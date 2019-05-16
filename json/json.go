package json

import (
	"github.com/ajeddeloh/vcontext/tree"
	// todo: rewrite this dep
	json "github.com/ajeddeloh/go-json"
)

func UnmarshalToContext(raw []byte) (tree.Node, error) {
	var ast json.Node
	if err := json.Unmarshal(raw, &ast); err != nil {
		return nil, err
	}
	node := fromJsonNode(ast)
	tree.FixLineColumn(node, raw)
	return node, nil
}

func fromJsonNode(n json.Node) tree.Node {
	m := tree.MarkerFromIndices(int64(n.Start), int64(n.End))

	switch v := n.Value.(type) {
	case map[string]json.Node:
		ret := tree.MapNode{
			Marker:   m,
			Children: make(map[string]tree.Node, len(v)),
			Keys:     make(map[string]tree.Leaf, len(v)),
		}
		for key, child := range v {
			ret.Children[key] = fromJsonNode(child)
			ret.Keys[key] = tree.Leaf{
				Marker: tree.MarkerFromIndices(int64(child.KeyStart), int64(child.KeyEnd)),
			}
		}
		return ret
	case []json.Node:
		ret := tree.SliceNode{
			Marker:   m,
			Children: make([]tree.Node, 0, len(v)),
		}
		for _, child := range v {
			ret.Children = append(ret.Children, fromJsonNode(child))
		}
		return ret
	default:
		return tree.Leaf{
			Marker: m,
		}
	}
}
