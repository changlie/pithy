package service

import (
	pt "changlie/pithy"
	"fmt"
)
// `rootpath#employee`
// `get#list`
func EmployeeList(r pt.Req) {
	fmt.Println("over### EmployeeList")
	r.RespJson("service#EmployeeList1111111")
}

// `post#create`
func EmployeeCreate(r pt.Req) {
	fmt.Println("over### EmployeeCreate")
	r.RespJson("service#EmployeeCreate222222222")
}


