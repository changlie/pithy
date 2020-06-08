package pithy

import (
  // "io"
  "log"
  "net/http"
  "os/exec"
  "time"
)

func Start() {
  http.HandleFunc("/", mainHandler)
  log.Fatal(http.ListenAndServe(":8888", nil))
}

func OpenBrowser(url string) {
  go func() {
    time.Sleep(time.Second)

    cmd := exec.Command("xdg-open", url)
    cmd.Run()
  }()
}
