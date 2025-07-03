package main

import (
	"akv/api"
	"akv/controller"
)

func main() {
	podController := controller.NewPodController()
	go podController.Run()

	api.StartServer(podController)
	stopch := make(chan struct{})
	<-stopch

}
