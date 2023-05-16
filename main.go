package main

import (
	"flag"
	"fmt"
	"gogen_grpc/application"
	"gogen_grpc/shared/gogen"
)

var Version = "0.0.1"

func main() {

	appMap := map[string]gogen.Runner{
		//
		"client": application.NewAppClient(),
		"server": application.NewAppServer(),
	}

	flag.Parse()

	app, exist := appMap[flag.Arg(0)]
	if !exist {
		fmt.Printf("You may try :\n\n")
		for appName := range appMap {
			fmt.Printf("    go run main.go %s\n", appName)
		}
		fmt.Printf("\n")
		return
	}

	//fmt.Printf("Version %s\n", Version)

	err := app.Run()
	if err != nil {
		return
	}

}
