# Volley
[![Build Status](https://dev.azure.com/csedevops/volley/_apis/build/status/TsuyoshiUshio.volley?branchName=master)](https://dev.azure.com/csedevops/volley/_build/latest?definitionId=228&branchName=master)

Volley is a command line tool for create/destroy Stress testing environment and help you to run the senario and getting a log with a command line. 

The key feature of Volley are: 

* **Privsioning/Deprovisioning:** Provision/Deplovision Master/Slave cluster of JMeter. Currently, we support Virtual Machine for Azure. However, you can contribute to add other providers. 
* **Send a scenario:** Send JMX file and data files to the Master.
* **Run:** Run Stress/Load testing.
* **Fetch a report:** Fetch a report from the Master. 
* **Server:** Worked as an API Server on the JMeter master side. Execute Remote JMeter Server from the client request. 

## Motivation 
Cloud based load testing was a cool service that enable us not to worry about the Stress Testing enviornment and has a good integration with CI tools. However, it was depricated. I'd like to create a command that does the same thing for us. I create this tool with go lang that is enable us to 
install/download with a single binary. All you need is just download a binary of your platform and add it to the PATH.

## Install

Go to [Release](https://github.com/TsuyoshiUshio/volley/releases) and find your platform binary. Download it and set the PATH to the binary.

## Refernce

### Provision (TODO)
Create a JMeter Cluster. Currently we support Azure VM. In the future, we can add other providers. 

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
#### sample

```bash
$ volley provision --cluster-name tsushi --slave 2
```

### Server
Start the volley API Server. The server port is `38080` by default. It provide several REST API. It is usally used from volley subcommand. 
You need to start server on the JMeter master server. `provision` subcommand will do it automatically.

```
NAME:
   volley server - API Server for uploading/receiving files

USAGE:
   volley server [command options] [arguments...]

OPTIONS:
   --help, -h  show help (default: false)

```

#### REST API

* **Upload config**: `POST  /config`: Upload JMeter Config files. It create a new UUID and create a folder `config/${UUID}` then put the files under the directory.  `parameter`: none `body`: multipart files. `return`: {"id":"${config_id}"}
* **Run JMeter Job**: `POST  /job` : Start Job that JMeter execution using config_id. It doesn't wait whole execution. It generate job_id start JMeter with the log under `job/${UUID}` directory.  `parameter`: none `body`: {"id":"${config_id}"} `return`:  {"job_id": "${job_id}", "config_id":"${config_id}"}
* **Update JMeter Property**: `POST /property` : Upload Remote IP hosts (slave) then it will update the default `jmeter.properties` file with the uploaded ip. You can find the modified `jmeter.property` file at the current directory of the `volley server`.  `parameter`: none `body`: {"remote_host_ips": ["${ip_address_1}", "${ip_address_2"}]} `return` : The same structure as the request body.
* **Check Job Status**: `GET /job/:job_id`: Get the status of a job. status will be found `job/job_id/status.json`. This api return the status. Possible value is `Running`, `Completed`, and `Failed`. `return`: {"status": "${status}"}.
* **Download Report**: `GET /asset/:job_id`: Download the result of the job execution. It downloads as a zip file. the fileName will be `${job_id}.zip`

#### sample

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

#### sample
Upload the jmx and csv file from the current directory. 
**NOTE:** Don't put two jmx file on the target directory. It doesn't error. However, volley server execute only one jmx file on `run` subcommand. 

```bash
$ volley config --directory . --master http://localhost 
{"id":"c0234dff-1b18-11ea-bd0d-00155d7fe159"}
```

### Run
Run the JMeter with the configuration you uploaded on the server side. 

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

#### sample

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

For enabling distributed testing, you need to post request to `/properties` with RemoteIP address of slaves. 
Then execute `volley run` command. 

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
Fetch the log file and report from Server side. It is extracted on the current directory with the sub directory with job_id. 

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

#### sample

```bash
$ volley log --job-id 4d93ea23-1b19-11ea-bd0d-00155d7fe159 --master http://localhost 
$ cd 4d93ea23-1b19-11ea-bd0d-00155d7fe159/
$ ls
report  status.json  stress.log
```

### Breaker

Break the build once if the result of the job execution doesn't meet the success criteria. 
This sub command resturn exit status 1 if it fails. 

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

#### sample

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

Execute the volley using sample log file and success_criteria.json.
The log file's result is 

> 2019/12/17 00:29:11 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 % Upto 250 Request Per Second.


This sample fails 

```bash
$ volley breaker -l pkg/model/test-data/success-criteria/avg-time-error-on-rps/stress.log -c pkg/model/test-data/success-criteria/config/success_criteria.json
2019/12/17 02:26:25 TotalRequest: 3916, Average Latency: 11593, ErrorRatio: 1 %
2019/12/17 02:26:25 Request Per Second Up to: 250, Target Average Letency Less than: 10000, Target Error Ratio Less than: 10
2019/12/17 02:26:25 Validation failed.
$ echo $?
1
```
If you change the `succss_criteria.json`'s avarage latency from 10000 to 20000, it will success. 

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
If you want to have Distributed Testing, you need to run this server on slave machines. This server accept API. Accept `POST /csv` on `38081` port. It receives multipart csv files and saves under `csv` directory. JMeter Slave has some restriction that requires to put csv files on PATH that enable for jmeter.  

```bash
NAME:
   volley slave-server - API Server for receive csv files on slave server

USAGE:
   volley slave-server [command options] [arguments...]

OPTIONS:
   --help, -h  show help (default: false)
```

#### sample

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

We are planning these feature as TODO.

* Provision/Deprovision cluster
* Wait option for Run for the CI. Wait until the JMeter execution finished on the server side. 
* Able to omit master parameter. It automatically fetched from the culster deployment. 
* Able to omit config_id parameter. It automatically set once upload the config. 
* Secure connection with Server/Client.

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

### Build For All Platform

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
