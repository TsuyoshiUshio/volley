# Volley
[![Build Status](https://dev.azure.com/csedevops/volley/_apis/build/status/TsuyoshiUshio.volley?branchName=master)](https://dev.azure.com/csedevops/volley/_build/latest?definitionId=228&branchName=master)

Volley is a command line tool for create/destroy Stress testing environment and help you to run the senario and getting a log with a command line. 

The key features of Volley are: 

* **Provisioning/Deprovisioning:** Provision/Deprovision a JMeter cluster with Master/Slave nodes. Currently, we support deployment of a Virtual Machine cluster in Azure Azure. However, you can contribute to add other providers. 
* **Send a scenario:** Send a JMX file and data files to the Master.
* **Run:** Run Stress/Load testing.
* **Fetch a report:** Fetch a report from the Master. 
* **Server:** Works as an API Server on the JMeter master which facilitates the execution of Remote JMeter Server requests from the client machine. 

## Motivation 
Cloud based load testing was a cool service that enabled to create a worry free Stress Testing enviornment which also had good integration with existing CI tools. However, with it's depracation as a service, I decided to create a command line interface (CLI) that offers similar features which is easy to configure, use and is opensource. The tool was created using [GO Lang](https://golang.org/) which enables us to install/download and run the application with a single binary. 

All that is required is downloading a binary (specific to your platform) and adding it to the PATH.

## Getting Started

* [Getting Started Volley](doc/getting-started.md)

## Install

Go to [Release](https://github.com/TsuyoshiUshio/volley/releases) and find your platform binary. Download it and set the PATH to the binary. 

### get_volley.sh

The sample script, `get_volley.sh` shows how to install volley to a linux based machine. It downloads the latest binary and moves the binary into `/usr/bin` directory. The script
requires that the user has privileges to execute the `sudo` command. 

```bash
$ GET_VOLLEY_SCRIPT=get_volley.sh
$ curl -fsSL https://raw.githubusercontent.com/TsuyoshiUshio/volley/master/script/$get_volley.sh -o $GET_VOLLEY_SCRIPT
$ /bin/bash ${GET_VOLLEY_SCRIPT}
```

If you use Microsoft Windows and have [GitBash](https://git-scm.com/download/win), you can use this script as well. However, it will just download the latest version. Please make sure to move the binary to a suitable directory and add its location to your PATH variables. 

## CI Sample 

This is an sexample of a CI pipeline. This example uses [`Azure Pipelines`](https://azure.microsoft.com/en-us/services/devops/pipelines/)  however, we don't use any `Azure Pipeline` specific features. So you can easily to modify it to work on other CI systems. 

* [JMeter CI Sample](doc/ci-sample.md)

## Reference

### Provision (TODO)
Create a JMeter Cluster. (We currenlty support provisioning the infrastructure on Azure VMs. In the future, other providers will be added to support provisioning on `Azure Container Instances`). 

```
NAME:
   volley provision - Provision JMeter cluster

USAGE:
   volley provision [command options] [arguments...]

OPTIONS:
   --cluster-name value, -c value  Specify Cluster Name. Should be uniq.
   --slave value, -s value         Specify the number of slaves of JMeter cluster (default: 1)
   --help, -h                      show help (default: false)
```
#### Sample

```bash
$ volley provision --cluster-name tsushi --slave 2
```

### Server
The Volley API Server command can be used to start the volley API Server. The default server port is set to `38080`. The API Server provides several REST APIs which are usually used by volley subcommands. 

The API Server needs to run on the JMeter master server. The volley `provision` subcommand will start an instance for you automatically.

```
NAME:
   volley server - API Server for uploading/receiving files

USAGE:
   volley server [command options] [arguments...]

OPTIONS:
   --help, -h  show help (default: false)

```

#### REST API

* **Upload config**: `POST  /config`: Is used to upload JMeter Config files. It creates a new UUID and creates a folder `config/${UUID}` and moves the configuration into the newly created directory. `parameter`: none `body`: multipart files. `return`: {"id":"${config_id}"}
* **Run JMeter Job**: `POST  /job` : Starts the Job on the JMeter Server using a config config_id. It runs asynchronously and does not wait the execution to complete. It generates a job_id starts JMeter with the log under `job/${UUID}` directory.  `parameter`: none `body`: {"id":"${config_id}"} `return`:  {"job_id": "${job_id}", "config_id":"${config_id}"}
* **Update JMeter Property**: `POST /property` : Upload Remote IP hosts (slave) then it will update the default `jmeter.properties` file with the uploaded ip. You can find the modified `jmeter.property` file at the current directory of the `volley server`.  `parameter`: none `body`: {"remote_host_ips": ["${ip_address_1}", "${ip_address_2"}]} `return` : The same structure as the request body.
* **Check Job Status**: `GET /job/:job_id`: Gets the status of a job. status will be found `job/job_id/status.json`. This api returns the status. Possible value is `Running`, `Completed`, and `Failed`. `return`: {"status": "${status}"}.
* **Download Report**: `GET /asset/:job_id`: Download the result of the job execution. It downloads as a zip file. the fileName will be `${job_id}.zip`

#### Sample

```bash
$ volley server
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> github.com/TsuyoshiUshio/volley/pkg/command.(*ServerCommand).Start.func1 (3 handlers)
[GIN-debug] POST   /config                   --> github.com/TsuyoshiUshio/volley/pkg/controller.CreateNewConfig (3 handlers)
[GIN-debug] POST   /job                      --> github.com/TsuyoshiUshio/volley/pkg/controller.Start (3 handlers)
[GIN-debug] POST   /property                 --> github.com/TsuyoshiUshio/volley/pkg/controller.UpdateJMeterConfig (3 handlers)
[GIN-debug] GET    /job/:job_id              --> github.com/TsuyoshiUshio/volley/pkg/controller.StatusCheck (3 handlers)
[GIN-debug] GET    /asset/:job_id            --> github.com/TsuyoshiUshio/volley/pkg/controller.Download (3 handlers)
```

### Config
Upload JMeter Config file and csv file to the server to be ready to run the JMeter. On the serverside, this command create config/config_id directory and upload files into it. 

```
NAME:
   volley config - Upload jmx, csv files to the server. Return value is config-id.

USAGE:
   volley config [command options] [arguments...]

OPTIONS:
   --directory value, -d value  Specify directory that contains jmx and csv files that you want to upload
   --master value, -m value     Specify master ip address or domain name.
   --port value, -p value       Specify master port. 38080 by default (default: "38080")
   --help, -h                   show help (default: false)
```

#### Sample
Upload the jmx and csv files from the current directory. 
**NOTE:** Do not put two jmx fils in the target directory. It does not report an error. However, the volley server will execute only one jmx file on `run` subcommand. 

```bash
$ volley config --directory . --master http://localhost 
{"id":"c0234dff-1b18-11ea-bd0d-00155d7fe159"}
```

### Run
Runs JMeter with the configuration you uploaded on the server side. 

```
$ ./volley run --help
NAME:
   volley run - Run JMeter

USAGE:
   volley run [command options] [arguments...]

OPTIONS:
   --config-id value, -c value          Specify config-id that is created by config command.
   --master value, -m value             Specify master ip address or domain name.
   --port value, -p value               Specify master port. 38080 by default (default: "38080")
   --wait, -w                           Make this subcommand wait for completion (default: false)
   --timeout value, -t value            Specify the default timeout in minutes if you use --wait (-w) flag (default: 30)
   --output-type value, -o value        Specify the how to output the job_id. Possible value is 'stdout', 'file', 'both', if you choose file or both, it will output as file. The file name will respect outpus-filename flag (default: 
"stdout")
   --output-filename value, --of value  Specify the output filename when you specify --output-type flag (default: "job.json")
   --distributed-testing, -d            Enable distributed testing (default: false)
   --help, -h                           show help (default: false)
```

#### Sample

Non broking run mode.

```bash
$ volley run --config-id c0234dff-1b18-11ea-bd0d-00155d7fe159 --master http://localhost 
{"config_id":"c0234dff-1b18-11ea-bd0d-00155d7fe159","job_id":"4d93ea23-1b19-11ea-bd0d-00155d7fe159"}
```

Broking run mode. 

```bash
./volley run -c 929715c6-2454-11ea-a403-00249b32d3f7  -m http://localhost -w -o both -of myjob.json
{"config_id":"929715c6-2454-11ea-a403-00249b32d3f7","job_id":"a4568d8b-2454-11ea-a403-00249b32d3f7"}

Waiting for Job completion...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 5.0016115s second ...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 10.0037235s second ...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 15.0054866s second ...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 20.0075305s second ...
```

Run with distributed testing

To enable distributed testing, you need to post requests to `/properties` with RemoteIP address of the configured JMeter slaves nodes. 
Then execute the `volley run` command. If you want to enable `-d` flag, you need to insure that you start `volley --slave-server` on slave machine. For distributed testing, send csv files to the slave server/s before running the `jmeter` command with distiruted testing options specified. 

```bash
$ curl -X POST -H "Content-Type: application/json" -d '{"remote_host_ips":["10.0.0.4", "10.0.0.5"]}' http://${MASTER_IP}:38080/property
$ volley run -c 929715c6-2454-11ea-a403-00249b32d3f7  -m http://localhost -w -o both -of myjob.json -d
{"config_id":"929715c6-2454-11ea-a403-00249b32d3f7","job_id":"a4568d8b-2454-11ea-a403-00249b32d3f7"}

Waiting for Job completion...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 5.0016115s second ...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 10.0037235s second ...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 15.0054866s second ...
Polling status for JobID: a4568d8b-2454-11ea-a403-00249b32d3f7 Status: running at 20.0075305s second ...
```


### Log
Fetches the log file and report from the JMeter Server and extracts the content into current directory. The subdirectory will be the job_id of the job in question. 

```
NAME:
   volley log - fetch log of the JMeter run.

USAGE:
   volley log [command options] [arguments...]

OPTIONS:
   --job-id value, -j value  Specify job-id that run sub command returns.
   --master value, -m value  Specify master ip address or domain name.
   --port value, -p value    Specify master port. 38080 by default (default: "38080")
   --help, -h                show help (default: false)
```

#### Sample

```bash
$ volley log --job-id 4d93ea23-1b19-11ea-bd0d-00155d7fe159 --master http://localhost 
$ cd 4d93ea23-1b19-11ea-bd0d-00155d7fe159/
$ ls
report  status.json  stress.log
```

### Breaker

The Breaker breaks the build in the event that the test results for the executing job does not meet the success criteria. 
This sub command returns exit status 1 if it fails. 

```bash
NAME:
   volley breaker - Build

USAGE:
   volley breaker [command options] [arguments...]

OPTIONS:
   --log-file value, -l value  File path for JMeter execution log file.
   --config value, -c value    Config file path of success_criteria (default: "success_criteria.json")
   --help, -h                  show help (default: false)
```

#### Sample

_success_criteria.json_

```
{
    "criteria":"average_time_error_on_rps",
    "Parameters":{
        "avg_latency":10000,
        "error_ratio":10,
        "rps":250
    }
}
```

Execute volley using the sample log file and success_criteria.json.
The log file's result resembles the following:

> 2019/12/17 00:29:11 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 % Upto 250 Request Per Second.


This sample demonstrates the console output in the case of a failure 

```bash
$ volley breaker -l pkg/model/test-data/success-criteria/avg-time-error-on-rps/stress.log -c pkg/model/test-data/success-criteria/config/success_criteria.json
2019/12/17 02:26:25 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 %
2019/12/17 02:26:25 Request Per Second Up to: 250, Target Average Letency Less than: 10000, Target Error Ratio Less than: 10
2019/12/17 02:26:25 Validation failed.
$ echo $?
1
```
If you change the `succss_criteria.json`'s avarage latency from 10000 to 20000, it will succeed. 

```bash
$ ./volley breaker -l pkg/model/test-data/success-criteria/avg-time-error-on-rps/stress.log -c pkg/model/test-data/success-criteria/config/success_c
riteria.json
2019/12/17 02:40:57 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 %
2019/12/17 02:40:57 Request Per Second Up to: 250, Target Average Letency Less than: 20000, Target Error Ratio Less than: 10
2019/12/17 02:40:57 Validation succeed.

$ echo $?
0
```
### Slave Server
If you want to have Distributed Testing, you need to run this server on slave machines. This server API accepts requests using  `POST /csv` on port `38081`. It is able to receive multipart csv files and saves them under the `csv` directory. (JMeter Slave servers have some restrictions in that it requires copying csv files in the path that is recorded in the PATH variable created for jmeter.  

```bash
NAME:
   volley slave-server - API Server for receive csv files on slave server

USAGE:
   volley slave-server [command options] [arguments...]

OPTIONS:
   --help, -h  show help (default: false)
```

#### Sample

```bash
$ volley ss 
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> github.com/TsuyoshiUshio/volley/pkg/command.(*SlaveServerCommand).Start.func1 (3 handlers)
[GIN-debug] POST   /csv                      --> github.com/TsuyoshiUshio/volley/pkg/controller.UploadCSV (3 handlers)
```

### Destroy (TODO)

### HELP

## TODO 

We are planning these feature as TODO  items.

* Provision/Deprovision cluster
* Wait option for Run for the CI. Wait until the JMeter execution completed on the server side. 
* Ability to omit master parameter. It should automatically fetch configuration from the cluster deployment. 
* Ability to omit config_id parameter. It should automatically be set once the config file is uploaded. 
* Secure connection between Server/Client.

## Developing Volley

### Build Project

_linux and mac_

```bash
./script/build.sh
```

_windows_

```cmd
script\build.bat
```

### Build For All Platforms

If you need to build multi platform build, execute this command. 
Currently only support bash. It works on linux, mac, and windows(git bash)

```bash
./script/build_all_platform.sh
```

### Run Unit Test

_linux and mac_

```bash
./script/test.sh
```

_windows_

```cmd
script\test.bat
```
