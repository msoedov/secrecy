package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/msoedov/environ/vars"
)

func main() {
	var (
		_ = flag.Bool("verbose", false, "")
		_ = flag.Bool("export", false, "Print into stdout export VAR=VALUE")
		_ = flag.Bool("ignore", false, "Don't fail if error retrieving parameter")
	)
	flag.Parse()
	mapping, err := vars.InstrumentEnv(true)
	exitIf(err)
	fmt.Printf("Instrumented %#v\n", mapping)
	// exitIf(syscall.Exec(args[1], args[1:], os.Environ()))
}

func exitIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ssm-env: %v\n", err)
		os.Exit(1)
	}
}
