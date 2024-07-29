package main

import (
	"context"
	"log"

	"github.com/spf13/viper"
)

type commandTIAddException struct {
	baseCommand
}

func newCommandAddIT() *commandTIAddException {
	c := &commandTIAddException{}
	c.Setup(cmdAddEception)
	return c
}

func (c *commandTIAddException) Setup(name string) {
	c.baseCommand.Setup(name, "Add IoC exception")
	c.fs.String(flagSOType, "", "IoC type (url, domain, ip, senderMailAddress, fileSha1, fileSha256)")
	c.fs.String(flagSO, "", "IoC value")
	c.fs.String(flagDescription, "", "IoC description")
}

func (c *commandTIAddException) Execute() error {
	so := viper.GetString(flagSO)
	if so == "" {
		log.Fatalf("--%s parameter can not be empty", flagSO)
	}
	addException := c.visionOne.AddExceptions()
	description := viper.GetString(flagDescription)
	soType := viper.GetString(flagSOType)
	switch soType {
	case "url":
		addException.AddURL(so, description)
	case "domain":
		addException.AddDomain(so, description)
	case "ip":
		addException.AddIP(so, description)
	case "senderMailAddress":
		addException.AddSenderMailAddress(so, description)
	case "fileSha1":
		addException.AddFileSHA1(so, description)
	case "fileSha256":
		addException.AddFileSHA256(so, description)
	default:
		log.Fatalf("--%s parameter has wrong value: '%s'. It should one of: url, domain, ip, senderMailAddress, fileSha1, fileSha256", flagSO, soType)
	}

	response, err := addException.Do(context.TODO())
	if err != nil {
		return err
	}
	for _, r := range *response {
		if r.Status/100 == 2 {
			log.Println("Ok")
			continue
		}
		log.Printf("%s: %s", r.Body.Error.Code, r.Body.Error.Message)
		//return fmt.Errorf("%d: %s", response.Body.Error.Code, response.Body.Error.Message)
	}
	//log.Printf("Responce: %v", response)
	return nil
}
