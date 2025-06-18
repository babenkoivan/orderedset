package orderedset_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/babenkoivan/orderedset"
)

func TestAdd(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(1)
	s.Add(2)
	s.Add(2)

	expected := []int{1, 2}
	if !reflect.DeepEqual(s.Values(), expected) {
		t.Errorf("Add failed: got %v, want %v", s.Values(), expected)
	}
}

func TestHas(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(1)

	if !s.Has(1) {
		t.Error("Has failed: expected true for 1")
	}
	if s.Has(2) {
		t.Error("Has failed: expected false for 2")
	}
}

func TestRemove(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(1)
	s.Add(2)

	s.Remove(1)
	if s.Has(1) {
		t.Error("Remove failed: 1 should be removed")
	}

	s.Remove(3)
	if s.Len() != 1 {
		t.Error("Remove failed: length should remain 1 after removing non-existing")
	}
}

func TestRemoveAt(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(10)
	s.Add(20)

	val, ok := s.RemoveAt(1)
	if !ok || val != 20 {
		t.Errorf("RemoveAt failed: got (%v, %v), want (20, true)", val, ok)
	}

	val, ok = s.RemoveAt(5)
	if ok {
		t.Error("RemoveAt failed: expected false for invalid index")
	}
}

func TestAt(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(5)

	val, ok := s.At(0)
	if !ok || val != 5 {
		t.Errorf("At failed: got (%v, %v), want (5, true)", val, ok)
	}

	_, ok = s.At(1)
	if ok {
		t.Error("At failed: expected false for out of range index")
	}
}

func TestIndexOf(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(5)
	s.Add(10)

	if idx := s.IndexOf(10); idx != 1 {
		t.Errorf("IndexOf failed: got %d, want 1", idx)
	}

	if idx := s.IndexOf(15); idx != -1 {
		t.Errorf("IndexOf failed: got %d, want -1", idx)
	}
}

func TestLen(t *testing.T) {
	s := orderedset.New[int]()
	if l := s.Len(); l != 0 {
		t.Errorf("Len failed: got %d, want 0", l)
	}

	s.Add(1)
	if l := s.Len(); l != 1 {
		t.Errorf("Len failed: got %d, want 1", l)
	}
}

func TestValues(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(1)
	s.Add(2)
	expected := []int{1, 2}

	if vals := s.Values(); !reflect.DeepEqual(vals, expected) {
		t.Errorf("Values failed: got %v, want %v", vals, expected)
	}
}

func TestClone(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(1)
	s.Add(2)

	clone := s.Clone()
	if !reflect.DeepEqual(s.Values(), clone.Values()) {
		t.Error("Clone failed: cloned values differ")
	}

	clone.Add(3)
	if s.Has(3) {
		t.Error("Clone failed: original modified after clone changed")
	}
}

func TestUnion(t *testing.T) {
	s1 := orderedset.New[int]()
	s2 := orderedset.New[int]()

	s1.Add(1)
	s1.Add(2)

	s2.Add(2)
	s2.Add(3)

	union := s1.Union(s2)
	expected := []int{1, 2, 3}

	if !reflect.DeepEqual(union.Values(), expected) {
		t.Errorf("Union failed: got %v, want %v", union.Values(), expected)
	}
}

func TestIntersect(t *testing.T) {
	s1 := orderedset.New[int]()
	s2 := orderedset.New[int]()

	s1.Add(1)
	s1.Add(2)

	s2.Add(2)
	s2.Add(3)

	intersect := s1.Intersect(s2)
	expected := []int{2}

	if !reflect.DeepEqual(intersect.Values(), expected) {
		t.Errorf("Intersect failed: got %v, want %v", intersect.Values(), expected)
	}
}

func TestDifference(t *testing.T) {
	s1 := orderedset.New[int]()
	s2 := orderedset.New[int]()

	s1.Add(1)
	s1.Add(2)

	s2.Add(2)
	s2.Add(3)

	diff := s1.Difference(s2)
	expected := []int{1}

	if !reflect.DeepEqual(diff.Values(), expected) {
		t.Errorf("Difference failed: got %v, want %v", diff.Values(), expected)
	}
}

func TestSlice(t *testing.T) {
	s := orderedset.New[int]()
	for i := 1; i <= 5; i++ {
		s.Add(i)
	}

	for n, tc := range map[string]struct {
		from    int
		to      int
		want    []int
		wantErr bool
	}{
		"valid slice 0-3":   {from: 0, to: 3, want: []int{1, 2, 3}},
		"valid slice 1-5":   {from: 1, to: 5, want: []int{2, 3, 4, 5}},
		"valid empty 3-3":   {from: 3, to: 3, want: []int{}},
		"invalid from < 0":  {from: -1, to: 3, wantErr: true},
		"invalid to > len":  {from: 2, to: 6, wantErr: true},
		"invalid from > to": {from: 4, to: 2, wantErr: true},
	} {
		t.Run(n, func(t *testing.T) {
			got, err := s.Slice(tc.from, tc.to)
			if (err != nil) != tc.wantErr {
				t.Errorf("Slice(%d, %d) error = %v, wantErr %v", tc.from, tc.to, err, tc.wantErr)
				return
			}
			if err == nil && !reflect.DeepEqual(got.Values(), tc.want) {
				t.Errorf("Slice(%d, %d) = %v, want %v", tc.from, tc.to, got.Values(), tc.want)
			}
		})
	}
}

func TestSortBy(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(3)
	s.Add(1)
	s.Add(2)

	s.SortBy(func(a, b int) bool { return a < b })
	expectedAsc := []int{1, 2, 3}
	if !reflect.DeepEqual(s.Values(), expectedAsc) {
		t.Errorf("SortBy ascending failed: got %v, want %v", s.Values(), expectedAsc)
	}

	s.SortBy(func(a, b int) bool { return a > b })
	expectedDesc := []int{3, 2, 1}
	if !reflect.DeepEqual(s.Values(), expectedDesc) {
		t.Errorf("SortBy descending failed: got %v, want %v", s.Values(), expectedDesc)
	}
}

func TestMarshalJSON(t *testing.T) {
	s := orderedset.New[int]()
	s.Add(1)
	s.Add(2)

	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	wantJSON := `[1,2]`
	if string(b) != wantJSON {
		t.Errorf("Marshal JSON failed: got %s, want %s", string(b), wantJSON)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	input := `[1,2]`
	s := orderedset.New[int]()
	err := json.Unmarshal([]byte(input), s)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	expected := []int{1, 2}
	if !reflect.DeepEqual(s.Values(), expected) {
		t.Errorf("Unmarshal JSON failed: got %v, want %v", s.Values(), expected)
	}
}
