package pithy

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(r Req)

type Handler struct {
	method string
	Paths []string
	Action HandleFunc
}

var urlmap = make(map[string]Handler)

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
	log.Printf("register service %v:%v \n", url, paths)
	urlmap[url] = Handler{Paths:paths, Action:hf}
}

func setHandler1(method, url string, hf HandleFunc) {
	paths := strings.Split(url, "/")
	log.Printf("register service[%v] %v:%v \n", method, url, paths)
	urlmap[url] = Handler{method:method, Paths:paths, Action:hf}
}

func SetGetHandler(url string, hf HandleFunc) {
	setHandler1(http.MethodGet, url, hf)
}

func SetPostHandler(url string, hf HandleFunc) {
	setHandler1(http.MethodPost, url, hf)
}

func mainHandler(w http.ResponseWriter, r *http.Request)  {
	url := r.URL.Path

	// if Accessing local file, return local file directly
	fp, err := localFile(url)
	if err == nil {
		returnLocalFile(w, fp)
		return
	}

	// if resources exists in cache, retrun it from cache
	rc := GetResource(url)
	if rc != nil {
		w.Write(rc)
		return
	}

	queryArgs := r.URL.Query()
	headers := r.Header

	req := &DefaultReq{Args:queryArgs, Header:headers}

	h := getHandler(url)
	if h == nil {
		NotFound(w)
		return
	}
	h(req)
	if req.resp != nil {
		handleResp(w, req.resp)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}










