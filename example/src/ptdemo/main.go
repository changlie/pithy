package main

import _ "internel_pithy_gen"
import (
	ptpkg "changlie/pithy/pkg"
	pt "changlie/pithy"
	"changlie/pithy/common"
	"fmt"
	"os"
)

func gen() {
	ptpkg.GenerateInitHandlerFile()
}

func serverStartup() {
	pt.Start()
}

func main() {
	if !common.GetConfigBool("dev.mode") {
		serverStartup()
		return
	}
	if len(os.Args) < 2 {
		return
	}

	action := os.Args[1]
	fmt.Println("action:", action)
	switch action {
	case "gen":
		fmt.Println("generate source file")
		gen()

	case "run":
		fmt.Println("start up server!")
		serverStartup()
	}
}





