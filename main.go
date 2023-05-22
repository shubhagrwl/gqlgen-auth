package main

import (
	"fmt"
	"net/http"
	"time"
	"todo/internal/app/api/router"
	"todo/internal/app/api/service"
	"todo/internal/app/constants"

	"github.com/spf13/viper"
)

func main() {
	services := service.Init()

	routerV1 := router.Init(services)
	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", viper.GetString(constants.ServerPort)),
		Handler:      routerV1,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.ListenAndServe()
}
