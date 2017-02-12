package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var (
	Tomato   = Topping{Value: true}
	Mushroom = Topping{Value: false}
)

type Topping struct {
	Value, visited bool
	sliceSize      int
}

type PizzaCutter struct {
	MinSliceToppingCount int
	MaxSliceSize         int
	Pizza                [][]Topping
}

type Cuts []Cut

type Cut struct {
	RowA, ColumnA, RowB, ColumnB int
}

func NewTopping(input rune) Topping {
	if input == 'T' {
		return Tomato
	}

	return Mushroom
}

func NewPizzaCutter(input string) PizzaCutter {
	inputRows := strings.Split(input, "\n")
	inputHeaders := strings.Split(inputRows[0], " ")
	inputRows = inputRows[1:]

	rowCount, _ := strconv.Atoi(inputHeaders[0])
	columnCount, _ := strconv.Atoi(inputHeaders[1])
	minSliceToppingCount, _ := strconv.Atoi(inputHeaders[2])
	maxSliceSize, _ := strconv.Atoi(inputHeaders[3])

	pizza := make([][]Topping, rowCount, rowCount)
	for i := range pizza {
		pizza[i] = make([]Topping, columnCount, columnCount)
		for j, toppingRune := range inputRows[i] {
			pizza[i][j] = NewTopping(toppingRune)
		}
	}

	return PizzaCutter{
		MinSliceToppingCount: minSliceToppingCount,
		MaxSliceSize:         maxSliceSize,
		Pizza:                pizza,
	}
}

func (t Topping) String() string {
	if t == Tomato {
		return "T"
	}

	return "M"
}

func (p PizzaCutter) String() string {
	return fmt.Sprintf("%d:%d:%s", p.MinSliceToppingCount, p.MaxSliceSize, p.Pizza)
}

func (c Cuts) String() string {
	output := strconv.Itoa(len(c)) + "\n"
	for _, slice := range c {
		output += strconv.Itoa(slice.RowA) + " " + strconv.Itoa(slice.ColumnA) + " " + strconv.Itoa(slice.RowB) + " " + strconv.Itoa(slice.ColumnB) + "\n"
	}

	return output
}

func NewPizzaCutterFromFile(path string) PizzaCutter {
	return NewPizzaCutter(readFile(path))
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func (p PizzaCutter) GetPerfectCuts() Cuts {
	perfectCuts := Cuts{}

	for rowIndex := 0; rowIndex < len(p.Pizza); rowIndex++ {
		row := p.Pizza[rowIndex]
		for toppingIndex := 0; toppingIndex < len(row); toppingIndex++ {
			topping := row[toppingIndex]
			if !topping.visited {
				for sliceSizeCounter := 0; sliceSizeCounter < p.MaxSliceSize; sliceSizeCounter++ {

				}
			}
		}
	}

	return perfectCuts
}

func main() {
	fmt.Println(NewPizzaCutterFromFile("/Users/danshu/Gocode/src/github.com/DanShu93/golang-playground/algorithms/pizza/input/big.in"))
}
