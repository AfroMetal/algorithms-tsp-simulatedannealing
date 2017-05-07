package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"./tsp"
)

var debugFlag = flag.Bool("d", false, "print debug info")

const (
	minItrs     = 5000000 // how many times to iterate (minimum)
	meltPointF  = 0.7    // fraction of melting point for starting temperature
	targetTempF = 0.01   // fraction of melting point for ending temperature
	stagItrsF   = 0.1    // fraction of stagnant minItrs allowed before reheating
)

var (
	n, k      int
	meltPoint float64
	t0, t     float64
	route     *tsp.Route
)
var m *tsp.Map
var optimum *tsp.Route

func randRange(from, to int) int {
	if from > to {
		panic("Range can't be negative")
	}
	return rand.Intn(to-from+1) + from
}

func iterate(t float64, x *tsp.Route) *tsp.Route {
	y := x.GetCopy()
	
	t1 := randRange(1, n-3)
	t2 := randRange(t1+2, n-1)
	y.Swap(m, t1, t2)

	if y.Cost <= x.Cost {
		return y
	} else if rand.Float64() < math.Exp((x.Cost-y.Cost)/t) {
		return y
	} else {
		return x
	}
}

func iterateP(p float64, x *tsp.Route) *tsp.Route {
	y := x.GetCopy()
	
	t1 := randRange(1, n-3)
	t2 := randRange(t1+2, int(math.Min(float64(t1+k), float64(n-1))))

	y.Swap(m, t1, t2)

	if y.Cost <= x.Cost {
		return y
	} else if rand.Float64() < p {
		return y
	} else {
		return x
	}
}

// http://www.bookstaber.com/david/opinion/SATSP.pdf

func main() {
	flag.Parse()

	if !*debugFlag {
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetFlags(log.Lshortfile)
	}

	rand.Seed(time.Now().UnixNano())

	scanner := bufio.NewScanner(os.Stdin)
	m = tsp.ParseFrom(scanner)

	n = m.Size()
	k = int(math.Max(0.01*float64(n), 25.0))

	route = tsp.InitialRoute(m)

	optimum = route.GetCopy()
	initialCost := route.Cost
	log.Printf("Initial cost: %.5f\n", route.Cost)

	testRoute := route.GetCopy()
	var minT, maxT float64
	var i int

	for i, maxT, minT = 0, 0, math.Pow(10.0, 10.0); i < int(math.Max(0.01*float64(minItrs), 2.0)); i++ {
		testRoute = iterateP(0.001, testRoute)
		if testRoute.Cost < minT {
			minT = testRoute.Cost
		}
		if testRoute.Cost > maxT {
			maxT = testRoute.Cost
		}
	}

	meltPoint = (maxT - minT) * math.Pow(10.0, -1.0*math.Log10(float64(n)))

	t0 = meltPoint * meltPointF
	// t0 = 50
	log.Printf("T0 = %.2f\n", t0)
	t = t0

	tDecay := math.Pow(targetTempF, 1.0/minItrs)

	var itr, optItr, stagItrs = 1, 1, 0
	for ; t > t0*targetTempF; itr++ {
		route = iterate(t, route)
		if route.Cost < optimum.Cost {
			optimum = route.GetCopy()
			optItr = itr
		} else {
			stagItrs++
		}

		if stagItrs == int(stagItrsF*minItrs) {
			stagItrs = 0
			log.Println("reheating")
			t *= 1.0 + 0.7*((float64(minItrs)-float64(itr))/float64(minItrs))
		}

		t *= tDecay

		if *debugFlag && itr%(minItrs/5) == 0 {
			log.Printf("T  = %.2f\n", t)
			log.Printf("cost:\t%.5f\n", route.Cost)
			log.Printf("optimum:\t%.5f\n", optimum.Cost)
		}
	}

	if !*debugFlag {
		optimum.PrintResult()
	} else {
		log.Printf("T  = %.2f\n", t)
		log.Printf("Optimum: %.5f\n", optimum.Cost)
		log.Printf("%.2f%% improvement\n", 100.0-100*optimum.Cost/initialCost)
		log.Printf("after %.2f%% of iterations\n", float64(100*optItr)/float64(itr))
		log.Printf("done %.2f%% of minimum iterations\n", float64(100*itr)/float64(minItrs))
	}
}
