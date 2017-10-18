package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"os"
	"encoding/json"
	"text/template"
	"log"
)

var settings struct {
	ServerMode bool `json:"serverMode"`
	SourceDir  string `json:"sourceDir"`
	TargetDir  string `json:"targetDir"`
	Local_token 	   string `json:"local_token"`
}

//this method is used to create application
func creteebApplication(){
	//svc := elasticbeanstalk.New(session.New())
	// Create a ElasticBeanstalk client with additional configuration
	svc := elasticbeanstalk.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))
	input := &elasticbeanstalk.CreateApplicationInput{
		ApplicationName: aws.String("my-app-from-go"),
		Description:     aws.String("my application creatinhg using go"),
	}

	result, err := svc.CreateApplication(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elasticbeanstalk.ErrCodeTooManyApplicationsException:
				fmt.Println(elasticbeanstalk.ErrCodeTooManyApplicationsException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func main() {
	t, err := template.ParseFiles( "envconfig.yaml.template")
	if err != nil {
		log.Print(err)
		return
	}

	f, err := os.Create("envconfig.yaml")
	if err != nil {
		log.Println("create file: ", err)
		return
	}

	// then config file settings

	configFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&settings); err != nil {
		fmt.Println("parsing config file", err.Error())
	}

	//fmt.Printf("%v %s %s %s", settings.ServerMode, settings.SourceDir, settings.TargetDir, settings.Locale)

	err = t.Execute(f, settings)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	return

}
