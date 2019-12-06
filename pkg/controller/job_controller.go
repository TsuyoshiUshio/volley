package controller

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Start(c *gin.Context) {
	var json model.Config
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job_id, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// start with initialize RunContext then execute with go routine. 

	path := filepath.Join(".", model.JobDir, job_id.String())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModeDir)
	}

	configDirPath := filepath.Join(".", model.ConfigDir, json.ID)

	statusPath := filepath.Join(".", model.JobDir, job_id.String())
	status := model.Status{
		Status: model.StatusRunning,
	}
	status.Write(statusPath)

	go func() {
		// docker run -d -P --name master -v /home/ushio/Codes/DevSecOps/EpiServer/StressTesting/Temp/:/jmeter_log tsuyoshiushio/jmeter jmeter -n -t /jmeter_log/MessageApi.jmx -l /jmeter_log/current2.jtl -e -o /jmeter_log/report2 -Jthreads=100 -Jduration=60 -Jport=6400
		out, err := exec.Command("docker", "run", "").Output
		fmt.Printf("%s", out)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

	}()

}
