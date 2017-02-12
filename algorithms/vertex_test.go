package algorithms

import (
	"reflect"
	"testing"
)

var graph = &Vertex{
	Value: "A",
	Children: []*Vertex{
		&Vertex{
			Value: "B",
			Children: []*Vertex{
				&Vertex{
					Value: "C",
				},
				&Vertex{
					Value: "D",
					Children: []*Vertex{
						&Vertex{
							Value: "E",
						}},
				},
			},
		},
		&Vertex{
			Value: "F",
		},
		&Vertex{
			Value: "G",
			Children: []*Vertex{
				&Vertex{
					Value: "H",
				},
			},
		},
	},
}

func TestBreadthsFirstSearch(t *testing.T) {
	expectedOrder := []string{"A", "B", "F", "G", "C", "D", "H", "E"}

	actualOrder := []string{}
	graph.BreadthsFirstSearch(func(value string) {
		actualOrder = append(actualOrder, value)
	})

	if !reflect.DeepEqual(actualOrder, expectedOrder) {
		t.Fatalf("BreadthsFirstSearch is not traversing graph in correct order. order= %s, want %s", actualOrder, expectedOrder)
	}
}

func TestDepthFirstSearch(t *testing.T) {
	expectedOrder := []string{"A", "B", "C", "D", "E", "F", "G", "H"}

	actualOrder := []string{}
	graph.DepthFirstSearch(func(value string) {
		actualOrder = append(actualOrder, value)
	})

	if !reflect.DeepEqual(actualOrder, expectedOrder) {
		t.Fatalf("DepthFirstSearch is not traversing graph in correct order. order= %s, want %s", actualOrder, expectedOrder)
	}
}
