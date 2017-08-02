package compose

import (
	"beji/grpcapi/conf"
	"beji/grpcapi/grpc"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/yaml"
	"github.com/lytics/logrus"
	"github.com/spf13/pflag"

	"github.com/docker/libcompose/project/options"

	y "gopkg.in/yaml.v2"
)

type Project struct {
	serverConf *conf.ServerConf
	clientConf *conf.ClientConf
}

func IfNotEmptyAssign(a *string, b string) {
	if b != "" {
		*a = b
	}
}

func (p *Project) CreateClientCmd() yaml.Command {
	if p.clientConf.Lang == "" && p.clientConf.Command == "" {
		return nil
	}

	var byteBuffer bytes.Buffer
	byteBuffer.WriteString("command: ")
	byteBuffer.WriteString("bash -c \"")

	if p.clientConf.Lang != "" {
		langs := strings.Split(p.clientConf.Lang, ",")
		cmds := grpc.NewGrpcCmd(p.serverConf.ApiDir, p.serverConf.ApiDir, p.serverConf.ApiFile, langs...)

		for _, cmd := range cmds {
			byteBuffer.WriteString(cmd)
			byteBuffer.WriteString(" && ")
		}
	}

	if p.clientConf.Command != "" {
		tmp := strings.Trim(p.clientConf.Command, "bash -c")
		tmp = strings.Trim(tmp, "\"")
		byteBuffer.WriteString(tmp)
	}
	byteBuffer.WriteString("\"")

	s := &conf.StructCommand{}
	_ = y.Unmarshal(byteBuffer.Bytes(), s)

	return s.Command
}

func (p *Project) CreateServerCmd() yaml.Command {
	if p.serverConf.Command == "" {
		return nil
	}

	var byteBuffer bytes.Buffer
	byteBuffer.WriteString("command: ")
	byteBuffer.WriteString("bash -c \"")

	tmp := strings.Trim(p.serverConf.Command, "bash -c")
	tmp = strings.Trim(tmp, "\"")
	byteBuffer.WriteString(tmp)
	byteBuffer.WriteString("\"")

	s := &conf.StructCommand{}
	_ = y.Unmarshal(byteBuffer.Bytes(), s)

	return s.Command
}

func NewProject(s, c map[string]string) *Project {
	p := Project{serverConf: conf.NewServerService(), clientConf: conf.NewClientService()}

	IfNotEmptyAssign(&p.serverConf.DockerDir, s["dockerdir"])
	IfNotEmptyAssign(&p.serverConf.Port, s["apiport"])
	IfNotEmptyAssign(&p.serverConf.Env, s["env"])
	IfNotEmptyAssign(&p.serverConf.ApiDir, s["apidir"])
	IfNotEmptyAssign(&p.serverConf.ApiFile, s["apifile"])

	IfNotEmptyAssign(&p.clientConf.DockerDir, c["dockerdir"])
	IfNotEmptyAssign(&p.clientConf.Env, c["env"])
	IfNotEmptyAssign(&p.clientConf.Lang, c["lang"])
	IfNotEmptyAssign(&p.clientConf.Command, c["command"])

	return &p
}

func (p *Project) GetServerServiceConfig() *config.ServiceConfig {
	c := config.ServiceConfig{}

	c.Build = yaml.Build{Context: p.serverConf.DockerDir}
	c.EnvFile = []string{p.serverConf.Env}
	c.Expose = []string{p.serverConf.Port}
	c.Volumes = &yaml.Volumes{
		Volumes: []*yaml.Volume{
			{
				Destination: p.serverConf.ApiDir,
			},
		}}

	c.Command = p.CreateServerCmd()

	for _, x := range c.Command {
		fmt.Printf("ser: %+v\n", x)
	}

	return &c
}
func (p *Project) GetClientServiceConfig() *config.ServiceConfig {
	c := config.ServiceConfig{}

	c.Build = yaml.Build{Context: p.clientConf.DockerDir}
	c.EnvFile = []string{p.clientConf.Env}
	c.Environment = []string{"API_HOST=" + p.serverConf.Host, "API_PORT=" + p.serverConf.Port}
	c.Links = []string{p.serverConf.Name + ":" + p.serverConf.Host}
	c.VolumesFrom = []string{p.serverConf.Name}

	c.Command = p.CreateClientCmd()

	for _, x := range c.Command {
		fmt.Printf("cli: %+v\n", x)
	}

	return &c
}

func (p *Project) Up(flags *pflag.FlagSet) error {

	pf := &ProjectFactory{}
	project, err := pf.Create()

	if err != nil {
		log.Fatal(err)
		return err
	}
	project.AddConfig(p.serverConf.Name, p.GetServerServiceConfig())
	project.AddConfig(p.clientConf.Name, p.GetClientServiceConfig())

	debug, _ := flags.GetBool("debug")
	nr, _ := flags.GetBool("no-recreate")
	fr, _ := flags.GetBool("force-recreate")
	nb, _ := flags.GetBool("no-build")
	build, _ := flags.GetBool("build")
	timeout, _ := flags.GetInt("timeout")

	options := options.Up{
		Create: options.Create{
			NoRecreate:    nr,
			ForceRecreate: fr,
			NoBuild:       nb,
			ForceBuild:    build,
		},
	}
	ctx, cancelFun := context.WithCancel(context.Background())

	err = project.Up(ctx, options)
	if err != nil {
		log.Fatal(err)
		return err
	}

	if !debug {
		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		errChan := make(chan error)
		go func() {
			errChan <- project.Log(ctx, true)
		}()
		go func() {
			select {
			case <-signalChan:
				fmt.Printf("\nGracefully stopping...\n")
				cancelFun()
				//projectStop
				err := project.Stop(context.Background(), timeout)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				cleanupDone <- true
			case err := <-errChan:
				if err != nil {
					fmt.Println(err)
					logrus.Fatal(err)
				}
				cleanupDone <- true
			}
		}()
		<-cleanupDone
		return nil
	}

	return nil
}
