package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func main() {
	clearTerminal()
	chars := []string{"X", "Y", "Z"}
	width := getTerminalWidth()
	for num := 0; ; num++ {
		for i, char := range chars {
			steps, _ := rand.Int(rand.Reader, big.NewInt(5))
			for j := int64(0); j < steps.Int64(); j++ {
				if len(chars[i]) >= width {
					fmt.Println(strconv.Itoa(i) + " WON!")
					return
				}
				chars[i] = " " + char
			}

			resetCursor()
			for _, printedChar := range chars {

				// for colorCode := 30; colorCode < 38; colorCode++ {
				fmt.Println(bold(color(31+i%6, printedChar)))
				// }
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}

func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	widthRegex := regexp.MustCompile(".* (.*)\n")
	width, _ := strconv.Atoi(string(widthRegex.ReplaceAll(out, []byte("$1"))))
	return width
}

func clearTerminal() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func resetCursor() {
	fmt.Printf("\033[0;0H")
}

func bold(input string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", input)
}

func color(color int, input string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, input)
}
