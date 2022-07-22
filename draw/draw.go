package draw

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/dominikbraun/graph"
)

const dotTemplate = `
strict {{.GraphType}} {
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

func Graph[K comparable, T any](g graph.Graph[K, T], options ...func(*config)) error {
	c := defaultConfig()

	for _, option := range options {
		option(&c)
	}

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
			statement := statement{
				Source: vertex,
			}
			desc.Statements = append(desc.Statements, statement)
			continue
		}

		for adjacency, edge := range adjacencies {
			statement := statement{
				Source: vertex,
				Target: adjacency,
				Weight: edge.Weight,
				Label:  edge.Label,
			}
			desc.Statements = append(desc.Statements, statement)
		}
	}

	return renderDOT(desc, &c)
}

func renderDOT(data description, c *config) error {
	tpl, _ := template.New("dotTemplate").Parse(dotTemplate)

	if c.writer != nil {
		return tpl.Execute(c.writer, data)
	}

	name := filepath.Join(c.directory, c.filename)

	file, err := os.Create(name)
	if err != nil {
		return err
	}

	return tpl.Execute(file, data)
}
