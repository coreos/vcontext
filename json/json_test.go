// Copyright 2019 Red Hat, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.)

package json

import (
	"reflect"
	"testing"

	"github.com/ajeddeloh/vcontext/tree"

	json "github.com/ajeddeloh/go-json"
)

func TestUnmarshalToContext(t *testing.T) {
	tests := []struct {
		in  json.Node
		out tree.Node
	}{
		// leaf
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: "foo",
			},
			tree.Leaf{
				Marker: tree.MarkerFromIndices(1, 2),
			},
		},
		// map
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: map[string]json.Node{},
			},
			tree.MapNode{
				Marker:   tree.MarkerFromIndices(1, 2),
				Keys:     map[string]tree.Leaf{},
				Children: map[string]tree.Node{},
			},
		},
		//slice
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: []json.Node{},
			},
			tree.SliceNode{
				Marker:   tree.MarkerFromIndices(1, 2),
				Children: []tree.Node{},
			},
		},
		// map of map, slice, leaf
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: map[string]json.Node{
					"foo": json.Node{
						KeyStart: 3,
						KeyEnd:   4,
						Start:    5,
						End:      6,
						Value:    map[string]json.Node{},
					},
					"bar": json.Node{
						KeyStart: 7,
						KeyEnd:   8,
						Start:    9,
						End:      10,
						Value:    []json.Node{},
					},
					"baz": json.Node{
						KeyStart: 11,
						KeyEnd:   12,
						Start:    13,
						End:      14,
						Value:    "quux",
					},
				},
			},
			tree.MapNode{
				Marker: tree.MarkerFromIndices(1, 2),
				Keys: map[string]tree.Leaf{
					"foo": tree.Leaf{
						Marker: tree.MarkerFromIndices(3, 4),
					},
					"bar": tree.Leaf{
						Marker: tree.MarkerFromIndices(7, 8),
					},
					"baz": tree.Leaf{
						Marker: tree.MarkerFromIndices(11, 12),
					},
				},
				Children: map[string]tree.Node{
					"foo": tree.MapNode{
						Marker:   tree.MarkerFromIndices(5, 6),
						Children: map[string]tree.Node{},
						Keys:     map[string]tree.Leaf{},
					},
					"bar": tree.SliceNode{
						Marker:   tree.MarkerFromIndices(9, 10),
						Children: []tree.Node{},
					},
					"baz": tree.Leaf{
						Marker: tree.MarkerFromIndices(13, 14),
					},
				},
			},
		},
		// slice of leaf
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: []json.Node{
					json.Node{
						Start: 3,
						End:   4,
						Value: "foo",
					},
				},
			},
			tree.SliceNode{
				Marker: tree.MarkerFromIndices(1, 2),
				Children: []tree.Node{
					tree.Leaf{
						Marker: tree.MarkerFromIndices(3, 4),
					},
				},
			},
		},
		// slice of slice
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: []json.Node{
					json.Node{
						Start: 3,
						End:   4,
						Value: []json.Node{},
					},
				},
			},
			tree.SliceNode{
				Marker: tree.MarkerFromIndices(1, 2),
				Children: []tree.Node{
					tree.SliceNode{
						Marker:   tree.MarkerFromIndices(3, 4),
						Children: []tree.Node{},
					},
				},
			},
		},
		// slice of map
		{
			json.Node{
				Start: 1,
				End:   2,
				Value: []json.Node{
					json.Node{
						Start: 3,
						End:   4,
						Value: map[string]json.Node{},
					},
				},
			},
			tree.SliceNode{
				Marker: tree.MarkerFromIndices(1, 2),
				Children: []tree.Node{
					tree.MapNode{
						Marker:   tree.MarkerFromIndices(3, 4),
						Children: map[string]tree.Node{},
						Keys:     map[string]tree.Leaf{},
					},
				},
			},
		},
	}
	for i, test := range tests {
		n := fromJsonNode(test.in)
		if !reflect.DeepEqual(test.out, n) {
			t.Errorf("test %d failed: expected: %v, got %v", i, test.out, n)
		}
	}
}
