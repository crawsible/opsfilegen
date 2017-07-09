package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

type OpDefinition struct {
	Type  string
	Path  string
	Value string `yaml:",omitempty"`
}

var sourceMarkup = []byte(`
---
consistent_value:
  nested_consistent_value: "value"
  nested_removed_value: "value"
removed_value: "value"
`)

var targetMarkup = []byte(`
---
consistent_value:
  nested_consistent_value: "value"
`)

func main() {
	currentPath := "/"

	var source interface{}
	var target interface{}

	_ = yaml.Unmarshal(sourceMarkup, &source)
	_ = yaml.Unmarshal(targetMarkup, &target)

	opDefs := compareObjects(source, target, currentPath)

	opsfile, _ := yaml.Marshal(opDefs)
	fmt.Printf("%s", string(opsfile))
}

func compareObjects(source, target interface{}, currentPath string) (opDefs []OpDefinition) {
	switch source := source.(type) {
	case map[interface{}]interface{}:
		target := target.(map[interface{}]interface{})

		for key, value := range source {
			key := key.(string)
			if target[key] == nil {
				opDefs = append(opDefs, buildOpDefinition(currentPath, key))
				continue
			}

			switch value := value.(type) {
			case map[interface{}]interface{}:
				newPath := currentPath + key + "/"
				opDefs = append(opDefs, compareObjects(value, target[key], newPath)...)
			}
		}
	}

	return
}

func buildOpDefinition(location, key string) OpDefinition {
	return OpDefinition{
		Type: "remove",
		Path: location + key,
	}
}
