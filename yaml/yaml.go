package json

import (
	"github.com/ajeddeloh/vcontext/tree"

	"gopkg.in/yaml.v3"
)

func UnmarshalToContext(raw []byte) (tree.Node, error) {
	var ast yaml.Node
	if err := yaml.Unmarshal(raw, &ast); err != nil {
		return nil, err
	}
	return fromYamlNode(ast), nil
}

func fromYamlNode(n yaml.Node) tree.Node {
	m := tree.Marker{
		StartP: &tree.Pos{
			Line:   int64(n.Line),
			Column: int64(n.Column),
		},
	}
	switch n.Kind {
	case 0:
		// empty
		return nil
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return nil
		}
		return fromYamlNode(*n.Content[1])
	case yaml.MappingNode:
		ret := tree.MapNode{
			Marker:   m,
			Children: make(map[string]tree.Node, len(n.Content)/2),
			Keys:     make(map[string]tree.Leaf, len(n.Content)/2),
		}
		// MappingNodes list keys and values like [k, v, k, v...]
		for i := 0; i < len(n.Content); i += 2 {
			key := *n.Content[i]
			value := *n.Content[i+1]
			ret.Keys[key.Value] = tree.Leaf{
				Marker: tree.Marker{
					StartP: &tree.Pos{
						Line:   int64(key.Line),
						Column: int64(key.Column),
					},
				},
			}
			ret.Children[key.Value] = fromYamlNode(value)
		}
		return ret
	case yaml.SequenceNode:
		ret := tree.SliceNode{
			Marker:   m,
			Children: make([]tree.Node, 0, len(n.Content)),
		}
		for _, child := range n.Content {
			ret.Children = append(ret.Children, fromYamlNode(*child))
		}
		return ret
	default: // scalars and aliases
		return tree.Leaf{
			Marker: m,
		}
	}
}
