package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"cuelang.org/go/cue"
)

func makeDepsMap(rootVal cue.Value) map[string][]string {
	dependencyMap := map[string][]string{}
	walkExpr(rootVal, rootVal, &dependencyMap)
	return reverseMap(dependencyMap)
}

func walkExpr(child cue.Value, rootVal cue.Value, dependencyMap *map[string][]string) {
	if op, args := child.Expr(); op != cue.NoOp && len(args) > 0 {
		// if has operator, walk operator
		walkOperator(op, child, args, rootVal, dependencyMap)
	} else {
		switch child.Kind() {
		case cue.StructKind:
			// if is struct, walk fields
			iter, _ := child.Fields()
			for iter.Next() {
				child := iter.Value()
				walkExpr(child, rootVal, dependencyMap)
			}
		case cue.ListKind:
			// if is list, walk elements
			iterList, _ := child.List()
			for iterList.Next() {
				child := iterList.Value()
				walkExpr(child, rootVal, dependencyMap)
			}

		case cue.StringKind, cue.FloatKind, cue.NumberKind, cue.BoolKind, cue.BytesKind, cue.IntKind:
			// do nothing (is literal)
		case cue.BottomKind:
			// do nothing (is type)
		case cue.NullKind:
			// do nothing (is null)

		case cue.TopKind:
			// i don't know what to do with this
			fmt.Println("this Kind without operator is not implemented")
			fmt.Println("  Kind:", child.Kind())
			fmt.Println("  operator:", op)
			fmt.Println("  args:", args)
			fmt.Println()
		default:
			fmt.Println("unkown kind with no operator")
			fmt.Println("  Kind:", child.Kind())
			fmt.Println("  operator:", op)
			fmt.Println("  args:", args)
			fmt.Println()
			panic("unkown kind with no operator")
		}
	}
}

func walkOperator(op cue.Op, node cue.Value, args []cue.Value, rootVal cue.Value, dependencyMap *map[string][]string) {
	switch op {
	case cue.NoOp:
		walkExpr(node, rootVal, dependencyMap)
	case cue.SelectorOp:
		_, refPath := node.ReferencePath()
		nodePath := node.Path()
		refSelectors := refPath.Selectors()
		if len(refSelectors) == 0 {
			// if no selector, just add the path to the dependency map
			(*dependencyMap)[node.Path().String()] = append((*dependencyMap)[node.Path().String()], cue.Dereference(node).Path().String())
			return
		}
		refLabel := refSelectors[len(refSelectors)-1]
		if strings.HasPrefix(refLabel.String(), "#") {
			// if the last selector has prefix #, treat it as a function
			child := rootVal.LookupPath(refPath)
			nestedDependencyMap := map[string][]string{}
			walkExpr(child, rootVal, &nestedDependencyMap)
			for rootPathString, deps := range nestedDependencyMap {
				rootPathString := cue.ParsePath(rootPathString).String()
				firstSelector := cue.ParsePath(rootPathString).Selectors()[0].String()
				nestedPathString := strings.Replace(rootPathString, firstSelector, nodePath.String(), 1)
				for _, depKey := range deps {
					nestedDepKey := strings.Replace(depKey, firstSelector, nodePath.String(), 1)
					(*dependencyMap)[nestedPathString] = append((*dependencyMap)[nestedPathString], nestedDepKey)
				}
			}
		} else if strings.HasPrefix(refLabel.String(), "_") {
			// if the last selector has prefix _, treat it as a hidden reference
			child := rootVal.LookupPath(refPath)
			walkExpr(child, rootVal, dependencyMap)
			(*dependencyMap)[nodePath.String()] = append((*dependencyMap)[nodePath.String()], refPath.String())
		} else {
			// if the last selector has no prefix, treat it as a normal reference
			(*dependencyMap)[nodePath.String()] = append((*dependencyMap)[nodePath.String()], refPath.String())
		}
	case
		cue.AndOp,
		cue.OrOp,
		cue.IndexOp,
		cue.SliceOp,
		cue.CallOp,
		cue.BooleanAndOp,
		cue.BooleanOrOp,
		cue.EqualOp,
		cue.NotOp,
		cue.NotEqualOp,
		cue.LessThanOp,
		cue.LessThanEqualOp,
		cue.GreaterThanOp,
		cue.GreaterThanEqualOp,
		cue.RegexMatchOp,
		cue.NotRegexMatchOp,
		cue.AddOp,
		cue.SubtractOp,
		cue.MultiplyOp,
		cue.FloatQuotientOp,
		cue.InterpolationOp,
		cue.IntQuotientOp,
		cue.IntRemainderOp,
		cue.IntDivideOp,
		cue.IntModuloOp:
		// if the operator is one of these, walk the args
		for _, arg := range args {
			argOp, argArgs := arg.Expr()
			walkOperator(argOp, arg, argArgs, rootVal, dependencyMap)
		}
	default:
		fmt.Printf("  unkown op: %s", op)
		fmt.Println()
	}
}

func extendDependencyMap(dependencyMap map[string][]string) map[string][]string {
	// Step 1: Initialize base dependency map
	expanded := make(map[string]map[string]bool)
	for k, deps := range dependencyMap {
		expanded[k] = makeSet(deps)
	}

	// Step 2: Field mapping propagation (e.g., myStruct → #SomeStruct)
	for structUsage, structDeps := range dependencyMap {
		for _, dep := range structDeps {
			if !strings.HasPrefix(dep, "#") {
				continue
			}
			for depField, fieldDeps := range dependencyMap {
				if strings.HasPrefix(depField, dep+".") {
					suffix := strings.TrimPrefix(depField, dep)
					newKey := structUsage + suffix
					if expanded[newKey] == nil {
						expanded[newKey] = map[string]bool{}
					}
					for _, d := range fieldDeps {
						expanded[newKey][d] = true
					}
				}
			}
		}
	}

	// Step 3: Recursive dependency resolution
	changed := true
	for changed {
		changed = false
		for key, deps := range expanded {
			newDeps := make(map[string]bool)
			for dep := range deps {
				if transitive, ok := expanded[dep]; ok {
					for t := range transitive {
						if !deps[t] {
							newDeps[t] = true
							changed = true
						}
					}
				}
			}
			for d := range newDeps {
				expanded[key][d] = true
			}
		}
	}

	// Step 4: Add parent chain dependencies
	for key := range expanded {
		if !strings.Contains(key, ".") {
			continue
		}
		segments := strings.Split(key, ".")
		for i := len(segments) - 1; i > 0; i-- {
			parent := strings.Join(segments[:i], ".")
			for dep := range expanded[parent] {
				if !expanded[key][dep] {
					expanded[key][dep] = true
				}
			}
		}
	}

	// Step 5: Convert back to map[string][]string
	final := make(map[string][]string)
	for k, v := range expanded {
		for dep := range v {
			final[k] = append(final[k], dep)
		}
	}

	return final
}

func makeSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

func exportToDot(dependencyMap map[string][]string) {
	var sb strings.Builder
	sb.WriteString("digraph G {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape=box];\n")

	for from, deps := range dependencyMap {
		if from == "#Component" || from == "#Config" {
			continue
		}
		for _, to := range deps {
			sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", from, to))
		}
	}

	sb.WriteString("}")
	dot := sb.String()

	fileError := os.WriteFile("deps.dot", []byte(dot), 0644)
	if fileError != nil {
		log.Fatal(fileError)
	}
	fmt.Println("Graph written to deps.dot")
}

func reverseMap(original map[string][]string) map[string][]string {
	reversed := make(map[string][]string)

	for key, values := range original {
		for _, value := range values {
			reversed[value] = append(reversed[value], key)
		}
	}

	return reversed
}

func printDependencyMap(dependencyMap map[string][]string) {
	for path, deps := range dependencyMap {
		fmt.Println(path)
		fmt.Println("  depends on: ")
		for _, dep := range deps {
			fmt.Print("  - ")
			fmt.Print(dep)
			fmt.Println()
		}
	}
}
