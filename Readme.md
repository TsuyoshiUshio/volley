# Volley
Volley is a command line tool for managing JMeter Cluster.

The key feature of Volley are: 

* **Privsioning/Deprovisioning:** Provision/Deplovision Master/Slave cluster of JMeter. Currently, we support Virtual Machine for Azure. However, you can contribute to add other providers. 
* **Send a scenario:** Send JMX file and data files to the Master.
* **Run:** Run Stress/Load testing.
* **Fetch a report:** Fetch a report from the Master. 
* **Server:** Worked as an API Server on the JMeter master side. Execute Remote JMeter Server from the client request.  

## Getting Started 

```bash
./volley provision --slave 10
Provisioning JMeter Environment Master: 1, Slave: 10
```

## Developing Volley

TODO