package pithy

import (
  "bufio"
  "net/http"
  "fmt"
  "strings"
)

func mainHandler(w http.ResponseWriter, r *http.Request)  {
  printReqInfo(r)
  url := r.URL.RawPath
  queryArgs := r.URL.Query()
  headers := r.Header

  req := Req{queryArgs, headers}

  h := getHandler(url)
  if h == nil {
    w.WriteHeader(404)
    fmt.Fprintln(w, "service is not found!")
    return
  }
  resp := h(req)
  handleResp(w, resp)

}

func handleResp(w http.ResponseWriter, resp *Resp) {
  if resp.isJson() {
    var res string
    switch val := resp.data.(type) {
    case string:
      res = val
    default:
      res = "not support temp"
    }
    fmt.Fprint(w, res)
    return
  }
  bufio.NewWriter(w).Write([]byte("not support temp"))
}

func printReqInfo(r *http.Request) {
  fmt.Println(r.URL)
  fmt.Println("head:---------------------")
  for k, v := range r.Header {
    fmt.Println(k, v)
  }
}

type uri struct {
  url string
  // subUrls []string
}

var (
  Urlmap = make(map[uri]Handler)
)

func getHandler(url string) Handler {
  u := NewUri(url)
  for p, h := range Urlmap {
    if p.equal(u) {
      return h
    }
  }
  return nil
}

func SetHandler(url string, h Handler) {
  u := NewUri(url)
  Urlmap[u] = h
}


func NewUri(url string) uri {
  paths := strings.Split(url, "/")
  return uri{url, paths}
}

func (u *uri) equal(target uri) bool {
  subUrls1 := target.subUrls
  for i, subUrl := range u.subUrls {
    if i >= len(subUrls1) {
      return false
    }
    subUrl1 := subUrls1[i]
    if subUrl != subUrl1 && !strings.HasPrefix(subUrl, ":") {
      return false
    }
  }
  return true
}

type Handler func(r Req) *Resp

type Req struct {
  Args map[string][]string
  Header http.Header
}

func (r *Req) getArg(name string) []string {
  return r.Args[name]
}

type Resp struct {
  data interface{}
  t respType
}

type respType int

const (
  respJson respType = 1 << iota
  respFile
)

func RespJson(d interface{}) *Resp {
  return &Resp{
    data: d,
    t: respJson,
  }
}

func RespFile(path string) *Resp {
  return &Resp{
    data: path,
    t: respFile,
  }
}

func (resp *Resp) isJson() bool {
  return (resp.t & respJson) != 0
}

func (resp *Resp) isFile() bool {
  return (resp.t & respFile) != 0
}

type Json interface {
  JsonString() string
}
