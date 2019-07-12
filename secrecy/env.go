package secrecy

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

const (
	BatchSize = 10
)

type Env map[string]string

func InstrumentEnv(export bool) (Env, error) {
	ssm := ssm.New(session.Must(awsSession()))
	variableNames, ssmVariables := defaultVariableDiscovery()
	if len(variableNames) == 0 {
		return *new(Env), nil
	}
	mapping, err := SmmFetcher(ssm, variableNames, true, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load param")
	}
	switch {
	case export:
		for name, val := range ssmVariables {
			fmt.Printf("export %s=%s\n", name, mapping[val])
		}
		return mapping, nil
	default:
		for name, val := range ssmVariables {
			os.Setenv(name, mapping[val])
		}
		return mapping, nil
	}
	return mapping, nil
}

func awsSession() (*session.Session, error) {
	sess := session.Must(session.NewSession())
	return sess, nil
}

// SmmFetcher fetch ssm params from a given mapping
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
		return variableMapping, errors.Wrapf(err, "Input=%v", input)
	}

	if len(resp.InvalidParameters) > 0 {
		return variableMapping, fmt.Errorf("InvalidParameters:=%v and input=%v", resp.InvalidParameters, input)
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

func MaskValues(secrets Env) Env {
	secured := make(Env)
	for name, secret := range secrets {
		secured[name] = maskText(secret)
	}
	return secured
}

func maskText(secret string) string {
	switch {
	case len(secret) < 3:
		return "***"
	default:
		return secret[:1] + "******" + secret[len(secret)-2:]
	}
}
