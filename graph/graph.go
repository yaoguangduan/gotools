package graph

import (
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
)

type Graph struct {
	gi *graphImpl
}

func Directed(allowSelfLoop, allowParallelEdge bool) IGraph {
	return createGraph(true, allowSelfLoop, allowParallelEdge)
}

func Undirected(allowSelfLoop, allowParallelEdge bool) IGraph {
	return createGraph(false, allowSelfLoop, allowParallelEdge)
}

func createGraph(directed, allowSelfLoop, allowParallelEdge bool) IGraph {
	gi := &graphImpl{make(map[int]nodeView), make(map[int]edgeView), directed, allowParallelEdge, allowSelfLoop, 0}
	return gi
}

func (g *graphImpl) IsDirected() bool {
	return g.isDirected
}
func (g *graphImpl) AllowSelfLoop() bool {
	return g.allowSelfLoop
}
func (g *graphImpl) AllowParallelEdge() bool {
	return g.allowParallelEdge
}
func (g *graphImpl) Nodes() []INode {
	nodes := make([]INode, 0, len(g.nodes))
	for _, node := range g.nodes {
		nodes = append(nodes, node.self)
	}
	return nodes
}
func (g *graphImpl) Edges() []IEdge {
	edges := make([]IEdge, 0, len(g.edges))
	for _, edge := range g.edges {
		edges = append(edges, edge.self)
	}
	return edges
}
func (g *graphImpl) AddNode(n INode) bool {
	_, b := g.addNodeInner(n)
	return b
}
func (g *graphImpl) addNodeInner(n INode) (nodeView, bool) {
	g.checkAndInit(n)
	old, ok := g.nodes[n.nodeId()]
	if !ok {
		nv := nodeView{self: n, incident: make(map[int]endpoint)}
		g.nodes[n.nodeId()] = nv
		return nv, true
	} else {
		return old, false
	}
}
func (g *graphImpl) AddEdge(u, v INode, e IEdge) error {
	uav, _ := g.addNodeInner(u)
	vav, _ := g.addNodeInner(v)
	if uav.self.nodeId() == vav.self.nodeId() && !g.allowSelfLoop {
		return ErrSelfLooped
	}
	if !g.allowParallelEdge {
		for _, end := range uav.incident {
			if end.peer.nodeId() == vav.self.nodeId() {
				return ErrParallelEdged
			}
		}
	}
	g.checkAndInit(e)
	_, ok := g.edges[e.edgeId()]
	if ok {
		return ErrAlreadyExists
	}
	ev := edgeView{e, u, v}
	g.edges[e.edgeId()] = ev
	uav.incident[e.edgeId()] = endpoint{e, v, false}
	vav.incident[e.edgeId()] = endpoint{e, u, true}
	return nil
}
func (g *graphImpl) AdjacentNodes(n INode) []INode {
	if g.invalid(n) {
		return nil
	}
	nv := g.nodes[n.nodeId()]
	var nodes = make(map[int]INode)
	for _, ev := range nv.incident {
		if ev.peer.nodeId() != n.nodeId() {
			nodes[ev.peer.nodeId()] = ev.peer
		}
	}
	return slices.Collect(maps.Values(nodes))
}
func (g *graphImpl) PredecessorNodes(n INode) []INode {
	if g.invalid(n) {
		return nil
	}
	if g.isDirected {
		nv := g.nodes[n.nodeId()]
		var nodes = make(map[int]INode)
		for _, ev := range nv.incident {
			if ev.in && ev.peer.nodeId() != n.nodeId() {
				nodes[ev.peer.nodeId()] = ev.peer
			}
		}
		return slices.Collect(maps.Values(nodes))
	} else {
		return g.AdjacentNodes(n)
	}
}
func (g *graphImpl) SuccessorNodes(n INode) []INode {
	if g.invalid(n) {
		return nil
	}
	if g.isDirected {
		nv := g.nodes[n.nodeId()]
		var nodes = make(map[int]INode)
		for _, ev := range nv.incident {
			if !ev.in && ev.peer.nodeId() != n.nodeId() {
				nodes[ev.peer.nodeId()] = ev.peer
			}
		}
		return slices.Collect(maps.Values(nodes))
	} else {
		return g.AdjacentNodes(n)
	}
}
func (g *graphImpl) IncidentEdges(n INode) []IEdge {
	if g.invalid(n) {
		return nil
	}
	var edges []IEdge
	nv := g.nodes[n.nodeId()]
	for _, ev := range nv.incident {
		edges = append(edges, ev.edge)
	}
	return edges
}
func (g *graphImpl) InEdges(n INode) []IEdge {
	if g.invalid(n) {
		return nil
	}
	if g.isDirected {
		var edges []IEdge
		nv := g.nodes[n.nodeId()]
		for _, ev := range nv.incident {
			if ev.in {
				edges = append(edges, ev.edge)
			}
		}
		return edges
	} else {
		return g.IncidentEdges(n)
	}
}
func (g *graphImpl) OutEdges(n INode) []IEdge {
	if g.invalid(n) {
		return nil
	}
	if g.isDirected {
		var edges []IEdge
		nv := g.nodes[n.nodeId()]
		for _, ev := range nv.incident {
			if !ev.in {
				edges = append(edges, ev.edge)
			}
		}
		return edges
	} else {
		return g.IncidentEdges(n)
	}
}
func (g *graphImpl) Degree(n INode) int {
	if g.invalid(n) {
		return 0
	}
	nv := g.nodes[n.nodeId()]
	var degree int
	for _, ev := range nv.incident {
		if ev.peer.nodeId() == n.nodeId() {
			degree += 2
		} else {
			degree++
		}
	}
	return degree
}
func (g *graphImpl) InDegree(n INode) int {
	if g.invalid(n) {
		return 0
	}
	if g.isDirected {
		var degree int
		nv := g.nodes[n.nodeId()]
		for _, ev := range nv.incident {
			if ev.in {
				degree++
			}
		}
		return degree
	}
	return g.Degree(n)
}
func (g *graphImpl) OutDegree(n INode) int {
	if g.invalid(n) {
		return 0
	}
	if g.isDirected {
		var degree int
		nv := g.nodes[n.nodeId()]
		for _, ev := range nv.incident {
			if !ev.in {
				degree++
			}
		}
		return degree
	}
	return g.Degree(n)
}
func (g *graphImpl) EdgesOf(u, v INode) []IEdge {
	if g.invalid(u) || g.invalid(v) {
		return nil
	}
	var edges []IEdge
	ui := g.nodes[u.nodeId()]
	for _, end := range ui.incident {
		if end.peer.nodeId() == v.nodeId() {
			edges = append(edges, end.edge)
		}
	}
	return edges
}
func (g *graphImpl) EdgeCount(u, v INode) int {
	return len(g.EdgesOf(u, v))
}
func (g *graphImpl) IncidentNodes(e IEdge) [2]INode {
	if g.invalid(e) {
		return [2]INode{}
	}
	ev := g.edges[e.edgeId()]
	return [2]INode{ev.u, ev.v}
}
func (g *graphImpl) AdjacentEdges(e IEdge) []IEdge {
	if g.invalid(e) {
		return nil
	}
	ev := g.edges[e.edgeId()]
	uv := g.nodes[ev.u.nodeId()]
	vv := g.nodes[ev.v.nodeId()]
	var edges = map[int]IEdge{}
	for _, end := range uv.incident {
		edges[end.edge.edgeId()] = end.edge
	}
	for _, end := range vv.incident {
		edges[end.edge.edgeId()] = end.edge
	}
	delete(edges, e.edgeId())
	return slices.Collect(maps.Values(edges))
}

type endpoint struct {
	edge IEdge
	peer INode
	in   bool
}
type nodeView struct {
	self     INode
	incident map[int]endpoint
}
type edgeView struct {
	self IEdge
	u    INode
	v    INode
}

func (ev *edgeView) eq(another *edgeView) bool {
	return ev.self.edgeId() == another.self.edgeId() && ev.u.nodeId() == another.u.nodeId() && ev.v.nodeId() == another.v.nodeId()
}

type graphImpl struct {
	nodes             map[int]nodeView
	edges             map[int]edgeView
	isDirected        bool
	allowParallelEdge bool
	allowSelfLoop     bool
	idSeq             int
}

func (g *graphImpl) checkAndGet(n identify) (reflect.Value, bool) {
	ptr := reflect.ValueOf(n)
	if ptr.Kind() != reflect.Ptr {
		panic("n must be a pointer")
	}
	nOrE := ptr.Elem()
	typ := nOrE.Type()
	_, exist := typ.FieldByName("INode")
	if exist {
		in := nOrE.FieldByName("INode")
		if !in.CanSet() {
			panic("cannot set INode")
		} else {
			return in, true
		}
	} else {
		_, exist = typ.FieldByName("IEdge")
		if exist {
			in := nOrE.FieldByName("IEdge")
			if !in.CanSet() {
				panic("cannot set IEdge")
			} else {
				return in, false
			}
		} else {
			panic("must impl INode or IEdge")
		}
	}
}

func (g *graphImpl) checkAndInit(ne identify) {
	v, node := g.checkAndGet(ne)
	if !v.IsNil() {
		return
	}
	if node {
		g.idSeq++
		v.Set(reflect.ValueOf(n{k: g.idSeq}))
	} else {
		g.idSeq++
		v.Set(reflect.ValueOf(e{k: g.idSeq}))
	}
}

func (g *graphImpl) String() string {
	sb := new(strings.Builder)
	sb.WriteString("NODE\tEDGE\tNODE\n")
	for _, v := range g.edges {
		sb.WriteString(fmt.Sprintf("%s\t%s\t%s\n", fmtSpec(v.u), fmtSpec(v.self), fmtSpec(v.v)))
	}
	return sb.String()
}

func (g *graphImpl) invalid(n identify) bool {
	v, _ := g.checkAndGet(n)
	return v.IsNil()
}
func fmtSpec(n interface{}) string {
	is := fmt.Sprintf("%+v", n)
	is = strings.TrimPrefix(is, "&")
	ib := strings.Index(is, "INode:")
	if ib != -1 {
		ie := ib + 12
		is = is[0:ib] + is[ie:]
	}
	ib = strings.Index(is, "IEdge:")
	if ib != -1 {
		ie := ib + 12
		is = is[0:ib] + is[ie:]
	}
	return is
}

type n struct {
	k int
}

func (ni n) nodeId() int {
	return ni.k
}

type e struct {
	k int
}

func (ni e) edgeId() int {
	return ni.k
}
