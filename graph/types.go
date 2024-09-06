package graph

import "errors"

var ErrSelfLooped = errors.New("self looped")
var ErrParallelEdged = errors.New("parallel edged")
var ErrAlreadyExists = errors.New("already exists")

type IGraph interface {
	IsDirected() bool
	AllowSelfLoop() bool
	AllowParallelEdge() bool
	Nodes() []INode
	Edges() []IEdge
	AddNode(n INode) bool
	AddEdge(u, v INode, e IEdge) error
	AdjacentNodes(INode) []INode
	PredecessorNodes(INode) []INode
	SuccessorNodes(INode) []INode
	IncidentEdges(INode) []IEdge
	InEdges(INode) []IEdge
	OutEdges(INode) []IEdge
	Degree(INode) int
	InDegree(INode) int
	OutDegree(INode) int
	EdgesOf(u, v INode) []IEdge
	EdgeCount(u, v INode) int
	IncidentNodes(IEdge) [2]INode
	AdjacentEdges(IEdge) []IEdge
}
type identify interface {
}
type INode interface {
	identify
	nodeId() int
}
type IEdge interface {
	identify
	edgeId() int
}
