package jsonmapping_test

import (
	"encoding/json"
	"testing"

	jsonmapping "jsonmaping"
)

func TestJsonMapper_MapModelToReference(t *testing.T) {
	type mapJsonCase struct {
		name          string
		referenceJson string
		stateJson     string
		expected      string
	}

	cases := []mapJsonCase{
		{
			name:          "string",
			referenceJson: `{"foo1":"bar"}`,
			stateJson:     `{"foo2":"bar"}`,
			expected:      `{"/foo2":["/response/foo1"]}`,
		},
		{
			name:          "int",
			referenceJson: `{"foo1":123}`,
			stateJson:     `{"foo2":123}`,
			expected:      `{"/foo2":["/response/foo1"]}`,
		},
		{
			name:          "float64",
			referenceJson: `{"foo1":1.23}`,
			stateJson:     `{"foo2":1.23}`,
			expected:      `{"/foo2":["/response/foo1"]}`,
		},
		{
			name:          "bool",
			referenceJson: `{"foo1":true}`,
			stateJson:     `{"foo2":true}`,
			expected:      `{"/foo2":["/response/foo1"]}`,
		},
		{
			name:          "array",
			referenceJson: `{"foo1":["bar1","bar2"]}`,
			stateJson:     `{"foo2":["bar1","bar2"]}`,
			expected:      `{"/foo2/0":["/response/foo1/0","/response/foo1/0"]}`,
		},
		{
			name:          "object",
			referenceJson: `{"foo1":{"bar1":"test1","bar2":"test2"}}`,
			stateJson:     `{"foo2":{"baz1":"test1","baz2":"test2"}}`,
			expected:      `{"/foo2/baz1":["/response/foo1/bar1"],"/foo2/baz2":["/response/foo1/bar2"]}`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var ref interface{}
			err := json.Unmarshal([]byte(c.referenceJson), &ref)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			refMap := make(map[string]interface{}, 0)
			refMap["response"] = ref

			mapper, err := jsonmapping.NewMapper(refMap)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			var model interface{}
			err = json.Unmarshal([]byte(c.stateJson), &model)

			mapping, err := mapper.MapModelToReference(model)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			mappingStr, err := json.Marshal(mapping)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if string(mappingStr) != c.expected {
				t.Errorf("Expected to get %v, but got %v", c.expected, string(mappingStr))
			}
		})
	}
}
