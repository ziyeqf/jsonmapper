package jsontree_test

import (
	"encoding/json"
	"testing"

	"jsonmaping/jsontree"
)

func TestNode_AddChild(t *testing.T) {
	parent := jsontree.NewNode("parent")
	child := jsontree.NewNode("child")
	parent.AddChild(child)

	if len(parent.Children()) != 1 {
		t.Errorf("Expected parent to have 1 child, but got %d", len(parent.Children()))
	}

	if parent.GetChild("child") != child {
		t.Error("Expected parent's child to be the same as the added child")
	}

	if child.Parent() != parent {
		t.Error("Expected child's parent to be the same as the parent")
	}
}

func TestUnmarshal_basic(t *testing.T) {
	type unmarshalCase struct {
		name     string
		jsonStr  string
		keyName  string
		expected jsontree.Node
	}

	cases := []unmarshalCase{
		{
			name:     "string",
			jsonStr:  `{"foo":"bar"}`,
			keyName:  "foo",
			expected: *jsontree.NewNode("$").AddChild(jsontree.NewNode("foo").SetValue("bar")),
		},
		{
			name:     "int",
			jsonStr:  `{"foo":123}`,
			keyName:  "foo",
			expected: *jsontree.NewNode("$").AddChild(jsontree.NewNode("foo").SetValue(float64(123))),
		},
		{
			name:    "float64",
			jsonStr: `{"foo":1.23}`,
			keyName: "foo",
			expected: func() jsontree.Node {
				return *jsontree.NewNode("$").AddChild(jsontree.NewNode("foo").SetValue(1.23))
			}(),
		},
		{
			name:     "bool",
			jsonStr:  `{"foo":true}`,
			keyName:  "foo",
			expected: *jsontree.NewNode("$").AddChild(jsontree.NewNode("foo").SetValue(true)),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var j map[string]interface{}
			err := json.Unmarshal([]byte(c.jsonStr), &j)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			var actual jsontree.Node
			actual.Unmarshal(j)

			if actual.GetChild(c.keyName).GetStringValue() != c.expected.GetChild(c.keyName).GetStringValue() {
				t.Errorf("Expected %s, but got %s", c.expected.GetChild(c.keyName).GetStringValue(), actual.GetChild(c.keyName).GetStringValue())
			}
		})
	}

}

func TestUnmarshal_array(t *testing.T) {
	jsonStr := `
{"foo":["bar","baz"]}
`
	expectedMap := make(map[string]string, 0)
	expectedMap["bar"] = "bar"
	expectedMap["baz"] = "baz"

	var j map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &j)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var actual jsontree.Node
	actual.Unmarshal(j)

	if len(actual.GetChild("foo").Children()) != 2 {
		t.Errorf("Expected 2 children, but got %d", len(actual.GetChild("foo").Children()))
	}

	for _, c := range actual.GetChild("foo").Children() {
		if _, ok := expectedMap[c.GetStringValue()]; !ok {
			t.Errorf("Expected value does not exist but exist: %v", c.GetStringValue())
		}
	}
}

func TestUnmarshal_map(t *testing.T) {
	jsonStr := `
{"foo":{"bar":"baz"}}
`
	var j map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &j)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	actual := jsontree.NewNode("$")
	actual.Unmarshal(j)

	if actual.GetChild("foo").GetChild("bar").GetStringValue() != "baz" {
		t.Errorf("Expected %v, but got %v", "baz", actual.GetChild("foo").GetChild("bar").GetStringValue())
	}
}

func TestNode_IsArrayNode(t *testing.T) {
	jsonStr := `
	[{"a":"1"},{"c":"d"}]
`
	node := jsontree.NewNode("$")
	if node.IsArrayNode() {
		t.Error("Expected node to not be an array node")
	}

	var j []interface{}
	err := json.Unmarshal([]byte(jsonStr), &j)
	if err != nil {
		t.Errorf("error unmarshalling input: %v", err)
	}
	node.Unmarshal(j)
	if !node.IsArrayNode() {
		t.Error("Expected node to be an array node")
	}
}

func TestIsMapToSameParent_map(t *testing.T) {
	mappedFoo := jsontree.NewNode("foo")
	mappedBar := jsontree.NewNode("bar")
	mapped := jsontree.NewNode("$").AddChild(mappedFoo).AddChild(mappedBar)

	actual := jsontree.NewNode("$").
		AddChild(jsontree.NewNode("foo").SetMappingRef(mappedFoo)).
		AddChild(jsontree.NewNode("bar").SetMappingRef(mappedBar))

	if n, yes := jsontree.IsMapToSameParent(actual.Children()); !yes || n != mapped {
		t.Error("Expected to be mapped to the same parent")
	}

	mappedFoo = jsontree.NewNode("foo")
	mappedBar = jsontree.NewNode("bar")
	mapped = jsontree.NewNode("$").AddArrChild(mappedFoo).AddArrChild(mappedBar)

	actual = jsontree.NewNode("$").
		AddArrChild(jsontree.NewNode("foo").SetMappingRef(mappedFoo)).
		AddArrChild(jsontree.NewNode("bar").SetMappingRef(mappedBar))
	if n, yes := jsontree.IsMapToSameParent(actual.Children()); !yes || n != mapped {
		t.Error("Expected to be mapped to the same parent")
	}
}

func TestGetJsonPath(t *testing.T) {
	jsonStr := `{
        "name": "John Doe",
        "age": 30,
        "address": {
            "city": "New York",
            "state": "NY"
        },
        "phoneNumbers": [
            {
                "type": "home",
                "number": "555-555-1234"
            },
            {
                "type": "work",
                "number": "555-555-5678"
            }
        ]
    }`

	var j map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &j)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	node := jsontree.NewNode("$")
	node.Unmarshal(j)

	expected := "$.address.state"
	actual := node.GetChild("address").GetChild("state").GetJsonPath()

	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}

	expected = "$.phoneNumbers.0.number"
	actual = node.GetChild("phoneNumbers").Children()[0].GetChild("number").GetJsonPath()
	if actual != expected {
		t.Errorf("Expected %v, but got %v", expected, actual)
	}

}
