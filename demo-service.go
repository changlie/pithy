package pithy

import (
  "fmt"
)

func init() {
  fmt.Println("service init...")

  SetHandler("users", users)
  SetHandler("userAdd", userAdd)
}


func users(r Req) *Resp {
  return RespJson("users")
}

func userAdd(r Req) *Resp {
  return RespJson("create user")
}
