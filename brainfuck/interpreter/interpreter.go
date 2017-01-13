// Package interpreter implements a function for running brainfuck code.
// The following operators can be used in brainfuck code:
// >	move data pointer to the next cell on the right
// <	move data pointer to the next cell on the left
// +	increment value of current cell
// -	decrement value of current cell
// .	output the byte of a currently pointed cell in ASCII code
// ,	read one byte from stdin and store its value at the current cell
// [	if the current cell is 0 then jump to the matching ]
// ]	jump to the matching [
//      all other characters are ignored
package interpreter

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Run runs the given brainfuck code with a tape of the given size.
func Run(code string, tapeSize int) {
	tape := make([]byte, tapeSize, tapeSize)
	tapePos := 0
	nestingLevel := 0
	for codePos := 0; codePos < len(code); {
		operator := code[codePos]
		if nestingLevel == 0 {
			switch operator {
			case '>':
				if tapePos == tapeSize-1 {
					tapePos = 0
				} else {
					tapePos++
				}
			case '<':
				if tapePos == 0 {
					tapePos = tapeSize - 1
				} else {
					tapePos--
				}
			case '+':
				if tape[tapePos] == math.MaxUint8 {
					tape[tapePos] = 0
				} else {
					tape[tapePos]++
				}
			case '-':
				if tape[tapePos] == 0 {
					tape[tapePos] = math.MaxUint8
				} else {
					tape[tapePos]--
				}
			case '.':
				fmt.Print(fmt.Sprintf("%c", tape[tapePos]))
			case ',':
				givenByte := promptByte()
				if givenByte == '\n' {
					givenByte = 0
				}
				tape[tapePos] = givenByte
			case '[':
				if tape[tapePos] == 0 {
					nestingLevel++
				}
			case ']':
				if tape[tapePos] != 0 {
					nestingLevel--
				}
			}
		} else {
			switch operator {
			case '[':
				nestingLevel++
			case ']':
				nestingLevel--
			}
		}

		if nestingLevel < 0 {
			codePos--
		} else {
			codePos++
		}
	}
}

func promptByte() byte {
	inputReader := bufio.NewReader(os.Stdin)
	givenByte, _ := inputReader.ReadByte()

	return givenByte
}
