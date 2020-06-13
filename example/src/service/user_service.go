package service

import (
	pt "changlie/pithy"
	"fmt"
)

type UserServiceImpl struct{}

func (u *UserServiceImpl) Users(r pt.Req) {

	r.RespJson("UserServiceImpl#Users")
}

func (u *UserServiceImpl) UserAdd(r pt.Req) {
	fmt.Println("create user successfully!")
	r.RespJson("UserServiceImpl#UserAdd1111")
}

func (u *UserServiceImpl) UserUpdate(r pt.Req) {
	r.RespJson("UserServiceImpl#UserUpdate")
}

func (u *UserServiceImpl) UserDel(r pt.Req) {
	r.RespJson("UserServiceImpl#UserDel")
}



