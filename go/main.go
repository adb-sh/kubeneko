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

	for path, deps := range makeDepsMap(v) {
		fmt.Println(path)
		fmt.Println("  depends on: ")
		for _, dep := range deps {
			fmt.Print("    ")
			fmt.Print(dep.String())
			fmt.Println()
		}
	}

	// // Lookup specific values
	// for path, deps := range makeDepsMap(v.LookupPath(cue.ParsePath("a.b.c"))) {
	// 	fmt.Println(path)
	// 	fmt.Println("  depends on: ")
	// 	for _, dep := range deps {
	// 		fmt.Print("    ")
	// 		fmt.Print(dep.String())
	// 		fmt.Println()
	// 	}
	// }
}

func makeDepsMap(val cue.Value) map[string][]cue.Path {
	dependencyMap := map[string][]cue.Path{}
	child := val.Value()
	walkExpr(child, &dependencyMap)
	return dependencyMap
}

func walkExpr(child cue.Value, dependencyMap *map[string][]cue.Path) {
	if op, args := child.Expr(); op != cue.NoOp && len(args) > 0 {
		pathsByOp(op, child, args, dependencyMap)
	} else {
		switch child.Kind() {
		case cue.StructKind:
			iter, _ := child.Fields()
			for iter.Next() {
				child := iter.Value()
				walkExpr(child, dependencyMap)
			}
		case cue.ListKind:
			iterList, _ := child.List()
			for iterList.Next() {
				child := iterList.Value()
				walkExpr(child, dependencyMap)
			}

		case cue.StringKind, cue.FloatKind, cue.NumberKind, cue.BoolKind, cue.BytesKind, cue.IntKind:
			// do nothing (is literal)
		case cue.BottomKind, cue.NullKind, cue.TopKind:
			fmt.Println("Kind with no op not implemented")
			fmt.Println("  Kind:", child.Kind())
			fmt.Println("  op:", op)
			fmt.Println("  args:", args)
			fmt.Println()
		default:
			fmt.Println("unkown kind with no op")
			fmt.Println("  Kind:", child.Kind())
			fmt.Println("  op:", op)
			fmt.Println("  args:", args)
			fmt.Println()
		}
	}
}

func pathsByOp(op cue.Op, node cue.Value, args []cue.Value, dependencyMap *map[string][]cue.Path) {
	switch op {
	case cue.SelectorOp:
		_, ref := node.ReferencePath()
		nodePath := node.Path()
		(*dependencyMap)[nodePath.String()] = append((*dependencyMap)[nodePath.String()], ref)
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
		for _, arg := range args {
			argOp, argArgs := arg.Expr()
			if argOp != cue.NoOp {
				pathsByOp(argOp, arg, argArgs, dependencyMap)
			}
		}
	default:
		fmt.Printf("  unkown op: %s", op)
		fmt.Println()
	}
}
