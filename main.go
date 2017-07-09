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
	sourceIDs := findUniqueIDs(source)
	targetIDs := findUniqueIDs(target)

	comparators := []Comparator{}
	handledSources := map[*interface{}]bool{}
	for id, sourceEl := range sourceIDs {
		if !handledSources[sourceEl] && targetIDs[id] != nil {
			comparators = append(comparators, Comparator{
				Source: *sourceEl,
				Target: *targetIDs[id],
				Path:   fmt.Sprintf("%s%s/", currentPath, id),
			})

			handledSources[sourceEl] = true
		}
	}

	for _, comparator := range comparators {
		opDefs = append(opDefs, compareObjects(comparator)...)
	}

	for id, sourceEl := range sourceIDs {
		if !handledSources[sourceEl] {
			opDefs = append(opDefs, buildOpDefinition(currentPath, id))
			handledSources[sourceEl] = true
		}
	}

	return
}

func findUniqueIDs(mapArray []interface{}) map[string]*interface{} {
	iDPresence := map[string][]*interface{}{}
	for i, el := range mapArray {
		for _, id := range getIDsForItem(el) {
			iDPresence[id] = append(iDPresence[id], &mapArray[i])
		}
	}

	uniqueIDs := map[string]*interface{}{}
	for id, mapArray := range iDPresence {
		if len(mapArray) == 1 {
			uniqueIDs[id] = mapArray[0]
		}
	}

	return uniqueIDs
}

func getIDsForItem(item interface{}) []string {
	switch item := item.(type) {
	case string:
		return []string{item}
	case int:
		return []string{string(item)}
	case map[interface{}]interface{}:
		ids := []string{}
		for key, value := range item {
			value, ok := value.(string)
			if !ok {
				continue
			}

			ids = append(ids, fmt.Sprintf("%s=%s", key.(string), value))
		}

		return ids
	default:
		return []string{}
	}
}

func buildOpDefinition(location, key string) OpDefinition {
	return OpDefinition{
		Type: "remove",
		Path: location + key,
	}
}
