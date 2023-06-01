package jsonmapping

import (
	"fmt"

	"jsonmaping/jsontree"
)

type DiffMapper struct {
	JsonMapper
	stateTreeRoot *jsontree.Node
}

func NewDiffMapper(response []byte, state []byte) (DiffMapper, error) {
	m, err := NewJsonMapper(response)
	if err != nil {
		return DiffMapper{}, fmt.Errorf("error unmarshalling response: %v", err)
	}

	s, err := m.MapJson(state)
	if err != nil {
		return DiffMapper{}, fmt.Errorf("error unmarshalling state json: %v", err)
	}

	return DiffMapper{m, s}, nil
}

// return a node in new tree which has different value with the old tree.
func (m *DiffMapper) diffTree(newTreeRoot *jsontree.Node, oldTreeRoot *jsontree.Node) (newNode *jsontree.Node, oldNode *jsontree.Node, err error) {
	// map[json_path]value, only keep value node
	pathMap := make(map[string]*jsontree.Node, 0)
	jsontree.Traverse(oldTreeRoot, func(node *jsontree.Node) {
		if node.IsValueNode() {
			pathMap[node.GetJsonPath()] = node
		}
	}, nil)

	diffCnt := 0
	jsontree.Traverse(newTreeRoot, func(node *jsontree.Node) {
		if node.IsValueNode() {
			tmp := pathMap[node.GetJsonPath()]
			if tmp.GetStringValue() != node.GetStringValue() {
				diffCnt++
				newNode = node
				oldNode = tmp
			}
		}
	}, nil)

	if diffCnt > 1 {
		return nil, nil, fmt.Errorf("more than one difference found")
	}

	return
}

// diffResponse will detect the difference between the original response and the new response
// will return the diff node in old tree
func (m *DiffMapper) diffResponse(response []byte) (newNode *jsontree.Node, oldNode *jsontree.Node, err error) {
	tree, err := unmarshalToTree(response)
	if err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling response: %v", err)
	}

	if ok, err := m.validateJsonStruct(tree, m.treeRoot); !ok {
		return nil, nil, fmt.Errorf("json struct is not the same: %v", err)
	}

	return m.diffTree(tree, m.treeRoot)
}

func (m *DiffMapper) DiffJson(newState []byte, newResponse []byte) (stateTree *jsontree.Node, err error) {
	_, oldResponseNode, err := m.diffResponse(newResponse)
	if err != nil {
		return nil, fmt.Errorf("error diffing response: %v", err)
	}

	newStateTree, err := unmarshalToTree(newState)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling new state: %v", err)
	}
	_, oldStateDiffNode, err := m.diffTree(newStateTree, m.stateTreeRoot)

	oldStateDiffNode.SetMappingRef(oldResponseNode)

	// in case some parent node is not mapped before.
	m.stateTreeRoot, err = m.mapTree(m.stateTreeRoot)
	if err != nil {
		return nil, fmt.Errorf("error mapping state tree: %v", err)
	}

	return m.stateTreeRoot, nil
}

func (m *DiffMapper) validateJsonStruct(tree1 *jsontree.Node, tree2 *jsontree.Node) (ok bool, err error) {
	// map[json_path]node
	map1 := make(map[string]*jsontree.Node, 0)
	map2 := make(map[string]*jsontree.Node, 0)

	jsontree.Traverse(tree1, func(n *jsontree.Node) {
		map1[n.GetJsonPath()] = n
	}, nil)
	jsontree.Traverse(tree2, func(n *jsontree.Node) {
		map2[n.GetJsonPath()] = n
	}, nil)

	if len(map1) != len(map2) {
		return false, fmt.Errorf("size of json are not same, map1: %v, map2: %v", len(map1), len(map2))
	}

	for k, _ := range map1 {
		if map2[k] == nil {
			return false, fmt.Errorf("json path %v not found in map2", k)
		}
	}

	return true, nil
}
