package main

import (
	"day17/internal/aocUtils"
	"day17/internal/day17"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {

	////////////////////////////////////////////////////////////////////
	// ENVIRONMENT SETUP
	////////////////////////////////////////////////////////////////////

	//os.Setenv("DEBUG", "true")
	INPUT_FILE := os.Getenv("INPUT_FILE")
	OUTPUT_FILE := os.Getenv("OUTPUT_FILE")
	PARALLELISM, err := strconv.Atoi(os.Getenv("PARALLELISM"))
	if PARALLELISM < 1 || err != nil {
		PARALLELISM = 1
	}
	fmt.Printf("PARALLELISM: %d\n\n", PARALLELISM)

	if INPUT_FILE == "" || OUTPUT_FILE == "" {
		fmt.Println("INPUT_FILE and OUTPUT_FILE environment variables not set")
		fmt.Println("Defaulting to input.txt and output.txt")
		INPUT_FILE = "input.txt"
		OUTPUT_FILE = "output.txt"
	}

	////////////////////////////////////////////////////////////////////
	// READ INPUT FILE
	////////////////////////////////////////////////////////////////////

	lines, err := aocUtils.ReadInput(INPUT_FILE)
	if err != nil {
		fmt.Printf("Error reading from %s: %v", INPUT_FILE, err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// SOLUTION LOGIC
	////////////////////////////////////////////////////////////////////

	input, err := parseLines(lines)
	if err != nil {
		fmt.Println("Error parsing input:", err)
		return
	}
	results, err := solve(input)
	if err != nil {
		fmt.Println("Error solving 1:", err)
		return
	}

	////////////////////////////////////////////////////////////////////
	// WRITE OUTPUT FILE
	////////////////////////////////////////////////////////////////////

	err = aocUtils.WriteToFile(OUTPUT_FILE, results)
	if err != nil {
		fmt.Printf("Error writing to %s: %v", OUTPUT_FILE, err)
		return
	}

	fmt.Printf("Successfully processed %s and created %s\n", INPUT_FILE, OUTPUT_FILE)
}

func parseLines(lines []string) (*day17.Computer, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	fmt.Println("Parsing Input...")

	if len(lines) == 0 {
		return nil, fmt.Errorf("input is empty")
	}

	computer := day17.NewComputer()

	regAString := strings.TrimSpace(lines[0])
	regAString = strings.TrimPrefix(regAString, "Register A: ")
	regA, ok := big.NewInt(0).SetString(regAString, 10)
	if !ok {
		return nil, fmt.Errorf("error parsing register A: %s", regAString)
	}
	computer.SetRegisterA(regA)

	regBString := strings.TrimSpace(lines[1])
	regBString = strings.TrimPrefix(regBString, "Register B: ")
	regB, ok := big.NewInt(0).SetString(regBString, 10)
	if !ok {
		return nil, fmt.Errorf("error parsing register B: %s", regBString)
	}
	computer.SetRegisterB(regB)

	regCString := strings.TrimSpace(lines[2])
	regCString = strings.TrimPrefix(regCString, "Register C: ")
	regC, ok := big.NewInt(0).SetString(regCString, 10)
	if !ok {
		return nil, fmt.Errorf("error parsing register C: %s", regCString)
	}
	computer.SetRegisterC(regC)

	if DEBUG {
		fmt.Printf("Parsed Registers as:\n")
		fmt.Printf("    (A) String: %-20s ParsedValue: %-20s Set Value: %s\n", regAString, regA, computer.GetRegisterA().String())
		fmt.Printf("    (B) String: %-20s ParsedValue: %-20s Set Value: %s\n", regBString, regB, computer.GetRegisterB().String())
		fmt.Printf("    (C) String: %-20s ParsedValue: %-20s Set Value: %s\n", regCString, regC, computer.GetRegisterC().String())
		fmt.Println()
	}

	// lines[3] is a blank line

	opcodes := make([]day17.Opcode, 0)
	programString := strings.TrimSpace(lines[4])
	programString = strings.TrimPrefix(programString, "Program: ")
	programParts := strings.Split(programString, ",")
	for _, part := range programParts {
		opcode, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("error parsing opcode: %s", part)
		}
		opcodes = append(opcodes, day17.Opcode(opcode))
	}

	err := computer.SetOpcodes(opcodes)
	if err != nil {
		return nil, fmt.Errorf("error setting opcodes on new computer: %v", err)
	}

	if DEBUG {
		fmt.Printf("Parsed Program as: %v\n", opcodes)
		fmt.Println()
	}

	return computer, nil
}

func solve(comp *day17.Computer) ([]string, error) {
	fmt.Println("Beginning single-threaded solve")
	results := findRangesBFS(comp.GetOpcodes())
	for i, r := range results {
		fmt.Printf("Result %d: %s\n", i, r.String())
	}

	return nil, nil
}

type Range struct {
	Match        string
	Start        *big.Int
	End          *big.Int
	Index        int
	OutputLength int
}

func (r Range) String() string {
	return fmt.Sprintf("{Start: %s, End: %s, Index: %d, Match: %s, OutputLength: %d}", r.Start.String(), r.End.String(), r.Index, r.Match, r.OutputLength)
}

type Increment struct {
	Input  *big.Int
	Output string
}

func findRanges(c *day17.Computer, r Range) []Range {
	DEBUG := os.Getenv("DEBUG") == "true"
	eight := big.NewInt(8)
	increment := big.NewInt(0).Exp(eight, big.NewInt(int64(r.Index)), nil)
	if DEBUG {
		fmt.Printf("Start: %s, End: %s, Increment: %s(8^%d), Match: %s\n", r.Start.String(), r.End.String(), increment.String(), r.Index, r.Match)
	}

	increments := []Increment{}
	// Collect the increments
	for i := r.End; i.Cmp(r.Start) >= 0; i = big.NewInt(0).Sub(i, increment) {
		cloneComp := c.Clone()
		cloneComp.SetRegisterA(i)
		output, err := SolveComputer(0, cloneComp)
		if err != nil {
			fmt.Printf("Error solving with value %s: %v\n", i.String(), err)
			continue
		}
		outputWithoutCommas := strings.ReplaceAll(output, ",", "")
		if len(outputWithoutCommas) < r.OutputLength {
			break
		}
		increments = append(increments, Increment{Input: i, Output: outputWithoutCommas})

		if DEBUG {
			octalI := fmt.Sprintf("%s", i.Text(8))
			outputString := fmt.Sprintf("Value: %-20s Octal: %-20s Output: %-20s", i.String(), octalI, output)
			// If i is divisible by 8
			if big.NewInt(0).Mod(i, big.NewInt(8)).Cmp(big.NewInt(0)) == 0 {
				outputString = fmt.Sprintf("\033[32m%s\033[0m", outputString)
			}
			fmt.Println(outputString)
		}
	}

	// Check the increments for matches
	ranges := []Range{}
	var rangeStart *big.Int = nil
	var rangeEnd *big.Int = nil
	for i, inc := range increments {
		compareString := inc.Output[r.Index-1:]
		if compareString == r.Match && rangeEnd == nil {
			fmt.Printf("End found at %d\n", i)
			if len(r.Match) == len(inc.Output) {
				rangeEnd = big.NewInt(0).Set(increments[i].Input)
				rangeStart = big.NewInt(0).Set(increments[i].Input)
			} else if i == 0 {
				rangeEnd = big.NewInt(0).Set(r.End)
			} else {
				rangeEnd = big.NewInt(0).Set(increments[i-1].Input)
			}
		}
		if compareString != r.Match && rangeEnd != nil {
			fmt.Printf("Start found: %s\n", inc)
			rangeStart = big.NewInt(0).Set(inc.Input)
		}
		if rangeStart != nil && rangeEnd != nil {
			ranges = append(ranges, Range{Start: rangeStart, End: rangeEnd, Index: r.Index - 1, Match: r.Match, OutputLength: r.OutputLength})
			rangeStart = nil
			rangeEnd = nil
		}
	}
	if rangeEnd != nil && rangeStart == nil {
		fmt.Printf("No start found, setting start to start of range: %s\n", r.Start)
		rangeStart = r.Start
		ranges = append(ranges, Range{Start: rangeStart, End: rangeEnd, Index: r.Index - 1, Match: r.Match, OutputLength: r.OutputLength})
	}

	return ranges
}

func findRangesBFS(opCodes []day17.Opcode) []*big.Int {
	// Initialize the result array for complete matches
	var completeMatches []*big.Int

	fmt.Println("Initializing computer with opcodes")
	comp := day17.NewComputer()
	comp.SetOpcodes(opCodes)

	index := len(opCodes)
	fmt.Printf("Opcode length: %d\n", index)

	initialEnd := big.NewInt(0).Exp(big.NewInt(8), big.NewInt(int64(index)), nil)
	initialEnd = big.NewInt(0).Sub(initialEnd, big.NewInt(8))
	fmt.Printf("Initial range: 0 to %s\n", initialEnd.String())

	// Initialize the heap (queue) with the initial range
	initialRange := Range{
		Start:        big.NewInt(0),
		End:          initialEnd,
		Index:        index,
		Match:        fmt.Sprintf("%d", opCodes[index-1]),
		OutputLength: len(opCodes),
	}
	queue := []Range{initialRange}
	fmt.Println("Initial range added to the queue")

	// Perform BFS
	for len(queue) > 0 {
		fmt.Printf("Queue size: %d\n", len(queue))
		// Dequeue the next range to process
		currentRange := queue[0]
		queue = queue[1:]

		currentRange.Match = ""
		for _, opcode := range opCodes[currentRange.Index-1:] {
			currentRange.Match += fmt.Sprintf("%d", opcode)
		}

		fmt.Printf("Processing range: %s\n", currentRange)

		// Call findRanges to get new ranges from the current range
		newRanges := findRanges(comp, currentRange)
		fmt.Printf("Found %d new ranges\n", len(newRanges))

		// Process each range found
		for _, r := range newRanges {
			fmt.Printf("Processing new range: %s\n", r)
			// Check if this range is a complete match
			if r.Start.Cmp(r.End) == 0 {
				fmt.Printf("Found complete match: %s\n", r.Start.String())
				// Add to complete matches
				completeMatches = append(completeMatches, r.Start)
			} else {
				// Enqueue for further processing
				fmt.Println("Enqueuing new range for processing")
				queue = append(queue, r)
			}
		}
	}

	fmt.Printf("Total complete matches found: %d\n", len(completeMatches))
	return completeMatches
}

func SolveComputer(workerId int, comp *day17.Computer) (string, error) {
	// DEBUG := os.Getenv("DEBUG") == "true"
	// if DEBUG {
	// 	fmt.Printf("Worker %d: Beginning solve\n", workerId)
	// 	fmt.Printf("Worker %d: Initial State: %s\n", workerId, comp)
	// }
	var loopDetection = 0

	// Get the opcodes from the computer
	opcodes := comp.GetOpcodes()
	// Get the instruction pointer from the computer
	ip := comp.GetInstructionPointer()

	var workerOutput string
	var workerWg sync.WaitGroup
	workerWg.Add(1)
	go func() {
		defer workerWg.Done()
		for out := range comp.Output {
			workerOutput = workerOutput + fmt.Sprintf("%s,", out)
		}
		// Remove the trailing comma
		if len(workerOutput) > 0 {
			workerOutput = strings.TrimSuffix(workerOutput, ",")
		}
	}()

	for ip < len(opcodes) {
		// Get the opcode at the instruction pointer
		opcode := opcodes[ip]

		// if DEBUG {
		// 	fmt.Printf("Worker %d: Executing", workerId)
		// 	fmt.Printf("Worker %d:     Computer State: %s\n", workerId, comp)
		// 	fmt.Printf("Worker %d:     Instruction pointer %d\n", workerId, ip)
		// 	fmt.Printf("Worker %d:     Opcode %d\n", workerId, opcode)
		// 	fmt.Printf("Worker %d:     Operand: %d\n", workerId, opcodes[ip+1])
		// }

		// Get the function associated with the opcode
		fn, err := opcode.GetInstruction()
		if err != nil {
			return "", fmt.Errorf("error in Worker %d: getting instruction for opcode %d at instruction pointer %d: %v", workerId, opcode, ip, err)
		}

		// if DEBUG {
		// 	fmt.Printf("Worker %d:     Executing Function\n", workerId)
		// }

		// Execute the function
		err = fn(comp, opcodes[ip+1])
		if err != nil {
			return "", fmt.Errorf("error in Worker %d: executing opcode %d: %v", workerId, opcode, err)
		}

		// if DEBUG {
		// 	fmt.Printf("Worker %d:     Fetching New Instruction Pointer\n", workerId)
		// }

		// Get the instruction pointer from the computer
		newIp := comp.GetInstructionPointer()
		// if DEBUG {
		// 	fmt.Printf("Worker %d:     New Instruction Pointer: %d\n", workerId, newIp)
		// }
		if newIp == ip {
			loopDetection++
			if loopDetection > 10 {
				return "", fmt.Errorf("Worker %d: loop detected", workerId)
			}
		} else {
			loopDetection = 0
		}
		ip = newIp
		if ip < 0 || ip >= len(opcodes) {
			close(comp.Output)
		}
	}
	workerWg.Wait()
	// if DEBUG {
	// 	fmt.Printf("Worker %d: Solve complete", workerId)
	// 	fmt.Printf("Worker %d: Final State: %s\n", workerId, comp)
	// 	fmt.Printf("Worker %d: Output: %s\n", workerId, workerOutput)
	// }
	return workerOutput, nil
}
