package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Topping struct {
	Value, visited bool
}

type PizzaCutter struct {
	MinSliceToppingCount int
	MaxSliceSize         int
	Pizza                Pizza
}

type Pizza [][]Topping

type Cuts []Cut

type Cut struct {
	RowA, ColumnA, RowB, ColumnB int
}

func NewTopping(input rune) Topping {
	if input == 'T' {
		return Topping{Value: true}
	}

	return Topping{Value: false}
}

func NewPizzaCutter(input string) PizzaCutter {
	inputRows := strings.Split(input, "\n")
	inputHeaders := strings.Split(inputRows[0], " ")
	inputRows = inputRows[1:]

	rowCount, _ := strconv.Atoi(inputHeaders[0])
	columnCount, _ := strconv.Atoi(inputHeaders[1])
	minSliceToppingCount, _ := strconv.Atoi(inputHeaders[2])
	maxSliceSize, _ := strconv.Atoi(inputHeaders[3])

	pizza := make(Pizza, rowCount, rowCount)
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
	if t.Value == true {
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

func (p PizzaCutter) getPerfectCutsR(startRow, startColumn, score int, perfectCuts Cuts) (Cuts, bool) {
	for rowIndex := startRow; rowIndex < len(p.Pizza); rowIndex++ {
		row := p.Pizza[rowIndex]
		for columnIndex := startColumn; columnIndex < len(row); columnIndex++ {
			topping := row[columnIndex]
			if !topping.visited {
				tomatoCountsPerColumn := make([]int, len(row), len(row))
				mushroomCountsPerColumn := make([]int, len(row), len(row))

				for sliceRowIndex := 0; sliceRowIndex < p.MaxSliceSize && rowIndex+sliceRowIndex < len(p.Pizza); sliceRowIndex++ {
					for sliceColumnIndex := 0; (sliceColumnIndex+1)*(sliceRowIndex+1) <= p.MaxSliceSize && columnIndex+sliceColumnIndex < len(row); sliceColumnIndex++ {
						cut := Cut{rowIndex, columnIndex, rowIndex + sliceRowIndex, columnIndex + sliceColumnIndex}

						checkedTopping := p.Pizza[cut.RowB][cut.ColumnB]
						if !p.Pizza.IsVisited(cut) {
							if checkedTopping.Value == true {
								tomatoCountsPerColumn[sliceColumnIndex]++
							} else {
								mushroomCountsPerColumn[sliceColumnIndex]++
							}
							if sliceSum(tomatoCountsPerColumn[:sliceColumnIndex+1]) >= p.MinSliceToppingCount && sliceSum(mushroomCountsPerColumn[:sliceColumnIndex+1]) >= p.MinSliceToppingCount {
								p.Pizza.SetVisited(true, cut)
								score += (sliceColumnIndex + 1) * (sliceRowIndex + 1)
								perfectCuts = append(perfectCuts, cut)

								if score == len(p.Pizza)*len(row) {
									fmt.Print("PERFECT")
									return perfectCuts, true
								}

								nextPosX, nextPosY, hasNextPos := p.Pizza.GetNextNonVisitedPosition(cut)
								if hasNextPos {
									cuts, isPerfect := p.getPerfectCutsR(nextPosX, nextPosY, score, perfectCuts)
									if isPerfect {
										return cuts, true
									}
									perfectCuts = cuts

								}
								p.Pizza.SetVisited(false, cut)
								score -= (sliceColumnIndex + 1) * (sliceRowIndex + 1)
								perfectCuts = perfectCuts[:len(perfectCuts)-1]
							}
						}
					}
				}
			}
		}
	}

	return perfectCuts, false
}

func (pizza Pizza) IsVisited(cut Cut) bool {
	for i := cut.RowA; i <= cut.RowB; i++ {
		for j := cut.ColumnA; j <= cut.ColumnB; j++ {
			if pizza[i][j].visited {
				return true
			}
		}
	}

	return false
}

func (pizza Pizza) SetVisited(visited bool, cut Cut) bool {
	for i := cut.RowA; i <= cut.RowB; i++ {
		for j := cut.ColumnA; j <= cut.ColumnB; j++ {
			pizza[i][j].visited = visited
		}
	}

	return false
}

func (pizza Pizza) GetNextNonVisitedPosition(cut Cut) (x int, y int, hasPosition bool) {
	for nextPosRowIndex, i := cut.RowA, 0; nextPosRowIndex < len(pizza); nextPosRowIndex, i = nextPosRowIndex+1, i+1 {
		nextPosColumnIndex := 0
		if i == 0 {
			nextPosColumnIndex = cut.ColumnB + 1
		}
		for ; nextPosColumnIndex < len(pizza[0]); nextPosColumnIndex++ {
			if pizza[nextPosRowIndex][nextPosColumnIndex].visited == false {
				return nextPosRowIndex, nextPosColumnIndex, true
			}
		}
	}

	return 0, 0, false
}

func sliceSum(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}

	return r
}

func main() {
	fmt.Println(NewPizzaCutterFromFile("/Users/danshu/Gocode/src/github.com/DanShu93/golang-playground/algorithms/pizza/input/small.in").GetPerfectCuts())
}
