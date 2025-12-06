package interval

// An interval search tree implementation in Go.

type Node struct {
	Start int   // The start of the interval
	End   int   // The end of the interval
	Max   int   // The maximum end value in the subtree
	Left  *Node // Left child
	Right *Node // Right child
}

type IntervalTree struct {
	Root *Node
}

func NewIntervalTree() *IntervalTree {
	return &IntervalTree{}
}

// A recursive function to insert a new interval into the tree.
func (t *IntervalTree) Insert(start, end int) {
	if t.Root == nil {
		t.Root = &Node{Start: start, End: end, Max: end}
		return
	}
	_insert(t.Root, start, end)
}

// A recursive function to insert a new interval into the tree.
func _insert(n *Node, start, end int) {
	newNode := &Node{Start: start, End: end, Max: end}
	if start < n.Start {
		if n.Left == nil {
			n.Left = newNode
		} else {
			_insert(n.Left, start, end)
		}
	} else {
		if n.Right == nil {
			n.Right = newNode
		} else {
			_insert(n.Right, start, end)
		}
	}

	// Update the max value of this ancestor node
	if n.Max < end {
		n.Max = end
	}
}

// A recursive function to insert a new interval into the tree.
// If the new interval overlaps with existing intervals, they are merged.
func (t *IntervalTree) InsertWithoutOverlap(start, end int) {
	if t.Root == nil {
		t.Root = &Node{Start: start, End: end, Max: end}
		return
	}
	t._insertWithoutOverlap(t.Root, start, end)
}

// A recursive function to insert a new interval into the tree.
// If the new interval overlaps with existing intervals, they are merged.
func (t *IntervalTree) _insertWithoutOverlap(n *Node, start, end int) {
	newNode := &Node{Start: start, End: end, Max: end}
	if start < n.Start && end < n.Start {
		if n.Left == nil {
			n.Left = newNode
		} else {
			t._insertWithoutOverlap(n.Left, start, end)
		}
	} else if start > n.End && end > n.End {
		if n.Right == nil {
			n.Right = newNode
		} else {
			t._insertWithoutOverlap(n.Right, start, end)
		}
	} else {
		// If it overlaps with the start we will add a new node
		// that covers the amount smaller than the current start
		if start < n.Start {
			t.InsertWithoutOverlap(start, n.Start-1)
		}
		// If it overlaps with the end we will add a new node
		// that covers the amount larger than the current end
		if end > n.End {
			t.InsertWithoutOverlap(n.End+1, end)
			return // Return here to avoid updating max twice
		}
	}

	// Update the max value of this ancestor node
	if n.Max < end {
		n.Max = end
	}
}

// Searches for all intervals in the tree that overlap with the given value.
func (t *IntervalTree) Search(value int) []*Node {
	var result []*Node
	t._searchIntervals(t.Root, value, &result)
	return result
}

func (t *IntervalTree) _searchIntervals(n *Node, value int, result *[]*Node) {
	if n == nil {
		return
	}
	// If the current interval overlaps with the value
	if n.Start <= value && value <= n.End {
		*result = append(*result, n)
	}
	// If the left child exists and its max is greater than or equal to the value, search left
	if n.Left != nil && n.Left.Max >= value {
		t._searchIntervals(n.Left, value, result)
	}
	// Always search right
	t._searchIntervals(n.Right, value, result)
}
