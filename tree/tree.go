package tree

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrBadPath = errors.New("invalid path")
)

// superset of Marker
type Node interface {
	Start() (int64, int64) // line, col
	End()   (int64, int64)
	Get(path []interface{}) (Node, error)
	GetMarker() Marker
	pos() []*Pos // just used for iterating through the markers to fill in line and column from index
}

func FixLineColumn(n Node, source []byte) {
	p := n.pos()
	sort.Slice(p, func(i, j int) bool {
		return p[i].Index < p[j].Index
	})
	pi := 0
	line, col := int64(1), int64(1)
	for i, c := range source {
		if pi == len(p) {
			return
		}
		if int64(i) == p[pi].Index {
			p[pi].Line = line
			p[pi].Column = col
			pi++
		}
		col++
		if c == '\n' {
			line++
			col = 1
		}
	}
}

// Key is used to differentiate leaves describing the start and end of where
// a key starts and where a value starts.
type Key string

type Pos struct {
	Index  int64
	Line   int64
	Column int64
}

func posString(p *Pos) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("line %d col %d", p.Line, p.Column)
}

func posLC(p *Pos) (int64, int64) {
	if p == nil {
		return 0, 0
	}
	return p.Line, p.Column
}

// IndexMarkers are composed of information regarding the start and
// end of where a Node exists in its source. Depending on where the marker
// came from (e.g. json or yaml) the *Idx or End* fields may not be populated.
type Marker struct {
	StartP *Pos
	EndP   *Pos
}

func (m Marker) Start() (int64, int64) {
	return posLC(m.StartP)
}

func (m Marker) End() (int64, int64) {
	return posLC(m.EndP)
}

func (m Marker) String() string {
	// Just do start for now, figure out end later
	return posString(m.StartP)
}

func (marker Marker) GetMarker() Marker {
	return marker
}

func MarkerFromIndices(start, end int64) Marker {
	m := Marker{}
	if start >= 0 {
		m.StartP = &Pos{Index: start}
	}
	if end >= 0 {
		m.EndP = &Pos{Index: end}
	}
	return m
}

func appendPos(l []*Pos, p *Pos) []*Pos {
	if p != nil {
		return append(l, p)
	}
	return l
}

type MapNode struct {
	Marker
	Children map[string]Node
	Keys     map[string]Leaf
}

func (m MapNode) Get(path []interface{}) (Node, error) {
	if len(path) == 0 {
		return m, nil
	}
	switch p := path[0].(type) {
	case string:
		if r, ok := m.Children[p]; ok {
			return r.Get(path[1:])
		} else {
			return nil, ErrBadPath
		}
	case Key:
		if r, ok := m.Keys[string(p)]; ok {
			return r.Get(path[1:])
		} else {
			return nil, ErrBadPath
		}
	default:
		return nil, ErrBadPath
	}
}

func (m MapNode) pos() []*Pos {
	ret := appendPos(nil, m.StartP)
	for _, v := range m.Children {
		ret = append(ret, v.pos()...)
	}
	for _, v := range m.Keys {
		ret = append(ret, v.pos()...)
	}
	ret = appendPos(ret, m.EndP)
	return ret
}

type Leaf struct {
	Marker
}

func (l Leaf) pos() []*Pos {
	return appendPos(appendPos(nil, l.StartP), l.EndP)
}

func (k Leaf) Get(path []interface{}) (Node, error) {
	if len(path) == 0 {
		return k, nil
	}
	return nil, ErrBadPath
}

type SliceNode struct {
	Marker
	Children []Node
}

func (s SliceNode) Get(path []interface{}) (Node, error) {
	if len(path) == 0 {
		return s, nil
	}
	if i, ok := path[0].(int); ok {
		if i >= len(s.Children) {
			return nil, ErrBadPath
		}
		return s.Children[i].Get(path[1:])
	}
	return nil, ErrBadPath
}

func (s SliceNode) pos() []*Pos {
	ret := appendPos(nil, s.StartP)
	for _, v := range s.Children {
		ret = append(ret, v.pos()...)
	}
	ret = appendPos(ret, s.EndP)
	return ret
}
