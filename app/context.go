package app

import (
	"encoding/json"
	"github.com/asaskevich/EventBus"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v1"
)

var (
	CommandHelpTemplate = `{{.cmd.Name}}{{if .cmd.Subcommands}} command{{end}}{{if .cmd.Flags}} [command options]{{end}} [arguments...]
{{if .cmd.Description}}{{.cmd.Description}}
{{end}}{{if .cmd.Subcommands}}
SUBCOMMANDS:
	{{range .cmd.Subcommands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
	{{end}}{{end}}{{if .categorizedFlags}}
{{range $idx, $categorized := .categorizedFlags}}{{$categorized.Name}} OPTIONS:
{{range $categorized.Flags}}{{"\t"}}{{.}}
{{end}}
{{end}}{{end}}`
)

func init() {
	cli.AppHelpTemplate = `{{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

	cli.CommandHelpTemplate = CommandHelpTemplate
}

// CommonConfig read before app run,this fuction shared by other moudles
type CommonConfig struct {
	HomeDir    string `json:"homeDir,omitempty"`
	ConfigFile string `json:"configFile,omitempty"`
}

// API describes the set of methods offered over the RPC interface
type API struct {
	Namespace string      // namespace under which the rpc methods of Service are exposed
	Version   string      // api version for DApp's
	Service   interface{} // receiver instance which holds the methods
	Public    bool        // indication if the methods must be considered safe for public use
}

// Services can customize their own configuration, command parameters, interfaces, services
type Service interface {
	Name() string                              // service  name must be unique
	Api() []API                                // Interfaces required for services
	CommandFlags() ([]cli.Command, []cli.Flag) // flags required for services
	//P2pMessages() map[int]interface{}
	//Receive(context actor.Context)
	Init(executeContext *ExecuteContext) error
	Start(executeContext *ExecuteContext) error
	Stop(executeContext *ExecuteContext) error
}

// ExecuteContext centralizes all the data and global parameters of application execution,
// and each service can read the part it needs.
type ExecuteContext struct {
	ConfigPath   string
	CommonConfig *CommonConfig //
	PhaseConfig  map[string]json.RawMessage
	Cli          *cli.Context
	LifeBus      EventBus.Bus
	Services     []Service

	GitCommit string
	Usage     string
	Quit      chan struct{}
}

// AddService add a service to context, The application then initializes and starts the service.
func (econtext *ExecuteContext) AddService(service Service) {
	econtext.Services = append(econtext.Services, service)
}

// GetService In addition, there is a dependency relationship between services.
// This method is used to find the dependency services you need in the context.
func (econtext *ExecuteContext) GetService(name string) Service {
	for _, service := range econtext.Services {
		if service.Name() == name {
			return service
		}
	}
	return nil
}

//	GetConfig Configuration is divided into several segments,
//	each service only needs to obtain its own configuration data,
//	and the parsing process is also controlled by each service itself.
func (econtext *ExecuteContext) GetConfig(phaseName string) json.RawMessage {
	phaseConfig, ok := econtext.PhaseConfig[phaseName]
	if ok {
		return phaseConfig
	} else {
		return nil
	}
}

// GetFlags aggregate command configuration items required for each service
func (econtext *ExecuteContext) AggerateFlags() ([]cli.Command, []cli.Flag) {
	allFlags := []cli.Flag{}
	allCommands := []cli.Command{}
	for _, service := range econtext.Services {
		commands, defaultFlags := service.CommandFlags()
		if commands != nil {
			allCommands = append(allCommands, commands...)
		}
		if defaultFlags != nil {
			allFlags = append(allFlags, defaultFlags...)
		}
	}
	return allCommands, allFlags
}

//	GetApis aggregate interface functions for each service to provide for use by RPC services
func (econtext *ExecuteContext) GetApis() []API {
	apis := []API{}
	for _, service := range econtext.Services {
		apis = append(apis, service.Api()...)
	}
	return apis
}

////	GetApis aggregate interface functions for each service to provide for use by RPC services
//func (econtext *ExecuteContext) GetMessages() (map[int]interface{}, error)  {
//	msg := map[int]interface{}{}
//	for _, service := range econtext.Services {
//		for k, v := range service.P2pMessages() {
//			if _, ok := msg[k]; ok {
//				return nil, errors.New("exist p2p message")
//			}
//			msg[k] = v
//		}
//	}
//	return msg, nil
//}

//	RequireService When a service depends on another service, RequireService is used to obtain the dependent service.
func (econtext *ExecuteContext) RequireService(name string) Service {
	for _, service := range econtext.Services {
		if service.Name() == name {
			return service
		}
	}
	panic(errors.Wrap(ErrServiceNotFound, name))
}

func (econtext *ExecuteContext) UnmashalConfig(serviceName string, config interface{}) error {
	service := econtext.GetService(serviceName)
	if service == nil {
		return errors.Wrapf(ErrServiceNotFound, "service name:%s", serviceName)
	}
	phase := econtext.GetConfig(service.Name())
	if phase == nil {
		return nil
	}
	err := json.Unmarshal(phase, config)
	if err != nil {
		return err
	}
	return nil
}
