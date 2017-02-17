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
	sliceStatus    int
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
	cuts, _ := p.getPerfectCutsR(0, 0, 0, Cuts{})

	return cuts
}

func (p PizzaCutter) getPerfectCutsR(posX, posY, score int, perfectCuts Cuts) (Cuts, bool) {
	for rowIndex := posX; rowIndex < len(p.Pizza); rowIndex++ {
		row := p.Pizza[rowIndex]
		for toppingIndex := posY; toppingIndex < len(row); toppingIndex++ {
			topping := row[toppingIndex]
			if !topping.visited {
				tomatoCounts := make([]int, len(row), len(row))
				mushroomCounts := make([]int, len(row), len(row))

				for sliceRowIndex := 0; sliceRowIndex < p.MaxSliceSize && rowIndex+sliceRowIndex < len(p.Pizza); sliceRowIndex++ {
					for sliceColumnIndex := 0; (sliceColumnIndex+1)*(sliceRowIndex+1) <= p.MaxSliceSize && toppingIndex+sliceColumnIndex < len(row); sliceColumnIndex++ {
						checkedTopping := p.Pizza[rowIndex+sliceRowIndex][toppingIndex+sliceColumnIndex]
						if !checkedTopping.visited {
							if checkedTopping.Value == true {
								tomatoCounts[sliceColumnIndex]++
							} else {
								mushroomCounts[sliceColumnIndex]++
							}
							if sliceSum(tomatoCounts[:sliceColumnIndex+1]) >= p.MinSliceToppingCount && sliceSum(mushroomCounts[:sliceColumnIndex+1]) >= p.MinSliceToppingCount {
								for i := 0; i <= sliceRowIndex; i++ {
									for j := 0; j <= sliceColumnIndex; j++ {
										p.Pizza[rowIndex+i][toppingIndex+j].visited = true
									}
								}
								score += (sliceColumnIndex + 1) * (sliceRowIndex + 1)
								perfectCuts = append(perfectCuts, Cut{rowIndex, toppingIndex, rowIndex + sliceRowIndex, toppingIndex + sliceColumnIndex})

								if score == len(p.Pizza)*len(row) {
									return perfectCuts, true
								}

								nextPosX := 0
								nextPosY := 0
								hasNextPos := false
							OUTER:
								for nextPosRowIndex, i := rowIndex, 0; nextPosRowIndex < len(p.Pizza); nextPosRowIndex, i = nextPosRowIndex+1, i+1 {
									nextPosColumnIndex := 0
									if i == 0 {
										nextPosColumnIndex = toppingIndex + sliceColumnIndex + 1
									}
									for ; nextPosColumnIndex < len(row); nextPosColumnIndex++ {
										if p.Pizza[nextPosRowIndex][nextPosColumnIndex].visited == false {
											nextPosX = nextPosRowIndex
											nextPosY = nextPosColumnIndex
											hasNextPos = true
											break OUTER
										}
									}
								}
								if hasNextPos {
									cuts, isPerfect := p.getPerfectCutsR(nextPosX, nextPosY, score, perfectCuts)
									if isPerfect {
										return cuts, true
									}
									perfectCuts = cuts

									for i := 0; i <= sliceRowIndex; i++ {
										for j := 0; j <= sliceColumnIndex; j++ {
											p.Pizza[rowIndex+i][toppingIndex+j].visited = false
										}
									}
									score -= (sliceColumnIndex + 1) * (sliceRowIndex + 1)
									perfectCuts = perfectCuts[:len(perfectCuts)-1]
								}

							}
						}
					}
				}
			}
		}
	}

	return perfectCuts, false
}

func sliceSum(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}

	return r
}

func main() {
	fmt.Println(NewPizzaCutterFromFile("/Users/danshu/Gocode/src/github.com/DanShu93/golang-playground/algorithms/pizza/input/example.in").GetPerfectCuts())
}
