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
	iterations        = 1000000 // how many times to iterate
	meltingpFraction  = 0.7     // fraction of melting point for starting temperature
	targetTemperature = 0.01    // fraction of melting point for ending temperature
	stagnantFactor    = 0.05    // fraction of stagnant iterations allowed before reheating
)

var (
	k            int
	n            int
	meltingPoint float64
	t0, t        float64
	route        *tsp.Route
)
var m *tsp.Map
var optimum *tsp.Route

func iterate(t float64, x *tsp.Route) *tsp.Route {
	y := x.GetCopy()
again:
	t1 := rand.Intn(n-1) + 1
	t2 := (t1+rand.Intn(k))%(n-2) + 1
	if t2-t1 <= 1 {
		goto again
	}

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
again:
	t1 := rand.Intn(n-1) + 1
	t2 := (t1+rand.Intn(k))%(n-2) + 1
	if t2-t1 <= 1 {
		goto again
	}

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
	k = int(math.Max(0.05*float64(n), 2.0))

	route = tsp.InitialRoute(m)

	optimum = route.GetCopy()
	initialCost := route.Cost
	log.Printf("Initial cost: %.5f\n", route.Cost)

	testRoute := route.GetCopy()
	var minf, maxf float64
	var i int
	for i, maxf, minf = 0, 0, math.Pow(10.0, 10.0); i < int(math.Max(0.01*float64(iterations), 2.0)); i++ {
		testRoute = iterateP(0.001, testRoute)
		if testRoute.Cost < minf {
			minf = testRoute.Cost
		}
		if testRoute.Cost > maxf {
			maxf = testRoute.Cost
		}
	}

	meltingPoint = (maxf - minf) * math.Pow(10.0, -1.0*math.Log10(float64(n)))

	t0 = meltingPoint * meltingpFraction
	//t0 = 100
	log.Printf("T0 = %.2f\n", t0)
	t = t0

	tDecay := math.Pow(targetTemperature, 1.0/iterations)

	var l, optItr = 1, 1
	for ; l < iterations; l++ {
		route = iterate(t, route)
		if route.Cost < optimum.Cost {
			optimum = route.GetCopy()
			optItr = l
		}

		if l-optItr == int(stagnantFactor*iterations) {
			optItr = l
			log.Println("reheating")
			t *= 1.0 + 0.5*((float64(iterations)-float64(l))/float64(iterations))
			//route = optimum.GetCopy()
		}

		t *= tDecay
		if l%(iterations/10) == 0 {
			log.Printf("T  = %.2f\n", t)
			log.Printf("cost:\t%.5f\n", route.Cost)
			log.Printf("optimum:\t%.5f\n", optimum.Cost)
		}
	}

	if !*debugFlag {
		optimum.PrintResult()
		return
	}
	log.Printf("T  = %.2f\n", t)
	log.Printf("Optimum: %.5f\n", optimum.Cost)
	log.Printf("%.2f%% improvement\n", 100.0-100*optimum.Cost/initialCost)
	log.Printf("after %.2f%% of iterations\n", float64(100*optItr)/float64(iterations))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
