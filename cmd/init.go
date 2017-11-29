// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
)

type Service struct {
	Name          string
	Version       string
	ImageRegistry string
}

func (s Service) ServiceDir() string {
	return "itom-" + s.Name + "-" + s.Version
}

var service Service

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a itsma service",
	Long: `Create a skeleton of the itsma service, including manifest.yaml
	deployer-controller.yaml pom.xml...`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		if err := CreateService(service); err != nil {
			log.Fatal(err)
		}
	},
}

func GenerateService(service Service) (err error) {
	sourceDir := "./templates/itsma_service"
	targetDir := service.ServiceDir()

	tmplFiles := []string{"Dockerfile", "settings.xml", "manifest.yaml", "service_property.json", "start.sh", "tag.properties", "pom.xml"}
	for _, tmplFile := range tmplFiles {
		if tmplFile == "settings.xml" {
			if err = ExecuteCommand(true, "cp", "-f", sourceDir+"/settings.xml", targetDir+"/settings.xml"); err != nil {
				log.Fatal("Copy settings.xml failed: ", err)
			}
		} else {
			err = GenerateFileFromTemplate(sourceDir+"/"+tmplFile, targetDir+"/"+tmplFile, service)
			if err != nil {
				log.Fatal("Generate docker files failed: ", err)
			}
		}
	}

	if err = os.MkdirAll(targetDir+"/yaml_templates", 0770); err != nil {
		log.Fatal("Create directory /yaml_templates failed: ", err)
	}

	return nil
}

func CreateService(service Service) error {
	if service.Name == "" {
		return errors.New("Service name can NOT be empty")
	}

	err := os.Mkdir(service.ServiceDir(), 0770)
	if err != nil {
		if os.IsExist(err) {
			log.Printf("Service: %s with version: %s already exists.\n", service.Name, service.Version)
		}
		return err
	}

	return GenerateService(service)
}

func MkdirIfNotExist(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.MkdirAll(path, 0770); err != nil {
				log.Fatal(err)
				return err
			}
		} else {
			return err
		}
	} else {
		if !fi.IsDir() {
			return errors.New("Mkdir failed because of file with same name exists")
		}
	}

	return err
}

func init() {
	RootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	initCmd.Flags().StringVarP(&service.Name, "servicename", "n", "", "(required) servie name")
	initCmd.Flags().StringVarP(&service.Version, "serviceversion", "v", "1.0.0", "(optional) servie version")
	initCmd.Flags().StringVarP(&service.ImageRegistry, "imageregistry", "i", "shc-harbor-dev.hpeswlab.net/itsma", "(required) registry where images to store")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
