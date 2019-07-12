# Environ

Ultralight AWS Parameter Store variables instrumentation

## Usage

```shell
export FOO=ssm:SSM_PARAM
`secrecy`
echo $FOO
SSM_VALUE:)
```

How this works:

```shell
export FOO=ssm:SSM_PARAM
environ
# stdout
export FOO=SSM_VALUE:)
```

## Install

```shell
go get -u github.com/msoedov/secrecy
```

## Configuration

```shell
# Optional as per aws/aws-sdk-go#configuring-credentials
AWS_SDK_LOAD_CONFIG=1
```
