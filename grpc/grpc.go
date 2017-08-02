package grpc

import (
	"fmt"
	"strings"
)

type grpcCmd struct {
	cmd       []string
	inputDir  string
	outputDir string
	proto     string
}

func NewGrpcCmd(in, out, proto string, langs ...string) []string {
	c := grpcCmd{inputDir: in, outputDir: out, proto: proto}
	for _, lang := range langs {
		switch strings.TrimSpace(lang) {
		case "go":
			c.Go()
		case "python":
			c.Python()
		case "":
			return nil
		default:
			fmt.Printf("No support language: %s\n", lang)
		}
	}

	return c.cmd
}

func (c *grpcCmd) Python() {
	command := "python -m grpc_tools.protoc -I=" + c.inputDir + " --python_out=" + c.outputDir + " --grpc_python_out=" + c.outputDir + " " + c.proto
	c.cmd = append(c.cmd, command)
}

func (c *grpcCmd) Go() {
	command := "protoc -I=" + c.inputDir + " --go_out=plugins=grpc:" + c.outputDir + " " + c.proto
	c.cmd = append(c.cmd, command)

	command = "cp -r " + c.outputDir + " /go/src"
	c.cmd = append(c.cmd, command)
}
