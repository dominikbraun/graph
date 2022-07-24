// Package draw provides functions for visualizing graph structures. At this time, draw supports
// the DOT language which can be interpreted by Graphviz, Grappa, and others.
package draw

import (
	"fmt"
	"io"
	"text/template"

	"github.com/dominikbraun/graph"
)

const dotTemplate = `strict {{.GraphType}} {
{{range $s := .Statements}}
	{{.Source}} {{if .Target}}{{$.EdgeOperator}} {{.Target}} [ weight={{.Weight}}, label="{{.Label}}" ]{{end}};
{{end}}
}
`

type description struct {
	GraphType    string
	EdgeOperator string
	Statements   []statement
}

type statement struct {
	Source interface{}
	Target interface{}
	Weight int
	Label  string
}

// Graph renders the given graph structure in DOT language into an io.Writer, for example a file.
// The generated output can be passed to Graphviz or other visualization tools supporting DOT.
//
// The following example renders a directed graph into a file mygraph.gv:
//
//	g := graph.New(graph.IntHash, graph.Directed())
//
//	g.Vertex(1)
//	g.Vertex(2)
//	g.Vertex(3)
//
//	_ = g.Edge(1, 2)
//	_ = g.Edge(1, 3)
//
//	file, _ := os.Create("./mygraph.gv")
//	_ = graph.Draw(g, file)
//
// To generate an SVG from the created file using Graphviz, use a command such as the following:
//
//	dot -Tsvg -O mygraph.gv
//
// Another possibility is to use os.Stdin as an io.Writer, print the DOT output to stdin, and pipe
// it as follows:
//
//	go run main.go | dot -Tsvg > output.svg
func Graph[K comparable, T any](g graph.Graph[K, T], w io.Writer) error {
	description := generateDOT(g)

	return renderDOT(w, description)
}

func generateDOT[K comparable, T any](g graph.Graph[K, T]) description {
	description := description{
		GraphType:    "graph",
		EdgeOperator: "--",
		Statements:   make([]statement, 0),
	}

	if g.Traits().IsDirected {
		description.GraphType = "digraph"
		description.EdgeOperator = "->"
	}

	for vertex, adjacencies := range g.AdjacencyMap() {
		if len(adjacencies) == 0 {
			statement := statement{
				Source: vertex,
			}
			description.Statements = append(description.Statements, statement)
			continue
		}

		for adjacency, edge := range adjacencies {
			statement := statement{
				Source: vertex,
				Target: adjacency,
				Weight: edge.Weight,
				Label:  edge.Label,
			}
			description.Statements = append(description.Statements, statement)
		}
	}

	return description
}

func renderDOT(w io.Writer, d description) error {
	tpl, err := template.New("dotTemplate").Parse(dotTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.Execute(w, d)
}
