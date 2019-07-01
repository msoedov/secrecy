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
	// AWS_SDK_LOAD_CONFIG=1

	// args := flag.Args()
	// scope := os.Environ()

	// ssm := ssm.New(session.Must(awsSession()))

	// ssmVariables := make(map[string]string)

	// variableNames := []string{}
	// for _, pair := range scope {
	// 	if strings.Contains(pair, "ssm:") {
	// 		localName, ssmVarName := unPack(pair)
	// 		ssmVarName = ssmVarName[len("ssm:"):]
	// 		ssmVariables[localName] = ssmVarName
	// 		variableNames = append(variableNames, ssmVarName)
	// 	}
	// }

	// fmt.Printf("scope %#v\n", variableNames)

	mapping := vars.Do(true)
	fmt.Printf("mapping %#v\n", mapping)
	// exitIf(syscall.Exec(args[1], args[1:], os.Environ()))
}

func exitIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ssm-env: %v\n", err)
		os.Exit(1)
	}
}
