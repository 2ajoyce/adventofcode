package interval

import (
	"fmt"
	"testing"
)

func TestInsertRoot(t *testing.T) {
	tree := NewIntervalTree()
	tree.Insert(5, 10)
	if tree.Root == nil {
		t.Fatal("expected root to be set")
	}
	if tree.Root.Start != 5 || tree.Root.End != 10 {
		t.Fatalf("unexpected root interval: got start=%d end=%d", tree.Root.Start, tree.Root.End)
	}
	if tree.Root.Left != nil || tree.Root.Right != nil {
		t.Fatalf("expected no children for single-node tree; got left=%v right=%v", tree.Root.Left, tree.Root.Right)
	}
}

func TestInsertAndSearch(t *testing.T) {
	tree := NewIntervalTree()
	tree.Insert(5, 10)
	tree.Insert(1, 2)
	tree.Insert(20, 30)

	if tree.Root == nil {
		t.Fatal("root nil")
	}
	if tree.Root.Max != 30 {
		t.Fatalf("expected root.Max 30, got %d", tree.Root.Max)
	}

	res := tree.Search(6)
	if len(res) != 1 {
		t.Fatalf("expected 1 result for value 6, got %d", len(res))
	}
	if res[0].Start != 5 || res[0].End != 10 {
		t.Fatalf("unexpected interval for value 6: %+v", res[0])
	}

	res2 := tree.Search(15)
	if len(res2) != 0 {
		t.Fatalf("expected 0 results for value 15, got %d", len(res2))
	}

	res3 := tree.Search(25)
	if len(res3) != 1 || res3[0].Start != 20 || res3[0].End != 30 {
		t.Fatalf("expected [20,30] for value 25, got %+v", res3)
	}
}

func TestSearch(t *testing.T) {
	cases := []struct {
		name    string
		inserts [][2]int
		value   int
		want    []string // intervals in the form "start:end"
	}{
		{
			name:    "overlapping_and_point",
			inserts: [][2]int{{1, 5}, {3, 7}, {5, 5}, {6, 10}},
			value:   5,
			want:    []string{"1:5", "3:7", "5:5"},
		},
		{
			name:    "no_hits",
			inserts: [][2]int{{1, 2}, {4, 4}, {10, 12}},
			value:   3,
			want:    []string{},
		},
		{
			name:    "same_start_different_end",
			inserts: [][2]int{{5, 6}, {5, 10}, {2, 4}},
			value:   6,
			want:    []string{"5:6", "5:10"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tree := NewIntervalTree()
			for _, it := range tc.inserts {
				tree.Insert(it[0], it[1])
			}

			res := tree.Search(tc.value)
			got := make(map[string]struct{}, len(res))
			for _, n := range res {
				got[fmt.Sprintf("%d:%d", n.Start, n.End)] = struct{}{}
			}

			if len(got) != len(tc.want) {
				t.Fatalf("%s: expected %d results, got %d (got=%v)", tc.name, len(tc.want), len(got), res)
			}
			for _, w := range tc.want {
				if _, ok := got[w]; !ok {
					t.Fatalf("%s: missing expected interval %s (got=%v)", tc.name, w, res)
				}
			}
		})
	}
}
