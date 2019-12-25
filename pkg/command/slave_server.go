package command

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/TsuyoshiUshio/volley/pkg/controller"
	"github.com/TsuyoshiUshio/volley/pkg/model"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

type SlaveServerCommand struct {
}

func (s *SlaveServerCommand) Start(c *cli.Context) error {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello slave server",
		})
	})

	router.POST("/csv", controller.UploadCSV)

	srv := &http.Server{
		Addr:    ":" + model.SlaveDefaultPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	return nil
}
