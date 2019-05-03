package json

import (
	"github.com/ajeddeloh/vcontext"
	// todo: rewrite this dep
	json "github.com/ajeddeloh/go-json"
)

func UnmarshalToContext(raw []byte) (vcontext.Node, error) {
	var ast json.Node
	if err := json.Unmarshal(raw, &ast); err != nil {
		return nil, err
	}
	return fromJsonNode(ast), nil
}

func fromJsonNode(n json.Node) vcontext.Node {
	m := vcontext.Marker{
		StartIdx: int64(n.Start),
		EndIdx:   int64(n.End),
	}
	switch v := n.Value.(type) {
	case map[string]json.Node:
		ret := vcontext.MapNode{
			Marker: m,
			Children: make(map[string]vcontext.Node, len(v)),
			Keys:     make(map[string]vcontext.Leaf, len(v)),
		}
		for key, child := range v {
			ret.Children[key] = fromJsonNode(child)
			ret.Keys[key] = vcontext.Leaf{
				Marker: vcontext.Marker{
					StartIdx: int64(child.KeyStart),
					EndIdx:   int64(child.KeyEnd),
				},
			}
		}
		return ret
	case []json.Node:
		ret := vcontext.SliceNode{
			Marker: m,
			Children: make([]vcontext.Node, 0, len(v)),
		}
		for _, child := range v {
			ret.Children = append(ret.Children, fromJsonNode(child))
		}
		return ret
	default:
		return vcontext.Leaf{
			Marker: m,
		}
	}
}

