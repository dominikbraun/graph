# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] - 2020-07-22

### Added
* Added `AdjacencyMap` function for retrieving an adjancency map for all vertices.

### Removed
* Removed the `AdjacencyList` function.

## [0.5.0] - 2020-07-21

### Added
* Added `AdjacencyList` function for retrieving an adjacency list for all vertices.
  
### Changed
* Updated the examples in the documentation.

## [0.4.0] - 2020-07-01

### Added
* Added `ShortestPath` function for computing shortest paths.

### Changed
* Changed the term "properties" to "traits" in the code and documentation.
* Don't traverse all vertices in disconnected graphs by design.

## [0.3.0] - 2020-06-27

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
