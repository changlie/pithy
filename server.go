package pithy

import (
	"io"
	"log"
	"net/http"

)

func Start() {
  // Hello world, the web server

  helloHandler := func(w http.ResponseWriter, req *http.Request) {
    io.WriteString(w, "Hello, world!\n")
  }

  http.HandleFunc("/hello", helloHandler)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

func OpenBrowser(url string) {
  go func() {
    time.Sleep(time.Second)

    cmd := exec.Command("xdg-open", url)
    cmd.Run()
  }()
}
