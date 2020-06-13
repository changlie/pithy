package _interface

import (
	pt "changlie/pithy"
	s "service"
)

type UserService interface {
	Users(r pt.Req)
	UserAdd(r pt.Req)
	UserUpdate(r pt.Req)
	UserDel(r pt.Req)
}

var userServiceImpl s.UserServiceImpl
func NewUserService() UserService {
	return &userServiceImpl
}
