package main

import (
	"fmt"
	"log"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/load"
	"cuelang.org/go/mod/modconfig"
)

func main() {
	ctx := cuecontext.New()

	// Create a registry client. Passing a nil config
	// will give us client that behaves like the cue command.
	reg, err := modconfig.NewRegistry(nil)
	if err != nil {
		log.Fatal(err)
	}

	// Load the package from the current directory.
	// We don't need to specify a Config in this example.
	insts := load.Instances([]string{"."}, &load.Config{
		Registry: reg,
	})

	// The current directory just has one file without any build tags,
	// and that file belongs to the example package, so we get a single
	// instance as a result.
	v := ctx.BuildInstance(insts[0])
	if err := v.Err(); err != nil {
		log.Fatal(err)
	}

	walk("", v)
}

func walk(path string, val cue.Value) {
	iter, _ := val.Fields()
	for iter.Next() {
		label := iter.Selector()
		child := iter.Value()
		// kind := child.Kind()

		deps := []cue.Path{}
		walkExpr(child, &deps)
		fmt.Printf("%s.%s depends on: ", path, label)
		fmt.Println(deps)

		walk(path+"."+label.String(), child)
	}
}

func walkExpr(val cue.Value, deps *[]cue.Path) {
	if _, ref := val.ReferencePath(); len(ref.Selectors()) > 0 {
		*deps = append(*deps, ref)
		return
	}

	op, args := val.Expr()
	if op == cue.NoOp {
		return
	}

	for _, arg := range args {
		walkExpr(arg, deps)
	}
}
