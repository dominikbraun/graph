package draw

import "github.com/dominikbraun/graph"

const dotTemplate = `
{{.GraphType}} {
	{{range $s := .Statements}}
		{{$s.Source}} {{.EdgeOperator}} {{$s.Target}} [ weight={{$s.Weight}}, label="{{$s.Label}}" ];
	{{end}}
}
`

type dot struct {
	GraphType    string
	EdgeOperator string
	Statements   []dotStatement
}

type dotStatement struct {
	Source string
	Target string
	Weight int
	Label  string
}

func (d dotStatement) Attributes()

func Graph[K comparable, T any](g graph.Graph[K, T], options ...func(*config)) {
	c := defaultConfig()

	for _, option := range options {
		option(&c)
	}

}
