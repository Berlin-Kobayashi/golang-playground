package main

import (
	"reflect"
	"testing"
)

var pizzaCutter = PizzaCutter{
	MinSliceCellCount: 1,
	MaxSliceSize:      6,
	Pizza: [][]Cell{
		{Cell{Value: true}, Cell{Value: true}, Cell{Value: true}, Cell{Value: true}, Cell{Value: true}},
		{Cell{Value: true}, Cell{Value: false}, Cell{Value: false}, Cell{Value: false}, Cell{Value: true}},
		{Cell{Value: true}, Cell{Value: true}, Cell{Value: true}, Cell{Value: true}, Cell{Value: true}},
	},
}

var perfectCut = Cuts{
	{RowA: 0, ColumnA: 0, RowB: 2, ColumnB: 1},
	{RowA: 0, ColumnA: 2, RowB: 2, ColumnB: 2},
	{RowA: 0, ColumnA: 3, RowB: 2, ColumnB: 4},
}

func TestNewCell(t *testing.T) {
	tomato := Cell{Value: true}
	mushroom := Cell{Value: false}
	if NewCell('T') != tomato {
		t.Fatalf("NewCell() should return %s for %c", "Tomato", 'T')
	}
	if NewCell('M') != mushroom {
		t.Fatalf("NewCell() should return %s for %c", "Mushroom", 'M')
	}
}

func TestNewPizzaCutter(t *testing.T) {
	input := `3 5 1 6
TTTTT
TMMMT
TTTTT
`
	actualPizzaCutter := NewPizzaCutter(input)

	if !reflect.DeepEqual(actualPizzaCutter, pizzaCutter) {
		t.Fatalf("NewPizzaCutter() should return %s for %s but returned %s", pizzaCutter, input, actualPizzaCutter)
	}
}

func TestNewPizzaCutterFromFile(t *testing.T) {
	path := "./input/example.in"
	actualPizzaCutter := NewPizzaCutterFromFile(path)
	if !reflect.DeepEqual(actualPizzaCutter, pizzaCutter) {
		t.Fatalf("NewPizzaCutterFromFile() should return %s for %s but returned %s", pizzaCutter, path, actualPizzaCutter)
	}
}

func TestCutToString(t *testing.T) {
	expectedOutput := `3
0 0 2 1
0 2 2 2
0 3 2 4
`
	actualOutput := perfectCut.String()
	if actualOutput != expectedOutput {
		t.Fatalf("Cuts.String() should return %s but returned %s", expectedOutput, actualOutput)
	}
}

func TestIsValidCuts(t *testing.T) {
	validCuts := Cuts{Cut{RowA: 0, RowB: 1, ColumnA: 1, ColumnB: 1}}
	valid, message := pizzaCutter.isValidCuts(validCuts)
	if !valid {
		t.Fatalf("PizzaCutter.IsValidCuts() should return true but returned %s", message)
	}
	notEnoughTomatoesCuts := Cuts{Cut{RowA: 1, RowB: 1, ColumnA: 1, ColumnB: 2}}
	valid, message = pizzaCutter.isValidCuts(notEnoughTomatoesCuts)
	if valid {
		t.Fatal("PizzaCutter.IsValidCuts() should return false cause not enough tomatoes but returned true")
	}
	notEnoughMushroomsCuts := Cuts{Cut{RowA: 0, RowB: 0, ColumnA: 0, ColumnB: 1}}
	valid, message = pizzaCutter.isValidCuts(notEnoughMushroomsCuts)
	if valid {
		t.Fatal("PizzaCutter.IsValidCuts() should return false cause not enough mushrooms but returned true")
	}
	tooBigCuts := Cuts{Cut{RowA: 0, RowB: 0, ColumnA: 2, ColumnB: 2}}
	valid, message = pizzaCutter.isValidCuts(tooBigCuts)
	if valid {
		t.Fatal("PizzaCutter.IsValidCuts() should return false cause too big cut but returned true")
	}
}

func TestGetPerfectCutTiny(t *testing.T) {
	actualPerfectCut := pizzaCutter.GetPerfectCuts()
	if !reflect.DeepEqual(actualPerfectCut, perfectCut) {
		t.Fatalf("PizzaCutter.GetPerfectCuts() should return %s but returned %s", perfectCut, actualPerfectCut)
	}
}

func TestGetPerfectCutSmall(t *testing.T) {
	path := "./input/small.in"
	actualPizzaCutter := NewPizzaCutterFromFile(path)
	actualPerfectCut := actualPizzaCutter.GetPerfectCuts()
	valid, message := actualPizzaCutter.isValidCuts(actualPerfectCut)
	if !valid {
		t.Fatalf("PizzaCutter.GetPerfectCuts() returned an invalid valid cut: %s \n%s", message, actualPerfectCut)
	}
}
