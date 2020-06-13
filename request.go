package pithy

import (
	"net/http"
)

type Req interface {
	getArg(name string) []string
	RespJson(d interface{})
	RespFile(path string)
}

type DefaultReq struct {
	Args map[string][]string
	Header http.Header
	resp *Resp
}

func (r *DefaultReq) getArg(name string) []string {
	return r.Args[name]
}

func (r *DefaultReq) RespJson(d interface{}) {
	r.resp = &Resp{
		data: d,
		t: respJson,
	}
}

func (r *DefaultReq) RespFile(path string) {
	r.resp = &Resp{
		data: path,
		t: respFile,
	}
}


