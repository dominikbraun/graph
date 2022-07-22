package draw

import (
	"os"
	"text/template"

	"github.com/dominikbraun/graph"
)

const graphTemplate = `
graph {
{{range $s := .Statements}}
	{{.Source}} {{if .Target}}-- {{.Target}} [ weight={{.Weight}}, label="{{.Label}}" ]{{end}};
{{end}}
}
`

const digraphTemplate = `
digraph {
{{range $s := .Statements}}
	{{.Source}} {{if .Target}}-> {{.Target}} [ weight={{.Weight}}, label="{{.Label}}" ]{{end}};
{{end}}
}
`

type dot struct {
	GraphType  string
	Statements []dotStatement
}

type dotStatement struct {
	Source       interface{}
	Target       interface{}
	EdgeOperator string
	Weight       int
	Label        string
}

func Graph[K comparable, T any](g graph.Graph[K, T], options ...func(*config)) {
	c := defaultConfig()

	for _, option := range options {
		option(&c)
	}

	data := dot{
		Statements: make([]dotStatement, 0),
	}

	isUndirected := false

	for vertex, adjacencies := range g.AdjacencyMap() {
		if len(adjacencies) == 0 {
			statement := dotStatement{
				Source: vertex,
			}
			data.Statements = append(data.Statements, statement)
			continue
		}

		for adjacency, edge := range adjacencies {
			statement := dotStatement{
				Source: vertex,
				Target: adjacency,
				Weight: edge.Weight,
				Label:  edge.Label,
			}
			data.Statements = append(data.Statements, statement)
		}
	}

	textTemplate := digraphTemplate

	if isUndirected {
		textTemplate = graphTemplate
	}

	tpl, _ := template.New("dot").Parse(textTemplate)

	if err := tpl.Execute(os.Stdout, data); err != nil {
		panic(err)
	}
}
