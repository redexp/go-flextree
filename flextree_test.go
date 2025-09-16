package flextree

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestLayout(t *testing.T) {
	file, err := os.Open("test-data.json")

	if err != nil {
		t.Fatal(err)
		return
	}

	defer file.Close()

	var tests []Test

	if err = json.NewDecoder(file).Decode(&tests); err != nil {
		t.Fatal(err)
		return
	}

	if len(tests) == 0 {
		t.Fatal("tests == 0")
		return
	}

	for _, test := range tests {
		tree := Layout(
			test.Input,
			func(tree TestTree) []TestTree {
				return tree.Children
			},
			func(tree TestTree) (float64, float64) {
				return tree.Data[0], tree.Data[1]
			},
		)

		if !valid(tree, test.Output) {
			t.Error(test.Name)
		}
	}
}

func f(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

func valid(tree *Tree, test TestTree) bool {
	if f(tree.X) != f(test.Data[0]) || f(tree.Y) != f(test.Data[1]) || len(tree.Children) != len(test.Children) {
		return false
	}

	for i, child := range test.Children {
		if !valid(tree.Children[i], child) {
			return false
		}
	}

	return true
}

type Test struct {
	Name   string
	Input  TestTree
	Output TestTree
}

type TestTree struct {
	Data     []float64
	Children []TestTree
}
