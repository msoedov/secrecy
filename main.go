package main

import (
	"flag"
	"os"

	_ "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/msoedov/environ/secrecy"
	log "github.com/sirupsen/logrus"
)

var (
	verbose = flag.Bool("verbose", true, "")
	_       = flag.Bool("export", false, "Print into stdout export VAR=VALUE")
	_       = flag.Bool("ignore", false, "Don't fail if error retrieving parameter")
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stderr)

	flag.Parse()
	mapping, err := secrecy.InstrumentEnv(true)
	exitIf(err)
	if *verbose {
		log.Infof("Instrumented %#v\n", secrecy.MaskValues(mapping))
	}
}

func exitIf(err error) {
	if err != nil {
		log.Fatalf("secrecy-env: %v\n", err)
		os.Exit(1)
	}
}
