package euler

import (
	"github.com/stretchr/testify/assert"
	"gitlab.lrz.de/hm/goal-core/hmgraph"
	"testing"
)

func TestSmallCase(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(8)
	vs[0].CreateEdges([]*hmgraph.Vertex{vs[1], vs[5]})
	vs[3].CreateEdges([]*hmgraph.Vertex{vs[2], vs[4]})
	vs[5].CreateEdges([]*hmgraph.Vertex{vs[3], vs[1]})
	vs[6].CreateEdges([]*hmgraph.Vertex{vs[4], vs[2]})

	start, edges, err := EulerPath(g)
	assert.Nil(t, err)
	assertEulerPath(t, g, start, edges)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func TestUnconnectedCase(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(10)
	vs[0].CreateEdges([]*hmgraph.Vertex{vs[1], vs[5]})
	vs[3].CreateEdges([]*hmgraph.Vertex{vs[2], vs[4]})
	vs[5].CreateEdges([]*hmgraph.Vertex{vs[3], vs[1]})
	vs[6].CreateEdges([]*hmgraph.Vertex{vs[4], vs[2]})
	vs[7].CreateEdges([]*hmgraph.Vertex{vs[8], vs[9]})
	vs[8].CreateEdge(vs[9])

	_, _, err := EulerPath(g)
	assert.NotNil(t, err, "no error on unconnected edges")
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func TestNotEulerianCase(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(6)
	vs[0].CreateEdges([]*hmgraph.Vertex{vs[1], vs[2]})
	vs[3].CreateEdges([]*hmgraph.Vertex{vs[0], vs[4]})
	vs[1].CreateEdges([]*hmgraph.Vertex{vs[5], vs[3]})

	_, _, err := EulerPath(g)
	assert.NotNil(t, err, "no error despite non-Eulerian graph")
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func TestComplexCase(t *testing.T) {
	g := hmgraph.NewGraph()
	vs := g.CreateVertices(8)
	vs[0].CreateEdges([]*hmgraph.Vertex{vs[2], vs[4], vs[3], vs[7]})
	vs[1].CreateEdges([]*hmgraph.Vertex{vs[3], vs[2], vs[6]})
	vs[2].CreateEdges([]*hmgraph.Vertex{vs[3], vs[4]})
	vs[3].CreateEdge(vs[7])
	vs[4].CreateEdges([]*hmgraph.Vertex{vs[7], vs[5]})
	vs[4].CreateEdges([]*hmgraph.Vertex{vs[7], vs[6]})

	start, edges, err := EulerPath(g)
	assert.Nil(t, err)
	assertEulerPath(t, g, start, edges)
	x, y, z := g.MapCount()
	assert.True(t, x+y+z == 0, "Not all maps disposed.")
}

func assertEulerPath(t *testing.T, g *hmgraph.Graph, start *hmgraph.Vertex, edges []*hmgraph.Edge) {
	used := hmgraph.CreateEdgeMap(g, "used", false)
	defer used.Dispose()

	assert.True(t, len(edges) == g.EdgeCount(), "not all edges are used")
	current := start
	for _, edge := range edges {
		if !edge.IsIncident(current) {
			t.Logf("edges don't form a cycle.")
			t.FailNow()
		}
		assert.False(t, used.Get(edge), "edge used twice")
		used.Set(edge, true)
		current = edge.Opposite(current)
	}
}
