package tsp

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	floats "github.com/gonum/floats"
)

type coordinatesList [][2]float64
type distancesMatrix [][]float64

// Type Map is specification of TSP dataset, with city'es coordinates
// and distances between all of them, also contains amount of cities
type Map struct {
	size        int
	coordinates coordinatesList
	distances   distancesMatrix
}

func readXY(scan *bufio.Scanner) (coord [2]float64, err error) {
	scan.Scan()
	text := scan.Text()
	for text == "" {
		scan.Scan()
		text = scan.Text()
	}
	fields := strings.Fields(text)
	if l := len(fields); l != 3 {
		return coord, fmt.Errorf("coordinates have to be three values: `id x y`, found %d", l)
	}
	coord[0], err = strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return coord, err
	}
	coord[1], err = strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return coord, err
	}
	return coord, nil
}

// ParseFrom parses Map using provided Scanner, first line have to be number of cities,
// next lines have to contain index with two coordinates separated by whitespace.
func ParseFrom(scanner *bufio.Scanner) (m *Map) {
	scanner.Scan()
	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}
	m = &Map{
		size:        n,
		coordinates: make(coordinatesList, n),
		distances:   make(distancesMatrix, n),
	}

	for i := range m.coordinates {
		m.coordinates[i], err = readXY(scanner)
		if err != nil {
			panic(err)
		}
	}

	for i := range m.distances {
		m.distances[i] = make([]float64, n)
		for j := range m.distances[i] {
			if i != j {
				m.distances[i][j] = math.Hypot(
					m.coordinates[j][0]-m.coordinates[i][0],
					m.coordinates[j][1]-m.coordinates[i][1])
			} else {
				m.distances[i][j] = math.Inf(1)
			}
		}
	}

	return m
}

// Size returns amount of cities on map
func (m *Map) Size() int {
	return m.size
}

// Coordinates returns coordinates of city i
func (m *Map) Coordinates(i int) (x, y float64) {
	return m.coordinates[i][0], m.coordinates[i][1]
}

// Distance returns distance between cities i and j
func (m *Map) Distance(i, j int) float64 {
	return m.distances[i][j]
}

func (m *Map) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("size = %d\n", m.size))
	for i, c := range m.coordinates {
		buffer.WriteString(fmt.Sprintf("%5d: [%10.3f, %10.3f]\n", i+1, c[0], c[1]))
	}
	for i := range m.distances {
		for _, d := range m.distances[i] {
			buffer.WriteString(fmt.Sprintf("%10.2f", d))
		}
		buffer.WriteString("\n")
	}

	return buffer.String()
}

// Type Route stores route described by cities as next nodes (starts and ends in first city)
// and route travel cost.
type Route struct {
	Cost   float64
	Cities []int
}

// GetCopy returns separate copy of r.
func (r *Route) GetCopy() (copy *Route) {
	copy = &Route{
		Cost:   0.0,
		Cities: make([]int, len(r.Cities)),
	}
	CopyRoute(copy, r)

	return copy
}

// UpdateCost recalculates Cities route cost
func (r *Route) UpdateCost(m *Map) {
	r.Cost = 0.0
	for i, c := range r.Cities[0 : len(r.Cities)-1] {
		r.Cost += m.distances[c][r.Cities[i+1]]
	}
}

// String returns string representation of Route
func (r *Route) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("cost = %.5f\n[", r.Cost))
	for _, c := range r.Cities[0 : len(r.Cities)-1] {
		buffer.WriteString(fmt.Sprintf("%d ", c+1))
	}
	buffer.WriteString(fmt.Sprintf("%d]\n", 1))

	return buffer.String()
}

// InitialRoute returns initial route on Map m determined using GRASP algorithm,
// it starts and ends in first city, have properly calculated Cost and is zero indexed.
// However `func (r *Route) String() string` returns string with one indexed Cities.
func InitialRoute(m *Map) (r *Route) {
	r = &Route{
		Cost:   0.0,
		Cities: make([]int, m.size+1),
	}
	r.Cities[0] = 0
	r.Cities[m.size] = 0

	d := make(distancesMatrix, m.size)
	for i := range d {
		d[i] = make([]float64, m.size)
		copy(d[i], m.distances[i])
	}

	k := 0
	for i := 1; i < m.size; i++ {
		r.Cities[i] = floats.MinIdx(d[k])
		k = r.Cities[i]
		d[k][0] = math.Inf(1)
		for v := 0; v <= i; v++ {
			d[k][r.Cities[v]] = math.Inf(1)
		}
		r.Cost += m.distances[r.Cities[i-1]][k]
	}

	r.Cost += m.distances[k][0]

	return r
}

// CopyRoute copies r to dst.
// dst have to be size of r, otherwise it will panic.
func CopyRoute(dst *Route, r *Route) {
	if len(r.Cities) != len(dst.Cities) {
		panic("Routes r and dst have to be the same size.")
	}
	dst.Cost = r.Cost
	copy(dst.Cities, r.Cities)
}
