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

	// Lookup the 'config' field and print it out
	config := v.LookupPath(cue.ParsePath(""))
	fmt.Println(config)

	walk("", v)
}

func walk(path string, val cue.Value) {
	iter, _ := val.Fields()
	for iter.Next() {
		// label := iter.Label()
		child := iter.Value()
		fmt.Println("path", child.Path())
		fmt.Println("kind", child.Kind())
		op, args := child.Expr()
		fmt.Println(op)
		fmt.Printf("%#v ", args)
		fmt.Println()
		fmt.Printf("%#v ", child.Syntax())
		fmt.Println()
		_, ref := child.ReferencePath()
		fmt.Println("ref", ref)

		fmt.Println()

		// if op, args := child.Expr(); op != cue.NoOp && len(args) > 0 {
		// 	fmt.Printf("%s.%s depends on: ", path, label)
		// 	fmt.Print(args)
		// 	// for _, arg := range args {
		// 	// 	defer func() {
		// 	// 		syn := arg.Syntax(
		// 	// 			cue.Docs(false),
		// 	// 			cue.Definitions(true),
		// 	// 		)
		// 	// 		fmt.Print(syn.End().String())
		// 	// 	}()
		// 	// panic("yikes")
		// 	// switch expr := syn.(type) {
		// 	// case ast.Expr:
		// 	// 	fmt.Printf("%#v ", expr)
		// 	// default:
		// 	// 	fmt.Printf("<non-expr: %T> ", syn)
		// 	// 	continue
		// 	// }
		// 	// }
		// 	fmt.Println()
		// }

		// fmt.Println(path + "." + label)
		// walk(path+"."+label, child)
	}
}
