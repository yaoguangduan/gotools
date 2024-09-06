package graph

import (
	"fmt"
	"testing"
)

type person struct {
	INode
	name string
}
type relation struct {
	IEdge
	weight int
}

func TestInnerRef(t *testing.T) {
	p1 := &person{nil, "john"}
	p2 := &person{nil, "kitty"}
	p3 := &person{nil, "maka"}
	p4 := &person{nil, "nana"}
	r1 := &relation{nil, 23}
	r2 := &relation{nil, 11}
	gi := Undirected(false, false)
	err := gi.AddEdge(p1, p2, r1)
	if err != nil {
		panic(err)
	}
	err = gi.AddEdge(p2, p3, r2)
	if err != nil {
		panic(err)
	}
	err = gi.AddEdge(p1, p4, &relation{nil, 78})
	if err != nil {
		panic(err)
	}
	printNodes(gi.Nodes())
	fmt.Printf("%+v\n", gi.Edges())
	fmt.Printf("%+v\n", gi.AdjacentNodes(p1))
	fmt.Printf("%+v\n", gi.AdjacentNodes(p2))
	fmt.Printf("%+v\n", gi.PredecessorNodes(p1))
	fmt.Printf("%+v\n", gi.SuccessorNodes(p1))
	fmt.Printf("%+v\n", gi.SuccessorNodes(p2))
	fmt.Println(gi)
}
func printNodes(nodes []INode) {
	for _, node := range nodes {
		fmt.Printf("%+v ", node.(*person).name)
	}
	fmt.Println()
}

type test struct {
	f1 string
	f2 int
	INode
}
