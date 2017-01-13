package main

import (
	"flag"
	"github.com/DanShu93/golang-playground/brainfuck/interpreter"
	"strconv"
)

func main() {
	flag.Parse()
	var input = flag.Args()
	size, _ := strconv.Atoi(input[0])
	interpreter.Run(input[1], size)
}
