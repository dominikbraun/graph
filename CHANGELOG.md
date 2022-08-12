# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.0] - 2022-08-16

### Added
* Added the `Graph.AddVertex` method for adding a vertex. This replaces `Graph.Vertex`.
* Added the `Graph.AddEdge` method for creating an edge. This replaces `Graph.Edge`.
* Added the `Graph.Vertex` method for retrieving a vertex by its hash. This is not to be confused with the old `Graph.Vertex` function for adding vertices that got replaced with `Graph.AddVertex`.
* Added the `Graph.Edge` method for retrieving an edge. This is not to be confused with the old `Graph.Edge` function for creating an edge that got replaced with `Graph.AddEdge`.
* Added the `Graph.Predecessors` function for retrieving a vertex' predecessors.
* Added the `DFS` function.
* Added the `BFS` function.
* Added the `CreatesCycle` function.
* Added the `StronglyConnectedComponents` function.
* Added the `ShortestPath` function.

### Removed
* Removed the `Graph.EdgeByHashes` method. Use `Graph.AddEdge` instead.
* Removed the `Graph.GetEdgeByHashes` method. Use `Graph.Edge` instead.
* Removed the `Graph.DegreeByHash` method. Use `Graph.Degree` instead.
* Removed the `Graph.DFS` and `Graph.DFSByHash` methods. Use `DFS` instead.
* Removed the `Graph.BFS` and `Graph.BFSByHash` methods. Use `BFS` instead.
* Removed the `Graph.CreatesCycle` and `Graph.CreatesCycleByHashes` methods. Use `CreatesCycle` instead.
* Removed the `Graph.StronglyConnectedComponents` method. Use `StronglyConnectedComponents` instead.
* Removed the `Graph.ShortestPath` and `Graph.ShortestPathByHash` methods. Use `ShortestPath` instead.

## [0.8.0] - 2022-08-01

### Added
* Added the `EdgeWeight` and `EdgeAttribute` functional options.
* Added the `Properties` field to `Edge`.

### Changed
* Changed `Edge` to accept a variadic `options` parameter.
* Changed `EdgeByHashes` to accept a variadic `options` parameter.
* Renamed `draw.Graph` to `draw.DOT` for more clarity regarding the rendering format.

### Removed
* Removed the `WeightedEdge` function. Use `Edge` with the `EdgeWeight` functional option instead.
* Removed the `WeightedEdgeByHashes` function. Use `EdgeByHashes` with the `EdgeWeight` functional option instead.

### Fixed
* Fixed missing edge attributes when drawing a graph using `draw.DOT`.

## [0.7.0] - 2022-07-26

### Added
* Added `draw` package for graph visualization using DOT-compatible renderers.
* Added `Traits` function for retrieving the graph's traits.

## [0.6.0] - 2022-07-22

### Added
* Added `AdjacencyMap` function for retrieving an adjancency map for all vertices.

### Removed
* Removed the `AdjacencyList` function.

## [0.5.0] - 2022-07-21

### Added
* Added `AdjacencyList` function for retrieving an adjacency list for all vertices.
  
### Changed
* Updated the examples in the documentation.

## [0.4.0] - 2022-07-01

### Added
* Added `ShortestPath` function for computing shortest paths.

### Changed
* Changed the term "properties" to "traits" in the code and documentation.
* Don't traverse all vertices in disconnected graphs by design.

## [0.3.0] - 2022-06-27

### Added
* Added `StronglyConnectedComponents` function for detecting SCCs.
* Added various images to usage examples.

## [0.2.0] - 2022-06-20

### Added
* Added `Degree` and `DegreeByHash` functions for determining vertex degrees.
* Added cycle checks when adding an edge using the `Edge` functions.

## [0.1.0] - 2022-06-19

### Added
* Added `CreatesCycle` and `CreatesCycleByHashes` functions for predicting cycles.

## [0.1.0-beta] - 2022-06-17

### Changed
* Introduced dedicated types for directed and undirected graphs, making `Graph[K, T]` an interface.

## [0.1.0-alpha] - 2022-06-13

### Added
* Introduced core types and methods.
