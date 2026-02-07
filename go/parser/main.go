package parser

import (
	"fmt"
	"log"
	"strconv"

	"cuelang.org/go/cue"
	"github.com/adb-sh/kubeneko/core"
)

type Ref struct {
	From     cue.Value
	To       cue.Value
	Resource *Resource
}

type Dep struct {
	Refs  []Ref
	Value cue.Path
}

type Resource struct {
	Manifest core.KubeRef
	Value    cue.Value
	Deps     []Dep
}

func Parse(val cue.Value) []Resource {
	v := val.LookupPath(cue.ParsePath("config"))
	if !v.Exists() {
		log.Fatal("component not found")
	}

	var config core.Config
	if err := v.Decode(&config); err != nil {
		log.Fatal("decode error:", err)
	}

	fmt.Println("config:", config)

	resources := []Resource{}
	walkConfig(config, v, &resources)
	// fmt.Println(resources)

	resolveRefs(resources)

	return resources
}

func walk(val cue.Value, deps *[]Dep) {
	// path := val.Path().String()
	iter, _ := val.Fields(cue.Attributes(true))
	for iter.Next() {
		// label := iter.Selector()
		child := iter.Value()
		// kind := child.Kind()

		refs := []Ref{}
		walkExpr(child, &refs)
		dep := Dep{
			Value: child.Path(),
			Refs:  refs,
		}
		*deps = append(*deps, dep)

		walk(child, deps)
	}
}

func walkExpr(val cue.Value, refs *[]Ref) {
	if ref, path := val.ReferencePath(); len(path.Selectors()) > 0 {
		link := Ref{
			From: ref.LookupPath(path),
			To:   val,
		}
		*refs = append(*refs, link)
		return
	}

	op, args := val.Expr()
	if op == cue.NoOp {
		return
	}

	for _, arg := range args {
		walkExpr(arg, refs)
	}
}

func walkConfig(conf core.Config, inst cue.Value, resources *[]Resource) {
	for i, el := range conf.Resources {
		val := inst.LookupPath(cue.ParsePath("resources[" + strconv.Itoa(i) + "]"))
		deps := []Dep{}
		walk(val, &deps)
		res := Resource{
			Manifest: el,
			Value:    val,
			Deps:     deps,
		}
		*resources = append(*resources, res)
	}
	for i, el := range conf.Components {
		val := inst.LookupPath(cue.ParsePath("components[" + strconv.Itoa(i) + "].out"))
		walkConfig(el.Out, val, resources)
	}
}

func resolveRefs(resources []Resource) {
	for _, res := range resources {
		fmt.Println("\n# resource", res.Manifest.Metadata.Name, res.Manifest.Kind, res.Value.Path())
		for _, dep := range res.Deps {
			for _, ref := range dep.Refs {
				fmt.Println("ref", ref.To.Path(), "-->", ref.From.Path())
				for _, candidate := range resources {
					fmt.Println("candidate", candidate.Value.Path())
					if candidate.Value == res.Value {
						continue // skip self
					}

					// TODO: check if `res.Value` is a child of `candidate.Value`

					if searchValueTree(candidate.Value, ref.From) {
						// ref.Resource = *candidate
						fmt.Println("dep", candidate.Manifest.Metadata.Name, candidate.Manifest.Kind)
					}
				}
			}
		}
	}
}

func searchValueTree(root cue.Value, target cue.Value) bool {
	it, _ := root.Fields()
	for it.Next() {
		v := it.Value()
		if v.Equals(target) {
			fmt.Println(v.Path(), target.Path())
			return true
		}
		if searchValueTree(v, target) {
			return true
		}
	}
	return false
}
