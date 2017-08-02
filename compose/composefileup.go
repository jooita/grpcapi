package compose

import (
	"context"
	"log"

	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"
)

func ComposeFileUp(file string) error {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{file},
			ProjectName:  "keti-grpc-api",
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
		return err
	}

	err = project.Up(context.Background(), options.Up{})

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
