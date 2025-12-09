package dsu

type DSU struct {
	parent []int
	rank   []int
	size   []int
	count  int
}

func NewDSU(n int) *DSU {
	d := &DSU{
		parent: make([]int, n),
		rank:   make([]int, n),
		size:   make([]int, n),
		count:  n,
	}
	for i := range n {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	// Return the representative (root) of the set containing x
	// with path compression

	p := d.parent[x]
	// This element has itself as parent we have reached the root
	// we have reached the root
	if p == x {
		return p
	}
	// Otherwise continue climbing the tree
	root := d.Find(p)

	// Path compression
	d.parent[x] = root

	return root
}

func (d *DSU) Union(a, b int) {
	// Merge the sets containing a and b
	rootA := d.Find(a)
	rootB := d.Find(b)

	// They are already in the same set
	if rootA == rootB {
		return
	}
	// Decrease the count of disjoint sets
	if d.rank[rootA] < d.rank[rootB] {
		d.parent[rootA] = rootB
		d.size[rootB] += d.size[rootA]
	} else if d.rank[rootA] > d.rank[rootB] {
		d.parent[rootB] = rootA
		d.size[rootA] += d.size[rootB]
	} else {
		d.parent[rootB] = rootA
		d.size[rootA] += d.size[rootB]
		d.rank[rootA]++
	}
	d.count--
}

func (d *DSU) Same(a, b int) bool {
	// Return true if a and b are in the same set
	if d.Find(a) == d.Find(b) {
		return true
	}
	return false
}

func (d *DSU) Size(x int) int {
	// Return the size of the set containing x
	return d.size[d.Find(x)]
}

func (d *DSU) Count() int {
	// Return the current number of disjoint sets
	return d.count
}

func (d *DSU) SetSizes() map[int]int {
	// Return a map of root -> size for all distinct sets
	sizes := make(map[int]int)
	for i := range d.parent {
		root := d.Find(i)
		sizes[root] = d.size[root]
	}
	return sizes
}
