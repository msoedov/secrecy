# Secrecy

Ultralight AWS Parameter Store variables instrumentation on AWS instances

### Features

Written in simple Go

No installation necessary - binary is provided

Intuitive and easy to use

Efficient batching

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
secrecy
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

### Docker instrumentation

TBD
