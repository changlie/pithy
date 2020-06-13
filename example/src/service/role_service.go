package service

import pt "changlie/pithy"

type RoleServiceImpl struct {}


func (role *RoleServiceImpl) RoleList(r pt.Req) {
    r.RespJson("it's role list.")
}

func  CreateRole(r pt.Req) {
    r.RespJson("create role successfully.")
}


func (role *RoleServiceImpl) RoleDelete(r pt.Req) {
    r.RespJson("delete role finish!")
}