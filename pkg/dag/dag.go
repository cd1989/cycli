package dag

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/cd1989/cycli/pkg/console"
)

type Node struct {
	Name string
}

type Edge struct {
	From      string
	To        string
	Decorator func(format string, a ...interface{}) string
}

type AsciDAG struct {
	nodes       []*Node
	edges       []*Edge
	nodeElement string
}

func NewAsciDAGRender(nodeElement string) *AsciDAG {
	if len(nodeElement) == 0 {
		nodeElement = "oo"
	}

	return &AsciDAG{
		nodeElement: nodeElement,
	}
}

func (d *AsciDAG) AddNode(node *Node) {
	for _, n := range d.nodes {
		if n.Name == node.Name {
			return
		}
	}

	d.nodes = append(d.nodes, node)
}

func (d *AsciDAG) AddEdge(edge *Edge) {
	if edge.From == edge.To {
		return
	}

	for _, e := range d.edges {
		if e.From == edge.From && e.To == edge.To {
			return
		}

		if e.From == edge.To && e.To == edge.From {
			return
		}
	}

	d.edges = append(d.edges, edge)
}

func (d *AsciDAG) leveledNodes() [][]*Node {
	touched := make(map[string]struct{})

	var leveled [][]*Node
	for i := 0; i < len(d.nodes); i++ {
		var nodes []*Node
		for _, node := range d.nodes {
			if _, ok := touched[node.Name]; ok {
				continue
			}

			ready := true
			for _, edge := range d.edges {
				if edge.To == node.Name {
					if _, ok := touched[edge.From]; !ok {
						ready = false
						break
					}
				}
			}
			if !ready {
				continue
			}

			nodes = append(nodes, node)
		}

		if len(nodes) > 0 {
			leveled = append(leveled, nodes)
			for _, n := range nodes {
				touched[n.Name] = struct{}{}
			}
		} else {
			break
		}
	}

	return leveled
}

type Position struct {
	X int
	Y int
}

func (d *AsciDAG) Render() {
	//width, height := common.TermSize()
	width, height := 110, 80
	if width == 0 || height == 0 {
		console.Error("Get terminal width/height error")
		return
	}

	leveledNodes := d.leveledNodes()
	maxLevelSize := 0
	for _, level := range leveledNodes {
		if len(level) > maxLevelSize {
			maxLevelSize = len(level)
		}
	}

	// Create canvas for the asci DAG
	xStep := 16
	yStep := 3
	displayHeight := maxLevelSize*yStep - 1
	displayWidth := (len(leveledNodes)-1)*xStep + len(leveledNodes)
	canvas := make([][]string, displayHeight)
	for i := range canvas {
		canvas[i] = make([]string, displayWidth)
		for j := 0; j < displayWidth; j++ {
			canvas[i][j] = " "
		}
	}

	// Draw nodes in canvas
	nodesPositions := make(map[string]Position)
	for level, nodes := range leveledNodes {
		var x, y int
		x = (xStep + 1) * level

		// Put node placeholder in each row, so that edges can align correctly
		for y := 0; y < displayHeight; y++ {
			// Trim space for emoji
			canvas[y][x] = strings.Repeat(" ", runewidth.StringWidth(strings.TrimSpace(d.nodeElement)))
		}

		for i, n := range nodes {
			y = i * yStep
			canvas[y][x] = d.nodeElement
			nodesPositions[n.Name] = Position{
				X: x,
				Y: y,
			}
		}
	}

	// Draw edges in canvas
	for _, edge := range d.edges {
		from, ok1 := nodesPositions[edge.From]
		to, ok2 := nodesPositions[edge.To]
		if !ok1 || !ok2 {
			break
		}

		d.drawEdge(from, to, canvas)
	}

	for _, row := range canvas {
		fmt.Println(strings.Join(row, ""))
	}
}

func (d *AsciDAG) drawEdge(from, to Position, canvas [][]string) {
	// Draw vertical line
	if from.X == to.X {
		x := from.X
		y1, y2 := from.Y, to.Y
		if y1 < y2 {
			for y := y1; y < y2; y++ {
				canvas[y][x] = "│"
			}
			canvas[y2-1][x] = "v"
		} else {
			for y := y2; y < y1; y++ {
				canvas[y][x] = "│"
			}
			canvas[y2+1][x] = "^"
		}
		return
	}

	// Draw horizontal line
	if from.Y == to.Y {
		y := from.Y
		for x := from.X + 2; x < to.X-1; x++ {
			canvas[y][x] = "─"
		}
		canvas[y][to.X-2] = ">"
	}

	// Draw polyline
	// ───╮
	//    │
	//    │
	//    ╰─────>
	if from.Y < to.Y {
		turningX := from.X + (to.X-from.X-1)*4/10 + 1
		for x := from.X + 2; x < turningX; x++ {
			canvas[from.Y][x] = "─"
		}
		for y := from.Y; y < to.Y; y++ {
			canvas[y][turningX] = "│"
		}
		for x := turningX; x < to.X-1; x++ {
			canvas[to.Y][x] = "─"
		}
		canvas[to.Y][to.X-2] = ">"
		canvas[from.Y][turningX] = "╮"
		canvas[to.Y][turningX] = "╰"
	}

	// Draw polyline
	//    ╭─────>
	//    │
	//    │
	// ───╯
	if from.Y > to.Y {
		turningX := from.X + (to.X-from.X-1)*4/10 + 1
		for x := from.X + 2; x < turningX; x++ {
			canvas[from.Y][x] = "─"
		}
		for y := to.Y; y < from.Y; y++ {
			canvas[y][turningX] = "│"
		}
		for x := turningX; x < to.X-1; x++ {
			canvas[to.Y][x] = "─"
		}
		canvas[to.Y][to.X-2] = ">"
		canvas[from.Y][turningX] = "╯"
		canvas[to.Y][turningX] = "╭"
	}
}
