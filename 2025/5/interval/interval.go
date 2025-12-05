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

// Searches for all intervals in the tree that overlap with the given value.
func (t *IntervalTree) Search(value int) []*Node {
	var result []*Node
	_searchIntervals(t.Root, value, &result)
	return result
}

func _searchIntervals(n *Node, value int, result *[]*Node) {
	if n == nil {
		return
	}
	// If the current interval overlaps with the value
	if n.Start <= value && value <= n.End {
		*result = append(*result, n)
	}
	// If the left child exists and its max is greater than or equal to the value, search left
	if n.Left != nil && n.Left.Max >= value {
		_searchIntervals(n.Left, value, result)
	}
	// Always search right
	_searchIntervals(n.Right, value, result)
}
