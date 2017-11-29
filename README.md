# **Suite-kit**

## Linux OR Windows ?

We both support Windows and Linux . Your CPU has to support VT . For linux  [How to check.][1], for Windows enable VT in your BIOS


----------

1.Prepare
===
### 1.1 install VirtualBox

please follow https://www.virtualbox.org/wiki/Linux_Downloads

### 1.2 install Vagrant 

https://www.vagrantup.com/downloads.html

### 1.3 You can use suitectl now


----------


2.USAGE
===

download suitectl from https://github.houston.softwaregrp.net/SMA-RnD/suite-kit/releases 

```bash
$ ./suitectl --help
suitectl is a tool wrote in golang, for init and deploy itsma services, and core platform as well.

Usage:
  suitectl [command]

Available Commands:
  deploy      deploy itsma suite service
  help        Help about any command
  init        init a itsma service
  list        list all available services in itsma suite.


Flags:
      --config string   config file (default is $HOME/.suitectl.yaml)
  -h, --help            help for suitectl
  -t, --toggle          Help message for toggle

Use "suitectl [command] --help" for more information about a command.

$ ./suitectl deploy --help
deploy one or more itsma suite services on target environment. For example:
  ./suitectl deploy -m 16.155.199.110 -n "16.155.199.112 16.155.199.113" 

----------

Usage:
  suitectl deploy [flags]

Flags:
  -M, --Mode string                  (optional)  Choose install mode first, either H_MODE or X_MODE (default "H_MODE")
  -H, --VCenterHostName string       (optional,if -c true ,it's required) The hostname or IP address of the vSphere vCenter
  -D, --VsphereDataCenter string     (optional,if -c true ,it's required) datacenter of vsphere vms
  -F, --VsphereFolder string         (optional,if -c true ,it's required) Define instance folder location.
  -P, --VspherePassword string       (optional,if -c true ,it's required) password of vsphere admin
  -N, --VsphereSnapShotName string   (optional,if -c true ,it's required) snapshot name of vsphere vms
  -U, --VsphereUser string           (optional,if -c true ,it's required) username of vsphere admin
      --admin-password string        (optional) Admin password of mng-portal (default "cloud")
  -T, --auto-test-tag string         (optional) Automation test tag(separated by space)
      --backup string                (optional) Backup  suite service data  (default "false")
      --backup-pkg-dir string        (optional) backup package directory.It's located at ${global-pv}/backup
      --backup-pkg-name string       (optional) backup package name.It's located at ${global-pv}/backup/${backup-pkg-dir}
      --cdfurl string                (optional) download url of CDF (default "http://shc-nexus-repo.hpeswlab.net:8080/repository/itsma-releases/com/hpe/shared/services/suite-platform/2017.06.00156/suite-platform-2017.06.00156.zip")
      --controllerimgtags string     (optional) one service, one controller image tag(separated by space), using stable by default
  -c, --deepclean string             (optional) revert vm to snapshot (default "false")
  -i, --deployerregistry string      (optional) image of installer
  -t, --deployertag string           (optional) The docker image tag of itsma suite deployer
  -d, --depth int                    (optional) depth -d 1 or -d 3 (default 3)
      --ds string                    (optional) DeactivatedService list ,the service will not run in suite
  -h, --help                         help for deploy
      --install-type int             (optional) --install-type 1 or 2 , 1 is new install . 2 is install-from-backup. (Default 1)
  -m, --masterip string              (required) IP address of the master node
  -n, --nodeip string                (required) IP addresses of all worker nodes(separated by space)
      --patch                        (optional) Redeploy specified service in suite with new tag 
  -r, --registryurls string          (optional) one service, one registry(separated by space), using shc-harbor-dev.hpeswlab.net/itsma by default
  -s, --service string               (optional) Itsma services to install(separated by space), all services will be installed by default
      --suite-data-tag string        (optional) tag of suite-data (default "dev")
      --suite-size string            (optional)  Itom Suite Size optional choice 'demo' ,'xsmall', 'small', 'medium', 'large' default demo
      --thinpool string              (optional) The LVM thinpool device to use for the Docker devicemapper storage driver (default "false")
      --upgrade string               (optional) upgrade suite,value can be true or false (default "false")
      --upgrade-image-tag string     (optional) The docker image tag of itsma suite update (itsma/itom-itsma-update) (default "yymmdd")

Global Flags:
      --config string   config file 
```

## 3. EXAMPLE

```bash
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
  
  11. Do not want to use Thinpool for the Docker devicemapper storage driver(use Thinpool as default)
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
  19. Mix mode (for SRE team only!)
  19.1 deploy mode P1 and wait until Install finished.
   ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip" --ds "itom-cmdb itom-sm" 
    ./suitectl deploy -m masterip -n "node1_ip node2_ip node3_ip"  --mix-mode true
```
# Compatibility
Refer to [Compatibility.md][2]

## 4. Change Log
### v1.0.0.23
```bash
use latest CDF build as default
```

### v1.0.0.22
```bash
Use Chrony instead of ntpdate to sync time
Update ansible version  
Update vagrant base box image
```

### v1.0.0.21
```bash
Update CDF to 157
```

### v1.0.0.20
```bash
add arg --cdfurl string                (optional) download url of CDF (default "http://shc-nexus-repo.hpeswlab.net:8080/repository/itsma-releases/com/hpe/shared/services/suite-platform/2017.06.00156/suite-platform-2017.06.00156.zip")
add arg --suite-data-tag string        (optional) tag of suite-data (default "dev")
update CDF to 156
```

### v1.0.0.19
```bash
Update for new upgrade code merged in deployer 
##Suitekit will backup suite automatically , and use the latest backup dir and backup package name unless you specify them with arg  --backup-pkg-dir --backup-pkg-name  refer example 15 for more details/
```

### v1.0.0.18

```bash
Update for new code merged in deployer
```

### v1.0.0.17
```
update CDF to build 155

```

### v1.0.0.15
```bash
update CDF to 153
suite-data use "dev" as default tag
```
### v1.0.0.14
```bash
Update Mode arg to upper case for defect 25378 (H_Mode --> H_MODE, X_Mode -->X_MODE), no effect on use.
support aws harbor
AWS download ojdbc8.jar for S3
Update CDF to build 152
```

### v1.0.0.12
```bash
fix ojdbc download url error
add arg --admin-password string        (optional) Admin password of mng-portal (default "cloud") Example 16
update CDF to 150
update suite version from 2017.10 to 2017.11

```

### v1.0.0.11
```bash
rename ojdbc8.jar to ojdbc.jar
update suite version to 2017.10
add namespace check before upgrade suite
```

### v1.0.0.10
```bash
update resource limit for AGM defect 24790 
update CDF to build 147
reload backup itsma_suite_metadata.json as install body
support suite upgrade
      --upgrade string               (optional) upgrade suite,value can be true or false (default "false")
      --upgrade-image-tag string     (optional) The docker image tag of itsma suite update (itsma/itom-itsma-update) (default yymmdd --> 170922)
update ojdbc6.jar to ojdbc8.jar
add alias tool https://github.houston.softwaregrp.net/SMA-RnD/suite-tools/blob/master/alias-tool/alias.sh (supported by Liu Bin bin.liu8@hpe.com)

```

### v1.0.0.9
```bash
update CDF kubernetes to 1.6.9

```

### v1.0.0.8
```bash
fix aws ECR token expire problem
fix backup  bug
default deploy X_MODE on AWS
add bo_admin password and sysadmin password
update CDF to 141
add SM Oracle clinet to global-pv 
```

### v1.0.0.7
```bash
support backup suite
      --backup string                (optional) Backup  suite service data  (default "false")
support install from backup      
    --install-type int             (optional) --install-type 1 or 2 , 1 is new install . 2 is install-from-backup. (Default 1)  
      --backup-pkg-dir string        (optional) backup package directory.It's located at ${global-pv}/backup
      --backup-pkg-name string       (optional) backup package name.It's located at ${global-pv}/backup/${backup-pkg-dir}

add more args in config file
update kube-proxy resource limit

```
### v1.0.0.6
```bash
update install body for new deployer code
```

### v1.0.0.5
```bash
add arg -T, --auto-test-tag string     (optional) Automation test tag(separated by space)

make master node as a worker node

```

### v1.0.0.4

```bash
update cdf to 139

```

### v1.0.0.3

```bash
add arg --suite-size string   (optional)  Itom Suite Size optional choice 'demo' ,'xsmall', 'small', 'medium', 'large' default demo
disable thinpool as default
delete image list of itsma_suitefeatures.2017.10.json
```

### v1.0.0.2
```bash
update H and X mode service list
add mode column in the output of  suitectl list
add ojdbc6.jar to /var/vols/itom/core/suite-install/itsma/output/  and /var/vols/itom/itsma/itsma-itsma-global/jdbc
```

### v1.0.0.1

```bash
add arg --ds , DeactivatedService list ,the service will not run in suite (example 9)
add arg --patch Only redeploy specified service in suite with new tag. (example 10)
            --patch do not delete the volume data
add arg -M  Install mode selection, either H_MODE or X_MODE(use H_MODE as default)(example 9)
            H_MODE service: itom-sm","itom-cmdb","itom-service-portal","itom-chat","itom-smartanalytics","itom-service-portal-ui","itom-xservices","itom-xservices-infra","itom-landing-page
            X_MODE service: itom-cmdb","itom-xruntime","itom-xruntime-infra","itom-backoffice (TBD)"
add arg --thinpool  The LVM thinpool device to use for the Docker devicemapper storage driver  (default "true")(example 11)
update profile to demo from default
```

### v0.0.17
```bash
update CDF to 136
add checksum of CDF zip
```


### v0.0.16
```bash
update sysadmin password to Admin_1234
```

### v0.0.15
```bash
delete image list in itsma_suitefeatures.json of suite-data
update ntpdate command in crontab
use uuid to revert vm instead of vm name ,no longer limited by vm folder
upgrade CDF to 129
```

### v0.0.14
```bash
upgrade CDF to build 120

```


### v0.0.13
```bash
upgrade CDF to build 105
fix hangs on getting fact, for agm 22225
support mutiple PV in suite (global-volume, db-volume, smartanalytics-volume)
support read args form configuration file
add node ip connection test
write suitekit deploy log to file
add CDF version check
```
### v0.0.12
```bash
upgrade CDF to version 2017.06
for now please  specify the deployer tag  to "upgradeCDF2.3"  with  "-t upgradeCDF2.3"
```
### v0.0.11
```bash
update default vsphere user and password
list service stable version from new git repo
use pre-packaged VritualBox instead of download on line

```

### v0.0.10
```bash
add fqdn check , fail fast if your fqdn is not in domain hpeswlab.net, for agm 22289
add wait_time when delete itsma namespace , for agm 22342

```

### v0.0.9
```bash
update unite test of command 'init'
add examples for command 'deploy'
add command 'list' to list all itsma services
use current date yymmdd as deployer default tag

```


### v0.0.8

**support VM in other POOL to install suite.**

```bash
 
  if your vms are in /SHCITSMA pool ,you don't need to give all the infomations , run commadn below . use -c default to revert vms
  ./suitectl deploy -m 15.119.87.12 -n "15.119.87.124 15.119.87.125"  -c default 


  if you have your own pool ,please use -c true and give all the args list below to revert vms
./suitectl deploy -m 15.119.87.12 -n "15.119.87.124 15.119.87.125"  -c true -U "ASIAPACIFIC\tom" -P "password" -D "SHCITSMA" -F "/SHCITSMA" -N init -H selvc.hpeswlab.net

  -H, --VCenterHostName string       (optional,if -c true ,it's required) The hostname or IP address of the vSphere vCenter
  -D, --VsphereDataCenter string     (optional,if -c true ,it's required) datacenter of vsphere vms
  -F, --VsphereFolder string         (optional,if -c true ,it's required) Define instance folder location.
  -P, --VspherePassword string       (optional,if -c true ,it's required) password of vsphere admin
  -N, --VsphereSnapShotName string   (optional,if -c true ,it's required) snapshot name of vsphere vms
  -U, --VsphereUser string           (optional,if -c true ,it's required) username of vsphere admin


```

### v0.0.7   fix bugs




### v0.0.6
```bash
add -d, --depth int
  if -d 1 , suitctl will only install CorePlatform ,default 3 , install CorePlatform and all services

```
### v0.0.4 -- 0.0.5  fix bugs

### v0.0.3

```bash
add -i arg ,
 -i, --image-name string    (optional) image of installer (default "shc-harbor-dev.hpeswlab.net/itsma/itom-itsma-installer")

add retrue code for jenkins 

```


### v0.0.2

```bash
add -c arg ,
 -c, --deepclean string     (optional) revert vm to snapshot (default "false")

only support VM from datacenter SHCITSMA in ShangHai

```


  [1]: https://www.howtogeek.com/howto/linux/linux-tip-how-to-tell-if-your-processor-supports-vt/
  [2]: https://github.houston.softwaregrp.net/SMA-RnD/suite-kit/blob/master/Compatibility.md
