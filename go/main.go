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
	// config := v.LookupPath(cue.ParsePath(""))
	// fmt.Println(config)

	walk("", v)
}

func walk(path string, val cue.Value) {
	iter, _ := val.Fields()
	for iter.Next() {
		label := iter.Selector()
		child := iter.Value()
		fmt.Printf("Field %s.%s ", path, label)
		fmt.Println()

		kind := child.Kind()
		switch child.Kind() {
		case cue.StructKind:
			walk(path+"."+label.String(), child)
		case cue.ListKind:
			iterList, _ := child.List()
			for iterList.Next() {
				indexLabel := iterList.Selector()
				walk(path+"."+label.String()+"."+indexLabel.String(), iterList.Value())
			}

		case cue.StringKind, cue.FloatKind, cue.NumberKind, cue.BoolKind, cue.BytesKind, cue.IntKind:
			if op, args := child.Expr(); op != cue.NoOp && len(args) > 0 {
				fmt.Print("  depends on ")
				fmt.Println(op.String())
				paths := pathsByOp(op, child, args)
				for _, path := range paths {
					fmt.Print("    ")
					fmt.Print(path)
					fmt.Println()
				}
			} else {
				value := child.Value()
				fmt.Println("  literal", value)

			}
		case cue.BottomKind:
			fmt.Println("  bottom")
		case cue.NullKind:
			fmt.Println("  null")
		case cue.TopKind:
			fmt.Println("  top")
		default:
			fmt.Printf("  unkown kind: %s", kind)
			fmt.Println()
		}
	}
}

func pathsByOp(op cue.Op, node cue.Value, args []cue.Value) [][]cue.Selector {
	var paths [][]cue.Selector
	switch op {
	case cue.NoOp:
		// fmt.Println("cue.NoOp")
		// TODO:

	case cue.AndOp:
		// fmt.Println("cue.AndOp")
		// TODO:
	case cue.OrOp:
		// fmt.Println("cue.OrOp")
		// TODO:

	case cue.SelectorOp:
		_, ref := node.ReferencePath()
		paths = append(paths, ref.Selectors())

	case cue.IndexOp:
		// fmt.Println("cue.IndexOp")
		// TODO:
	case cue.SliceOp:
		// fmt.Println("cue.SliceOp")
		// TODO:
	case cue.CallOp:
		// fmt.Println("cue.CallOp")
		// TODO:

	case cue.BooleanAndOp:
		// fmt.Println("cue.BooleanAndOp")
		// TODO:
	case cue.BooleanOrOp:
		// fmt.Println("cue.BooleanOrOp")
		// TODO:

	case cue.EqualOp:
		// fmt.Println("cue.EqualOp")
		// TODO:
	case cue.NotOp:
		// fmt.Println("cue.NotOp")
		// TODO:
	case cue.NotEqualOp:
		// fmt.Println("cue.NotEqualOp")
		// TODO:
	case cue.LessThanOp:
		// fmt.Println("cue.LessThanOp")
		// TODO:
	case cue.LessThanEqualOp:
		// fmt.Println("cue.LessThanEqualOp")
		// TODO:
	case cue.GreaterThanOp:
		// fmt.Println("cue.GreaterThanOp")
		// TODO:
	case cue.GreaterThanEqualOp:
		// fmt.Println("cue.GreaterThanEqualOp")
		// TODO:

	case cue.RegexMatchOp:
		// fmt.Println("cue.RegexMatchOp")
		// TODO:
	case cue.NotRegexMatchOp:
		// fmt.Println("cue.NotRegexMatchOp")
		// TODO:

	case cue.AddOp:
		for _, arg := range args {
			argOp, argArgs := arg.Expr()
			if argOp != cue.NoOp {
				argPaths := pathsByOp(argOp, arg, argArgs)
				paths = append(paths, argPaths...)
			}
		}

	case cue.SubtractOp:
		// fmt.Println("cue.SubtractOp")
	case cue.MultiplyOp:
		// fmt.Println("cue.MultiplyOp")
	case cue.FloatQuotientOp:
		// fmt.Println("cue.FloatQuotientOp")
	case cue.IntQuotientOp:
		// fmt.Println("cue.IntQuotientOp")
	case cue.IntRemainderOp:
		// fmt.Println("cue.IntRemainderOp")
	case cue.IntDivideOp:
		// fmt.Println("cue.IntDivideOp")
	case cue.IntModuloOp:
		// fmt.Println("cue.IntModuloOp")

	case cue.InterpolationOp:
		// fmt.Println("cue.InterpolationOp")
	default:
		fmt.Printf("  unkown op: %s", op)
		fmt.Println()
	}
	return paths
}
