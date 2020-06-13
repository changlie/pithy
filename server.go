package pithy

import (
	"changlie/pithy/common"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

var DevMode bool

// server entry.
func Start() {
	DevMode = common.GetConfigBool("dev.mode")
	if !DevMode {
		os.Chdir(filepath.Dir(os.Args[0]))
	}

	http.HandleFunc("/", mainHandler)
	port := common.GetConfig("server.port", "8888")
	addr := fmt.Sprintf(":%s", port)
	log.Printf("server [%v] start up successfully!", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Pwd() string {
	dir, _ := os.Getwd()
	return dir
}

func OpenBrowser(url string) {
	go func() {
		time.Sleep(time.Second)

		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", "start", url)
		} else {
			cmd = exec.Command("xdg-open", url)
		}
		cmd.Run()
	}()
}
