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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

const Services_default_url = "https://raw.github.houston.softwaregrp.net/SMA-RnD/suite-deployer/master/docker/profiles/ci/itsmaProducts.json"

type ItsmaSuiteService struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Mode    []string `json:"modes"`
}

type ItsmaSuiteServices struct {
	Services []ItsmaSuiteService `json:"itsmaServices"`
}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all available services in itsma suite.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		services_url := viper.GetString("commands.list.servicesVersionInfoLocation")
		if services_url == "" {
			services_url = Services_default_url
		}

		fmt.Println("===ListCommand. Getting services version from: ", services_url)
		data, err := downloadFile(services_url)
		if err != nil {
			fmt.Println("===ListCommand. Failed. Download itsmaProducts.json error: ", err.Error())
			return
		}
		var services ItsmaSuiteServices
		if err = json.Unmarshal(data, &services); err != nil {
			fmt.Println("===ListCommand. Failed. Json unmarshal error: ", err.Error())
			return
		}

		fmt.Println("Itsma suite services: ")
		fmt.Println("-------------------------------------------------------")
		fmt.Println("Name \t\t\t  Mode \t\tVersion(Stable)")
		fmt.Println("-------------------------------------------------------")
		for _, service := range services.Services {
			fmt.Printf("%-25s %-15s %-15s\n", service.Name, strings.Join(service.Mode,","), service.Version)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

func downloadFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		errMsg := fmt.Sprintf("===DownloadFile: Failed. Http Get error: %s.", err.Error())
		return nil, errors.New(errMsg)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errMsg := fmt.Sprintf("===DownloadFile: Failed. Read http response body error: %s.", err.Error())
		return nil, errors.New(errMsg)
	}
	if resp.StatusCode == 200 {
		return body, nil
	} else {
		errMsg := fmt.Sprintf("===DownloadFile: Failed. Http rc expected: 200, actual: %d. Http rb: %s.",
			resp.StatusCode, string(body))
		return nil, errors.New(errMsg)
	}
}
