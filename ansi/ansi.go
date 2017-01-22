package main

import (
	"fmt"
	"time"
)

func main() {
	for i := 30; i < 38; i++ {
		fmt.Println(color(i, "hello"))
		fmt.Println(bold(color(i, "hello")))
	}
	fmt.Print("Hollaa")
	time.Sleep(time.Second * 1)
	fmt.Print("\rasdasda")
	fmt.Print("\n")
	time.Sleep(time.Second * 1)
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
	fmt.Println("test")
	fmt.Println("test2")
    time.Sleep(time.Second * 1)
	fmt.Printf("\033[0;0H")
	fmt.Println("testy")
    fmt.Println("test2")
    time.Sleep(time.Second * 1)
}

func bold(input string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", input)
}

func color(color int, input string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", color, input)
}
