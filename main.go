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
consistent_value: "value"
removed_value: "value"
another_removed_value: "new_value"
`)

var targetMarkup = []byte(`
---
consistent_value: "value"
`)

func main() {
	var opDefs []OpDefinition
	currentPath := "/"

	var sourceRaw interface{}
	var targetRaw interface{}

	_ = yaml.Unmarshal(sourceMarkup, &sourceRaw)
	_ = yaml.Unmarshal(targetMarkup, &targetRaw)

	switch source := sourceRaw.(type) {
	case map[interface{}]interface{}:
		target := targetRaw.(map[interface{}]interface{})

		for key, _ := range source {
			key := key.(string)
			if target[key] == nil {
				opDefs = append(opDefs, buildOpDefinition(currentPath, key))
			}
		}
	}

	opsfile, _ := yaml.Marshal(opDefs)
	fmt.Printf("%s", string(opsfile))
}

func buildOpDefinition(location, key string) OpDefinition {
	return OpDefinition{
		Type: "remove",
		Path: location + key,
	}
}
