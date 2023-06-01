package jsonmapping_test

import (
	"fmt"
	"testing"

	jsonmapping "jsonmaping"
	"jsonmaping/jsontree"
)

func TestJsonMapper_FindNodeByValue(t *testing.T) {
	type findNodeByValueCase struct {
		name    string
		jsonStr string
		value   string
	}

	cases := []findNodeByValueCase{
		{
			name:    "string",
			jsonStr: `{"foo":"bar"}`,
			value:   "bar",
		},
		{
			name:    "int",
			jsonStr: `{"foo":123}`,
			value:   fmt.Sprintf("%f", float64(123)),
		},
		{
			name:    "float64",
			jsonStr: `{"foo":1.23}`,
			value:   fmt.Sprintf("%f", 1.23),
		},
		{
			name:    "bool",
			jsonStr: `{"foo":true}`,
			value:   "true",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mapper, err := jsonmapping.NewJsonMapper([]byte(c.jsonStr))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			node := mapper.FindNodeByValue(c.value)

			if node == nil {
				t.Errorf("Expected to find node with value %s", c.value)
			}

		})
	}
}

func TestJsonMapper_MapJson(t *testing.T) {
	type mapJsonCase struct {
		name         string
		responseJson string
		stateJson    string
		mappedCnt    int
	}

	cases := []mapJsonCase{
		{
			name:         "string",
			responseJson: `{"foo1":"bar"}`,
			stateJson:    `{"foo2":"bar"}`,
			mappedCnt:    2,
		},
		{
			name:         "int",
			responseJson: `{"foo1":123}`,
			stateJson:    `{"foo2":123}`,
			mappedCnt:    2,
		},
		{
			name:         "float64",
			responseJson: `{"foo1":1.23}`,
			stateJson:    `{"foo2":1.23}`,
			mappedCnt:    2,
		},
		{
			name:         "bool",
			responseJson: `{"foo1":true}`,
			stateJson:    `{"foo2":true}`,
			mappedCnt:    2,
		},
		{
			name:         "array",
			responseJson: `{"foo1":["bar1","bar2"]}`,
			stateJson:    `{"foo2":["bar1","bar2"]}`,
			mappedCnt:    4,
		},
		{
			name:         "object",
			responseJson: `{"foo1":{"bar1":"test1","bar2":"test2"}}`,
			stateJson:    `{"foo2":{"baz1":"test1","baz2":"test2"}}`,
			mappedCnt:    4,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mapper, err := jsonmapping.NewJsonMapper([]byte(c.responseJson))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			mapTree, err := mapper.MapJson([]byte(c.stateJson))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			mappedCnt := 0
			jsontree.Traverse(mapTree, func(node *jsontree.Node) {
				if node.GetMappingRef() != nil {
					mappedCnt++
				}
			}, nil)
			if mappedCnt != c.mappedCnt {
				t.Errorf("Expected to find %d mapped nodes, but got %d", c.mappedCnt, mappedCnt)
			}
		})
	}
}
