package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
)

//>	move data pointer to the next cell on the right
//<	move data pointer to the next cell on the left
//+	increment value of current cell
//-	decrement value of current cell
//.	output the byte of a currently pointed cell in ASCII code
//,	read one byte from stdin and store its value at the current cell
//[	if the current cell is 0 then jump to the matching ]
//]	jump to the matching [
//everything else ignored
func main() {
	flag.Parse()
	var input = flag.Args()
	size, _ := strconv.Atoi(input[1])
	Run(input[0], size)
}

func Run(program string, size int) {
	tape := make([]byte, size, size)
	pos := 0
	opened := 0
	closed := 0
	for i:= 0; i < len(program);  {
		operator := program[i]
		if opened > 0 {
			switch operator {
		case '['  :
			opened++
		case ']' :
		 	opened--
		}
		i++
		} else if closed > 0 {
		switch operator {
		case '['  :
			closed--
		case ']' :
		 	closed++
		}
		if closed == 0 {
			i++
		}else {
   			i--
		}
		}else {
		switch operator {
		case '>':
			if pos+1 >= size {
				pos = 0
			} else {
				pos++
			}
				i++
		case '<':
			if pos-1 < 0 {
				pos = size - 1
			} else {
				pos--
			}
				i++
		case '+':
			if tape[pos] == 255 {
				tape[pos] = 0
			} else {
				tape[pos]++
			}
				i++
		case '-':
			if tape[pos] == 0 {
				tape[pos] = 255
			} else {
				tape[pos]--
			}
				i++
		case '.':
			fmt.Print(fmt.Sprintf("%c", tape[pos]))
				i++
		case ',':
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadByte()
			if text == 10 {
				text = 0
			}
			tape[pos] = text
				i++
		case '[' :
			if tape[pos] == 0 {
				opened++
			}
				i++
		case ']' :
			if tape[pos] != 0 {
				closed++
					i--
			}else {
	i++
			}

		}

		}
	}
}
