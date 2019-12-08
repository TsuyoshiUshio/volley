package command

import (
	"github.com/TsuyoshiUshio/volley/pkg/controller"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type ServerCommand struct {
}

func (s *ServerCommand) Start(c *cli.Context) error {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello server",
		})
	})

	router.POST("/job", controller.Start)
	router.GET("/job/:job_id", controller.StatusCheck)

	srv := &http.Server{
		Addr:    ":38080",
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
