# JMeter Docker

JMeter Docker images for Stress/Load testing. 
This docker image wraps the jmeter commands required to perform stress testing. 

_linux_

```bash
$ docker run -d -P --name master -v /home/ushio/Codes/DevSecOps/EpiServer/StressTesting/Temp/:/jmeter_log tsuyoshiushio/jmeter jmeter -n -t /jmeter_log/MessageApi.jmx -l /jmeter_log/current2.jtl -e -o /jmeter_log/report2 -Jthreads=100 -Jduration=60 -Jport=6400
```

_windows_

```powershell
PS > docker run -d -P --name master4 -v C:/Users/tsushi/JMeter:/jmeter_log tsuyoshiushio/jmeter jmeter -n -t /jmeter_log/MessageApi.jmx -l /jmeter_log/current2.jtl -e -o /jmeter_log/report -Jthreads=100 -Jduration=60 -Jport=6400
```

**NOTE:** If you want to use Windows Subsystem for Linux (WSL), you will need to Enable Shared Drives on Docker for Windows. To Enable Shared Drives if you are authenticating to Azure Active Directly, You will need to add another Local Admin account and use its credential to map the Share Drive. Then assign read/write/modify permissions to the targed directory to this account. For more details review the contents of [this issue](https://github.com/docker/for-win/issues/1801)

## Master/Slave configuration

Configuring an Inbound TCP PORT for `1099` and `30000-65535` will be required. 