package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

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

type OpDefinitions []OpDefinition

func (o OpDefinitions) Len() int {
	return len(o)
}

func (o OpDefinitions) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}

func (o OpDefinitions) Less(i, j int) bool {
	return o[i].Path < o[j].Path
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
	sort.Sort(opDefs)

	opsfile, _ := yaml.Marshal(opDefs)
	fmt.Printf("%s", string(opsfile))
}

func compareObjects(c Comparator) OpDefinitions {
	switch source := c.Source.(type) {
	case map[interface{}]interface{}:
		target, ok := c.Target.(map[interface{}]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Skipping path %s; replace not yet implemented\n", c.Path)
			return OpDefinitions{}
		}

		return compareMaps(source, target, c.Path)
	case []interface{}:
		target, ok := c.Target.([]interface{})
		if !ok {
			fmt.Fprintf(os.Stderr, "Skipping path %s; replace not yet implemented\n", c.Path)
			return OpDefinitions{}
		}

		return compareSlices(source, target, c.Path)
	default:
		return OpDefinitions{}
	}
}

func compareMaps(source, target map[interface{}]interface{}, currentPath string) (opDefs OpDefinitions) {
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

func compareSlices(source, target []interface{}, currentPath string) (opDefs OpDefinitions) {
	sourceIds := findUniqueIds(source)
	targetIds := findUniqueIds(target)

	nameSourceIds := map[string]int{}
	for id, sourceIndex := range sourceIds {
		if strings.HasPrefix(id, "name=") {
			nameSourceIds[id] = sourceIndex
		}
	}

	comparators := []Comparator{}
	matchedSourceIds := map[int]bool{}

	var targetIndex int
	var ok bool
	for id, sourceIndex := range nameSourceIds {
		targetIndex, ok = targetIds[id]
		if ok {
			comparators = append(comparators, Comparator{
				Source: source[sourceIndex],
				Target: target[targetIndex],
				Path:   fmt.Sprintf("%s%s/", currentPath, id),
			})

			matchedSourceIds[sourceIndex] = true
		}
	}

	for id, sourceIndex := range sourceIds {
		if matchedSourceIds[sourceIndex] {
			continue
		}

		targetIndex, ok = targetIds[id]
		if ok {
			comparators = append(comparators, Comparator{
				Source: source[sourceIndex],
				Target: target[targetIndex],
				Path:   fmt.Sprintf("%s%s/", currentPath, id),
			})

			matchedSourceIds[sourceIndex] = true
		}
	}

	for id, sourceIndex := range nameSourceIds {
		if !matchedSourceIds[sourceIndex] {
			opDefs = append(opDefs, buildOpDefinition(currentPath, id))
			matchedSourceIds[sourceIndex] = true
		}
	}

	for id, sourceIndex := range sourceIds {
		if !matchedSourceIds[sourceIndex] {
			opDefs = append(opDefs, buildOpDefinition(currentPath, id))
			matchedSourceIds[sourceIndex] = true
		}
	}

	for _, comparator := range comparators {
		opDefs = append(opDefs, compareObjects(comparator)...)
	}

	return
}

func findUniqueIds(items []interface{}) map[string]int {
	idPresence := map[string][]int{}
	for i, item := range items {
		for _, id := range getIDsForItem(item) {
			idPresence[id] = append(idPresence[id], i)
		}
	}

	uniqueIds := map[string]int{}
	for id, indices := range idPresence {
		if len(indices) == 1 {
			uniqueIds[id] = indices[0]
		}
	}

	return uniqueIds
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
