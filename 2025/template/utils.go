package main

import (
	"fmt"
	"strconv"
)

// Functions I don't want to redefine every problem
// Many of these swallow errors or panic, they're used here in
// the context of ridgidly defined inputs.

func StrToArrRune(s string) []rune {
	return []rune(s)
}

func ArrRuneToInt(r []rune) int {
	i, err := strconv.Atoi(string(r))
	if err != nil {
		panic(fmt.Sprintf("failed to convert rune slice %v to integer: %v", r, err))
	}
	return i
}

func IntToArrRune(i int) []rune {
	return []rune(strconv.Itoa(i))
}

func StrToInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Sprintf("failed to convert string %q to integer: %v", s, err))
	}
	return num
}

func ArrRuneToStr(r []rune) string {
	return string(r)
}
