package dsu

import "testing"

// AI generated tests with some tweaks
// If things go sideways, I'd be skeptical of these tests correctness
// This wasn't the focus of the exercise so I didn't spend too long here

func assertEqual(t *testing.T, name string, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: got %d, want %d", name, got, want)
	}
}

func assertTrue(t *testing.T, name string, cond bool) {
	t.Helper()
	if !cond {
		t.Fatalf("%s: expected true, got false", name)
	}
}

func assertFalse(t *testing.T, name string, cond bool) {
	t.Helper()
	if cond {
		t.Fatalf("%s: expected false, got true", name)
	}
}

func TestNewDSU(t *testing.T) {
	const n = 5
	d := NewDSU(n)

	assertEqual(t, "NewDSU count", d.Count(), n)

	for i := range n {
		assertEqual(t, "NewDSU size of singleton", d.Size(i), 1)

		for j := 0; j < n; j++ {
			if i == j {
				assertTrue(t, "NewDSU Same(i,i)", d.Same(i, j))
			} else {
				assertFalse(t, "NewDSU Same(i,j) for distinct i,j", d.Same(i, j))
			}
		}
	}
}
func TestUnionBasicGroups(t *testing.T) {
	d := NewDSU(6)

	// Initially: 6 singleton sets
	assertEqual(t, "start count", d.Count(), 6)

	// Step 1: union(0,1) -> {0,1}, {2}, {3}, {4}, {5}
	d.Union(0, 1)

	assertTrue(t, "after union(0,1) Same(0,1)", d.Same(0, 1))
	assertEqual(t, "after union(0,1) size(0)", d.Size(0), 2)
	assertEqual(t, "after union(0,1) count", d.Count(), 5)
	assertFalse(t, "after union(0,1) Same(0,2)", d.Same(0, 2))
	assertFalse(t, "after union(0,1) Same(1,2)", d.Same(1, 2))

	// Step 2: union(2,3) and union(4,5) -> {0,1}, {2,3}, {4,5}
	d.Union(2, 3)
	d.Union(4, 5)

	assertTrue(t, "after unions(2,3) Same(2,3)", d.Same(2, 3))
	assertTrue(t, "after unions(4,5) Same(4,5)", d.Same(4, 5))
	assertEqual(t, "after three unions count", d.Count(), 3)

	// Cross-group checks
	assertFalse(t, "groups distinct Same(0,2)", d.Same(0, 2))
	assertFalse(t, "groups distinct Same(0,4)", d.Same(0, 4))
	assertFalse(t, "groups distinct Same(2,4)", d.Same(2, 4))

	// Step 3: union(1,3) -> merge {0,1} and {2,3} into {0,1,2,3}
	d.Union(1, 3)

	assertTrue(t, "after union(1,3) Same(0,3)", d.Same(0, 3))
	assertTrue(t, "after union(1,3) Same(1,2)", d.Same(1, 2))
	assertEqual(t, "after union(1,3) size(0)", d.Size(0), 4)
	assertEqual(t, "after union(1,3) count", d.Count(), 2)
	assertFalse(t, "after union(1,3) groups {0..3} and {4,5} distinct", d.Same(0, 4))

	// Step 4: union(0,4) -> merge everything into {0,1,2,3,4,5}
	d.Union(0, 4)

	assertEqual(t, "after union(0,4) count", d.Count(), 1)
	for i := range 6 {
		assertTrue(t, "after union(0,4) Same(0,i)", d.Same(0, i))
	}
	assertEqual(t, "after union(0,4) final size(0)", d.Size(0), 6)
}

func TestUnionIdempotentAndCommutative(t *testing.T) {
	d := NewDSU(4)

	// Idempotent: repeated union doesn't change count or groups
	d.Union(0, 1)
	before := d.Count()
	d.Union(0, 1)
	assertEqual(t, "Union idempotent count", d.Count(), before)
	assertTrue(t, "Union idempotent Same(0,1)", d.Same(0, 1))

	// Commutativity: Union(a,b) vs Union(b,a) leads to same partition
	d2 := NewDSU(4)
	d2.Union(1, 0) // flipped order

	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if d.Same(i, j) != d2.Same(i, j) {
				t.Fatalf("Union commutative: partition mismatch for (%d,%d)", i, j)
			}
		}
	}
}

func TestPathCompression(t *testing.T) {
	const n = 10
	d := NewDSU(n)

	// Build a chain manually: 0 <- 1 <- 2 <- ... <- n-1
	for i := range n {
		d.parent[i] = i
		d.size[i] = 1
	}
	for i := 1; i < n; i++ {
		d.parent[i] = i - 1
	}

	// Sanity check: last node should not directly point to root yet.
	if d.parent[n-1] == 0 {
		t.Fatalf("setup: expected deep chain, but parent[n-1] already 0")
	}

	root := d.Find(n - 1)
	if root != 0 {
		t.Fatalf("Find: expected root 0, got %d", root)
	}

	// After path compression, all nodes should point directly to root.
	for i := range n {
		if d.parent[i] != 0 {
			t.Fatalf("PathCompression: expected parent[%d] == 0 after Find, got %d", i, d.parent[i])
		}
	}

	// And logically, Same() should report all elements in one set.
	for i := range n {
		if !d.Same(0, i) {
			t.Fatalf("PathCompression: expected Same(0,%d) == true", i)
		}
	}
}

func TestSetSizes_Distribution(t *testing.T) {
	d := NewDSU(5)
	d.Union(0, 1)
	d.Union(2, 3)

	sizes := d.SetSizes()

	if got := len(sizes); got != 3 {
		t.Fatalf("expected 3 distinct sets, got %d", got)
	}

	// Build a distribution of sizes (size -> count)
	dist := make(map[int]int)
	for _, sz := range sizes {
		dist[sz]++
	}

	if dist[2] != 2 || dist[1] != 1 {
		t.Fatalf("unexpected size distribution: %v", dist)
	}
}

func TestSetSizes_FullUnion(t *testing.T) {
	n := 6
	d := NewDSU(n)
	for i := 1; i < n; i++ {
		d.Union(0, i)
	}

	sizes := d.SetSizes()
	if got := len(sizes); got != 1 {
		t.Fatalf("expected 1 set after full union, got %d", got)
	}

	for _, sz := range sizes {
		if sz != n {
			t.Fatalf("expected size %d, got %d", n, sz)
		}
	}
}
