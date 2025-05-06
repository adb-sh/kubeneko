package main

import (
	"fmt"
	"log"

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
	// insts := load.Instances([]string{"."}, &load.Config{
	insts := load.Instances([]string{"./examples/simpleapp"}, &load.Config{
		Registry: reg,
	})

	// The current directory just has one file without any build tags,
	// and that file belongs to the example package, so we get a single
	// instance as a result.
	v := ctx.BuildInstance(insts[0])
	if err := v.Err(); err != nil {
		log.Fatal(err)
	}

	// build dependencyMap
	// Example:
	// { "field1": {"field2"}, "field2": {"field3"} }
	// field1 changes field2 and field2 changes field3
	dependencyMap := makeDepsMap(v)
	fmt.Println()
	fmt.Println("Dependency Map:")
	printDependencyMap(dependencyMap)
	fmt.Println()

	// extends the map so that all nested paths are included in the parent as a dependency.
	// Example:
	// { "field1": {"field2"}, "field2": {"field3"} }
	// field1 changes field2 and field2 changes field3
	// therfore field1 changes field3 too
	// the result of this function is:
	// { "field1": {"field2", "field3"}, "field2": {"field3"} }
	// field1 changes field2 and field3. field2 changes field3
	extendedMap := extendDependencyMap(dependencyMap)
	fmt.Println("Extended Dependency Map:")
	printDependencyMap(extendedMap)
	fmt.Println()

	// export dependency map to dot file (Graphviz)
	exportToDot(dependencyMap)
}
