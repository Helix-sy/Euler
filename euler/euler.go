package euler

import (
	"errors"
	"gitlab.lrz.de/hm/goal-core/hmgraph"
)

// UndirectedEulerPath computes an UndirectedEulerPath for an undirected graph
// if the graph is at least semi-Eulerian. Returns an error
// if the graph contains any Arcs or is not semi-Eulerian.
func UndirectedEulerPath(g *hmgraph.Graph) (start *hmgraph.Vertex, edges []*hmgraph.Edge, err error) {
	// Check if the graph contains any arcs
	if g.ArcCount() > 0 {
		return nil, nil, errors.New("graph contains arcs")
	}

	// Find a suitable starting vertex
	start, err = findStartVertex(g)
	if err != nil {
		return nil, nil, err
	}

	// Create a map to track visited edges
	visited := hmgraph.CreateEdgeMap(g, "visited", false)
	defer visited.Dispose()

	// Compute the Euler path
	path := make([]*hmgraph.Edge, 0, g.EdgeCount())

	findPath(start, nil, visited, &path)

	// Reverse the path to get the correct order
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	// Validate whether we used all edges
	if len(path) != g.EdgeCount() {
		return nil, nil, errors.New("graph is not semi-Eulerian")
	}

	return start, path, nil
}

// findStartVertex finds a suitable starting vertex for an Euler path
func findStartVertex(g *hmgraph.Graph) (*hmgraph.Vertex, error) {
	oddDegreeVertices := 0
	var candidateStart *hmgraph.Vertex

	// Check if graph is empty
	if g.VertexCount() == 0 {
		return nil, errors.New("graph has no vertices")
	}

	// First pass to count odd degree vertices and find a candidate start
	vertices := g.GetVertices()
	for _, v := range vertices {
		if v.Degree()%2 != 0 {
			oddDegreeVertices++
			candidateStart = v
		}
	}

	// For an Euler path, we should have 0 or 2 vertices with odd degree
	if oddDegreeVertices == 0 {
		// Eulerian circuit - any vertex can be the start
		var firstVertex *hmgraph.Vertex
		for _, v := range vertices {
			// Pick any vertex with edges
			if v.Degree() > 0 {
				firstVertex = v
				break
			}
		}

		// If no vertex with edges was found, the graph is either empty or has isolated vertices
		if firstVertex == nil {
			return nil, errors.New("graph has no edges")
		}

		return firstVertex, nil
	} else if oddDegreeVertices == 2 {
		// Semi-Eulerian path - start from one of the odd-degree vertices
		return candidateStart, nil
	} else {
		// Not Eulerian or semi-Eulerian
		return nil, errors.New("graph is not semi-Eulerian")
	}
}

// findPath implements the depth-first search approach from cp-algorithms.com
// but without modifying the graph structure
func findPath(current *hmgraph.Vertex, edge *hmgraph.Edge, visited *hmgraph.EdgeMap[bool], path *[]*hmgraph.Edge) {
	// Iterate through all edges of the current vertex
	for _, e := range current.GetEdges() {
		if !visited.Get(e) {
			visited.Set(e, true)

			// Get the next vertex
			next := e.Opposite(current)

			// Recursive DFS
			findPath(next, e, visited, path)

			// Add the edge to the path
			*path = append(*path, e)
		}
	}
}
