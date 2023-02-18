/*
	Trend Micro Vision One API SDK
	(c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)

	Submit - Cloud Sandbox Comman Line Interface

	main.go - command line utility to run cloud sandbox functions
*/

package main

import (
	"context"
	"fmt"
	"log"
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
	cmdSubmit = "submit"
	cmdQuota  = "quota"
)

const (
	flagAddress  = "address"
	flagToken    = "token"
	flagFileName = "filename"
	flagTimeout  = "timeout"
)

type command interface {
	Name() string
	Init(args []string) error
	Execute() error
}

type baseCommand struct {
	name      string
	visionOne *vone.VOne
	ctx       context.Context
	fs        *pflag.FlagSet
}

func (c *baseCommand) Setup(name string) {
	c.name = name
	c.fs = pflag.NewFlagSet(name, pflag.ExitOnError)
	c.fs.String(flagAddress, "", "Vision One entry point URL")
	c.fs.String(flagToken, "", "Vision One API Token")
}

func (c *baseCommand) Name() string {
	return c.name
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
	//	c.ctx = context.Background()
	//	if viper.GetBool(flagDryRun) {
	//		c.ctx = ddan.DryRunContext(context.Background(), func(line string) {
	//			fmt.Println(line)
	//		})
	//	}
	return nil
}

type commandQuota struct {
	baseCommand
}

func newCommandQuota() *commandQuota {
	c := &commandQuota{}
	c.Setup(cmdQuota)
	return c
}

func (c *commandQuota) Execute() error {
	quota, err := c.visionOne.SandboxDailyReserve().Do()
	if err != nil {
		return err
	}
	log.Printf("Submission Reserve Count: %d", quota.SubmissionReserveCount)
	log.Printf("Submission Remaining Count: %d", quota.SubmissionRemainingCount)
	log.Printf("Submission Count: %d", quota.SubmissionCount)
	log.Printf("Submission ExemptionCount: %d", quota.SubmissionExemptionCount)
	log.Printf("Submission Count Detail: File Count: %d", quota.SubmissionCountDetail.FileCount)
	log.Printf("Submission Count Detail: File Exemption Count: %d", quota.SubmissionCountDetail.FileExemptionCount)
	log.Printf("Submission Count Detail: URL Count: %d", quota.SubmissionCountDetail.URLCount)
	log.Printf("Submission Count Detail: URL Exemption Count: %d", quota.SubmissionCountDetail.URLExemptionCount)
	return nil
}

var commands = []command{
	newCommandSubmit(),
	newCommandQuota(),
}

func usage() {
	var commandNames []string
	for _, c := range commands {
		commandNames = append(commandNames, c.Name())
	}
	fmt.Printf("VOne\nUsage: %s%s {%s} [options]\n",
		name(), exe(), strings.Join(commandNames, "|"))
	os.Exit(2)
}

func pickCommand(args []string) error {
	subcommand := args[0]
	for _, cmd := range commands {
		if cmd.Name() == subcommand {
			err := cmd.Init(args[1:])
			if err != nil {
				return fmt.Errorf("Init error: %v", err)
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
