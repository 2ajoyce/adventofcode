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

	"github.com/schollz/progressbar/v3"
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
	results, err := solve(input, PARALLELISM)
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

func solve(comp *day17.Computer, WORKER_COUNT int) ([]string, error) {
	fmt.Printf("Beginning solve with %d workers\n", WORKER_COUNT)

	// Calculate the expected output based on the initial state of the computer
	expectedOutput := ""
	opcodes := comp.GetOpcodes()
	for _, opcode := range opcodes {
		expectedOutput = expectedOutput + fmt.Sprintf("%d,", opcode)
	}
	expectedOutput = strings.TrimSuffix(expectedOutput, ",")

	testRegAValue := big.NewInt(216134799294990)
	var optimalRegAValue *big.Int = nil
	bar := progressbar.Default(testRegAValue.Int64(), "Solving")

	// Create a channel to distribute tasks to workers
	tasks := make(chan *big.Int, WORKER_COUNT)
	type inputOutput struct {
		input  string
		output string
	}
	results := make(chan inputOutput, WORKER_COUNT)

	// Read results from workers
	go func() {
		header := "input,output"
		err := aocUtils.AppendToFile("data.csv", []string{header})
		if err != nil {
			panic(fmt.Errorf("unexpected error: %v", err))
		}
		comparisonString := strings.ReplaceAll(expectedOutput, ",", "")
		aocUtils.AppendToFile("data.csv", []string{comparisonString})
		for result := range results {
			if len(result.output) < len(comparisonString) {
				break
			}
			if !strings.HasSuffix(result.output, "155530") {
				continue
			}
			data := fmt.Sprintf("%s,%s", result.input, result.output)
			err := aocUtils.AppendToFile("data.csv", []string{data})
			if err != nil {
				panic(fmt.Errorf("unexpected error: %v", err))
			}
		}
	}()

	var wg sync.WaitGroup
	// Launch worker goroutines
	for i := 0; i < WORKER_COUNT; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for regA := range tasks {
				cloneComp := comp.Clone()
				cloneComp.SetRegisterA(regA)

				workerOutput, err := SolveComputer(workerId, cloneComp)
				if err != nil {
					panic(fmt.Errorf("unexpected error: %v", err))
				}

				if optimalRegAValue != nil {
					break
				}
				if workerOutput == expectedOutput {
					fmt.Printf("\nWorker %d found optimal RegA value: %s\n", workerId, regA)
					optimalRegAValue = big.NewInt(0).Set(regA)
				}
				// Send the result to the results channel
				regAString := regA.String()
				workerOutput = strings.ReplaceAll(workerOutput, ",", "")
				results <- inputOutput{input: regAString, output: workerOutput}
			}
		}(i)
	}

	// Distribute tasks to workers
	tasksDistributed := 0
	go func() {
		fmt.Printf("Ignoring testRegAValue(%s) for now", testRegAValue)
		i := 1
		for optimalRegAValue == nil {
			increment := big.NewInt(0).Exp(big.NewInt(8), big.NewInt(int64(6)), nil)
			testRegAValue = big.NewInt(0).Sub(testRegAValue, increment)
			tasks <- testRegAValue
			bar.Add(1)
			tasksDistributed++
			i = i * 10
		}
		close(tasks)
	}()

	wg.Wait()
	output := fmt.Sprintf("Lowest RegA Value: %s", optimalRegAValue)
	fmt.Println(output)
	return []string{output}, nil
}

func SolveComputer(workerId int, comp *day17.Computer) (string, error) {
	DEBUG := os.Getenv("DEBUG") == "true"
	if DEBUG {
		fmt.Printf("Worker %d: Beginning solve\n", workerId)
		fmt.Printf("Worker %d: Initial State: %s\n", workerId, comp)
	}
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

		if DEBUG {
			fmt.Printf("Worker %d: Executing", workerId)
			fmt.Printf("Worker %d:     Computer State: %s\n", workerId, comp)
			fmt.Printf("Worker %d:     Instruction pointer %d\n", workerId, ip)
			fmt.Printf("Worker %d:     Opcode %d\n", workerId, opcode)
			fmt.Printf("Worker %d:     Operand: %d\n", workerId, opcodes[ip+1])
		}

		// Get the function associated with the opcode
		fn, err := opcode.GetInstruction()
		if err != nil {
			return "", fmt.Errorf("error in Worker %d: getting instruction for opcode %d at instruction pointer %d: %v", workerId, opcode, ip, err)
		}

		if DEBUG {
			fmt.Printf("Worker %d:     Executing Function\n", workerId)
		}

		// Execute the function
		err = fn(comp, opcodes[ip+1])
		if err != nil {
			return "", fmt.Errorf("error in Worker %d: executing opcode %d: %v", workerId, opcode, err)
		}

		if DEBUG {
			fmt.Printf("Worker %d:     Fetching New Instruction Pointer\n", workerId)
		}

		// Get the instruction pointer from the computer
		newIp := comp.GetInstructionPointer()
		if DEBUG {
			fmt.Printf("Worker %d:     New Instruction Pointer: %d\n", workerId, newIp)
		}
		if newIp == ip {
			loopDetection++
			if loopDetection > 10 {
				return "", fmt.Errorf("Worker %d: loop detected", workerId)
			}
		} else {
			loopDetection = 0
		}
		ip = newIp
	}

	close(comp.Output)
	workerWg.Wait()
	if DEBUG {
		fmt.Printf("Worker %d: Solve complete", workerId)
		fmt.Printf("Worker %d: Final State: %s\n", workerId, comp)
		fmt.Printf("Worker %d: Output: %s\n", workerId, workerOutput)
	}
	return workerOutput, nil
}