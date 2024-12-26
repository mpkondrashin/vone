/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Submit - Cloud Sandbox Comman Line Interface

	main.go - command line utility to run cloud sandbox functions
*/

package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mpkondrashin/vone"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	EnvPrefix      = "VONE"
	ConfigFileName = "config"
	ConfigFileType = "yaml"
)

const (
	cmdCheck            = "check"
	cmdSubmit           = "submit"
	cmdQuota            = "quota"
	cmdPDF              = "pdf"
	cmdGetEndpointData  = "endpoint"
	cmdListGetEndpoints = "list"
	cmdHighRiskDevices  = "high_risk_devices"
	cmdPing             = "ping"
	cmdAddEception      = "it_exception"
	cmdGetOATEvents     = "oat"
)

const (
	flagAddress       = "address"
	flagToken         = "token"
	flagLog           = "log"
	flagFileName      = "filename"
	flagMask          = "mask"
	flagURL           = "url"
	flagURLsFile      = "urlfile"
	flagTimeout       = "timeout"
	flagID            = "id"
	flagQuery         = "query"
	flagFilter        = "filter"
	flagOrderBy       = "order_by"
	flagTop           = "top"
	flagProxy         = "proxy"
	flagProxyUser     = "proxy_user"
	flagProxyPassword = "proxy_password"
	flagProxyDomain   = "proxy_domain"
	flagSOType        = "so_type"
	flagSO            = "so"
	flagDescription   = "description"
	flagDetectedStart = "detected_start"
	flagDetectedEnd   = "detected_end"
	flagIngestedStart = "ingested_start"
	flagIngestedEnd   = "ingested_end"
)

type command interface {
	Name() string
	Description() string
	Init(args []string) error
	Execute() error
}

type baseCommand struct {
	name        string
	description string
	visionOne   *vone.VOne
	fs          *pflag.FlagSet
}

func (c *baseCommand) Setup(name, description string) {
	c.name = name
	c.description = description
	c.fs = pflag.NewFlagSet(name, pflag.ExitOnError)
	c.fs.String(flagAddress, "", "Vision One entry point URL")
	c.fs.String(flagToken, "", "Vision One API Token")
	c.fs.String(flagLog, "", "Log file path")

	c.fs.String(flagProxy, "", "Proxy URL (scheme://address:port)")
	c.fs.String(flagProxyUser, "", "Proxy username")
	c.fs.String(flagProxyPassword, "", "Proxy password")
	c.fs.String(flagProxyDomain, "", "Proxy domain (for NTLM auth)")

	c.fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\nAvailable options:\n", c.description)
		c.fs.PrintDefaults()
	}

}

func (c *baseCommand) Name() string {
	return c.name
}

func (c *baseCommand) Description() string {
	return c.description
}

func (c *baseCommand) String() string {
	return c.name
}

func (c *baseCommand) Init(args []string) error {
	err := c.fs.Parse(os.Args[2:])
	if err != nil {
		return err
	}
	if err := viper.BindPFlags(c.fs); err != nil {
		panic(err)
	}
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()

	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileType)
	path, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(path)
		viper.AddConfigPath(dir)
	}
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		notFoundErr, ok := err.(viper.ConfigFileNotFoundError)
		_ = notFoundErr
		if !ok {
			panic(err) //Fatal(RCConfigReadError, "ReadInConfig: %v", err)
		} else {
			log.Printf("%s: loaded", ConfigFileName)
		}
		//LogIt(Debug, "ReadInConfig: %v", notFoundErr)
	}
	c.visionOne = vone.NewVOne(
		viper.GetString(flagAddress),
		viper.GetString(flagToken),
	)
	rateLimiter := vone.NewAdaptiveRateLimiter(vone.VOneRateLimitSurpassedError, nil)
	c.visionOne.SetRateLimiter(rateLimiter)

	if viper.GetString(flagProxy) != "" {
		u, err := url.Parse(viper.GetString(flagProxy))
		if err != nil {
			log.Fatal(err)
		}
		proxy := vone.NewProxy(u)
		log.Println("Use proxy")
		if viper.GetString(flagProxyUser) != "" {
			if viper.GetString(flagProxyPassword) == "" {
				log.Fatal("missing proxy password")
			}
			if viper.GetString(flagProxyDomain) != "" {
				log.Println("Use NTLM proxy auth")
				proxy.NTLMAuth(
					viper.GetString(flagProxyUser),
					viper.GetString(flagProxyPassword),
					viper.GetString(flagProxyDomain),
				)
			} else {
				log.Println("Use basic proxy auth")
				proxy.BasicAuth(
					viper.GetString(flagProxyUser),
					viper.GetString(flagProxyPassword),
				)
			}
		}
		c.visionOne.AddTransportModifier(proxy.GetModifier())
	}
	//	c.ctx = context.Background()
	//	if viper.GetBool(flagDryRun) {
	//		c.ctx = ddan.DryRunContext(context.Background(), func(line string) {
	//			fmt.Println(line)
	//		})
	//	}
	logFilePath := viper.GetString(flagLog)
	if logFilePath != "" {
		logFile, err := os.Create(logFilePath)
		if err != nil {
			log.Fatal(err)
		}
		log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	}
	return nil
}

var commands = []command{
	newCommandCheck(),
	newCommandSubmit(),
	newCommandQuota(),
	newCommandPDF(),
	newCommandGetEndpointData(),
	newCommandListEndpoints(),
	newCommandHighRiskDevices(),
	newCommandAddIT(),
	newCommandGetOATEvents(),
}

func usage() {
	var commandNames []string
	for _, c := range commands {
		commandNames = append(commandNames, c.Name())
	}
	fmt.Fprintf(os.Stderr, "VOne — various Trend Micro Vision One API functions demo\nUsage: %s%s {%s} [options]\n",
		name(), exe(), strings.Join(commandNames, "|"))
	fmt.Fprintf(os.Stderr, "Commands available:\n")
	for _, c := range commands {
		fmt.Fprintf(os.Stderr, "\t%s - %s\n", c.Name(), c.Description())
	}
	fmt.Fprintf(os.Stderr, "For more details, run vone <command> --help\n")
	os.Exit(2)
}

func pickCommand(args []string) error {
	subcommand := args[0]
	for _, cmd := range commands {
		if cmd.Name() == subcommand {
			err := cmd.Init(args[1:])
			if err != nil {
				return fmt.Errorf("Init error: %w", err)
			}
			log.Printf("Command %s\n", cmd.Name())
			return cmd.Execute()
		}
	}
	return fmt.Errorf("unknown command: %s", subcommand)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	err := pickCommand(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Done")
	}
}

func exe() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func name() string {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		panic(fmt.Errorf("runtime.Caller() error"))
	}
	dir := filepath.Dir(path)
	folder := filepath.Base(dir)
	return folder
}
