package internel_pithy_gen

import pt "changlie/pithy"
import s "service"

func init() {
    pt.SetHandler("/employee/list", s.EmployeeList)
    pt.SetHandler("/employee/create", s.EmployeeCreate)


    var roleserviceimpl *s.RoleServiceImpl
    pt.SetHandler("/role/list", roleserviceimpl.RoleList)
    pt.SetHandler("/create/role", s.CreateRole)
    pt.SetHandler("/role/delete", roleserviceimpl.RoleDelete)


    var userserviceimpl *s.UserServiceImpl
    pt.SetHandler("/users", userserviceimpl.Users)
    pt.SetHandler("/user/add", userserviceimpl.UserAdd)
    pt.SetHandler("/user/update", userserviceimpl.UserUpdate)
    pt.SetHandler("/user/del", userserviceimpl.UserDel)


}
