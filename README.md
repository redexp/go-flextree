# go-flextree

[![Go Reference](https://pkg.go.dev/badge/github.com/redexp/go-flextree.svg)](https://pkg.go.dev/github.com/redexp/go-flextree)

[Go](https://pkg.go.dev/github.com/redexp/go-flextree) implementation of [d3-flextree](https://github.com/Klortho/d3-flextree)

## Usage

With `Layout` and map functions `getChildren` and `getSize`

```go
tree := Layout(
    root,
    func(node Node) []Node {
        return node.Children
    },
    func(node Node) (float64, float64) {
        return node.Width, node.Height
    },
)

for _, child := range tree.Children {
    fmt.Printf("x: %.1f, y: %.1f\n", child.X, child.Y)
}
```

Or create `*Tree` by yourself and call `Reset()` and `Update()` methods

```go
tree := &Tree{
	Width: 100,
	Height: 100,
	Children: []*Tree{
        {
			Width: 50,
			Height: 30,
        },
        {
			Width: 50,
			Height: 30,
        },
    },
}

tree.Reset()
tree.Update()

for _, child := range tree.Children {
    fmt.Printf("x: %.1f, y: %.1f\n", child.X, child.Y)
}
```