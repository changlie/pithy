package pithy

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func mainHandler(w http.ResponseWriter, r *http.Request)  {
	printReqInfo(r)
	url := r.URL.Path

	// if Accessing local file, return local file directly
	fp, err := localFile(url)
	if err == nil {
		returnLocalFile(w, fp)
		return
	}

	log.Printf("string(): %v; path: %v; rawpath: %v \n", r.URL, r.URL.Path, r.URL.RawPath)
	queryArgs := r.URL.Query()
	headers := r.Header

	req := Req{Args:queryArgs, Header:headers}

	h := getHandler(url)
	if h == nil {
		NotFound(w)
		return
	}
	resp := h(req)
	handleResp(w, resp)
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "service is not found!")
}

func InternalServerError(w http.ResponseWriter, e error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "server error: "+e.Error())
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

	if resp.isFile() {
		path, ok := resp.data.(string)
		if !ok {
			NotFound(w)
			return
		}
		fp, err := localFile(path)
		if err != nil {
			NotFound(w)
			return
		}
		returnLocalFile(w, fp)

		return
	}

	bufio.NewWriter(w).Write([]byte("not support temp"))
}

func localFile(path string) (string, error) {
	if path == "" || path == "/" {
		return "", errors.New("don't handle root path")
	}
	fp := filepath.Join(Pwd(), "static", path)
	_, err := os.Open(fp)
	return fp, err
}

func returnLocalFile(w http.ResponseWriter, fp string) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		InternalServerError(w, err)
		return
	}
	bw := bufio.NewWriter(w)
	bw.Write(data)
	bw.Flush()
}

func printReqInfo(r *http.Request) {
	fmt.Println(r.URL)
	fmt.Println("head:---------------------")
	for k, v := range r.Header {
		fmt.Println(k, v)
	}
}

var (
	urlmap = make(map[string]Handler)
)

func getHandler(url string) HandleFunc {
	for p, h := range urlmap {
		if p == url {
			return h.Action
		}
	}
	log.Printf("current urlmap: %v \n", urlmap)
	log.Printf("service is not found for %v \n", url)

	return nil
}

func SetHandler(url string, hf HandleFunc) {
	paths := strings.Split(url, "/")

	urlmap[url] = Handler{Paths:paths, Action:hf}
}

type HandleFunc func(r Req) *Resp

type Handler struct {
	Paths []string
	Action HandleFunc
}

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
