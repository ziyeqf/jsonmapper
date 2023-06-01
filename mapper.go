package jsonmapping

import (
	"encoding/json"
	"fmt"
	"strconv"

	"jsonmaping/jsontree"
)

type JsonMapper struct {
	valueKeyMap map[string]*jsontree.Node
	treeRoot    *jsontree.Node
}

func unmarshalToTree(jsonByte []byte) (*jsontree.Node, error) {
	unmarshalledResponse := make(map[string]interface{}, 0)

	err := json.Unmarshal(jsonByte, &unmarshalledResponse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	root := jsontree.NewNode("$")
	root.Unmarshal(unmarshalledResponse)

	return root, nil
}

func NewJsonMapper(response []byte) (m JsonMapper, err error) {
	m = JsonMapper{}
	m.treeRoot, err = unmarshalToTree(response)
	if err != nil {
		return m, fmt.Errorf("error unmarshalling response: %v", err)
	}
	// generate a map for finding the node by value
	m.valueKeyMap = make(map[string]*jsontree.Node, 0)
	jsontree.Traverse(m.treeRoot, func(n *jsontree.Node) {
		if n.IsValueNode() {
			m.valueKeyMap[n.GetStringValue()] = n
		}
	}, nil)

	return m, nil
}

// FindNodeByValue could only find string, float64 and bool value.
func (m *JsonMapper) FindNodeByValue(value interface{}) *jsontree.Node {
	key := ""
	switch value.(type) {
	case string:
		key = value.(string)
	case float64:
		key = fmt.Sprintf("%f", value.(float64))
	case bool:
		key = strconv.FormatBool(value.(bool))
	}
	return m.valueKeyMap[key]
}

// mapTree maps a json tree to the response json tree.
func (m *JsonMapper) mapTree(mapTree *jsontree.Node) (*jsontree.Node, error) {
	// traverse the map tree and find the corresponding node in the response tree
	jsontree.Traverse(mapTree, func(n *jsontree.Node) {
		if n.IsValueNode() {
			if mappedNode := m.FindNodeByValue(n.GetStringValue()); mappedNode != nil {
				n.SetMappingRef(mappedNode)
			}
		}
	}, func(n *jsontree.Node) {
		if !n.IsValueNode() {
			if parent, yes := jsontree.IsMapToSameParent(n.Children()); yes {
				n.SetMappingRef(parent)
			}
		}
	})

	return mapTree, nil
}

// MapJson maps a json to the response json tree.
func (m *JsonMapper) MapJson(input []byte) (*jsontree.Node, error) {
	mapTree, err := unmarshalToTree(input)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling state json: %v", err)
	}

	return m.mapTree(mapTree)
}
