package main

import (
	"reflect"
	"testing"
)

var pizzaCutter = PizzaCutter{
	MinSliceToppingCount: 1,
	MaxSliceSize:         6,
	Pizza: [][]Topping{
		[]Topping{Tomato, Tomato, Tomato, Tomato, Tomato},
		[]Topping{Tomato, Mushroom, Mushroom, Mushroom, Tomato},
		[]Topping{Tomato, Tomato, Tomato, Tomato, Tomato},
	},
}

var perfectCut = Cuts{
	{RowA: 0, ColumnA: 0, RowB: 2, ColumnB: 1},
	{RowA: 0, ColumnA: 2, RowB: 2, ColumnB: 2},
	{RowA: 0, ColumnA: 3, RowB: 2, ColumnB: 4},
}

func TestNewTopping(t *testing.T) {
	if NewTopping('T') != Tomato {
		t.Fatalf("NewTopping() should return %s for %c", "Tomato", 'T')
	}
	if NewTopping('M') != Mushroom {
		t.Fatalf("NewTopping() should return %s for %c", "Mushroom", 'M')
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

func TestGetPerfectCut(t *testing.T) {
	actualPerfectCut := pizzaCutter.GetPerfectCuts()
	if !reflect.DeepEqual(actualPerfectCut, perfectCut) {
		t.Fatalf("PizzaCutter.GetPerfectCuts() should return %s but returned %s", perfectCut, actualPerfectCut)
	}
}
