package main

import (
	"fmt"

	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/clients"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/config"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/log"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/proxy"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/server"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/utils"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func main() {
	// init logrus logger
	log.InitLogger()

	// log runtime attributes
	utils.LogRuntimeAttributes()

	// parse env variables
	if err := envconfig.Process("", &config.Config); err != nil {
		logrus.Fatalf("Failed to parse ENVs, %v", err)
	}

	// set proxy envs, if provided
	proxy.SetProxyENVs()

	// Initialize mongodb clients
	businessClient, jobClient, jobApplicationClient, initJobApplication, schemeClient, schemeApplicationClient := server.InitMongoDB()

	// Set up clients
	clients := clients.NewClients(jobClient, businessClient, jobApplicationClient, initJobApplication, schemeClient, schemeApplicationClient)

	// initialize the server
	server := server.SetupServer(clients)

	// start the server
	if err := server.Run(fmt.Sprintf(":%s", config.Config.HTTPPort)); err != nil {
		logrus.Fatal(err)
	}
}
