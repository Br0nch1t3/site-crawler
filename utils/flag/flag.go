package utilsflag

import (
	"flag"
	"fmt"
	"os"
)

// Function called when a flag or argument is encoutered
type LookupFn func(string) error

type varFn[T any] func(p *T, name string, value bool, usage string)

func NewVar[T any](fn varFn[T], p *T, name string, shorthand string, value bool, usage string) [2]string {
	fn(p, name, value, usage)
	fn(p, shorthand, value, "")

	return [2]string{name, shorthand}
}

func NewFunc(name string, shorthand string, usage string, fn LookupFn) [2]string {
	flag.Func(name, usage, fn)
	flag.Func(shorthand, "", fn)

	return [2]string{name, shorthand}
}

func NewBoolFunc(name string, shorthand string, usage string, fn LookupFn) [2]string {
	flag.BoolFunc(name, usage, fn)
	flag.BoolFunc(shorthand, "", fn)

	return [2]string{name, shorthand}
}

// Parses command line arguments, including positional ones
//
/* Parse("[OPTIONS] [ARG1] [ARG2]", [][2]string{NewVar(...)}, []LookupFn{...}) */
func Parse(commandLineUsage string, flags [][2]string, positionals []LookupFn) {
	flags = append(flags, NewBoolFunc("help", "h", "Prints this help message and exit", func(s string) error {
		flag.Usage()
		os.Exit(0)
		return nil
	}))
	flag.Usage = usageBuilder(commandLineUsage, flags)
	flag.Parse()
	if len(positionals) == flag.NArg() {
		for i, pos := range positionals {
			if err := pos(flag.Arg(i)); err != nil {
				fmt.Fprintln(os.Stderr, err)
				flag.Usage()
				os.Exit(1)
			}
		}
		return
	}
	fmt.Fprintln(os.Stderr, "invalid number of arguments")
	flag.Usage()
	os.Exit(1)
}

func usageBuilder(commandLineUsage string, flags [][2]string) func() {
	usage := fmt.Sprintf("Usage of %s:\n%s %s", flag.CommandLine.Name(), flag.CommandLine.Name(), commandLineUsage)

	for _, names := range flags {
		f := flag.Lookup(names[0])
		_, flagUsage := flag.UnquoteUsage(f)
		usage = fmt.Sprintf("%s\n  -%s, --%s\n\t%s", usage, names[1], f.Name, flagUsage)
		if len(f.DefValue) > 0 {
			usage = fmt.Sprintf("%s (default: %s)", usage, f.DefValue)
		}
	}

	return func() {
		fmt.Fprintln(os.Stderr, usage)
	}
}
