package conf

import "github.com/docker/libcompose/yaml"

type ServerConf struct {
	Name      string
	DockerDir string
	Env       string
	Host      string
	Port      string
	ApiDir    string
	ApiFile   string
	Command   string
}

type ClientConf struct {
	Name      string
	DockerDir string
	Env       string
	Lang      string
	Command   string
}

type StructCommand struct {
	Command yaml.Command `yaml:"command,flow,omitempty"`
}

type ProjectConfig struct {
	Server struct {
		Name        string               `yaml:"name"`
		DockerDir   string               `yaml:"dockerdir,omitempty"`
		Image       string               `yaml:"image,omitempty"`
		EnvFile     yaml.Stringorslice   `yaml:"env,omitempty"`
		Environment yaml.MaporEqualSlice `yaml:"environment,omitempty"`
		ApiHost     string               `yaml:"apihost"`
		ApiPort     string               `yaml:"apiport"`
		ApiDir      string               `yaml:"apidir"`
		ApiFile     string               `yaml:"apifile,omitempty"`
		Command     yaml.Command         `yaml:"command,flow,omitempty"`
	} `yaml:"api-server"`
	Client struct {
		Name        string               `yaml:"name"`
		DockerDir   string               `yaml:"dockerdir,omitempty"`
		Image       string               `yaml:"image,omitempty"`
		EnvFile     yaml.Stringorslice   `yaml:"env,omitempty"`
		Environment yaml.MaporEqualSlice `yaml:"environment,omitempty"`
		Lang        yaml.Stringorslice   `yaml:"lang,omitempty"`
		Command     yaml.Command         `yaml:"command,flow,omitempty"`
	} `yaml:"api-client"`
}

func NewServerService() *ServerConf {
	return &ServerConf{
		Name:    "api-server",
		Host:    "server",
		Port:    "50051",
		ApiDir:  "/api",
		ApiFile: "/api/main.proto",
	}
}

func NewClientService() *ClientConf {
	return &ClientConf{
		Name: "api-client",
	}
}
