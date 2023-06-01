package jsontree

import (
	"fmt"
	"strconv"
)

type Node struct {
	// only exists in non-root nodes.
	parent     *Node
	mappingRef *Node

	// arrChildren and children are mutually exclusive.
	arrChildren []*Node
	children    map[string]*Node
	// JsonName is the name of the node in the json struct.
	jsonName string
	// value only exists in leaf nodes.
	// possible type of value is string, float64, bool, nil
	value interface{}
}

func (n *Node) AddChild(child *Node) *Node {
	if n.children == nil {
		n.children = make(map[string]*Node, 0)
	}
	n.children[child.jsonName] = child
	child.parent = n
	return n
}

func (n *Node) AddArrChild(child *Node) *Node {
	if n.arrChildren == nil {
		n.arrChildren = make([]*Node, 0)
	}
	n.arrChildren = append(n.arrChildren, child)
	child.parent = n
	return n
}

func (n *Node) GetChild(jsonName string) *Node {
	return n.children[jsonName]
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) IsArrayNode() bool {
	return n.arrChildren != nil
}

func (n *Node) IsValueNode() bool {
	return n.arrChildren == nil && n.children == nil
}

// SetMappingRef adds a reference to another node in another tree.
func (n *Node) SetMappingRef(ref *Node) *Node {
	n.mappingRef = ref
	return n
}

func (n *Node) GetMappingRef() *Node {
	return n.mappingRef
}

// HasChild returns true if the node has named children.
// array children are not included.
func (n *Node) HasChild() bool {
	return n.children != nil && len(n.children) > 0
}

func (n *Node) Children() []*Node {
	if n.IsArrayNode() {
		return n.arrChildren
	}
	var children []*Node
	for _, v := range n.children {
		children = append(children, v)
	}
	return children
}

func (n *Node) GetJsonPath() string {
	if n.parent == nil {
		return n.jsonName
	}
	if n.IsArrayNode() {
		return n.parent.GetJsonPath() + "." + n.jsonName
	}
	return n.parent.GetJsonPath() + "." + n.jsonName
}

func (n *Node) SetValue(value interface{}) *Node {
	n.value = value
	return n
}

func (n *Node) GetStringValue() string {
	if !n.IsValueNode() {
		return ""
	}
	switch n.value.(type) {
	case string:
		return n.value.(string)
	case float64:
		return fmt.Sprintf("%f", n.value.(float64))
	case bool:
		return strconv.FormatBool(n.value.(bool))
	}
	return ""
}

func (n *Node) Unmarshal(j interface{}) {
	switch j.(type) {
	case map[string]interface{}:
		for k, v := range j.(map[string]interface{}) {
			child := NewNode(k)
			child.Unmarshal(v)
			child.jsonName = k
			n.AddChild(child)
		}

	case []interface{}:
		for i, v := range j.([]interface{}) {
			child := NewNode(fmt.Sprintf("%v", i))
			child.Unmarshal(v)
			n.AddArrChild(child)
		}
	case string:
		n.value = j
	case float64:
		n.value = j
	case bool:
		n.value = j
	}
}

func NewNode(jsonName string) *Node {
	c := Node{
		jsonName: jsonName,
	}
	return &c
}

// Traverse
func Traverse(node *Node, headFunc func(*Node), tailFunc func(*Node)) {
	if headFunc != nil {
		headFunc(node)
	}
	if node.IsArrayNode() {
		for _, child := range node.arrChildren {
			Traverse(child, headFunc, tailFunc)
		}
	} else {
		for _, child := range node.children {
			Traverse(child, headFunc, tailFunc)
		}
	}
	if tailFunc != nil {
		tailFunc(node)
	}
}

func IsMapToSameParent(nodes []*Node) (*Node, bool) {
	if len(nodes) == 0 {
		return nil, false
	}
	if len(nodes) == 1 {
		return nodes[0].parent, true
	}
	if nodes[0].mappingRef == nil {
		return nil, false
	}
	parent := nodes[0].mappingRef.parent
	for _, node := range nodes {
		if node.mappingRef == nil {
			return nil, false
		}
		if node.mappingRef.parent != parent {
			return nil, false
		}
	}
	return parent, parent != nil
}
