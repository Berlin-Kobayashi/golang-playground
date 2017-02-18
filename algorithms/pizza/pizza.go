package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"flag"
	"os"
	"os/signal"
	"runtime"
)

type Cell struct {
	Value, visited bool
}

type PizzaCutter struct {
	MinSliceCellCount int
	MaxSliceSize      int
	Pizza             Pizza
	BestResultChannel chan Result
	bestResult        Result
}

type Pizza [][]Cell

type Cuts []Cut

type Cut struct {
	RowA, ColumnA, RowB, ColumnB int
}

type Result struct {
	Score int
	Cuts  Cuts
}

func NewCell(input rune) Cell {
	if input == 'T' {
		return Cell{Value: true}
	}

	return Cell{Value: false}
}

func NewPizzaCutter(input string, bestResultChannel chan (Result)) PizzaCutter {
	inputRows := strings.Split(input, "\n")
	inputHeaders := strings.Split(inputRows[0], " ")
	inputRows = inputRows[1:]

	rowCount, _ := strconv.Atoi(inputHeaders[0])
	columnCount, _ := strconv.Atoi(inputHeaders[1])
	minSliceCellCount, _ := strconv.Atoi(inputHeaders[2])
	maxSliceSize, _ := strconv.Atoi(inputHeaders[3])

	pizza := make(Pizza, rowCount, rowCount)
	for i := range pizza {
		pizza[i] = make([]Cell, columnCount, columnCount)
		for j, cellRune := range inputRows[i] {
			pizza[i][j] = NewCell(cellRune)
		}
	}

	return PizzaCutter{
		MinSliceCellCount: minSliceCellCount,
		MaxSliceSize:      maxSliceSize,
		Pizza:             pizza,
		BestResultChannel: bestResultChannel,
	}
}

func (t Cell) String() string {
	if t.Value == true {
		return "T"
	}

	return "M"
}

func (p PizzaCutter) String() string {
	return fmt.Sprintf("%d:%d:%s", p.MinSliceCellCount, p.MaxSliceSize, p.Pizza)
}

func (c Cut) Rotate(rotations int) Cut {
	switch rotations % 4 {
	case 0:
		return Cut{RowA: c.RowA, RowB: c.RowB, ColumnA: c.ColumnA, ColumnB: c.ColumnB}
	case 1:
		return Cut{RowA: c.ColumnB, RowB: c.ColumnA, ColumnA: c.RowA, ColumnB: c.RowB}
	case 2:
		return Cut{RowA: c.RowB, RowB: c.RowA, ColumnA: c.ColumnB, ColumnB: c.ColumnA}
	default:
		return Cut{RowA: c.ColumnA, RowB: c.ColumnB, ColumnA: c.RowB, ColumnB: c.RowA}
	}
}

func (c Cuts) Rotate(rotations int) Cuts {
	result := make(Cuts, len(c), len(c))
	for i, cut := range c {
		result[i] = cut.Rotate(rotations)
	}

	return result
}

func (c Cuts) String() string {
	output := strconv.Itoa(len(c)) + "\n"
	for _, slice := range c {
		output += strconv.Itoa(slice.RowA) + " " + strconv.Itoa(slice.ColumnA) + " " + strconv.Itoa(slice.RowB) + " " + strconv.Itoa(slice.ColumnB) + "\n"
	}

	return output
}

func NewPizzaCutterFromFile(path string, bestResultChannel chan (Result)) PizzaCutter {
	return NewPizzaCutter(readFile(path), bestResultChannel)
}

func readFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}

func (p *PizzaCutter) GetPerfectCuts() {
	go p.getPerfectCutsR(0, 0, 0, Cuts{})
}

func (p *PizzaCutter) getPerfectCutsR(startRow, startColumn, score int, cuts Cuts) (Cuts, bool) {
	for rowIndex := startRow; rowIndex < len(p.Pizza); rowIndex++ {
		row := p.Pizza[rowIndex]
		for columnIndex := startColumn; columnIndex < len(row); columnIndex++ {
			cell := row[columnIndex]
			if !cell.visited {
				tomatoCountsPerColumn := make([]int, len(row), len(row))
				mushroomCountsPerColumn := make([]int, len(row), len(row))

				for sliceRowIndex := 0; sliceRowIndex < p.MaxSliceSize && rowIndex+sliceRowIndex < len(p.Pizza); sliceRowIndex++ {
					for sliceColumnIndex := 0; (sliceColumnIndex+1)*(sliceRowIndex+1) <= p.MaxSliceSize && columnIndex+sliceColumnIndex < len(row); sliceColumnIndex++ {
						cut := Cut{rowIndex, columnIndex, rowIndex + sliceRowIndex, columnIndex + sliceColumnIndex}

						checkedCell := p.Pizza[cut.RowB][cut.ColumnB]
						if !p.Pizza.IsVisited(cut) {
							if checkedCell.Value == true {
								tomatoCountsPerColumn[sliceColumnIndex]++
							} else {
								mushroomCountsPerColumn[sliceColumnIndex]++
							}
							if sliceSum(tomatoCountsPerColumn[:sliceColumnIndex+1]) >= p.MinSliceCellCount && sliceSum(mushroomCountsPerColumn[:sliceColumnIndex+1]) >= p.MinSliceCellCount {
								p.Pizza.SetVisited(true, cut)
								score += (sliceColumnIndex + 1) * (sliceRowIndex + 1)
								cuts = append(cuts, cut)
								if score > p.bestResult.Score {
									p.bestResult.Score = score
									cutsCopy := make(Cuts, len(cuts), len(cuts))
									copy(cutsCopy, cuts)
									p.bestResult.Cuts = cutsCopy
									p.BestResultChannel <- p.bestResult
									if score == len(p.Pizza)*len(row) {
										close(p.BestResultChannel)
										return cuts, true
									}
								}

								nextPosX, nextPosY, hasNextPos := p.Pizza.GetNextNonVisitedPosition(cut)
								if hasNextPos {
									cuts, isPerfect := p.getPerfectCutsR(nextPosX, nextPosY, score, cuts)
									if isPerfect {
										return cuts, true
									}
									cuts = cuts

								}
								p.Pizza.SetVisited(false, cut)
								score -= (sliceColumnIndex + 1) * (sliceRowIndex + 1)
								cuts = cuts[:len(cuts)-1]
							}
						}
					}
				}
			}
		}
	}

	return cuts, false
}

func (p *PizzaCutter) isValidCuts(cuts Cuts) (bool, string) {
	pizza := make([][]Cell, len(p.Pizza), len(p.Pizza))
	for i := range pizza {
		pizza[i] = make([]Cell, len(p.Pizza[0]), len(p.Pizza[0]))
	}

	for c, cut := range cuts {
		tomatoCount := 0
		mushroomCount := 0
		if (cut.RowB-cut.RowA)*(cut.ColumnB-cut.ColumnA) > p.MaxSliceSize {
			return false, fmt.Sprintf("To big cut at cut number %d", c)
		}

		for rowIndex, i := cut.RowA, 0; rowIndex < len(pizza); rowIndex, i = rowIndex+1, i+1 {
			columnIndex := 0
			if i == 0 {
				columnIndex = cut.ColumnA
			}
			for ; columnIndex < len(pizza[0]); columnIndex++ {
				if rowIndex >= cut.RowA && rowIndex <= cut.RowB && columnIndex >= cut.ColumnA && columnIndex <= cut.ColumnB {
					if pizza[rowIndex][columnIndex].visited {
						return false, fmt.Sprintf("Reused cell at cut number %d", c)
					}
					pizza[rowIndex][columnIndex].visited = true
					if p.Pizza[rowIndex][columnIndex].Value {
						tomatoCount++
					} else {
						mushroomCount++
					}
				}
			}
		}

		if tomatoCount < p.MinSliceCellCount {
			return false, fmt.Sprintf("Not enough tomatos at cut number %d", c)
		}
		if mushroomCount < p.MinSliceCellCount {
			return false, fmt.Sprintf("Not enough mushrooms at cut number %d", c)
		}
	}

	return true, ""
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

func (pizza Pizza) Rotate(rotations int) Pizza {
	previous := pizza
	for r := 0; r < rotations%4; r++ {
		result := Pizza{}
		for i, row := range previous {
			for j := len(row) - 1; j >= 0; j-- {
				cell := row[j]
				if len(result) == (len(row)-1)-j {
					result = append(result, make([]Cell, len(previous), len(previous)))
				}
				result[(len(row)-1)-j][i] = Cell{Value: cell.Value}
			}
		}
		previous = result
	}

	return previous
}

func sliceSum(slice []int) int {
	r := 0
	for _, v := range slice {
		r += v
	}

	return r
}

func main() {
	fmt.Printf("%d cores available\n", runtime.NumCPU())
	flag.Parse()
	var input = flag.Args()
	size := input[0]

	outputPath := fmt.Sprintf("./output/%s.out", size)
	bestResultChannel := make(chan Result)
	cutter := NewPizzaCutterFromFile(fmt.Sprintf("./input/%s.in", size), bestResultChannel)

	lastResult := Result{}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	interrupted := false
	go func() {
		<-signals
		interrupted = true
		fmt.Println("Saving results ...")
		saveResults(outputPath, lastResult.Cuts)
		fmt.Println("Results saved!")

		printResultValidationStatus(lastResult, cutter)

		os.Exit(1)
	}()

	cutter.GetPerfectCuts()

	for result := range cutter.BestResultChannel {
		if !interrupted {
			fmt.Printf("%d > %d\n", result.Score, lastResult.Score)
			cuts := make(Cuts, len(result.Cuts), len(result.Cuts))
			copy(cuts, result.Cuts)
			lastResult = Result{Score: result.Score, Cuts: result.Cuts}
		}
	}

	printResultValidationStatus(lastResult, cutter)

	saveResults(outputPath, lastResult.Cuts)
}

func saveResults(path string, results fmt.Stringer) {
	ioutil.WriteFile(path, []byte(results.String()), 0644)
}

func printResultValidationStatus(result Result, pizzaCutter PizzaCutter) {
	valid, _ := pizzaCutter.isValidCuts(result.Cuts)
	fmt.Printf("Valid %t\n", valid)
}
