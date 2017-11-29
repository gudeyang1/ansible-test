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
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
	"net"
	"log"
	"encoding/json"
	"encoding/xml"
)

var (
	MasterIP            string
	NodeIP              string
	Services            string
	RegistryUrls        string
	ControllerImgTags   string
	DeployerImgTag      string
	DeepClean           string
	Depth               int
	DeployerRegistry    string
	VsphereUser         string
	VspherePassword     string
	VsphereDataCenter   string
	VsphereFolder       string
	VsphereSnapShotName string
	VCenterHostName     string
	DeactivatedService  string
	DeactivatedServiceList	[]string
	Patch				bool
	Thinpool			string
	Mode				string
	ItomSuiteSize		string
	AutomationTag		string
	BackupSuite			string
	InstallType			int
	BackupPackageDir	string
	BackupPackageName	string
	Upgrade				string
	UpgradeImageTag		string
	AdminPassword		string
	CdfUrl				string
	SuiteDataImageTag	string
	MixMode				string



)

type ItsmaService struct {
	Name             string
	RegistryUrl      string
	ControllerImgTag string
}

type deploy_data struct {
	Masters         []string
	Workers         []string
	Services        []ItsmaService
	DeployerImgTag  string
	Registry        string
	DeepClean       string
	Depth           int
	VsphereUser     string
	VspherePassword string
	Patch			bool
	Thinpool		string
	Mode			string
	ItomSuiteSize	string
	AutomationTag	string
	BackupSuite		string
	InstallType		string
	BackupPackageDir	string
	BackupPackageName	string
	Upgrade				string
	UpgradeImageTag		string
	AdminPassword		string
	CdfUrl				string
	SuiteDataImageTag	string
	MixMode				string
}

// for generate cdf download url EOF
type MetadateOfNexus struct {
	XMLName     xml.Name `xml:"metadata"`
	GroupId     string   `xml:"groupId,attr"`
	Versioning         []Versioning `xml:"versioning"`
}

type Versioning struct {
	Latest 		string   	`xml:"latest"`
	Release   	string   	`xml:"release"`
}
// for generate cdf download url EOF

var XMode  = GetServiceNameFormJson("X_MODE")
var HMode  = GetServiceNameFormJson("H_MODE")

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploy itsma suite service",
	Long: `deploy one or more itsma suite services on target environment. 
	Examples:
	  0. Only deploy CDF
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -d 1

	  1. Deploy the whole itsma suite, including all stable services.
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip"

	  2. Deploy single itsma service of the stable version (such as itom-ucmdb)
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -s itom-ucmdb

	  3. Deploy single itsma service of the test/dev version (such as itom-ucmdb)
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -s itom-ucmdb --controllerimgtags test

	  4. Deploy more than one itsma services of the stable version (such as itom-ucmdb and itom-sm)
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -s "itom-ucmdb itom-sm"

	  5. Deploy more than one itsma services of the test/dev version (such as itom-ucmdb and itom-sm)
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -s "itom-ucmdb itom-sm" --controllerimgtags "test test"

	  6. Deploy more than one itsma services, including both the stable of the test/dev version (such as itom-ucmdb of stable version and itom-sm of test/dev version)
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -s "itom-sm itom-ucmdb" --controllerimgtags "test"

	  Note! The order of services are: -s "itom-sm(test/dev version) itom-xxx(test/dev version) itom-ucmdb(stable version) itom-xxx(stable version)", and only specify the tags in --controllerimgtags
	  7. based on 6 , but with specified registry
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" -s "itom-sm itom-ucmdb" --controllerimgtags "test" -r "sm-registry ucmdb-registry"
	  Note ! if you don't speclfy ucmdb-registry , suitekit will use default registry shc-harbor-dev.hpeswlab.net/itsma

	  8. Deploy with configuration reading from file (support combination of command parameters and configuration in file. command parameters precedes configuration in file)
	  ./suitectl --config .suitectl.yaml deploy -c default -d 1
	   please copy the default config file .suitectl.yaml.default , but DO NOT modify/add any key , just change the value .

	  9. Deploy all services except itom-cmdb in H_MODE
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip"  --ds "itom-cmdb"

	  10. Patch a new service tag in X_MODE
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip"  -s "patched_service1 patched_service2" --patch --controllerimgtags "patched_service1_tag patched_service2_tag" -M X_MODE

	  11. Disable Thinpool for the Docker devicemapper storage driver(use Thinpool as default)
          ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --thinpool false

	  12. Deploy suite with xsmall size
      	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --suite-size "xsmall"
	  13. Backup an exist suite
	   ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --backup true

	  14. Install from backup
	  14.1. install a new nev
		./suitectl deploy -m NEW_masterip -n "NEW_node1_ip NEW_node2_ip NEW_node3_ip" -d 1
	  14.2 copy backup package and DB data from old env to new installed CDF env
	  14.3 ./suitectl deploy -m NEW_masterip -n "NEW_node1_ip NEW_node2_ip NEW_node3_ip" --install-type 2 --backup-pkg-dir "1504687116129" --backup-pkg-name "ITSMA_v2017.07_1504687116129.zip"

	  15. Upgrade suite
	  # CDF version >= 147 , kit version >= 1.0.0.10
		15.1 Prepare suite env
		#backup automatically ,and use the latest backup dir and backup package name
		15.2 ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --upgrade true --upgrade-image-tag 170922 // tag is optional
		#backup automatically, but has backup records and specify the backup dir and backup package name
		15.3 ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --upgrade true --upgrade-image-tag 170922 --backup-pkg-dir "1508986223681" --backup-pkg-name "ITSMA_v2017.11_1508986223681.zip"
	  16. Specify the password of mng-protal
		#Run this command after you changed the password of mng-protal, DO NOT specify the password for the first running!
		./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --admin-password "your password"
	  17. Deploy CDF with specified version
	  ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --cdfurl {cdf download url}
	  18. Deploy with specified suite-data tag
	   ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --suite-data-tag "PR-XXX"
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Work your own magic here
		if err := checkFlags(); err != nil {
			return err
		}


		var deploy deploy_data
		deploy.Masters = strings.Fields(MasterIP)
		deploy.Workers = strings.Fields(NodeIP)
		DeactivatedServiceList = strings.Fields(DeactivatedService)
		Activated_Services_list := strings.Fields(Services)
		if DeactivatedService != ""{
			//All_service_list := GetServiceNameFormJson()
			All_service_list := GetAllServiceList()
			//the all service list before add --ds
			//fmt.Println("all service list is ", All_service_list)
			//fmt.Println("deactivated service:", DeactivatedServiceList)
			count := 0
			for _,value :=range DeactivatedServiceList{
				for k,v := range All_service_list{
					if value == v{
						// delete deactivated services from all service list
						All_service_list = append(All_service_list[:k],All_service_list[k+1:]...)
						count += 1
					}
				}
			}
			// the service list after add --ds
			//fmt.Println(All_service_list)
			if count != len(DeactivatedServiceList){
				return errors.New("DeactivatedService name must be the same as ItsmaProduct.json , please check service name again.")
				}
			Activated_Services_list = All_service_list
		}

		for k, name := range  Activated_Services_list{
			var s ItsmaService
			s.Name = name
			if RegistryUrls != "" {
				//fmt.Println(k,len(strings.Fields(RegistryUrls)))
				if len(strings.Fields(RegistryUrls)) >= k+1 {
					s.RegistryUrl = strings.Fields(RegistryUrls)[k]
				} else {
					s.RegistryUrl = ""
				}

			}

			if ControllerImgTags != "" {
				if len(strings.Fields(ControllerImgTags)) >= k+1 {
					s.ControllerImgTag = strings.Fields(ControllerImgTags)[k]
				} else {
					s.ControllerImgTag = ""
				}

			}

			deploy.Services = append(deploy.Services, s)
		}

		// generate install type form command line , 1 is new-install , 2 is install-from-backup
		InstallTypeString := GetInstallType(InstallType)
		// backup flag check ,if no input from command line , will read from config file
		InstallFromBackupFlagCheck(InstallTypeString)


		deploy.DeployerImgTag = DeployerImgTag
		deploy.DeepClean = DeepClean
		deploy.Registry = DeployerRegistry
		deploy.Depth = Depth
		deploy.Patch = Patch
		deploy.Thinpool = Thinpool
		deploy.Mode = Mode
		deploy.ItomSuiteSize = ItomSuiteSize
		deploy.AutomationTag = AutomationTag
		deploy.BackupSuite = BackupSuite
		deploy.InstallType = InstallTypeString
		deploy.BackupPackageDir = BackupPackageDir
		deploy.BackupPackageName = BackupPackageName
		deploy.Upgrade = Upgrade
		deploy.UpgradeImageTag	= UpgradeImageTag
		deploy.AdminPassword = AdminPassword
		deploy.CdfUrl = CdfUrl
		deploy.SuiteDataImageTag = SuiteDataImageTag
		deploy.MixMode = MixMode



	// test connection  EOF
		fmt.Println("===Connection Testing .....")
		all:=append(deploy.Workers,MasterIP)
		unreachable_ip:=[]string{}
		for _,ip :=range all{
			_,err :=net.Dial("tcp",ip + ":22")
			if err != nil{
				fmt.Printf("%s UNREACHABLE!! \n", ip)
				unreachable_ip=append(unreachable_ip,ip)
			}else {
			fmt.Printf("%s test connection ok\n", ip)
			}
		}
		if len(unreachable_ip) != 0{
				os.Exit(1)
			}
	// test connection EOF



		err := GenerateFileFromTemplate(filepath.FromSlash("ansible/host.tmpl"), filepath.FromSlash("ansible/host"), deploy)
		if err != nil {
			return err
		}

		curwd, _ := os.Getwd()

		os.Chdir("ansible")
		// if DeepClean == true , will  use the given vsphere username and password to revert vm, else use default writen in ansible
		if DeepClean == "true" {
			//generate vsphere user and password file  and load them as vars into ansible
			GenerateVsphereUserPass()
		}

		var cmdname string
		var cmdargs []string
		if runtime.GOOS == "windows" {
			cmdname = "cmd"
			cmdargs = append(cmdargs, "/C")
		} else {
			cmdname = "sh"
			cmdargs = append(cmdargs, "-c")
		}

		switch Depth {
		case 1:
			fmt.Println("Install the core platform, suite installer")
			cmdargs = append(cmdargs, "vagrant destroy -f && vagrant up --provision")
		case 3:
			fmt.Println("Install the core platform, suite installer and all services")
			cmdargs = append(cmdargs, "vagrant destroy -f && vagrant up --provision && vagrant destroy -f || vagrant destroy -f") //--provision
			//cmdargs = append(cmdargs, "cd") //--provision
		}
		err = ExecuteCommand(false, cmdname, cmdargs...)

		os.Chdir(curwd)

		if err != nil {
			return err
		}

		return nil
	},
}

func checkFlags() error {
	configFile := viper.ConfigFileUsed()
	PatchFlagCheck()
	ModeServiceCheck()
	//get cdf download url
	if CdfUrl == ""{
		CdfUrl = GetSuitePlatformLatestDownloadUrl()
		fmt.Printf("===CheckingParameters: Success. Using %s.\n", CdfUrl)
	}else {
		fmt.Printf("===CheckingParameters: Success. Using %s.\n", CdfUrl)
	}

	if DeactivatedService != ""{
		ModeDeactivatedServiceCheck()
	}
	fmt.Println("===H_MODE Service List:",HMode)
	fmt.Println("===X_MODE Service List:",XMode)
	Mode = strings.ToUpper(Mode)
	switch Mode {
	case "X_MODE":
		fmt.Printf("===CheckingParameters: Success. Using %s.\n", Mode)
	case "H_MODE":
		fmt.Printf("===CheckingParameters: Success. Using %s .\n", Mode)
	default:
		fmt.Println("===CheckingParameters: Failed. Please provide right Mode name,Either X_MODE or H_MODE.")
		os.Exit(1)
	}
	if MasterIP == "" {
		MasterIP = viper.GetString("commands.deploy.master")
		if MasterIP == "" {
			return errors.New("===CheckingParameters: Failed. Please provide IP address of Master.")
		} else {
			fmt.Printf("===CheckingParameters: Success. Using master: %s configured in %s.\n", MasterIP, configFile)
		}
	} else {
		fmt.Printf("===CheckingParameters: Success. master: %s.\n", MasterIP)
	}

	if NodeIP == "" {
		NodeIP = viper.GetString("commands.deploy.nodes")
		if NodeIP == "" {
			return errors.New("===CheckingParameters: Failed. Please provide IP address of worker nodes.")
		} else {
			fmt.Printf("===CheckingParameters: Success. Using worker nodes: %s configured in %s.\n", NodeIP, configFile)
		}
	} else {
		fmt.Printf("===CheckingParameters: Success. worker nodes: %s.\n", NodeIP)
	}

	if Services == "" {
		Services = viper.GetString("commands.deploy.services")
		if Services != "" {
			fmt.Printf("===CheckingParameters: Success. Using services: %s configured in %s.\n", Services, configFile)
		} else {
			fmt.Println("===CheckingParameters: Success. Full installation.")
		}
	} else {
		fmt.Printf("===CheckingParameters: Success. Services: %s.\n", Services)
	}

	if RegistryUrls == "" {
		RegistryUrls = viper.GetString("commands.deploy.registryUrls")
		if RegistryUrls != "" {
			fmt.Printf("===CheckingParameters: Success. Using RegistryUrls: %s configured in %s.\n", RegistryUrls, configFile)
		}
	} else {
		fmt.Printf("===CheckingParameters: Success. RegistryUrls: %s.\n", RegistryUrls)
	}

	if ControllerImgTags == "" {
		ControllerImgTags = viper.GetString("commands.deploy.controllerImgTags")
		if ControllerImgTags != "" {
			fmt.Printf("===CheckingParameters: Success. Using ControllerImgTags: %s configured in %s.\n", ControllerImgTags, configFile)
		}
	} else {
		fmt.Printf("===CheckingParameters: Success. ControllerImgTags: %s.\n", ControllerImgTags)
	}

	if DeployerRegistry == "" {
		DeployerRegistry = viper.GetString("commands.deploy.deployerRegistry")
		if DeployerRegistry == "" {
			DeployerRegistry = "shc-harbor-dev.hpeswlab.net/itsma"
			fmt.Printf("===CheckingParameters: Success. Using default deployerRegistry: %s.\n", DeployerRegistry)
		} else {
			fmt.Printf("===CheckingParameters: Success. Using deployerRegistry: %s configured in %s.\n", DeployerRegistry, configFile)
		}
	} else {
		fmt.Printf("===CheckingParameters: Success. deployerRegistry: %s.\n", DeployerRegistry)
	}

	if DeployerImgTag == "" {
		DeployerImgTag = viper.GetString("commands.deploy.deployerImgTag")
		if DeployerImgTag != "" {
			fmt.Printf("===CheckingParameters: Success. Using DeployerImgTag: %s configured in %s.\n", DeployerImgTag, configFile)
		} else {
			DeployerImgTag = Today("060102")
			fmt.Printf("===CheckingParameters: Success. Using the latest tag: %s.\n", DeployerImgTag)
		}

	} else {
		fmt.Printf("===CheckingParameters: Success. deployerImgTag: %s.\n", DeployerImgTag)
	}

	if ItomSuiteSize == "" {
		ItomSuiteSize = viper.GetString("commands.deploy.ItomSuiteSize")
		if ItomSuiteSize == "" {
			ItomSuiteSize = "demo"
			fmt.Printf("===CheckingParameters: Success. Using default ItomSuiteSize: %s.\n", ItomSuiteSize)
		} else {
			SuiteSizeFlagCheck()
		}
	} else {
		SuiteSizeFlagCheck()
		fmt.Printf("===CheckingParameters: Success. ItomSuiteSize: %s.\n", ItomSuiteSize)
	}

	if DeepClean == "true" &&
		(VsphereUser == "" || VspherePassword == "" || VsphereDataCenter == "" || VsphereFolder == "" || VsphereSnapShotName == "" || VCenterHostName == "") {
		return errors.New("===CheckingParameters: Failed. Please provide necessary info to revert VMs,VsphereUser,VspherePassword, etc.")
	}

	return nil
}

func PatchFlagCheck() error {
	//patched service and tags should not empty
	if Patch == true {
		if Services == "" {
			Services = viper.GetString("commands.deploy.services")
			if Services == "" {
				fmt.Println("===CheckingParameters: Failed. Patched Service name should not empty")
				os.Exit(1)
			}
		}
		if ControllerImgTags == "" {
			ControllerImgTags = viper.GetString("commands.deploy.controllerImgTags")
			if ControllerImgTags == "" {
				fmt.Println("===CheckingParameters: Failed. Patched service tag should not empty")
				os.Exit(1)
			}
		}
		// if patch more then one service , the number of service must be the same as the controller tags
		if len(strings.Fields(ControllerImgTags)) != len(strings.Fields(Services)){
			fmt.Println("===CheckingParameters: Failed. Patched service number should be equal service tags's ,one service one tag .")
			os.Exit(1)
		}

	}
	return nil
}

func init() {
	RootCmd.AddCommand(deployCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployCmd.Flags().StringVarP(&MasterIP, "masterip", "m", "", "(required) IP address of the master node")
	deployCmd.Flags().StringVarP(&NodeIP, "nodeip", "n", "", "(required) IP addresses of all worker nodes(separated by space)")
	deployCmd.Flags().StringVarP(&Services, "service", "s", "", "(optional) Itsma services to install(separated by space), all services will be installed by default")
	deployCmd.Flags().StringVarP(&RegistryUrls, "registryurls", "r", "", "(optional) one service, one registry(separated by space), using shc-harbor-dev.hpeswlab.net/itsma by default")
	deployCmd.Flags().StringVarP(&ControllerImgTags, "controllerimgtags", "", "", "(optional) one service, one controller image tag(separated by space), using stable by default")
	deployCmd.Flags().StringVarP(&DeployerImgTag, "deployertag", "t", "", "(optional) The docker image tag of itsma suite deployer")
	deployCmd.Flags().IntVarP(&Depth, "depth", "d", 3, "(optional) depth -d 1 or -d 3")
	deployCmd.Flags().StringVarP(&DeepClean, "deepclean", "c", "false", "(optional) revert vm to snapshot")
	deployCmd.Flags().StringVarP(&DeployerRegistry, "deployerregistry", "i", "", "(optional) image of installer")
	deployCmd.Flags().StringVarP(&VsphereUser, "VsphereUser", "U", "", "(optional,if -c true ,it's required) username of vsphere admin")
	deployCmd.Flags().StringVarP(&VspherePassword, "VspherePassword", "P", "", "(optional,if -c true ,it's required) password of vsphere admin")
	deployCmd.Flags().StringVarP(&VsphereDataCenter, "VsphereDataCenter", "D", "", "(optional,if -c true ,it's required) datacenter of vsphere vms")
	deployCmd.Flags().StringVarP(&VsphereFolder, "VsphereFolder", "F", "", "(optional,if -c true ,it's required) Define instance folder location.")
	deployCmd.Flags().StringVarP(&VsphereSnapShotName, "VsphereSnapShotName", "N", "", "(optional,if -c true ,it's required) snapshot name of vsphere vms")
	deployCmd.Flags().StringVarP(&VCenterHostName, "VCenterHostName", "H", "", "(optional,if -c true ,it's required) The hostname or IP address of the vSphere vCenter")
	deployCmd.Flags().StringVarP(&DeactivatedService, "ds", "", "", "(optional) DeactivatedService list ,the service will not run in suite")
	deployCmd.Flags().BoolVarP(&Patch, "patch", "", false, "(optional) Redeploy specified service in suite with new tag ")
	deployCmd.Flags().StringVarP(&Thinpool, "thinpool", "", "false", "(optional) The LVM thinpool device to use for the Docker devicemapper storage driver")
	deployCmd.Flags().StringVarP(&Mode, "Mode", "M", "H_MODE", "(optional)  Choose install mode first, either H_MODE or X_MODE")
	deployCmd.Flags().StringVarP(&ItomSuiteSize, "suite-size", "", "", "(optional)  Itom Suite Size optional choice 'demo' ,'xsmall', 'small', 'medium', 'large' default demo")
	deployCmd.Flags().StringVarP(&AutomationTag, "auto-test-tag", "T", "", "(optional) Automation test tag(separated by space)")
	deployCmd.Flags().StringVarP(&BackupSuite, "backup", "", "false", "(optional) Backup  suite service data ")
	deployCmd.Flags().IntVarP(&InstallType, "install-type", "", 0, "(optional) --install-type 1 or 2 , 1 is new install . 2 is install-from-backup. (Default 1)")
	deployCmd.Flags().StringVarP(&BackupPackageDir, "backup-pkg-dir", "", "", "(optional) backup package directory.It's located at ${global-pv}/backup")
	deployCmd.Flags().StringVarP(&BackupPackageName, "backup-pkg-name", "", "", "(optional) backup package name.It's located at ${global-pv}/backup/${backup-pkg-dir}")
	deployCmd.Flags().StringVarP(&Upgrade, "upgrade", "", "false", "(optional) upgrade suite,value can be true or false")
	deployCmd.Flags().StringVarP(&UpgradeImageTag, "upgrade-image-tag", "", Today("060102"), "(optional) The docker image tag of itsma suite update (itsma/itom-itsma-update)")
	deployCmd.Flags().StringVarP(&AdminPassword, "admin-password", "", "cloud", "(optional) Admin password of mng-portal")
	deployCmd.Flags().StringVarP(&CdfUrl, "cdfurl", "", "", "(optional) download url of CDF")
	deployCmd.Flags().StringVarP(&SuiteDataImageTag, "suite-data-tag", "", "dev", "(optional) tag of suite-data")
	deployCmd.Flags().StringVarP(&MixMode, "mix-mode", "", "false", "(optional) true or false , if true will apply external SM to suite. ")



}

func ExecuteCommand(quiet bool, name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// start the command after having set up the pipe
	if err = cmd.Start(); err != nil {
		return err
	}

	outerr := io.MultiReader(stdout, stderr)
	scanner := bufio.NewScanner(outerr)
	scanner.Buffer(make([]byte, bufio.MaxScanTokenSize*2), bufio.MaxScanTokenSize*2)
	if !quiet {
		fmt.Println("Command output message:")
		file_name := "../"+Today("2006-01-02_1504")+".log"
		log_file, err := os.Create(file_name)
		defer log_file.Close()
		debugLog := log.New(log_file,"",log.LstdFlags)
		if  err != nil {
			fmt.Printf("===CreateLogFile: Failed to create  file: %s. Error: %s\n", file_name, err.Error())
		}
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			vagrant_log := scanner.Text()

			debugLog.Println(vagrant_log)

			if strings.Index(vagrant_log, "site.retry") != -1 {
				fmt.Println("Cluster Info:\n\tMasterIp:", MasterIP, "\n\tNodeIP:  ", NodeIP)
				fmt.Println("\tInstall Result: FAILED")
				os.Exit(1)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Command output error: ", err)
	}
	fmt.Println("Cluster Info:\n\tMasterIp:", MasterIP, "\n\tNodeIP:  ", NodeIP)
	fmt.Println("\tInstall Result: SUCCESS")

	return err
}

func GenerateFileFromTemplate(templatefilename string, outputfilename string, data interface{}) error {
	tmpl := template.Must(template.New(filepath.Base(templatefilename)).ParseFiles(templatefilename))
	f, err := os.Create(outputfilename)
	if err != nil {
		return err
	}

	defer f.Close()

	err = tmpl.Execute(f, data)
	return err
}

func GenerateVsphereUserPass() {
	file := "revert-vm/vars/vsphere.yml"
	userYaml, err := os.Create(file)
	if err != nil {
		fmt.Printf("===GenerateVsphereUserPass: Failed. Failed to create file: %s.\n", file, " Rrror: ", err)
		return
	}

	defer userYaml.Close()

	content := fmt.Sprintf("vsphere_user: %s\nvsphere_pwd: %s\ndatacenter: %s\nfolder: %s\nsnapshot_name: %s\nhostname: %s\n",
		VsphereUser, VspherePassword, VsphereDataCenter, VsphereFolder, VsphereSnapShotName, VCenterHostName)
	if _, err = userYaml.WriteString(content); err != nil {
		fmt.Printf("===GenerateVsphereUserPass: Failed. Failed to write data: %s to file: %s. Error: %s\n", content, file, err.Error())
	}
}

// return the string date of today in format yymmdd

func Today(date_format string) string {
	//http://www.golangnote.com/topic/11.html
	t := time.Now().UTC()
	timestamp := t.Unix()
	_, offset := t.Zone()
	currentTime := time.Unix(timestamp+int64(offset), 0)

	return currentTime.Format(date_format)
}

func GetServiceNameFormJson(mode_name string) ([]string)  {
		//fmt.Println("===ListCommand. Getting services version from: ", Services_default_url)
	services_url := viper.GetString("commands.list.servicesVersionInfoLocation")
	if services_url == "" {
		services_url = Services_default_url
	}
	data, err := downloadFile(services_url)
	if err != nil {
		fmt.Println("===ListCommand. Failed. Download itsmaProducts.json error: ", err.Error())
		os.Exit(1)
	}
	var all_services ItsmaSuiteServices
	if err = json.Unmarshal(data, &all_services); err != nil {
		fmt.Println("===ListCommand. Failed. Json unmarshal error: ", err.Error())
		os.Exit(1)
	}
	all_service_list := []string{}
	ModeServiceList := []string{}
	for _, service := range all_services.Services {
		all_service_list = append(all_service_list,service.Name)
	}
	//fmt.Println(all_service_list)

	for _,service := range all_services.Services {
		if strings.Index(strings.Join(service.Mode,","), mode_name) != -1 || strings.Index(strings.Join(service.Mode,","), "INFRA") != -1{
			ModeServiceList = append(ModeServiceList,service.Name)
		}
	}
	//fmt.Println(ModeList)

	return ModeServiceList

}

func ModeServiceCheck() error {
	Activated_Services_list := strings.Fields(Services)
	if Mode == "X_MODE"{
		flag := 0
		for _,value :=range Activated_Services_list{
			for _,v :=range XMode{
				if value == v{
					flag += 1
				}
			}
		}
		if  len(Activated_Services_list) != 0 && len(Activated_Services_list) != flag{
			fmt.Printf("Some of your service don't belong to X_MODE : %s.\n", XMode)
			os.Exit(1)
		}
	}else if Mode == "H_MODE"{
		flag := 0
		for _,value :=range Activated_Services_list{
			for _,v :=range HMode{
				if value == v{
					flag += 1
				}
			}
		}
		if len(Activated_Services_list) != 0 &&  len(Activated_Services_list) != flag{
			fmt.Printf("Some of your service don't belong to H_MODE : %s.\n", HMode)
			os.Exit(1)
		}
	}else {
		errors.New("Please check your mode name, either X_MODE or H_MODE")
	}
	return nil
}

func ModeDeactivatedServiceCheck() error {
	DeactivatedServiceList := strings.Fields(DeactivatedService)
	if Mode == "X_MODE"{
		flag := 0
		for _,value :=range DeactivatedServiceList{
			for _,v :=range XMode{
				if value == v{
					flag += 1
				}
			}
		}
		if  len(DeactivatedServiceList) != 0 && len(DeactivatedServiceList) != flag{
			fmt.Printf("Some of your  deactivated service don't belong to X_MODE : %s.\n", XMode)
			os.Exit(1)
		}
	}else if Mode == "H_MODE"{
		flag := 0
		fmt.Println(DeactivatedServiceList)
		for _,value :=range DeactivatedServiceList{
			for _,v :=range HMode{
				if value == v{
					flag += 1
				}
			}
		}
		if len(DeactivatedServiceList) != 0 &&  len(DeactivatedServiceList) != flag{
			fmt.Printf("Some of your deactivated service don't belong to H_MODE : %s.\n", HMode)
			os.Exit(1)
		}
	}else {
		errors.New("Please check your mode name, either X_MODE or H_MODE")
	}
	return nil

}

func GetAllServiceList() []string {
	if Mode == "X_MODE"{
		return XMode
	}else {
		return HMode
	}
}

func SuiteSizeFlagCheck() error {
	// ItomSuiteSize jugement
	//suite size should  be 'demo' ,'xsmall', 'small', 'medium', 'large'
	suite_size_list := []string{"demo","xsmall","small","medium","large"}
	Flag := true
	for i :=0; i< len(suite_size_list); i++{
		if suite_size_list[i] == ItomSuiteSize{
			fmt.Printf("===CheckingParameters: Success. Using  ItomSuiteSize: %s.\n", ItomSuiteSize)
			Flag = false
			break
		}
	}
	if Flag{
		fmt.Printf("suite size must be one of demo,xsmall,small,medium,large , what you give is %s", ItomSuiteSize)
		os.Exit(1)
	}
	return nil
}

func GetInstallType(install_type_from_cmd_line int) string  {

	if install_type_from_cmd_line == 0{
		install_type_from_cmd_line = viper.GetInt("commands.deploy.InstallType")
	}


	var InstallTypeString string
	switch install_type_from_cmd_line {
	case 1:
		InstallTypeString = "new_install"
	case 2:
		InstallTypeString = "install_from_backup"
	default:
		InstallTypeString = "new_install"
	}
	fmt.Printf("===CheckingParameters: Success. Using InstallType: %s.\n",InstallTypeString)
	return InstallTypeString
}

func InstallFromBackupFlagCheck(install_type string)  {


	//read backup info from config file.

	if install_type == "install_from_backup" {
		if BackupPackageDir == "" {
			BackupPackageDir = viper.GetString("commands.deploy.BackupPackageDir")
			if BackupPackageDir != "" {
				fmt.Printf("===CheckingParameters: Success. Using BackupPackageDir: %s.\n", BackupPackageDir)
			} else {
				fmt.Println("===CheckingParameters: Failed . BackupPackageDir can't be null")
				os.Exit(1)
			}
		}else {
			fmt.Printf("===CheckingParameters: Success. Using BackupPackageDir: %s.\n", BackupPackageDir)
		}
		if BackupPackageName == "" {
			BackupPackageName = viper.GetString("commands.deploy.BackupPackageName")
			if BackupPackageName != "" {
				fmt.Printf("===CheckingParameters: Success. Using BackupPackageName: %s .\n", BackupPackageName)
			} else {
				fmt.Println("===CheckingParameters: Failed . BackupPackageName can't be null")
				os.Exit(1)
			}
		}else {
			fmt.Printf("===CheckingParameters: Success. Using BackupPackageName: %s .\n", BackupPackageName)
		}

// exit when specify backup info but not in install_form_back_up mode
	}else {
		if (BackupPackageName != "" || BackupPackageDir != "") && Upgrade != "true"{
			fmt.Println("===CheckingParameters: Failed. Only install_from_backup and Upgrade can specify BackUpPackageName and BackUpDir.")
			fmt.Println("Please remove the BackUpPackageName and BackUpDir or add --install-type 2 or add --upgrade true to continue.")
			os.Exit(1)
		}

	}



}

func GetSuitePlatformLatestDownloadUrl() string  {

	var cdf_metadata  string =  "http://shc-nexus-repo.hpeswlab.net:8080/repository/itsma-releases/com/hpe/shared/services/suite-platform/maven-metadata.xml"
	data,err1 := downloadFile(cdf_metadata)
	if err1 != nil {
		fmt.Printf("get cdf metadata url error: %s", err1)
		os.Exit(1)
	}
	v := MetadateOfNexus{}
	err := xml.Unmarshal([]byte(data), &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}

	//fmt.Println(v)

	var a string
	for _, i := range v.Versioning {
		a = fmt.Sprintf("http://shc-nexus-repo.hpeswlab.net:8080/"+
				"repository/itsma-releases/com/hpe/shared/services/suite-platform/%s/suite-platform-%s.zip", i.Release, i.Release)
		//a = "http://shc-nexus-repo.hpeswlab.net:8080/"+
			//"repository/itsma-releases/com/hpe/shared/services/suite-platform/"+i.Release +"/suite-platform-"+i.Release+".zip"
	}

	return  a

}

