package json

import (
	"github.com/ajeddeloh/vcontext"

	"gopkg.in/yaml.v3"
)

func UnmarshalToContext(raw []byte) (vcontext.Node, error) {
	var ast yaml.Node
	if err := yaml.Unmarshal(raw, &ast); err != nil {
		return nil, err
	}
	return fromYamlNode(ast), nil
}

func fromYamlNode(n yaml.Node) vcontext.Node {
	m := vcontext.Marker{
		StartIdx: int64(n.Line),
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
		ret := vcontext.MapNode{
			Marker: m,
			Children: make(map[string]vcontext.Node, len(n.Content)/2),
			Keys:     make(map[string]vcontext.Leaf, len(n.Content)/2),
		}
		// MappingNodes list keys and values like [k, v, k, v...]
		for i := 0; i < len(n.Content); i+=2 {
			key := *n.Content[i]
			value := *n.Content[i+1]
			ret.Keys[key.Value] = vcontext.Leaf{
				Marker: vcontext.Marker{
					StartIdx: int64(key.Line),
				},
			}
			ret.Children[key.Value] = fromYamlNode(value)
		}
		return ret
	case yaml.SequenceNode:
		ret := vcontext.SliceNode{
			Marker: m,
			Children: make([]vcontext.Node, 0, len(n.Content)),
		}
		for _, child := range n.Content {
			ret.Children = append(ret.Children, fromYamlNode(*child))
		}
		return ret
	default: // scalars and aliases
		return vcontext.Leaf{
			Marker: m,
		}
	}
}

