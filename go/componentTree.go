package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"cuelang.org/go/cue"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Component struct {
	cuePath   cue.Path
	waitFor   []*Component
	resources []*Resource
}
type Resource struct {
	kubeRef KubeRef
	cuePath cue.Path
	healthy bool
}
type KubeRef struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   metav1.ObjectMeta `json:"metadata"`
	Spec       interface{}       `json:"spec,omitempty"`   // loosely typed
	Status     interface{}       `json:"status,omitempty"` // loosely typed
}

func exportComponentsToDot(componentList []*Component) {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape=box];\n")

	depList := []string{}
	for _, component := range componentList {
		for _, dep := range component.waitFor {
			depString := "\"" + component.cuePath.String() + "\" -> \"" + dep.cuePath.String() + "\""
			if !_containsString(depList, depString) {
				sb.WriteString(fmt.Sprintf("    %s;\n", depString))
				depList = append(depList, depString)
			}
		}
	}

	sb.WriteString("}")
	dot := sb.String()

	fileError := os.WriteFile("components.dot", []byte(dot), 0644)
	if fileError != nil {
		log.Fatal(fileError)
	}
	fmt.Println("Graph written to components.dot")
}

func _containsString(depList []string, depString string) bool {
	for _, dep := range depList {
		if dep == depString {
			return true
		}
	}
	return false
}

func makeComponentTree(v cue.Value, rootVal cue.Value) []*Component {
	// make a component tree
	mainComponent := _makeComponentTree(v, nil)
	// convert tree-like structure to list
	componentList := []*Component{}
	_makeComponentList(mainComponent, &componentList)
	// add cue-dependencies as component dependencies
	_addCueDependencies(componentList, rootVal)
	return componentList
}

func _addCueDependencies(componentList []*Component, rootVal cue.Value) {
	for _, component := range componentList {
		fmt.Println(component.cuePath.String())
		dependencyMap := makeDepsMap(rootVal.LookupPath(component.cuePath).LookupPath(cue.ParsePath("in")))
		for _, deps := range dependencyMap {
			for _, dep := range deps {
				compPath, compType := getComponentPathFromResourcePath(dep)
				if compType == "out" {
					comp := _findComponentByPath(componentList, compPath)
					component.waitFor = append(component.waitFor, comp)
				}
			}
		}

		fmt.Println()
	}
}

func _findComponentByPath(components []*Component, path string) *Component {
	for i := range components {
		if components[i].cuePath.String() == path {
			return components[i]
		}
	}
	return nil
}

func getComponentPathFromResourcePath(path string) (string, string) {
	lastIn := strings.LastIndex(path, ".in.")
	lastOut := strings.LastIndex(path, ".out.")
	last := lastIn
	lastType := "in"
	if lastOut > lastIn {
		last = lastOut
		lastType = "out"
	}
	if last == -1 {
		return "", "" // no ".in." or ".out." found
	}
	return path[:last], lastType
}

func _makeComponentList(c *Component, list *[]*Component) {
	*list = append(*list, c)
	for _, child := range c.waitFor {
		if !_containsComponent(list, child) {
			_makeComponentList(child, list)
		}
	}
}

func _containsComponent(list *[]*Component, child *Component) bool {
	for _, component := range *list {
		if component.cuePath.String() == child.cuePath.String() {
			return true
		}
	}
	return false
}

func _makeComponentTree(v cue.Value, parent *Component) *Component {
	componentsCue := v.LookupPath(cue.ParsePath("out.components"))
	resourcesCue := v.LookupPath(cue.ParsePath("out.resources"))
	if !resourcesCue.Exists() && !componentsCue.Exists() {
		fmt.Print("No resources or component found in " + v.Path().String())
		panic("Integrity check failed! No resources or component found in " + v.Path().String())
	}
	resourcesIter, errResources := resourcesCue.Fields()
	if errResources != nil {
		fmt.Print("Error getting resources: " + errResources.Error())
		panic("Error getting resources: " + errResources.Error())
	}

	myResources := []*Resource{}
	for resourcesIter.Next() {
		kuberef := KubeRef{}
		resourcesIter.Value().Decode(&kuberef)
		resource := &Resource{
			kubeRef: kuberef,
			cuePath: resourcesIter.Value().Path(),
			healthy: false,
		}
		myResources = append(myResources, resource)
	}

	waitFor := []*Component{}
	components, errComponents := componentsCue.Fields()
	if errComponents != nil {
		fmt.Print("Error getting components: " + errComponents.Error())
		panic("Error getting components: " + errComponents.Error())
	}
	for components.Next() {
		childCueComp := components.Value()
		childComponent := _makeComponentTree(childCueComp, parent)
		waitFor = append(waitFor, childComponent)
	}

	return &Component{
		cuePath:   v.Path(),
		waitFor:   waitFor,
		resources: myResources,
	}
}
