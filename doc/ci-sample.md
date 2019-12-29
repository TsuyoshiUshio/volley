# JMeter CI Sample

This document explain the sample of CI Pipeline sample. 
Before starting this pipeline, please deploy JMeter Cluster in advance. We provide ARM template for deploying the cluster. Please refer [Getting Started Volley](getting-started.md)

## Azure Pipeline 

Example for Azure Pipeline. For more details of the configration, please refer [Getting Started Volley](getting-started.md) 

This example assume that there is `JmxSample` on the repo and it includes `jmx` file and `csv` file for the `jmx` file data. 

Also, it assume `success_criteria.json` that specify the threshold of the JMeter result if the pipeline will fail or not' 


```yaml

trigger:
- master

pool:
  vmImage: 'ubuntu-latest'

variables:
  masterIP: 'YOUR_MASTER_IP_ADDRESS_HERE'
  slaveIP1: 'YOUR_SLAVE_IP_ADDRESS_1_HERE'
  slaveIP2: 'YOUR_SLAVE_IP_ADDRESS_2_HERE'

steps:
- task: CmdLine@2
  displayName: 'install volley'
  inputs:
    script: |
      GET_VOLLEY_SCRIPT=get_volley.sh
      curl -fsSL https://raw.githubusercontent.com/TsuyoshiUshio/volley/master/script/get_volley.sh -o $GET_VOLLEY_SCRIPT
      /bin/bash ${GET_VOLLEY_SCRIPT}
      which volley

- task: CmdLine@2
  displayName: 'run stress testing'
  inputs:
    script: |
      MASTER_IP=$(masterIP)
      SLAVE_IP_1=$(slaveIP1)
      SLAVE_IP_2=$(slaveIP2)
      
      cd JmxSample

      CONFIG_ID=`volley config -d . -m http://${MASTER_IP} | jq .id | xargs`
      echo $CONFIG_ID
      
      curl -X POST -H "Content-Type: application/json" -d "{\"remote_host_ips\":[\"${SLAVE_IP_1}\", \"${SLAVE_IP_2}\"]}" http://${MASTER_IP}:38080/property
      
      volley run -c $CONFIG_ID -m http://${MASTER_IP} -w -o both -of myjob.json -d
      
      JOB_ID=`cat myjob.json | jq .job_id | xargs`
      volley log -j $JOB_ID -m http://${MASTER_IP}
      
      cd $JOB_ID
      ls -l

      volley breaker -l ./stress.log -c ../../success_criteria.json
      status=$?
      if [ $status -ne 0 ]; then
        exit $status
      fi 
      cd ..
      cd ..

- task: PublishBuildArtifacts@1
  displayName: 'publish JMeter logs'
  condition: succeededOrFailed()
  inputs:
    PathtoPublish: '$(Build.SourcesDirectory)/JmxSample'
    ArtifactName: 'jmeter'
    publishLocation: 'Container'
```

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