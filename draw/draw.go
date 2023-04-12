// Package draw provides functions for visualizing graph structures. At this
// time, draw supports the DOT language which can be interpreted by Graphviz,
// Grappa, and others.
package draw

import (
	"fmt"
	"io"
	"text/template"

	"github.com/dominikbraun/graph"
)

// ToDo: This template should be simplified and split into multiple templates.
const dotTemplate = `strict {{.GraphType}} {
{{range $k, $v := .Attributes}}
	{{$k}}="{{$v}}";
{{end}}
{{range $s := .Statements}}
	"{{.Source}}" {{if .Target}}{{$.EdgeOperator}} "{{.Target}}" [ {{range $k, $v := .EdgeAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.EdgeWeight}} ]{{else}}[ {{range $k, $v := .SourceAttributes}}{{$k}}="{{$v}}", {{end}} weight={{.SourceWeight}} ]{{end}};
{{end}}
}
`

type description struct {
	GraphType    string
	Attributes   map[string]string
	EdgeOperator string
	Statements   []statement
}

type statement struct {
	Source           interface{}
	Target           interface{}
	SourceWeight     int
	SourceAttributes map[string]string
	EdgeWeight       int
	EdgeAttributes   map[string]string
}

// DOT renders the given graph structure in DOT language into an io.Writer, for
// example a file. The generated output can be passed to Graphviz or other
// visualization tools supporting DOT.
//
// The following example renders a directed graph into a file my-graph.gv:
//
//	g := graph.New(graph.IntHash, graph.Directed())
//
//	_ = g.AddVertex(1)
//	_ = g.AddVertex(2)
//	_ = g.AddVertex(3, graph.VertexAttribute("style", "filled"), graph.VertexAttribute("fillcolor", "red"))
//
//	_ = g.AddEdge(1, 2, graph.EdgeWeight(10), graph.EdgeAttribute("color", "red"))
//	_ = g.AddEdge(1, 3)
//
//	file, _ := os.Create("./my-graph.gv")
//	_ = draw.DOT(g, file)
//
// To generate an SVG from the created file using Graphviz, use a command such
// as the following:
//
//	dot -Tsvg -O my-graph.gv
//
// Another possibility is to use os.Stdout as an io.Writer, print the DOT output
// to stdout, and pipe it as follows:
//
//	go run main.go | dot -Tsvg > output.svg
//
// DOT also accepts the [GraphAttribute] functional option, which can be used to
// add global attributes when rendering the graph:
//
//	_ = draw.DOT(g, file, draw.GraphAttribute("label", "my-graph"))
func DOT[K comparable, T any](g graph.Graph[K, T], w io.Writer, options ...func(*description)) error {
	desc, err := generateDOT(g, options...)
	if err != nil {
		return fmt.Errorf("failed to generate DOT description: %w", err)
	}

	return renderDOT(w, desc)
}

// GraphAttribute is a functional option for the [DOT] method.
func GraphAttribute(key, value string) func(*description) {
	return func(d *description) {
		d.Attributes[key] = value
	}
}

func generateDOT[K comparable, T any](g graph.Graph[K, T], options ...func(*description)) (description, error) {
	desc := description{
		GraphType:    "graph",
		Attributes:   make(map[string]string),
		EdgeOperator: "--",
		Statements:   make([]statement, 0),
	}

	for _, option := range options {
		option(&desc)
	}

	if g.Traits().IsDirected {
		desc.GraphType = "digraph"
		desc.EdgeOperator = "->"
	}

	adjacencyMap, err := g.AdjacencyMap()
	if err != nil {
		return desc, err
	}

	for vertex, adjacencies := range adjacencyMap {
		_, sourceProperties, err := g.VertexWithProperties(vertex)
		if err != nil {
			return desc, err
		}

		stmt := statement{
			Source:           vertex,
			SourceWeight:     sourceProperties.Weight,
			SourceAttributes: sourceProperties.Attributes,
		}
		desc.Statements = append(desc.Statements, stmt)

		for adjacency, edge := range adjacencies {
			stmt := statement{
				Source:         vertex,
				Target:         adjacency,
				EdgeWeight:     edge.Properties.Weight,
				EdgeAttributes: edge.Properties.Attributes,
			}
			desc.Statements = append(desc.Statements, stmt)
		}
	}

	return desc, nil
}

func renderDOT(w io.Writer, d description) error {
	tpl, err := template.New("dotTemplate").Parse(dotTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	return tpl.Execute(w, d)
}
