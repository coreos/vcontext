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

package tree

import (
	"reflect"
	"testing"
)

func TestFixLineColumn(t *testing.T) {
	tests := []struct {
		in  []int64
		src string
		out [][]int64 //list of index, line, col
	}{
		{
			in:  []int64{0, 1, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			src: "01\n3\n567\n9",
			out: [][]int64{
				{0, 1, 1},
				{1, 1, 2},
				{1, 1, 2},
				{2, 1, 3},
				{3, 2, 1},
				{4, 2, 2},
				{5, 3, 1},
				{6, 3, 2},
				{7, 3, 3},
				{8, 3, 4},
				{9, 4, 1},
			},
		},
	}

	for i, test := range tests {
		p := make([]*Pos, len(test.in), len(test.in))
		expected := make(map[int64]Pos, len(test.in))
		for j, index := range test.in {
			p[j] = &Pos{Index: index}
			expected[test.out[j][0]] = Pos{
				Index:  test.out[j][0],
				Line:   test.out[j][1],
				Column: test.out[j][2],
			}
		}
		fixLineColumn(p, []byte(test.src))
		for _, pos := range p {
			exp := expected[pos.Index]
			if !reflect.DeepEqual(*pos, exp) {
				t.Errorf("#%d: expected %+v, got %+v:", i, exp, *pos)
			}
		}
	}
}
