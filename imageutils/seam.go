package imageutils

import (
	"fmt"
	"image"
	"image/color"
)

type SeamNode struct {
	// if null = root node
	PreviosNode      *SeamNode
	X                int
	Y                int
	ThisCost         uint
	AccumulatedCosts uint
}

func (s SeamNode) String() string {
	if s.PreviosNode == nil {
		return fmt.Sprintf("(%d, %d; %d)", s.X, s.Y, s.AccumulatedCosts)
	} else {
		return fmt.Sprintf("(%d, %d; %d) -> %v", s.X, s.Y, s.AccumulatedCosts, s.PreviosNode)
	}
}

func (s SeamNode) WriteSeamChainToImage(img *image.RGBA) {
	img.Set(
		s.X,
		s.Y,
		color.RGBA{255, 0, 0, 1},
	)

	if s.PreviosNode != nil {
		s.PreviosNode.WriteSeamChainToImage(img)
	}
}

type Position struct {
	X int
	Y int
}

type ImageSeams struct {
	// map of <y row> -> SeamNodes
	Seams [][]*SeamNode
}

func NewImageSeams() ImageSeams {
	return ImageSeams{
		Seams: make([][]*SeamNode, 0),
	}
}

func (i *ImageSeams) CreateSeamsFromRectangle(rect image.Rectangle) {
	for y := range rect.Max.Y {
		rowSeams := make([]*SeamNode, 0)

		for x := range rect.Max.X {
			rowSeams = append(rowSeams, &SeamNode{
				X: x,
				Y: y,
			})
		}

		i.Seams = append(i.Seams, rowSeams)
	}
}

// Get the nodes above a given x, y coordinate.
// Get the left and right neighbors, if there are any.
func (i *ImageSeams) GetNodesAbove(x int, y int) []*SeamNode {
	if y == 0 {
		return nil
	}

	rowSeams := i.Seams[y-1]

	seams := make([]*SeamNode, 0, 3)

	// Left node
	if x != 0 {
		seams = append(seams, rowSeams[x-1])
	}
	// Middle node
	seams = append(seams, rowSeams[x])
	// Right node
	if x < (len(rowSeams) - 1) {
		seams = append(seams, rowSeams[x+1])
	}

	return seams
}

// Find the best node that's above the current one at `(x, y)`
func (i *ImageSeams) FindBestNodeAbove(x int, y int) *SeamNode {
	nodes := i.GetNodesAbove(x, y)
	lowestCostNode := nodes[0]

	for _, node := range nodes {
		if node.AccumulatedCosts < lowestCostNode.AccumulatedCosts {
			lowestCostNode = node
		}
	}

	return lowestCostNode
}

func (i *ImageSeams) CreateOptimizedRoutesForRow(y int) {
	row := i.Seams[y]

	for x, node := range row {
		bestNode := i.FindBestNodeAbove(x, y)

		node.PreviosNode = bestNode
		node.AccumulatedCosts = bestNode.AccumulatedCosts + node.ThisCost
	}
}

func (i *ImageSeams) CreateOptimizedRoutes() {
	for y := 1; y < len(i.Seams); y++ {
		i.CreateOptimizedRoutesForRow(y)
	}
}

func (i *ImageSeams) GetLowestSeam() *SeamNode {
	lastRow := i.Seams[len(i.Seams)-1]
	lowestNode := lastRow[0]

	for _, node := range lastRow {
		if node.AccumulatedCosts < lowestNode.AccumulatedCosts {
			lowestNode = node
		}
	}

	return lowestNode
}

func (i *ImageSeams) SetCostForNode(x int, y int, cost uint) {
	node := i.Seams[y][x]
	node.ThisCost = cost
	node.AccumulatedCosts = cost
}
