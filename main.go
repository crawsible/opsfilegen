package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Comparator struct {
	Source interface{}
	Target interface{}
	Path   string
}

type OpDefinition struct {
	Type  string
	Path  string
	Value string `yaml:",omitempty"`
}

func main() {
	sourceFilename := os.Args[1]
	targetFilename := os.Args[2]

	sourceBytes, _ := ioutil.ReadFile(sourceFilename)
	targetBytes, _ := ioutil.ReadFile(targetFilename)

	c := Comparator{Path: "/"}
	_ = yaml.Unmarshal(sourceBytes, &c.Source)
	_ = yaml.Unmarshal(targetBytes, &c.Target)

	opDefs := compareObjects(c)

	opsfile, _ := yaml.Marshal(opDefs)
	fmt.Printf("%s", string(opsfile))
}

func compareObjects(c Comparator) []OpDefinition {
	switch source := c.Source.(type) {
	case map[interface{}]interface{}:
		target, ok := c.Target.(map[interface{}]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Skipping path %s; replace not yet implemented\n", c.Path)
			return []OpDefinition{}
		}

		return compareMaps(source, target, c.Path)
	case []interface{}:
		target, ok := c.Target.([]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Skipping path %s; replace not yet implemented\n", c.Path)
			return []OpDefinition{}
		}

		return compareSlices(source, target, c.Path)
	default:
		return []OpDefinition{}
	}
}

func compareMaps(source, target map[interface{}]interface{}, currentPath string) (opDefs []OpDefinition) {
	for key, value := range source {
		key := key.(string)
		if target[key] == nil {
			opDefs = append(opDefs, buildOpDefinition(currentPath, key))
			continue
		}

		comparator := Comparator{
			Source: value,
			Target: target[key],
			Path:   currentPath + key + "/",
		}
		opDefs = append(opDefs, compareObjects(comparator)...)
	}

	return
}

func compareSlices(source, target []interface{}, currentPath string) (opDefs []OpDefinition) {
	sourceKVs := findUniqueKVs(source)
	targetKVs := findUniqueKVs(target)

	comparators := []Comparator{}
	handledSources := map[*interface{}]bool{}
	for kv, sourceEl := range sourceKVs {
		if !handledSources[sourceEl] && targetKVs[kv] != nil {
			comparators = append(comparators, Comparator{
				Source: *sourceEl,
				Target: *targetKVs[kv],
				Path:   fmt.Sprintf("%s%s/", currentPath, kv),
			})

			handledSources[sourceEl] = true
		}
	}

	for _, comparator := range comparators {
		opDefs = append(opDefs, compareObjects(comparator)...)
	}

	for kv, sourceEl := range sourceKVs {
		if !handledSources[sourceEl] {
			opDefs = append(opDefs, buildOpDefinition(currentPath, kv))
			handledSources[sourceEl] = true
		}
	}

	return
}

func findUniqueKVs(mapArray []interface{}) map[string]*interface{} {
	kVPresence := map[string][]*interface{}{}
	for i, el := range mapArray {
		el := el.(map[interface{}]interface{})
		for key, value := range el {
			value, ok := value.(string)
			if !ok {
				continue
			}

			kv := fmt.Sprintf("%s=%s", key.(string), value)
			kVPresence[kv] = append(kVPresence[kv], &mapArray[i])
		}
	}

	uniqueKVs := map[string]*interface{}{}
	for kv, mapArray := range kVPresence {
		if len(mapArray) == 1 {
			uniqueKVs[kv] = mapArray[0]
		}
	}

	return uniqueKVs
}

func buildOpDefinition(location, key string) OpDefinition {
	return OpDefinition{
		Type: "remove",
		Path: location + key,
	}
}
