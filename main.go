/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 17-01-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/LightsPlatform/vActuator/actuator"
	"github.com/LightsPlatform/vActuator/stateManager"
	"encoding/json"
	"github.com/LightsPlatform/vSensor/sensor"
)

var actuators map[string]*actuator.Actuator

// init initiates global variables
func init() {
	actuators = make(map[string]*actuator.Actuator)
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)
		api.POST("/actuator/:id", actuatorCreateHandler)
		api.POST("/actuator/:id/trigger", actuatorTriggerHandler)
		api.GET("/actuator/:id/state", actuatorDataHandler)
		api.GET("/actuator/", actuatorListHandler)
		api.DELETE("/actuator/:id", acuatorDeleteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func acuatorDeleteHandler(c *gin.Context) {
	log.Println("delete")
	id := c.Param("id")

	actuator, ok := actuators[id]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Sensor %s was not found on vSensor", id)})
		return
	}

	actuator.Stop()
	state := actuator.State
	delete(actuators, id)

	c.JSON(http.StatusOK, state)
}

func actuatorListHandler(c *gin.Context) {
	output := make([]string, 0)
	for _, actuator := range actuators {
		output = append(output, actuator.Name)
	}
	c.JSON(http.StatusOK, output)
}

func actuatorTriggerHandler(c *gin.Context){

	id := c.Param("id")
	if actuators[id] == nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Actuator not found!",
		})
		return
	}
	action,_ := c.GetPostForm("action")
	if(action == ""){
		action = "null"
	}
	result := actuators[id].Trigger(action);
	c.JSON(http.StatusOK,result)
}

func actuatorDataHandler(c *gin.Context){
	id := c.Param("id")
	if actuators[id] == nil{
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Actuator not found!",
		})
		return
	}

	c.JSON(http.StatusOK, actuators[id].State)
}

func actuatorCreateHandler(c *gin.Context) {
	id := c.Param("id")
	code, ok := c.GetPostForm("code")
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "code must send!",
		})
		return
	}
	configData, ok := c.GetPostForm("config")
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "config must send!",
		})
		return
	}

	config := &stateManager.Config{
		StateType: map[string][]string{},
	}

	error := json.Unmarshal([]byte(configData), config)
	if error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "config is not correct! ",
			"description" : error.Error(),
		})
		return
	}

	actuator, error := actuator.New(id, []byte(code), *config)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	if _, ok := actuators[id]; !ok {
		go actuator.Run()
	}

	actuators[id] = actuator
	c.String(http.StatusOK, id)
}

func main() {
	fmt.Println("vActuator Light @ 2018")

	srv := &http.Server{
		Addr:    ":8181",
		Handler: handle(),
	}

	go func() {
		fmt.Printf("vActuator Listen: %s\n", srv.Addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Listen Error:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("vActuator Shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown Error:", err)
	}
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "I'll keep Light in my heart ❤")
}
