# JMeter Docker

JMeter Docker images for Stress/Load testing. 
This docker image wraps the jmeter commands required to perform stress testing. 

_linux_

```bash
$ docker run -d -P --name master -v /home/ushio/Codes/DevSecOps/EpiServer/StressTesting/Temp/:/jmeter_log tsuyoshiushio/jmeter jmeter -n -t /jmeter_log/MessageApi.jmx -l /jmeter_log/current2.jtl -e -o /jmeter_log/report2 -Jthreads=100 -Jduration=60 -Jport=6400
```

_windows_

On windows, WSL and Windows both using Docker for Windows. You can't use localhost as a host. For more details you can refer [this stack overflow post](https://stackoverflow.com/questions/40746453/how-to-connect-to-docker-host-from-container-on-windows-10-docker-for-windows). Also, in this sample, you need to give permission for `C:/Docker` directory. For example, I give `Everyone` to `full access` only on this directory. Don't forget to trun Share Directory on Docker for Windows on. 

```powershell
PS > docker run -v C:/Docker/config:/jmeter_config -v C:/Docker/report:/jmeter_report some jmeter -n -t /jmeter_config/DocumentApi.jmx -l /jmeter_report/current.jtl -e -o /jmeter_report/report -Jhost=docker.for.win.localhost -Jthreads=100 -Jduration=60 -Jport=6400
```


**NOTE:** If you want to use Windows Subsystem for Linux (WSL), you will need to Enable Shared Drives on Docker for Windows. To Enable Shared Drives if you are authenticating to Azure Active Directly, You will need to add another Local Admin account and use its credential to map the Share Drive. Then assign read/write/modify permissions to the targed directory to this account. For more details review the contents of [this issue](https://github.com/docker/for-win/issues/1801)

## Master/Slave configuration

Configuring an Inbound TCP PORT for `1099` and `30000-65535` will be required. 


docker run -v C:/Users/tsushi/Codes/EpiStorageService/StressTest/config:/jmeter_config -v C:/Users/tsushi/Codes/EpiStorageService/StressTest/report:/jmeter_report some jmeter -n -t /jmeter_log/DocumentApi.jmx -l /jmeter_report/current.jtl -e -o /jmeter_report/report