package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"

	"./tsp"
)

var debugFlag = flag.Bool("d", false, "print debug info")

const (
	ITS        = 1000000
	PMELT      = 0.7
	TARGT      = 0.1
	STAGFACTOR = 0.01
)

var (
	K int
	N int
)
var m *tsp.Map

func iterate(temp float64, x *tsp.Route) *tsp.Route {
	var pt, sh int
	y := x.GetCopy()
	pt = rand.Int()%(N-1) + 1
	sh = (pt+(rand.Int()%K)+1)%(N-1) + 1
	y.Cities[pt] = y.Cities[pt] ^ y.Cities[sh]
	y.Cities[sh] = y.Cities[sh] ^ y.Cities[pt]
	y.Cities[pt] = y.Cities[pt] ^ y.Cities[sh]
	y.UpdateCost(m)
	if y.Cost <= x.Cost {
		return y
	} else if rand.Float64() < math.Exp((x.Cost-y.Cost)/temp) {
		return y
	} else {
		return x
	}
}

// http://www.bookstaber.com/david/opinion/SATSP.pdf

func main() {
	flag.Parse()
	var i, l, optgen int
	var minf, maxf, rangee, Dtemp float64
	var x, optimum *tsp.Route

	scanner := bufio.NewScanner(os.Stdin)
	m = tsp.ParseFrom(scanner)
	//debugf("Map: %v", m)

	N = m.Size()
	K = int(0.7 * float64(N))

	x = tsp.InitialRoute(m)
	debugf("Initial cost: %.5f\n", x.Cost)

	xCopy := x.GetCopy()
	for i, maxf, minf = 0, 0, math.Pow(10.0, 10.0); i < int(math.Max(0.01*float64(N), 2.0)); i++ {
		xCopy = iterate(math.Pow(10.0, 10.0), x)
		if xCopy.Cost < minf {
			minf = xCopy.Cost
		}
		if xCopy.Cost > maxf {
			maxf = xCopy.Cost
		}
	}

	rangee = (maxf - minf) * PMELT
	debugf("range = %f\n", rangee)
	Dtemp = math.Pow(TARGT, 1.0/ITS)
	debugf("Dtemp = %f\n", Dtemp)

	for l, optgen, optimum = 1, 1, x.GetCopy(); l < ITS; l++ {
		x = iterate(rangee*math.Pow(Dtemp, float64(l)), x)
		if x.Cost < optimum.Cost {
			optimum = x.GetCopy()
			optgen = l
		}

		if l-optgen == int(STAGFACTOR*ITS) {
			Dtemp = math.Pow(Dtemp, .05*float64(l)/float64(ITS))
		}
	}

	debugf("Optimum: %.5f\nFound in iteration: %d of %d\n", optimum.Cost, optgen, ITS)
}

func debug(s string) {
	if *debugFlag {
		fmt.Println(s)
	}
}

func debugf(format string, args ...interface{}) {
	if *debugFlag {
		fmt.Printf(format, args...)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
