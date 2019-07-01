package vars

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	BatchSize = 10
)

type Env map[string]string

func Do(export bool) Env {
	ssm := ssm.New(session.Must(awsSession()))
	variableNames, ssmVariables := defaultVariableDiscovery()
	mapping, err := SmmFetcher(ssm, variableNames, false, nil)
	if err != nil {
		fmt.Printf("Fail %#v\n", err)
		return nil
	}
	switch {
	case export:
		for name, val := range ssmVariables {
			fmt.Printf("export %s=%s\n", name, mapping[val])
		}
		return mapping
	default:
		for name, val := range ssmVariables {
			os.Setenv(name, mapping[val])
		}
		return mapping
	}
	return mapping
}

func awsSession() (*session.Session, error) {
	sess := session.Must(session.NewSession())
	return sess, nil
}

func SmmFetcher(session *ssm.SSM,
	names []string,
	shouldDecrypt bool,
	variableMapping Env) (Env, error) {
	if variableMapping == nil {
		variableMapping = make(Env)
	}

	input := ssm.GetParametersInput{
		WithDecryption: aws.Bool(shouldDecrypt),
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
		return variableMapping, err
	}

	if len(resp.InvalidParameters) > 0 {
		return variableMapping, fmt.Errorf("InvalidParameters:=%v", resp.InvalidParameters)
	}

	for _, p := range resp.Parameters {
		variableMapping[*p.Name] = *p.Value
	}
	if len(tail) > 0 {
		return SmmFetcher(session, tail, shouldDecrypt, variableMapping)
	}
	return variableMapping, nil
}

func unPack(v string) (key, val string) {
	parts := strings.Split(v, "=")
	return parts[0], parts[1]
}

func defaultVariableDiscovery() ([]string, Env) {
	ssmVariables := make(Env)
	scope := os.Environ()
	variableNames := []string{}
	for _, pair := range scope {
		if strings.Contains(pair, "ssm:") {
			localName, ssmVarName := unPack(pair)
			ssmVarName = ssmVarName[len("ssm:"):]
			ssmVariables[localName] = ssmVarName
			variableNames = append(variableNames, ssmVarName)
		}
	}
	return variableNames, ssmVariables
}
