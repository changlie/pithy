package pithy

import (
	"fmt"
	"net/http"
)

type Resp struct {
	data interface{}
	t respType
}

type respType int

const (
	respJson respType = 1 << iota
	respFile
)


func (resp *Resp) isJson() bool {
	return (resp.t & respJson) != 0
}

func (resp *Resp) isFile() bool {
	return (resp.t & respFile) != 0
}

func (resp *Resp) String() string {
	return fmt.Sprintf("%v", resp.data)
}

type Json interface {
	JsonString() string
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

	w.Write([]byte("request is not support temporarily!"))
}
