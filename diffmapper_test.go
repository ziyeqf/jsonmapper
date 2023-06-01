package jsonmapping_test

import (
	"testing"

	jsonmapping "jsonmaping"
	"jsonmaping/jsontree"
)

func TestDiffMapper_DiffJson(t *testing.T) {
	type testCase struct {
		name            string
		responseJson    string
		stateJson       string
		newResponseJson string
		newStateJson    string
		mappedCnt       int
	}

	cases := []testCase{
		{
			name:            "string",
			responseJson:    `{"foo":"bar"}`,
			stateJson:       `{"foo1":"bar1"}`,
			newResponseJson: `{"foo":"baz"}`,
			newStateJson:    `{"foo1":"baz1"}`,
			mappedCnt:       2,
		},
		{
			name:            "int",
			responseJson:    `{"foo":1}`,
			stateJson:       `{"foo1":2}`,
			newResponseJson: `{"foo":3}`,
			newStateJson:    `{"foo1":4}`,
			mappedCnt:       2,
		},
		{
			name:            "float64",
			responseJson:    `{"foo":1.0}`,
			stateJson:       `{"foo1":2.0}`,
			newResponseJson: `{"foo":3.0}`,
			newStateJson:    `{"foo1":4.0}`,
			mappedCnt:       2,
		},
		{
			name:            "bool",
			responseJson:    `{"foo":true}`,
			stateJson:       `{"foo1":false}`,
			newResponseJson: `{"foo":false}`,
			newStateJson:    `{"foo1":true}`,
			mappedCnt:       2,
		},
		{
			name:            "complex",
			responseJson:    `{"foo":{"bar":{"baz":1,"qux":4}}}`,
			stateJson:       `{"foo1":{"bar1":{"baz1":3,"qux1":4}}}`,
			newResponseJson: `{"foo":{"bar":{"baz":5,"qux":4}}}`,
			newStateJson:    `{"foo1":{"bar1":{"baz1":7,"qux1":4}}}`,
			mappedCnt:       5,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mapper, err := jsonmapping.NewDiffMapper([]byte(c.responseJson), []byte(c.stateJson))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			mapTree, err := mapper.DiffJson([]byte(c.newStateJson), []byte(c.newResponseJson))
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
				t.Errorf("mappedCnt: %v, want: %v", mappedCnt, c.mappedCnt)
			}

		})
	}
}
