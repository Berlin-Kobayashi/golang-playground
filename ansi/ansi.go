package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	var input = flag.Args()
	characterAmount, _ := strconv.Atoi(input[0])
	clearTerminal()
	chars := createCharacters(characterAmount)
	width := getTerminalWidth()
	border := createBoarder(width)
	speed, _ := strconv.Atoi(input[1])
	for {
		fmt.Println(border)
		for i, char := range chars {
			steps, _ := rand.Int(rand.Reader, big.NewInt(5))
			for j := int64(0); j < steps.Int64(); j++ {
				if len(chars[i]) >= width {
					fmt.Println(string(chars[i][len(chars[i])-1]) + " WON!")
					return
				}
				chars[i] = " " + char
			}
			fmt.Println(color(32, bold(char)))
			fmt.Println(border)
		}

		time.Sleep(time.Millisecond * time.Duration(speed))
		resetCursor()
	}
}

func createCharacters(characterAmount int) []string {
	characters := []string{}

	for i := 0; i < characterAmount; i++ {
		characters = append(characters, string('!'+byte(i)))
	}

	return characters
}

func createBoarder(width int) string {
	border := ""

	for i := 0; i < width; i++ {
		border += "-"
	}

	return border
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
