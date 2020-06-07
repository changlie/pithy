package pithy

import (
	"io"
	"log"
	"net/http"
  "os/exec"
  "time"
)

func Start() {
  // Hello world, the web server

  helloHandler := func(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "Hello, world!\n")
  }

  http.HandleFunc("/hello", helloHandler)
  go openBrowser()
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func openBrowser() {
  time.Sleep(time.Second)

  cmd := exec.Command("xdg-open", "http://as:8080/hello")
  cmd.Run()
}
