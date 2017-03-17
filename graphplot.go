// package graphplot is used for plotting gonum graph
package graphplot

import (
	"image/color"
	"math"

	"github.com/gonum/graph"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg/draw"
)

type XY struct {
	X, Y float64
}

type NodeLocationer interface {
	NodeLocation(n graph.Node) XY
}

type NodeStyler interface {
	NodeStyle(n graph.Node) draw.GlyphStyle
}

type EdgeStyler interface {
	EdgeStyle(e graph.Edge) draw.LineStyle
}

// NodeCircle plots nodes in a circle
// TODO(btracey): Lots of false assumptions
type NodeCircle struct {
	Nodes int
}

func (n NodeCircle) NodeLocation(node graph.Node) XY {
	id := node.ID()
	return XY{
		X: math.Cos(float64(id) / float64(n.Nodes) * 2 * math.Pi),
		Y: math.Sin(float64(id) / float64(n.Nodes) * 2 * math.Pi),
	}
}

// Probably need to add a "ChevronGlyph" or ArrowHeadGlyph

type GraphPlot struct {
	NodeLocationer
	NodeStyler
	EdgeStyler

	// Add something like NodeLabel and EdgeLabel?
	g graph.Graph
}

// NewGraphPlot implements the plotter interface. It does not copy the current
// graph, but retains a reference to it.
// Currently only works when there is only one (directed) edge between u and v.
func NewDirected(g graph.Graph) *GraphPlot {
	return &GraphPlot{g: g}
}

func (p *GraphPlot) Plot(c draw.Canvas, plt *plot.Plot) {
	graph := p.g
	nodes := graph.Nodes()
	nNodes := len(nodes)
	loc := p.NodeLocationer
	if loc == nil {
		loc = NodeCircle{nNodes}
	}

	// First, loop over all the edges and add them to the plot
	for _, start := range nodes {
		ends := graph.From(start)
		startxy := loc.NodeLocation(start)
		for _, end := range ends {
			endxy := loc.NodeLocation(end)
			edge := graph.Edge(start, end)
			xys := plotter.XYs{startxy, endxy}
			//fmt.Println(xys)
			l, err := plotter.NewLine(xys)
			if err != nil {
				panic(err)
			}
			//fmt.Println("p = ", p)
			if p.EdgeStyler != nil {
				l.LineStyle = p.EdgeStyle(edge)
			} else {
				l.LineStyle.Color = color.RGBA{R: 100, G: 100, B: 100, A: 55}
			}
			l.Plot(c, plt)
			if err != nil {
				panic(err)
			}
		}
	}

	// Now add the Nodes to the plot
	for _, node := range nodes {
		xy := loc.NodeLocation(node)
		scat, err := plotter.NewScatter(plotter.XYs{xy})
		if err != nil {
			panic(err)
		}
		if p.NodeStyler == nil {
			scat.GlyphStyle.Shape = draw.CircleGlyph{}
			scat.GlyphStyle.Color = color.RGBA{}
		} else {
			scat.GlyphStyle = p.NodeStyle(node)
		}
		scat.Plot(c, plt)
	}
}
