# Getting Started

Welcome to volley!

Volley help to execute JMeter testing with the JMeter master/slave cluster and fetch the log and report. Also, It help us to execute JMeter Stress testing on CI. 

This getting started guide introduces how to use volley manually. I asusme you use bash based terminal for linux and mac. For windows, Please download [GitBash](https://gitforwindows.org/).

## Deploy JMeter Cluster 

You can deploy JMeter Cluster with ARM template. In the near future, it will be wrapped with volley sub command. Currenlty please deploy via Azure CLI. If you don't have Azure CLI, please install it from [here](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest). 

### Create a Resource Group
The sample resrouce group name is "JMeterCluster".

```bash
$ az group create -n "JMeterCluster" -l westus
```

### Deploy the JMeter Cluster
This ARM template will create JMeter Cluster with JMeter Server, volley server on master/slave. If you encounter the deployment issue, please execute this command twice (this is known issue.) If you want to increse/decrease the number of slave, please use `slaveCount` parameter.

```bash
$ cd template/template/azure/vm
$ az group deployment create --parameters sshKeyData="<YOUR_SSH_PUBLIC_KEY_HERE>" -g JMeterCluster --template-file azuredeploy.json
```

### Install Volley locally
Install volley on your local machine. If you are using Linux(Include WSL) or Mac, you can use this script. It automatically download volley binary and put on `/usr/bin`. If you are using windows, you can use this script with GitBash. Unfrotunatelly, on windows, it doesn't automatically set binary to the PATH. Please get the binary on the current directry and copy it to a directory that is on the PATH environment variables. 
Other way to install volley is download binary from [Release]() then put the binary into a directory that is on the PATH envrionment variables. 

```bash
$ GET_VOLLEY_SCRIPT=get_volley.sh
$ curl -fsSL https://raw.githubusercontent.com/TsuyoshiUshio/volley/master/script/$get_volley.sh -o $GET_VOLLEY_SCRIPT
$ /bin/bash ${GET_VOLLEY_SCRIPT}
```

### Configure volley

#### Upload JMeter files 
volley help you to upload JMX file and csv files to the JMeter Cluster. Get the JMeter Master IP address. 

You can find JMeter Master IP Address by this command. Or you can find the IP Address on the Azure Portal. `JMeterCluster` is the resoruce group of the JMeter Cluster. 

```bash
$ MASTER_IP=`az vm list-ip-addresses -g JMeterCluster -n jmeterVM1 | jq -r  '.[0].virtualMachine.network.publicIpAddresses[0].ipAddress'`
```

For distributed testing, you need to give them the slave private IP address to the master server. Or you can find the private IP Address on the Azure Portal. 

```bash
$ SLAVE_IP_1=`az vm list-ip-addresses -g JMeterCluster -n jmeterVM2 | jq -r  '.[0].virtualMachine.network.privateIpAddresses[0]'`
$ SLAVE_IP_2=`az vm list-ip-addresses -g JMeterCluster -n jmeterVM2 | jq -r  '.[0].virtualMachine.network.privateIpAddresses[0]'`
```

Call REST API to configure Slave Servers. In the near feature, it will includes into the `provision` sub command. 

```bash
$ curl -X POST -H "Content-Type: application/json" -d "{\"remote_host_ips\":[\"${SLAVE_IP_1}\", \"${SLAVE_IP_2}\"]}" http://${MASTER_IP}:38080/property      
```

Then go to the directory that you have JMX file and csv files and upload it. `config` sub command upload 'jmx' and 'csv' file to the cluster. The jmx file should be one. csv file is eventually transfer to the slave server on the current directory that volley slave server working. JMeter Server requires csv file on each machine. That is limitation of the JMeter distributed testing. 

This command returns `Config ID` to the client. `-d` means `directory` that includes `jmx` and `csv` file. `-m` means `master` server's URL. 

```bash 
$ CONFIG_ID=`volley config -d . -m http://${MASTER_IP} | jq -r .id`
```  

Now ready to run stress/load testing. 

### Run Stress/Load Testing

Start stress/load testing using this command. This command execute JMeter distributed testing on JMeter Cluster side and wait until the execution is finished. It generate, `job.json` that include `job_id`.

`-c` means `config id`. `-m` means `master` server's URL. `-o` means `output type`. You can choose the output(job_id) to the file or stdout. You can choose `stdout`, `file`, `both`. If you choose `file` or `both`, the `-d` is `distributed testing`, and `-w` is `wait` until the test is finished. If you don't specify the `-w` flag, this command returns immediately. 

```bash
$ volley run -c $CONFIG_ID -m http://${MASTER_IP} -o both -d -w
```

### Get Report

To get the log file and rerpot of the distributed testing, you can use `log` sub command. `run` sub command generate `job.json` to the current directory. `log` sub command use it and donwload the artifact on the `$JOB_ID` directory. 

```bash
$ JOB_ID=`cat job.json | jq -r .job_id`
$ ./volley log -j $JOB_ID -m http://${MASTER_IP}
$ cd $JOB_ID
$ ls
report  status.json  stress.log 
```

Currently, we need to pass `MASTER_IP`, `CONFIG_ID`, and `JOB_ID` however, we have a plan to remove these parameter with default configration file. 

### (Optional) Build Breaker for CI

You might want to make the script fail if it doesn't reach the threshold. 
Create this json file. This config file means, `Average latency is less than `30000`, Error ratio is less than `30` at the Request Per Second less than `250`. If the actual value exceed this codition, the script will fail. 

_success_criteria.json_

```json
{
    "criteria":"average_time_error_on_rps",
    "Parameters":{
        "avg_latency":30000,
        "error_ratio":30,
        "rps":250
    }
}
```

_pipeline.sh_

```bash
      volley breaker -l ./stress.log -c ../../success_criteria.json
      status=$?
      if [ $status -ne 0 ]; then
        exit $status
      fi 
```

