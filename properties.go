package graph

type Properties struct {
	isDirected bool
	isAcyclic  bool
	isWeighted bool
	isRooted   bool
}

func Directed() func(*Properties) {
	return func(p *Properties) {
		p.isDirected = true
	}
}

func Acyclic() func(*Properties) {
	return func(p *Properties) {
		p.isAcyclic = true
	}
}

func Weighted() func(*Properties) {
	return func(p *Properties) {
		p.isWeighted = true
	}
}

func Rooted() func(*Properties) {
	return func(p *Properties) {
		p.isRooted = true
	}
}

func Tree() func(*Properties) {
	return func(p *Properties) {
		Acyclic()(p)
		Rooted()(p)
	}
}
