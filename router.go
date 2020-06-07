package pithy


var (
  Urlmap = make(map[string]Handler)
)

type Handler func(r Req) Resp

type Req struct {

}

type Resp struct {

}
