package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/drep-project/DREP-Chain/common"
	"github.com/drep-project/DREP-Chain/common/fileutil"
	"gopkg.in/urfave/cli.v1"
)

var (
	// General settings
	ConfigFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TODO add config description",
	}

	HomeDirFlag = common.DirectoryFlag{
		Name:  "homedir",
		Usage: "Home directory for the datadir logdir and keystore",
	}

	PprofFlag = cli.StringFlag{
		Name:  "pprof",
		Usage: "ppfof for debug performance, --pprof 1",
	}
)

type Option func()

// DrepApp based on the cli.App, the module service operation is encapsulated.
// The purpose is to achieve the independence of each module and reduce dependencies as far as possible.
type DrepApp struct {
	Context *ExecuteContext
	*cli.App
	options []Option
}

// NewApp create a new app
func NewApp() *DrepApp {
	return &DrepApp{
		Context: &ExecuteContext{
			Quit: make(chan struct{}),
		},
		App:     cli.NewApp(),
		options: []Option{},
	}
}

// AddService add a server into context
func (mApp DrepApp) addServiceInstance(service Service) {
	mApp.Context.AddService(service)
}

// AddServiceType add many services
func (mApp DrepApp) AddServices(serviceInstances ...interface{}) error {
	nilService := reflect.TypeOf((*Service)(nil)).Elem()

	for _, serviceInstance := range serviceInstances {
		serviceType := reflect.TypeOf(serviceInstance)
		if serviceType.Kind() == reflect.Ptr {
			serviceType = serviceType.Elem()
		}
		serviceVal := reflect.New(serviceType)
		if !serviceVal.Type().Implements(nilService) {
			return ErrNotMatchedService
		}
		mApp.addService(serviceVal)
	}
	return nil
}

// addServiceType add a service and iterator all service that has added in and fields in current service,
// if exist , set set service in the field
func (mApp DrepApp) addService(serviceValue reflect.Value) {
	serviceType := serviceValue.Type()
	serviceNumFields := serviceType.Elem().NumField()
	for i := 0; i < serviceNumFields; i++ {
		serviceTypeField := serviceType.Elem().Field(i)
		refServiceName := GetServiceTag(serviceTypeField)
		if refServiceName != "" {
			preAddServices := mApp.Context.Services
			hasService := false
			for _, addedService := range preAddServices {

				if addedService.Name() == refServiceName {
					//TODO the filed to be set must be set public field, but it wiil be better to set it as a private field ,
					//TODO There are still some technical difficulties that need to be overcome.
					//TODO UnsafePointer may help
					serviceValue.Elem().Field(i).Set(reflect.ValueOf(addedService))
					hasService = true
				}
			}

			if !hasService {
				fmt.Println(fmt.Sprintf("service not exist %s require %s", serviceValue.Interface().(Service).Name(), refServiceName))
				//dlog.Debug("service not exist",  "Service", addedService.Name() ,"RefService", refServiceName)
			}
		}
	}
	mApp.addServiceInstance(serviceValue.Interface().(Service))
}

// GetServiceTag read service tag name to match service that has added in
func GetServiceTag(field reflect.StructField) string {
	serviceTagStr := field.Tag.Get("service")
	if serviceTagStr == "" {
		return ""
	}
	serviceName := strings.Split(serviceTagStr, ",")
	if len(serviceName) == 0 {
		return ""
	} else {
		return serviceName[0]
	}
}

// Init do something before app run
func (mApp *DrepApp) Option(option func()) {
	mApp.options = append(mApp.options, option)
}

//TODO need a more graceful  command supporter
//Run read the global configuration, parse the global command parameters,
// initialize the main process one by one, and execute the service before the main process starts.
func (mApp *DrepApp) Run() error {
	for _, op := range mApp.options {
		op()
	}

	mApp.Before = mApp.before
	mApp.Flags = append(mApp.Flags, ConfigFileFlag)
	mApp.Flags = append(mApp.Flags, HomeDirFlag)
	mApp.Flags = append(mApp.Flags, PprofFlag)

	allCommands, allFlags := mApp.Context.AggerateFlags()
	for i := 0; i < len(allCommands); i++ {
		allCommands[i].Flags = append(allCommands[i].Flags, allFlags...)
		allCommands[i].Action = mApp.action
	}
	mApp.Flags = append(mApp.Flags, allFlags...)
	mApp.App.Commands = allCommands
	mApp.Action = mApp.action
	if err := mApp.App.Run(os.Args); err != nil {
		return err
	}
	return nil
}

// action used to init and run each services
func (mApp *DrepApp) action(ctx *cli.Context) error {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			fmt.Println("app action err", err)
		}
		length := len(mApp.Context.Services)
		for i := length; i > 0; i-- {
			err := mApp.Context.Services[i-1].Stop(mApp.Context)
			if err != nil {
				return
			}
		}
	}()
	mApp.Context.Cli = ctx //NOTE this set is for different commmands-
	for _, service := range mApp.Context.Services {
		//config
		serviceValue := reflect.ValueOf(service)
		serviceType := serviceValue.Type()
		var config reflect.Value
		fieldValue := reflect.ValueOf(service).Elem().FieldByName("Config")
		if hasMethod(serviceType, "DefaultConfig") {
			defaultConfigVal := serviceValue.MethodByName("DefaultConfig").Call([]reflect.Value{})
			if len(defaultConfigVal) > 0 && !defaultConfigVal[0].IsNil() {
				config = defaultConfigVal[0]
			} else {
				t := fieldValue.Type()
				config = reflect.New(t.Elem())
			}
		} else {
			t := fieldValue.Type()
			config = reflect.New(t.Elem())
		}
		fieldValue.Set(config)
		//flag
		err := mApp.Context.UnmashalConfig(service.Name(), fieldValue.Interface())
		if err != nil {
			fmt.Println("unmashalConfig", service.Name(), err)
			return err
		}
		err = service.Init(mApp.Context)
		if err != nil {
			return err
		}
	}

	for _, service := range mApp.Context.Services {
		err := service.Start(mApp.Context)
		if err != nil {
			return err
		}
	}
	exit := make(chan struct{})
	exitSignal(exit)
	select {
	case <-exit:
	case <-mApp.Context.Quit:
	}
	return nil
}

//  read global config before main process
func (mApp *DrepApp) before(ctx *cli.Context) error {
	mApp.Context.Cli = ctx

	homeDir := ""
	if ctx.GlobalIsSet(HomeDirFlag.Name) {
		homeDir = ctx.GlobalString(HomeDirFlag.Name)
	} else {
		homeDir = common.AppDataDir(mApp.Name, false)
	}
	mApp.Context.ConfigPath = homeDir

	mApp.Context.CommonConfig = &CommonConfig{
		HomeDir: homeDir,
	}
	phaseConfig, err := loadConfigFile(ctx, homeDir)

	if err != nil {
		fmt.Println("before(),loadConfigFile err:", err)
		return err
	}
	mApp.Context.PhaseConfig = phaseConfig

	if ctx.GlobalIsSet(PprofFlag.Name) {
		go func() {
			fmt.Println("http://localhost:8080/debug/pprof")
			http.ListenAndServe("0.0.0.0:8080", nil)
		}()
	}

	return nil
}

//	loadConfigFile sed to read configuration files
func loadConfigFile(ctx *cli.Context, homeDir string) (map[string]json.RawMessage, error) {
	configFile := filepath.Join(homeDir, "config.json")

	if ctx.GlobalIsSet(ConfigFileFlag.Name) {
		file := ctx.GlobalString(ConfigFileFlag.Name)
		if fileutil.IsFileExists(file) {
			//report error when user specify
			return nil, ErrConfigiNotFound
		}
		configFile = file
	}

	if !fileutil.IsFileExists(configFile) {
		//use default
		cfg := &CommonConfig{
			HomeDir:    homeDir,
			ConfigFile: configFile,
		}
		originConfigBytes, err := json.MarshalIndent(cfg, "", "\t")
		if err != nil {
			return nil, err
		}
		fileutil.EnsureFile(configFile)
		file, err := os.OpenFile(configFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
		if err != nil {
			return nil, err
		}
		file.Write(originConfigBytes)
	}
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(content))
	jsonPhase := map[string]json.RawMessage{}
	err = json.Unmarshal(content, &jsonPhase)
	if err != nil {
		return nil, err
	}
	return jsonPhase, nil
}

func hasMethod(t reflect.Type, name string) (hashMehod bool) {
	defer func() {
		if err := recover(); err != nil {
			hashMehod = false
		}
	}()
	_, hashMehod = t.MethodByName(name)
	return
}
