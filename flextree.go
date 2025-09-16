package flextree

func Layout[T any](
	data T,
	getChildren func(T) []T,
	getSize func(T) (float64, float64),
) *Tree {
	tree := NewTree(data, getChildren, getSize)

	tree.Update()

	return tree
}

type Tree struct {
	Input    any
	Parent   *Tree
	Children []*Tree

	X      float64
	Y      float64
	Width  float64
	Height float64

	Depth   float64
	Spacing float64

	lExt     *Tree
	rExt     *Tree
	lThr     *Tree
	rThr     *Tree
	relX     float64
	prelim   float64
	lExtRelX float64
	rExtRelX float64
	shift    float64
	change   float64
}

func NewTree[T any](
	data T,
	getChildren func(T) []T,
	getSize func(T) (float64, float64),
) *Tree {
	var walk func(node T, parent *Tree) *Tree

	walk = func(data T, parent *Tree) *Tree {
		width, height := getSize(data)

		node := &Tree{
			Input:  data,
			Parent: parent,
			Width:  width,
			Height: height,
		}

		if parent != nil {
			node.Depth = parent.Depth + 1
		}

		node.lExt = node
		node.rExt = node

		children := getChildren(data)

		node.Children = make([]*Tree, len(children))

		for i, child := range children {
			node.Children[i] = walk(child, node)
		}

		return node
	}

	return walk(data, nil)
}

func (tree *Tree) Update() {
	layoutChildren(tree, 0)
	resolveX(tree, -tree.relX-tree.prelim, 0)
}

func (tree *Tree) FirstChild() *Tree {
	count := len(tree.Children)

	if count > 0 {
		return tree.Children[0]
	}

	return nil
}

func (tree *Tree) LastChild() *Tree {
	count := len(tree.Children)

	if count > 0 {
		return tree.Children[count-1]
	}

	return nil
}

func (tree *Tree) Bottom() float64 {
	return tree.Y + tree.Height
}

func layoutChildren(tree *Tree, y float64) {
	tree.Y = y

	prevLows := &Lows{}
	var lowY float64

	for i, child := range tree.Children {
		layoutChildren(child, tree.Bottom())

		if i == 0 {
			lowY = child.lExt.Bottom()
		} else {
			lowY = child.rExt.Bottom()
		}

		if i > 0 {
			separate(tree, i, prevLows)
		}

		prevLows = updateLows(lowY, i, prevLows)
	}

	shiftChange(tree)
	positionRoot(tree)
}

func resolveX(tree *Tree, prevSum, parentX float64) {
	sum := prevSum + tree.relX
	tree.relX = sum + tree.prelim - parentX
	tree.prelim = 0
	tree.X = parentX + tree.relX

	for _, child := range tree.Children {
		resolveX(child, sum, tree.X)
	}
}

func separate(tree *Tree, i int, lows *Lows) {
	lSib := tree.Children[i-1]
	curSubtree := tree.Children[i]
	rContour := lSib
	rSumMods := lSib.relX
	lContour := curSubtree
	lSumMods := curSubtree.relX
	isFirst := true

	for rContour != nil && lContour != nil {
		if rContour.Bottom() > lows.lowY {
			lows = lows.next
		}

		dist := (rSumMods + rContour.prelim) - (lSumMods + lContour.prelim) +
			rContour.Width/2 + lContour.Width/2 +
			rContour.Spacing

		if dist > 0 || (dist < 0 && isFirst) {
			lSumMods += dist
			moveSubtree(curSubtree, dist)
			distributeExtra(tree, i, lows.index, dist)
		}

		isFirst = false

		rightBottom := rContour.Bottom()
		leftBottom := lContour.Bottom()

		if rightBottom <= leftBottom {
			rContour = nextRContour(rContour)
			if rContour != nil {
				rSumMods += rContour.relX
			}
		}
		if rightBottom >= leftBottom {
			lContour = nextLContour(lContour)
			if lContour != nil {
				lSumMods += lContour.relX
			}
		}
	}

	if rContour == nil && lContour != nil {
		setLThr(tree, i, lContour, lSumMods)
	} else if rContour != nil && lContour == nil {
		setRThr(tree, i, rContour, rSumMods)
	}
}

type Lows struct {
	lowY  float64
	index int
	next  *Lows
}

func updateLows(lowY float64, index int, lows *Lows) *Lows {
	for lows != nil && lowY >= lows.lowY {
		lows = lows.next
	}

	return &Lows{
		lowY:  lowY,
		index: index,
		next:  lows,
	}
}

func moveSubtree(subtree *Tree, distance float64) {
	subtree.relX += distance
	subtree.lExtRelX += distance
	subtree.rExtRelX += distance
}

func distributeExtra(tree *Tree, curSubtreeI int, leftSibI int, dist float64) {
	curSubtree := tree.Children[curSubtreeI]
	n := curSubtreeI - leftSibI

	if n > 1 {
		delta := dist / float64(n)
		tree.Children[leftSibI+1].shift += delta
		curSubtree.shift -= delta
		curSubtree.change -= dist - delta
	}
}

func nextRContour(tree *Tree) *Tree {
	if len(tree.Children) > 0 {
		return tree.LastChild()
	} else {
		return tree.rThr
	}
}

func nextLContour(tree *Tree) *Tree {
	if len(tree.Children) > 0 {
		return tree.FirstChild()
	} else {
		return tree.lThr
	}
}

func setLThr(tree *Tree, i int, lContour *Tree, lSumMods float64) {
	firstChild := tree.FirstChild()
	lExt := firstChild.lExt
	curSubtree := tree.Children[i]
	lExt.lThr = lContour

	diff := lSumMods - lContour.relX - firstChild.lExtRelX

	lExt.relX += diff
	lExt.prelim -= diff

	firstChild.lExt = curSubtree.lExt
	firstChild.lExtRelX = curSubtree.lExtRelX
}

func setRThr(tree *Tree, i int, rContour *Tree, rSumMods float64) {
	curSubtree := tree.Children[i]
	rExt := curSubtree.rExt
	lSib := tree.Children[i-1]
	rExt.rThr = rContour
	diff := rSumMods - rContour.relX - curSubtree.rExtRelX
	rExt.relX += diff
	rExt.prelim -= diff
	curSubtree.rExt = lSib.rExt
	curSubtree.rExtRelX = lSib.rExtRelX
}

func shiftChange(tree *Tree) {
	var lastShiftSum float64
	var lastChangeSum float64

	for _, child := range tree.Children {
		shiftSum := lastShiftSum + child.shift
		changeSum := lastChangeSum + shiftSum + child.change
		child.relX += changeSum
		lastShiftSum = shiftSum
		lastChangeSum = changeSum
	}
}

func positionRoot(tree *Tree) {
	if len(tree.Children) == 0 {
		return
	}

	k0 := tree.FirstChild()
	kf := tree.LastChild()

	prelim := (k0.prelim + k0.relX - k0.Width/2 +
		kf.relX + kf.prelim + kf.Width/2) / 2

	tree.prelim = prelim
	tree.lExt = k0.lExt
	tree.lExtRelX = k0.lExtRelX
	tree.rExt = kf.rExt
	tree.rExtRelX = kf.rExtRelX
}
