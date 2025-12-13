package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	path := "./input.txt"
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(data, []byte("\n"))
	machines := []Machine{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		machine, err := parseLine(line)
		if err != nil {
			panic(err)
		}

		machines = append(machines, machine)
	}

	totalPresses := 0
	for i, machine := range machines {
		fmt.Printf("processing machine #%d\n", i+1)
		presses := machine.Configure()
		totalPresses += presses
	}

	elapsed := time.Since(start).Seconds()
	fmt.Printf("total presses to configure all machines: %d\n", totalPresses)
	fmt.Printf("time elapsed: %vs\n", elapsed)
}

type Machine struct {
	Indicator        Indicator
	Buttons          []Button
	Joltage          Joltage
	Cache            map[CacheKey]int
	CounterToButtons [][]int
}

type CacheKey struct {
	i    int
	mask uint64
}

func (m *Machine) Configure() int {
	// Ensure dp cache populated
	_ = m.dp(0, 0)

	minParity := m.dp(0, 0)
	if minParity > len(m.Buttons) {
		panic("no parity solution")
	}

	bestTotal := math.MaxInt

	// Small slack is usually enough; bump if needed.
	// You can start with 0,2,4,... and stop early using bestTotal.
	for slack := 0; slack <= len(m.Buttons); slack += 2 {
		maxParity := minParity + slack

		ps := m.allParityVectorsUpTo(maxParity)

		// Try cheaper parity vectors first (helps bound early)
		sort.Slice(ps, func(i, j int) bool {
			return m.parityCost(ps[i]) < m.parityCost(ps[j])
		})

		for _, p := range ps {
			pCost := m.parityCost(p)

			// Early prune: even with kCost=0 we'd be worse
			if pCost >= bestTotal {
				continue
			}

			r := m.remainingVector(p)

			kCost, ok := m.solveK(r)
			if !ok {
				continue
			}

			total := pCost + 2*kCost
			if total < bestTotal {
				bestTotal = total
			}
		}

		// Early stop: next slack increases parity by at least 2 presses.
		// If we already have a solution with total == min possible for this slack,
		// further slack cannot beat it.
		if bestTotal != math.MaxInt && minParity+slack >= bestTotal {
			break
		}
	}

	if bestTotal == math.MaxInt {
		panic("no solution for machine")
	}
	return bestTotal
}

func (m *Machine) dp(i int, mask uint64) int {
	if v, ok := m.Cache[CacheKey{i, mask}]; ok {
		return v
	}

	if mask == m.Joltage.ParityMask {
		m.Cache[CacheKey{i, mask}] = 0
		return 0
	}

	if i == len(m.Buttons) {
		return len(m.Buttons) + 1
	}

	skip := m.dp(i+1, mask)
	use := 1 + m.dp(i+1, m.Buttons[i].IndicatorPress(mask))

	result := min(skip, use)
	m.Cache[CacheKey{i, mask}] = result
	return result
}

func (m *Machine) allParityVectorsUpTo(maxCost int) [][]int {
	// Ensure DP cache populated
	_ = m.dp(0, 0)

	n := len(m.Buttons)
	curr := make([]int, n)
	var out [][]int

	var rec func(i int, mask uint64)
	rec = func(i int, mask uint64) {
		// If already at target, force tail to 0 and record solution.
		if mask == m.Joltage.ParityMask {
			sol := make([]int, n)
			copy(sol, curr)
			for t := i; t < n; t++ {
				sol[t] = 0
			}
			out = append(out, sol)
			return
		}

		if i == n {
			return
		}

		// Prune: if even the optimal dp from here exceeds maxCost, stop.
		if m.dp(i, mask) > maxCost {
			return
		}

		// Option 1: skip
		curr[i] = 0
		if m.dp(i+1, mask) <= maxCost {
			rec(i+1, mask)
		}

		// Option 2: use
		nextMask := m.Buttons[i].IndicatorPress(mask)
		curr[i] = 1
		if 1+m.dp(i+1, nextMask) <= maxCost {
			rec(i+1, nextMask)
		}

		curr[i] = 0
	}

	rec(0, 0)
	return out
}

func (m *Machine) parityCost(p []int) int {
	c := 0
	for _, v := range p {
		c += v
	}
	return c
}

func (m *Machine) solveK(r []int) (int, bool) {
	assigned := make([]bool, len(m.Buttons))

	state := SolverState{
		r:        append([]int(nil), r...),
		upper:    m.upperVector(r, assigned),
		assigned: assigned,
		cost:     0,
	}

	bestCost := math.MaxInt
	search(m, state, &bestCost)

	if bestCost == math.MaxInt {
		return 0, false
	}
	return bestCost, true
}

func (m *Machine) remainingVector(p []int) []int {
	c := make([]int, len(m.Joltage.Values))
	for i, b := range m.Buttons {
		if p[i] == 1 {
			for _, idx := range b.Counters {
				c[idx]++
			}
		}
	}

	r := make([]int, len(m.Joltage.Values))
	for i := range r {
		diff := m.Joltage.Values[i] - c[i]
		if diff < 0 || diff%2 != 0 {
			panic(fmt.Errorf("invalid parity choice: %+v", p))
		}

		r[i] = diff / 2
	}

	return r
}

func (m *Machine) upperVector(r []int, assigned []bool) []int {
	u := make([]int, len(m.Buttons))

	for i, b := range m.Buttons {
		if len(b.Counters) == 0 {
			u[i] = 0
			continue
		}

		if assigned[i] {
			u[i] = 0
			continue
		}

		ub := math.MaxInt
		for _, ctr := range b.Counters {
			ub = min(ub, r[ctr])
		}

		u[i] = ub
	}

	return u
}

type SolverState struct {
	r        []int
	upper    []int
	assigned []bool
	cost     int
}

func (s *SolverState) assign(m *Machine, j int, v int) bool {
	for _, idx := range m.Buttons[j].Counters {
		s.r[idx] -= v
		if s.r[idx] < 0 {
			return false
		}
	}

	s.cost += v
	s.assigned[j] = true

	s.upper = m.upperVector(s.r, s.assigned)

	return true
}

func (s *SolverState) solved() bool {
	for _, need := range s.r {
		if need != 0 {
			return false
		}
	}

	return true
}

func (s *SolverState) chooseButton() int {
	best := -1
	bestUB := math.MaxInt

	for j := range s.upper {
		if s.assigned[j] {
			continue
		}

		ub := s.upper[j]
		if ub <= 0 {
			continue
		}

		if ub < bestUB {
			bestUB = ub
			best = j
		}
	}

	return best
}

func (s *SolverState) cloneState() SolverState {
	r2 := make([]int, len(s.r))
	copy(r2, s.r)

	u2 := make([]int, len(s.upper))
	copy(u2, s.upper)

	a2 := make([]bool, len(s.assigned))
	copy(a2, s.assigned)

	return SolverState{
		r:        r2,
		upper:    u2,
		assigned: a2,
		cost:     s.cost,
	}
}

func propagate(m *Machine, s *SolverState) bool {
	for {
		changed := false

		for i, need := range s.r {
			if need == 0 {
				continue
			}

			feasible := -1
			for _, j := range m.CounterToButtons[i] {
				if !s.assigned[j] && s.upper[j] > 0 {
					if feasible == -1 {
						feasible = j
					} else {
						feasible = -2
						break
					}
				}
			}

			if feasible == -1 {
				return false
			}

			if feasible >= 0 {
				v := need
				if v > s.upper[feasible] {
					return false
				}
				if !s.assign(m, feasible, v) {
					return false
				}
				changed = true
				break
			}
		}

		if !changed {
			return true
		}
	}
}

func search(m *Machine, s SolverState, bestCost *int) {
	if s.cost >= *bestCost {
		return
	}

	if !propagate(m, &s) {
		return
	}

	if s.solved() {
		if s.cost < *bestCost {
			*bestCost = s.cost
		}
		return
	}

	j := s.chooseButton()
	if j == -1 {
		return
	}

	for v := 0; v <= s.upper[j]; v++ {
		next := s.cloneState()
		if !next.assign(m, j, v) {
			continue
		}
		if next.cost >= *bestCost {
			continue
		}
		search(m, next, bestCost)
	}
}

type Indicator struct {
	BinaryString string
	BinaryValue  uint64
}

func (i Indicator) Length() int {
	return len(i.BinaryString)
}

type Button struct {
	BinaryValue uint64
	Counters    []int
}

func (b Button) IndicatorPress(input uint64) uint64 {
	return input ^ b.BinaryValue
}

type Joltage struct {
	Values     []int
	ParityMask uint64
}

func parseLine(input []byte) (Machine, error) {
	tokens := bytes.Split(input, []byte(" "))

	n := len(tokens)
	indicatorToken := tokens[0]
	buttonTokens := tokens[1 : n-1]
	joltageToken := tokens[n-1]

	indicator, err := parseIndicatorToken(indicatorToken)
	if err != nil {
		return Machine{}, err
	}

	buttons := []Button{}
	for _, buttonToken := range buttonTokens {
		button, err := parseButtonToken(buttonToken)
		if err != nil {
			return Machine{}, err
		}

		buttons = append(buttons, button)
	}

	joltage, err := parseJoltageToken(joltageToken)
	if err != nil {
		return Machine{}, err
	}

	cache := make(map[CacheKey]int)

	counterToButtons := make([][]int, len(joltage.Values))
	for i, b := range buttons {
		for _, idx := range b.Counters {
			counterToButtons[idx] = append(counterToButtons[idx], i)
		}
	}

	return Machine{indicator, buttons, joltage, cache, counterToButtons}, nil
}

func parseIndicatorToken(token []byte) (Indicator, error) {
	n := len(token)
	if token[0] != '[' || token[n-1] != ']' {
		return Indicator{}, fmt.Errorf("invalid indicator token: %s", string(token))
	}

	var b strings.Builder
	for _, c := range token {
		switch c {
		case '.':
			b.WriteString("0")
		case '#':
			b.WriteString("1")
		default:
			continue
		}
	}

	binaryString := b.String()
	binaryVal, err := strconv.ParseUint(binaryString, 2, 0)
	if err != nil {
		return Indicator{}, err
	}

	return Indicator{binaryString, binaryVal}, nil
}

func parseButtonToken(token []byte) (Button, error) {
	n := len(token)
	if token[0] != '(' || token[n-1] != ')' {
		return Button{}, fmt.Errorf("invalid button token: %s", string(token))
	}

	indexes := []int{}
	for _, c := range token {
		if c >= '0' && c <= '9' {
			v, err := strconv.Atoi(string(c))
			if err != nil {
				return Button{}, err
			}

			indexes = append(indexes, v)
		}
	}

	sort.Ints(indexes)

	binaryVal := uint64(0)
	for _, idx := range indexes {
		binaryVal |= uint64(1) << uint(idx)
	}

	return Button{BinaryValue: binaryVal, Counters: indexes}, nil
}

func parseJoltageToken(token []byte) (Joltage, error) {
	n := len(token)
	if token[0] != '{' || token[n-1] != '}' {
		return Joltage{}, fmt.Errorf("invalid joltage token: %s", string(token))
	}

	token = token[1 : n-1]
	valBytes := bytes.Split(token, []byte(","))

	values := make([]int, len(valBytes))
	for i, b := range valBytes {
		v, err := strconv.Atoi(string(b))
		if err != nil {
			return Joltage{}, err
		}

		values[i] = v
	}

	parity := uint64(0)
	for i, v := range values {
		if v&1 == 1 {
			parity |= 1 << i
		}
	}

	return Joltage{Values: values, ParityMask: parity}, nil
}

func sum(input []int) int {
	s := 0
	for _, v := range input {
		s += v
	}

	return s
}
