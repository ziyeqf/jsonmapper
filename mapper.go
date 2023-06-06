package jsonmapping

import (
	"fmt"
	"strconv"
)

type JsonMapper struct {
	valuePathMap map[string]string
}

func NewMapper(referenceModels map[string]interface{}) (m JsonMapper, err error) {
	m = JsonMapper{}
	if len(referenceModels) == 0 {
		return m, fmt.Errorf("the reference model is empty")
	}

	m.valuePathMap = generateValuePathMap(referenceModels, "")

	return m, nil
}

func (m *JsonMapper) MapModelToReference(model interface{}) (mapping map[string][]string, err error) {
	mapping = make(map[string][]string, 0)
	//if len(model) == 0 {
	//	return mapping, fmt.Errorf("the model is empty")
	//}

	valuePathMap := generateValuePathMap(model, "")
	for value, path := range valuePathMap {
		if refPath, ok := m.valuePathMap[value]; ok {
			mapping[path] = append(mapping[path], refPath)
		}
	}

	return mapping, nil
}

func generateValuePathMap(input interface{}, path string) map[string]string {
	result := make(map[string]string)
	switch input.(type) {
	case map[string]interface{}:
		for k, v := range input.(map[string]interface{}) {
			tmpMap := generateValuePathMap(v, path+"/"+k)
			for k, v := range tmpMap {
				result[k] = v
			}
		}
	case []interface{}:
		for _, v := range input.([]interface{}) {
			tmpMap := generateValuePathMap(v, path+"/0")
			for k, v := range tmpMap {
				result[k] = v
			}
		}
	default:
		result[ToString(input)] = path
	}
	return result
}

func ToString(input interface{}) string {
	switch input.(type) {
	case string:
		return input.(string)
	case float64:
		return fmt.Sprintf("%f", input.(float64))
	case bool:
		return strconv.FormatBool(input.(bool))
	}
	return ""
}
