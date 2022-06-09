package code

import error2 "github.com/quanxiang-cloud/cabin/error"

func init() {
	error2.CodeTable = CodeTable
}

const (
	// ErrSuperUpdate 超级管理不允许修改
	ErrSuperUpdate = 60014090001
	// InvalidParams 无效的参数
	InvalidParams = 20014000002
	// UnknownRole 未知角色
	UnknownRole = 60014040003
	// InvalidRoleOwner 无效的角色拥有者
	InvalidRoleOwner = 60014000004
	// RoleNotExist 角色不存在
	RoleNotExist = 60014040005

	// ErrNameUsed 名称已经被使用
	ErrNameUsed = 20014000006
	//ErrNoFuncs 没有系统功能，请补充
	ErrNoFuncs = 20014090007
	//ErrRoleNoOwners 角色没有关联拥有者
	ErrRoleNoOwners = 20014090008
)

// CodeTable 码表
var CodeTable = map[int64]string{
	ErrSuperUpdate:   "超级管理不允许修改.",
	UnknownRole:      "未知角色.",
	InvalidRoleOwner: "无效的角色拥有者.",
	RoleNotExist:     "角色不存在.",
	ErrNameUsed:      "名称已经被使用",
	ErrNoFuncs:       "没有系统功能，请补充",
	ErrRoleNoOwners:  "角色没有关拥有者，请补充",
}
