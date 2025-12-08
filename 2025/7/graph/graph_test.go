package graph

import (
	"slices"
	"testing"
)

func TestNewGraph(t *testing.T) {
	g := NewGraph()
	if g == nil {
		t.Fatal("NewGraph returned nil")
	}
	if g.Nodes == nil {
		t.Fatal("NewGraph.Nodes is nil; expected empty map")
	}
	if len(g.Nodes) != 0 {
		t.Fatalf("expected 0 nodes, got %d", len(g.Nodes))
	}
}

func TestAddNode(t *testing.T) {
	g := NewGraph()
	g.AddNode("a")
	if _, ok := g.Nodes["a"]; !ok {
		t.Fatal("node 'a' not created")
	}
	if g.Nodes["a"] == nil {
		t.Fatal("node 'a' slice is nil; expected empty slice")
	}
	if len(g.Nodes["a"]) != 0 {
		t.Fatalf("expected node 'a' to have 0 edges, got %d", len(g.Nodes["a"]))
	}

	// idempotent
	g.AddNode("a")
	if len(g.Nodes) != 1 {
		t.Fatalf("expected 1 node after adding 'a' twice, got %d", len(g.Nodes))
	}
}

func TestAddEdgeCreatesNodesAndPreventsDuplicates(t *testing.T) {
	g := NewGraph()

	// Add edge when neither node exists
	g.AddEdge("a", "b")
	if _, ok := g.Nodes["a"]; !ok {
		t.Fatal("expected node 'a' to be created by AddEdge")
	}
	if _, ok := g.Nodes["b"]; !ok {
		t.Fatal("expected node 'b' to be created by AddEdge")
	}
	if !slices.Contains(g.Nodes["a"], "b") {
		t.Fatalf("expected 'a' -> 'b' edge, got %v", g.Nodes["a"])
	}
	if len(g.Nodes["b"]) != 0 {
		t.Fatalf("expected 'b' to have 0 outgoing edges, got %d", len(g.Nodes["b"]))
	}

	// Add duplicate edge and ensure it's not added twice
	g.AddEdge("a", "b")
	count := 0
	for _, v := range g.Nodes["a"] {
		if v == "b" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected single 'b' in adjacency of 'a', found %d", count)
	}

	// Add more edges and verify ordering/contents
	g.AddEdge("a", "c")
	g.AddEdge("a", "d")
	if !slices.Contains(g.Nodes["a"], "c") || !slices.Contains(g.Nodes["a"], "d") {
		t.Fatalf("expected 'a' to contain 'c' and 'd', got %v", g.Nodes["a"])
	}
}

func TestAddEdgeHandlesNilSlice(t *testing.T) {
	g := NewGraph()
	// create node with nil slice explicitly
	g.Nodes["x"] = nil
	g.AddEdge("x", "y")
	if !slices.Contains(g.Nodes["x"], "y") {
		t.Fatalf("expected 'x' -> 'y' after AddEdge on nil slice, got %v", g.Nodes["x"])
	}
	// adding again should not duplicate
	g.AddEdge("x", "y")
	cnt := 0
	for _, v := range g.Nodes["x"] {
		if v == "y" {
			cnt++
		}
	}
	if cnt != 1 {
		t.Fatalf("expected single 'y' in adjacency of 'x', found %d", cnt)
	}
}

func TestCountPathsFrom_NonExistent(t *testing.T) {
	g := NewGraph()
	if got := g.CountPathsFrom("z"); got != 0 {
		t.Fatalf("expected 0 for non-existent node, got %d", got)
	}
}

func TestCountPathsFrom_LeafAndNilSlice(t *testing.T) {
	g := NewGraph()
	g.AddNode("b")
	if got := g.CountPathsFrom("b"); got != 1 {
		t.Fatalf("expected 1 for leaf node, got %d", got)
	}
}

func TestCountPathsFrom_Chain(t *testing.T) {
	g := NewGraph()
	// 0 -> 1 -> 2 (single path from 0 to leaf)
	g.AddEdge("a", "b")
	g.AddEdge("b", "c")

	if got := g.CountPathsFrom("a"); got != 1 {
		t.Fatalf("expected 1 path for chain starting at 'a', got %d", got)
	}
	if got := g.CountPathsFrom("b"); got != 1 {
		t.Fatalf("expected 1 path for chain starting at 'b', got %d", got)
	}
	if got := g.CountPathsFrom("c"); got != 1 {
		t.Fatalf("expected 1 path for leaf node 'c', got %d", got)
	}
}

func TestCountPathsFrom_BranchingAndMerging(t *testing.T) {
	g := NewGraph()
	// Graph:
	// 0 -> 1 -> 3
	// 0 -> 2 -> 3
	// Expect two distinct paths from 0 to leaf 3
	g.AddEdge("a", "b")
	g.AddEdge("a", "c")
	g.AddEdge("b", "d")
	g.AddEdge("c", "d")

	if got := g.CountPathsFrom("a"); got != 2 {
		t.Fatalf("expected 2 distinct paths from 'a', got %d", got)
	}
	if got := g.CountPathsFrom("b"); got != 1 {
		t.Fatalf("expected 1 path from 'b', got %d", got)
	}
	if got := g.CountPathsFrom("d"); got != 1 {
		t.Fatalf("expected 1 path from leaf 'd', got %d", got)
	}
}
