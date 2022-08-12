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
	{{.Source}} {{if .Target}}{{$.EdgeOperator}} {{.Target}} [ {{range $k, $v := .Attributes}}{{$k}}="{{$v}}", {{end}} weight={{.Weight}} ]{{end}};
{{end}}
}
`

type description struct {
	GraphType    string
	EdgeOperator string
	Statements   []statement
}

type statement struct {
	Source     interface{}
	Target     interface{}
	Attributes map[string]string
	Weight     int
}

// DOT renders the given graph structure in DOT language into an io.Writer, for example a file. The
// generated output can be passed to Graphviz or other visualization tools supporting DOT.
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
// Another possibility is to use os.Stdout as an io.Writer, print the DOT output to stdout, and
// pipe it as follows:
//
//	go run main.go | dot -Tsvg > output.svg
func DOT[K comparable, T any](g graph.Graph[K, T], w io.Writer) error {
	desc := generateDOT(g)

	return renderDOT(w, desc)
}

func generateDOT[K comparable, T any](g graph.Graph[K, T]) description {
	desc := description{
		GraphType:    "graph",
		EdgeOperator: "--",
		Statements:   make([]statement, 0),
	}

	if g.Traits().IsDirected {
		desc.GraphType = "digraph"
		desc.EdgeOperator = "->"
	}

	for vertex, adjacencies := range g.AdjacencyMap() {
		if len(adjacencies) == 0 {
			stmt := statement{
				Source: vertex,
			}
			desc.Statements = append(desc.Statements, stmt)
			continue
		}

		for adjacency, edge := range adjacencies {
			stmt := statement{
				Source:     vertex,
				Target:     adjacency,
				Weight:     edge.Properties.Weight,
				Attributes: edge.Properties.Attributes,
			}
			desc.Statements = append(desc.Statements, stmt)
		}
	}

	return desc
}

func renderDOT(w io.Writer, d description) error {
	tpl, err := template.New("dotTemplate").Parse(dotTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.Execute(w, d)
}
