package graphplot

import (
	"math"
	"testing"

	"github.com/gonum/graph/simple"
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg"
)

func TestGraphPlot(t *testing.T) {
	g := simple.NewDirectedGraph(1, math.NaN())
	for i := 0; i < 5; i++ {
		g.AddNode(simple.Node(i))
	}
	for i := 0; i < 4; i++ {
		edge := simple.Edge{
			F: simple.Node(i),
			T: simple.Node(i + 1),
			W: 1,
		}
		g.SetEdge(edge)
	}
	gp := NewDirected(g)
	plot, err := plot.New()
	if err != nil {
		panic("plot fail")
	}
	plot.Add(gp)
	plot.HideAxes()
	err = plot.Save(4*vg.Inch, 4*vg.Inch, "test.pdf")
	if err != nil {
		t.Fatalf(err.Error())
	}

}
