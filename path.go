package vcontext

import (
	"errors"
)

var (
	ErrBadPath = errors.New("invalid path")
)

type Node interface {
	Start() int64
	End()   int64
	Get(path ...interface{}) (Node, error)
}

type Key string

type Marker struct {
	StartIdx int64
	EndIdx   int64
}

func (m Marker) Start() int64 {
	return m.StartIdx
}

func (m Marker) End() int64 {
	return m.EndIdx
}

type MapNode struct {
	Marker
	Children map[string]Node
	Keys     map[string]Leaf
}

type Leaf struct {
	Marker
}

func (k Leaf) Get(path ...interface{}) (Node, error) {
	if len(path) == 0 {
		return k, nil
	}
	return nil, ErrBadPath
}

func (m MapNode) Get(path ...interface{}) (Node, error) {
	if len(path) == 0 {
		return m, nil
	}
	switch p := path[0].(type) {
	case string:
		if r, ok := m.Children[p]; ok {
			return r.Get(path[1:]...)
		} else {
			return nil, ErrBadPath
		}
	case Key:
		if r, ok := m.Keys[string(p)]; ok {
			return r.Get(path[1:]...)
		} else {
			return nil, ErrBadPath
		}
	default:
		return nil, ErrBadPath
	}
}
		
type SliceNode struct {
	Marker
	Children []Node
}

func (s SliceNode) Get(path ...interface{}) (Node, error) {
	if len(path) == 0 {
		return s, nil
	}
	if i, ok := path[0].(int); ok {
		if i >= len(s.Children) {
			return nil, ErrBadPath
		}
		return s.Children[i].Get(path[1:]...)
	}
	return nil, ErrBadPath
}
