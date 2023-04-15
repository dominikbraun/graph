package graph

import (
	"fmt"
)

// Union combines two given graphs into a new graph. The vertex hashes in both
// graphs are expected to be unique. The two input graphs will remain unchanged.
//
// Both graphs should be either directed or undirected. All traits for the new
// graph will be derived from g.
func Union[K comparable, T any](g, h Graph[K, T]) (Graph[K, T], error) {
	union, err := g.Clone()
	if err != nil {
		return union, fmt.Errorf("failed to clone g: %w", err)
	}

	adjacencyMap, err := h.AdjacencyMap()
	if err != nil {
		return union, fmt.Errorf("failed to get adjacency map: %w", err)
	}

	addedEdges := make(map[K]map[K]struct{})

	for currentHash, _ := range adjacencyMap {
		vertex, properties, err := h.VertexWithProperties(currentHash)
		if err != nil {
			return union, fmt.Errorf("failed to get vertex %v: %w", currentHash, err)
		}

		err = union.AddVertex(vertex, copyVertexProperties(properties))
		if err != nil {
			return union, fmt.Errorf("failed to add vertex %v: %w", currentHash, err)
		}
	}

	for _, adjacencies := range adjacencyMap {
		for _, edge := range adjacencies {
			if _, sourceOK := addedEdges[edge.Source]; sourceOK {
				if _, targetOK := addedEdges[edge.Source][edge.Target]; targetOK {
					// If the edge addedEdges[source][target] exists, the edge
					// has already been created and thus can be skipped here.
					continue
				}
			}

			err = union.AddEdge(edge.Source, edge.Target, copyEdgeProperties(edge.Properties))
			if err != nil {
				return union, fmt.Errorf("failed to add edge (%v, %v): %w", edge.Source, edge.Target, err)
			}

			if _, ok := addedEdges[edge.Source]; !ok {
				addedEdges[edge.Source] = make(map[K]struct{})
			}
			addedEdges[edge.Source][edge.Target] = struct{}{}
		}
	}

	return union, nil
}

func copyVertexProperties(source VertexProperties) func(*VertexProperties) {
	return func(p *VertexProperties) {
		for k, v := range source.Attributes {
			p.Attributes[k] = v
		}
		p.Weight = source.Weight
	}
}

func copyEdgeProperties(source EdgeProperties) func(properties *EdgeProperties) {
	return func(p *EdgeProperties) {
		for k, v := range source.Attributes {
			p.Attributes[k] = v
		}
		p.Weight = source.Weight
		p.Data = source.Data
	}
}
