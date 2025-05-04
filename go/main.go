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
		fmt.Printf("Field %s.%s (%s)", path, label, child.Kind())
		fmt.Println()

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
				fmt.Println("  depends on ")
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
			fmt.Printf("  unkown kind")
			fmt.Println()
		}
	}
}

func pathsByOp(op cue.Op, node cue.Value, args []cue.Value) [][]cue.Selector {
	var paths [][]cue.Selector
	switch op {
	case cue.NoOp:
		fmt.Println("    TODO: cue.NoOp")

	case cue.AndOp:
		fmt.Println("    TODO: cue.AndOp")
	case cue.OrOp:
		fmt.Println("    TODO: cue.OrOp")

	case cue.SelectorOp:
		_, ref := node.ReferencePath()
		paths = append(paths, ref.Selectors())

	case cue.IndexOp:
		fmt.Println("    TODO: cue.IndexOp")
	case cue.SliceOp:
		fmt.Println("    TODO: cue.SliceOp")
	case cue.CallOp:
		fmt.Println("    TODO: cue.CallOp")

	case cue.BooleanAndOp:
		fmt.Println("    TODO: cue.BooleanAndOp")
	case cue.BooleanOrOp:
		fmt.Println("    TODO: cue.BooleanOrOp")

	case cue.EqualOp:
		fmt.Println("    TODO: cue.EqualOp")
	case cue.NotOp:
		fmt.Println("    TODO: cue.NotOp")
	case cue.NotEqualOp:
		fmt.Println("    TODO: cue.NotEqualOp")
	case cue.LessThanOp:
		fmt.Println("    TODO: cue.LessThanOp")
	case cue.LessThanEqualOp:
		fmt.Println("    TODO: cue.LessThanEqualOp")
	case cue.GreaterThanOp:
		fmt.Println("    TODO: cue.GreaterThanOp")
	case cue.GreaterThanEqualOp:
		fmt.Println("    TODO: cue.GreaterThanEqualOp")

	case cue.RegexMatchOp:
		fmt.Println("    TODO: cue.RegexMatchOp")
	case cue.NotRegexMatchOp:
		fmt.Println("    TODO: cue.NotRegexMatchOp")

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
