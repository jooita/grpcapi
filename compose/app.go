package compose

import (
	"github.com/docker/libcompose/cli/logger"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/lytics/logrus"
)

// ProjectFactory is a struct that holds the app.ProjectFactory implementation.
type ProjectFactory struct {
}

// Populate updates the specified docker context based on command line arguments and subcommands.
func Populate(context *ctx.Context) {
	opts := client.Options{}
	clientFactory, err := client.NewDefaultFactory(opts)
	if err != nil {
		logrus.Fatalf("Failed to construct Docker client: %v", err)
	}

	context.ClientFactory = clientFactory
}

// Create implements ProjectFactory.Create using docker client.
func (p *ProjectFactory) Create() (project.APIProject, error) {
	context := &ctx.Context{}
	context.LoggerFactory = logger.NewColorLoggerFactory()
	Populate(context)
	return docker.NewProject(context, nil)
}
