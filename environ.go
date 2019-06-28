package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	BatchSize = 10
)

func main() {
	var (
		decrypt = flag.Bool("verbose", false, "")
		export  = flag.Bool("export", false, "Print into stdout export VAR=VALUE")
		_       = flag.Bool("ignore", false, "Don't fail if error retrieving parameter")
	)
	flag.Parse()
	args := flag.Args()
	scope := os.Environ()

	ssm := ssm.New(session.Must(awsSession()))

	ssmVariables := make(map[string]string)

	variableNames := []string{}
	for _, pair := range scope {
		if strings.Contains(pair, "ssm:") {
			localName, ssmVarName := unPack(pair)
			ssmVarName = ssmVarName[len("ssm:"):]
			ssmVariables[localName] = ssmVarName
			variableNames = append(variableNames, ssmVarName)
		}
	}

	// fmt.Printf("scope %#v\n", variableNames)

	mapping, err := smmFetcher(ssm, variableNames, *decrypt, nil)
	// fmt.Printf("mapping %#v\n", mapping)
	fmt.Printf("err %#v\n", err)
	switch {
	case *export:
		for name, val := range ssmVariables {
			fmt.Printf("export %s=%s\n", name, mapping[val])
		}
		return
	default:
		for name, val := range ssmVariables {
			os.Setenv(name, mapping[val])
		}
		return
	}
	exitIf(syscall.Exec(args[1], args[1:], os.Environ()))
}

func awsSession() (*session.Session, error) {
	sess := session.Must(session.NewSession())
	return sess, nil
}

func smmFetcher(session *ssm.SSM,
	names []string,
	decrypt bool,
	mapping map[string]string) (map[string]string, error) {
	if mapping == nil {
		mapping = make(map[string]string)
	}

	input := ssm.GetParametersInput{
		WithDecryption: aws.Bool(decrypt),
	}
	head := names
	tail := names[:0]
	if len(head) > BatchSize {
		head = names[:BatchSize]
		tail = names[BatchSize:]
	}
	for _, paramName := range head {
		input.Names = append(input.Names, aws.String(paramName))
	}

	resp, err := session.GetParameters(&input)
	if err != nil {
		return mapping, err
	}

	if len(resp.InvalidParameters) > 0 {
		return mapping, fmt.Errorf("InvalidParameters:=%v", resp.InvalidParameters)
	}

	for _, p := range resp.Parameters {
		mapping[*p.Name] = *p.Value
	}
	if len(tail) > 0 {
		return smmFetcher(session, tail, decrypt, mapping)
	}
	return mapping, nil
}

func unPack(v string) (key, val string) {
	parts := strings.Split(v, "=")
	return parts[0], parts[1]
}

func exitIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ssm-env: %v\n", err)
		os.Exit(1)
	}
}
