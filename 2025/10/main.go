package main

import (
	"2ajoyce/adventofcode/2025/10/equation"
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// First Problem
	input := make(chan *equation.Equation)
	go ReadInput("input1.txt", input)
	result, err := Solve1(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)

	// Second Problem
	input = make(chan *equation.Equation)
	go ReadInput("input2.txt", input)
	result, err = Solve2(input)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
}

// ReadInput reads the input from the filepath and sends each line to the provided channel.
func ReadInput(filepath string, c chan *equation.Equation) {
	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		c <- ParseInput(line)
	}
	close(c)
}

// Split the input line into
func ParseInput(input string) *equation.Equation {
	// [.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
	eq := equation.NewEquation(strings.TrimSpace(input))
	return &eq
}

func Solve1(input chan *equation.Equation) (string, error) {
	total := 0

	// Collect all equations
	equations := []*equation.Equation{}
	for eq := range input {
		equations = append(equations, eq)
	}

	// Solve each equation
	for _, eq := range equations {
		minButtonPushes, err := minButtonPressesBFS(eq)
		if err != nil {
			return "", err
		}
		total += minButtonPushes
	}

	return fmt.Sprintf("%d", total), nil
}

// minButtonPressesBFS returns the minimum number of button presses to
// get from the start state (all zeros) to eq.Target using BFS or an
// error if the target is unreachable.
func minButtonPressesBFS(eq *equation.Equation) (int, error) {
	// The  initial state is all zeros
	// Since State is a uint16, all zeros is just 0
	var start equation.State = 0

	// If start is already the target return 0
	if start == eq.Target {
		return 0, nil
	}

	// The queue holds states to explore
	queue := make([]equation.State, 0)
	// visited tracks states which have been seen
	visited := make(map[equation.State]bool)
	// dist tracks the number of button presses to reach each state
	dist := make(map[equation.State]int)

	// Add the start state to the queue, mark it visited, and set its distance to 0
	queue = append(queue, start)
	visited[start] = true
	dist[start] = 0

	head := 0
	for head < len(queue) {
		// Dequeue the next state
		cur := queue[head]
		head++
		// Try pressing each button
		for _, btn := range eq.Buttons {
			next := cur.PressButton(btn)
			// If we reached the target, return the distance
			if next == eq.Target {
				return dist[cur] + 1, nil
			}
			// If already visited, skip
			if !visited[next] {
				// Mark as visisted and record distance
				visited[next] = true
				dist[next] = dist[cur] + 1
				// Enqueue the new state
				queue = append(queue, next)
			}
		}
	}

	return 0, fmt.Errorf("no solution for equation: %v", eq)
}

func Solve2(input chan *equation.Equation) (string, error) {
	total := 0

	// Collect all equations
	equations := []*equation.Equation{}
	for eq := range input {
		equations = append(equations, eq)
	}

	if len(equations) == 0 {
		return "0", nil
	}

	bar := progressbar.Default(int64(len(equations)))

	type result struct {
		val int
		err error
	}

	results := make(chan result, len(equations))

	const maxWorkers = 8
	sem := make(chan struct{}, maxWorkers) // semaphore to limit concurrency
	var wg sync.WaitGroup

	// Launch workers
	for _, eq := range equations {
		sem <- struct{}{} // acquire a slot
		wg.Add(1)

		go func(eq *equation.Equation) {
			defer wg.Done()
			defer func() { <-sem }() // release slot

			val, err := minButtonPressesVoltage(eq)
			results <- result{val: val, err: err}
		}(eq)
	}

	// Close results channel when all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	var firstErr error

	// Collect results
	for res := range results {
		bar.Add(1)

		if res.err != nil && firstErr == nil {
			firstErr = res.err
			// We don't early-return here to avoid leaving goroutines
			// running; we just remember the first error.
		}
		if res.err == nil {
			total += res.val
		}
	}

	if firstErr != nil {
		return "", firstErr
	}

	return fmt.Sprintf("%d", total), nil
}

func minButtonPressesVoltage(eq *equation.Equation) (int, error) {
	m := len(eq.TargetVoltage) // number of indices / rows
	n := len(eq.Buttons)       // number of buttons / columns

	// 1. Build matrix A and RHS vector V
	A := make([][]int, m) // size m x n
	V := make([]int, m)   // size m

	// For each index
	for i := range m {
		A[i] = make([]int, n)
		// For each button
		for j := range n {
			btn := eq.Buttons[j]
			// If the button affects index i set A[i][j] = 1
			if (btn>>i)&1 == 1 {
				A[i][j] = 1
			} else { // else set A[i][j] = 0
				A[i][j] = 0
			}
		}
		// Set V[i] equal to our target voltage at index i
		V[i] = int(eq.TargetVoltage[i])
	}

	// 2. Solve A * x = V over non-negative integers
	x, err := solveNonNegativeIntegerSystem(A, V)
	if err != nil {
		return 0, fmt.Errorf("no integer solution for equation: %v", err)
	}

	// 3. Compute total button presses from solution vector
	totalPresses := 0
	for j := 0; j < len(x); j++ {
		totalPresses += x[j]
	}

	return totalPresses, nil
}

// solveIntegerSystem solves the integer linear system A * x = v
// for non-negative integer vector x using Gaussian elimination with
// back substitution. It returns an error if no solution exists or
// if the solution is not integral/non-negative.
// This failed because it finds solutions with negative values
// I couldn't figure out how to enforce non-negativity constraints
// and moved on to meet the 24-hour deadline.
func solveIntegerSystemGaussian(A [][]int, v []int) ([]int, error) {
	var epsPivot float64 = 1e-9
	var epsInt float64 = 1e-6
	// Dimensions:
	//   m rows (equations)    = len(A)
	//   n columns (variables) = len(A[0]) if m > 0
	m := len(A)
	if m == 0 {
		return nil, fmt.Errorf("empty system")
	}
	n := len(A[0])
	if len(v) != m {
		return nil, fmt.Errorf("dimension mismatch: v has %d rows, A has %d", len(v), m)
	}

	// 1. Build augmented matrix M as float64
	//   M[i][0..n-1] = A[i][j]
	//   M[i][n]      = v[i]
	// M is basically the same as A, but it's easier to work with floats
	M := make([][]float64, m)
	for i := range m {
		if len(A[i]) != n {
			return nil, fmt.Errorf("row %d of A has length %d; want %d", i, len(A[i]), n)
		}
		M[i] = make([]float64, n+1)
		for j := range n {
			M[i][j] = float64(A[i][j])
		}
		M[i][n] = float64(v[i])
	}

	// 2. Forward elimination (row echelon form)
	// Begin pivot on row 0
	row := 0
	for col := 0; col < n && row < m; col++ {
		// Loop until we find a non-zero entry to use as pivot
		pivot := -1
		for r := row; r < m; r++ {
			if math.Abs(M[r][col]) > epsPivot {
				pivot = r
				break
			}
		}

		// If no pivot found in this column, move to next column (free variable)
		if pivot == -1 {
			continue
		}

		// If a new pivot was found, swap it into the current row
		if pivot != row {
			M[row], M[pivot] = M[pivot], M[row]
		}

		// Normalize the pivot row so that M[row][col] == 1
		// In other words, divide entire row by the pivot value
		pivotVal := M[row][col]
		for k := range n + 1 { // 0..n inclusive
			M[row][k] /= pivotVal
		}

		// Eliminate this column from all rows *below* the pivot row.
		for r := row + 1; r < m; r++ {
			// I found this part tricky to understand, these are the best explainations I found

			// Explaination 1
			// After we normalize the pivot row, M[row][col] is 1.
			// In some row r below the pivot row, the entry in that same column is M[r][col] = factor.
			// We want that entry to become 0, so that column is "cleaned" below the pivot.

			// Explaination 2
			// // factor is the current value in this row at the pivot column.
			// Subtracting factor * pivotRow will zero out the pivot column in this row.
			factor := M[r][col]
			if math.Abs(factor) < epsPivot {
				continue
			}
			for k := range n + 1 {
				M[r][k] -= factor * M[row][k]
			}
		}

		// Move to next pivot row.
		row++
	}

	// 3. Check for inconsistency
	// Any row with all zeros is fine
	// Any row with all zeros except last entry is inconsistent
	// An inconsistent system has no solutions
	for i := range m {
		allZero := true
		for j := range n {
			if math.Abs(M[i][j]) > epsPivot {
				allZero = false
				break
			}
		}
		if allZero && math.Abs(M[i][n]) > epsPivot {
			return nil, fmt.Errorf("inconsistent system at row %d", i)
		}
	}

	// 4. Back substitution to get one solution xs (floats)
	// xs will hold the solution
	xs := make([]float64, n)
	// xs[i] defaults to 0, don't need to initialize explicitly

	// Loop rows from bottom to top:
	for i := m - 1; i >= 0; i-- {
		// 1. Find the leftmost non-zero coefficient in this row.
		pivotCol := -1

		// Loop until we find a non-zero entry to use as pivot
		for col := range n {
			if math.Abs(M[i][col]) > epsPivot {
				pivotCol = col
				break
			}
		}
		// If we didn't find a pivot, skip this row.
		if pivotCol == -1 {
			continue
		}

		// accumulate every variable to the right of the pivot column
		// multiplied by its known value in xs
		var sum float64
		for k := pivotCol + 1; k < n; k++ {
			sum += M[i][k] * xs[k]
		}

		// All variables to the right of pivotCol (x_{pivotCol+1}…x_{n-1}) have
		// already been computed and stored in xs.
		// Variables left of pivotCol either don't appear or are are in rows above the pivot
		// and will be handled when we reach them (moving bottom to top)
		// This step moves the sum of known variables to the right side of the equation in
		// this row
		xs[pivotCol] = (M[i][n] - sum) / M[i][pivotCol]
	}

	// 5. Round floats to integer solution and enforce non-negativity
	// x will hold the final integer solution
	x := make([]int, n)
	// For each variable(column)
	for j := range n {
		// Round to nearest integer
		nearest := math.Round(xs[j])
		// if not close enough, error
		if math.Abs(xs[j]-nearest) > epsInt {
			return nil, fmt.Errorf("non-integer solution for variable %d: got %f", j, xs[j])
		}
		// if negative, error
		if nearest < 0 {
			return nil, fmt.Errorf("negative solution for variable %d: got %f", j, xs[j])
		}
		// Otherwise, set x[j]
		x[j] = int(nearest)
	}

	// 6. Verify A * x == v exactly
	for i := range m {
		sum := 0
		for j := range n {
			sum += A[i][j] * x[j]
		}
		if sum != v[i] {
			return nil, fmt.Errorf("solution check failed: row %d, got %d, want %d", i, sum, v[i])
		}
	}

	return x, nil
}

// solveNonNegativeIntegerSystem solves A * x = v with x_j >= 0 integers,
// and among all such x, returns one with minimal sum(x_j)
//
// A is m x n (0/1 entries), v is length m
// This uses a depth-first search with pruning (branch-and-bound),
// and reorders variables to reduce branching
func solveNonNegativeIntegerSystem(A [][]int, v []int) ([]int, error) {
	m := len(A)
	if m == 0 {
		return nil, fmt.Errorf("empty system")
	}
	n := len(A[0])
	if len(v) != m {
		return nil, fmt.Errorf("dimension mismatch: v has %d rows, A has %d", len(v), m)
	}
	for i := range m {
		if len(A[i]) != n {
			return nil, fmt.Errorf("row %d of A has length %d; want %d", i, len(A[i]), n)
		}
	}

	// 1. Compute simple upper bounds ub[j] for each variable x_j
	ub := make([]int, n)
	for j := range n {
		minVal := -1
		for i := range m {
			if A[i][j] == 1 {
				if minVal == -1 || v[i] < minVal {
					minVal = v[i]
				}
			}
		}
		if minVal == -1 {
			ub[j] = 0
		} else {
			ub[j] = minVal
		}
	}

	// 2. Choose a good variable order for DFS
	//    We want to assign "tighter" variables first:
	//      - smaller ub[j] (more constrained)
	//      - if tie, larger degree (touches more rows)
	type varInfo struct {
		idx    int
		ub     int
		degree int
	}
	vars := make([]varInfo, n)
	for j := range n {
		deg := 0
		for i := 0; i < m; i++ {
			if A[i][j] == 1 {
				deg++
			}
		}
		vars[j] = varInfo{idx: j, ub: ub[j], degree: deg}
	}

	sort.Slice(vars, func(a, b int) bool {
		if vars[a].ub != vars[b].ub {
			return vars[a].ub < vars[b].ub
		}
		// For equal ub, sort by higher degree first
		if vars[a].degree != vars[b].degree {
			return vars[a].degree > vars[b].degree
		}
		// Stable tie-breaker on original index
		return vars[a].idx < vars[b].idx
	})

	// order[d] = actual variable index j we assign at DFS depth d
	order := make([]int, n)
	for d := range n {
		order[d] = vars[d].idx
	}

	// 3. Precompute rowMaxRemainingDFS[d][i] = max extra we can add to row i
	//    using variables at DFS depths d..n-1 (i.e. variables order[d],...,order[n-1]),
	//    assuming each x_j <= ub[j]
	rowMaxRemainingDFS := make([][]int, n+1)
	for d := range n + 1 {
		rowMaxRemainingDFS[d] = make([]int, m)
	}

	// Base case: with no remaining variables, no extra capacity
	for i := range m {
		rowMaxRemainingDFS[n][i] = 0
	}

	// Fill backwards in DFS depth
	for d := n - 1; d >= 0; d-- {
		j := order[d] // actual variable index at this depth
		for i := range m {
			rowMaxRemainingDFS[d][i] = rowMaxRemainingDFS[d+1][i]
			if A[i][j] == 1 {
				rowMaxRemainingDFS[d][i] += ub[j]
			}
		}
	}

	// 4. DFS state
	rowSums := make([]int, m)  // current partial sums per row
	xCurrent := make([]int, n) // xCurrent[j] is current assignment for variable j (original index)

	bestX := make([]int, n)
	bestFound := false
	bestSum := 0

	// dfs(d, currentSum) - assign variable at DFS depth d
	var dfs func(d int, currentSum int)
	dfs = func(d int, currentSum int) {
		// If all variables assigned, check exact match
		if d == n {
			for i := range m {
				if rowSums[i] != v[i] {
					return
				}
			}
			if !bestFound || currentSum < bestSum {
				bestFound = true
				bestSum = currentSum
				copy(bestX, xCurrent)
			}
			return
		}

		j := order[d] // actual variable index we are assigning now

		// Early feasibility check: overshoot or impossible to reach target
		for i := range m {
			if rowSums[i] > v[i] {
				return
			}
			if rowSums[i]+rowMaxRemainingDFS[d][i] < v[i] {
				return
			}
		}

		// Compute max value for x_j based on remaining capacity in rows it touches
		maxX := ub[j]
		for i := range m {
			if A[i][j] == 1 {
				remaining := v[i] - rowSums[i]
				if remaining < 0 {
					return
				}
				if remaining < maxX {
					maxX = remaining
				}
			}
		}

		for xj := 0; xj <= maxX; xj++ {
			newSum := currentSum + xj
			if bestFound && newSum >= bestSum {
				// Can't beat current best
				continue
			}

			valid := true
			if xj > 0 {
				for i := range m {
					if A[i][j] == 1 {
						rowSums[i] += xj
						if rowSums[i] > v[i] {
							valid = false
						}
					}
				}
			}

			if valid {
				xCurrent[j] = xj
				dfs(d+1, newSum)
			}

			// undo rowSums
			if xj > 0 {
				for i := range m {
					if A[i][j] == 1 {
						rowSums[i] -= xj
					}
				}
			}
		}
	}

	// 5. Run DFS
	dfs(0, 0)

	if !bestFound {
		return nil, fmt.Errorf("no non-negative integer solution")
	}
	return bestX, nil
}
