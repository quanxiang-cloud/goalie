package restful

import (
	error2 "github.com/quanxiang-cloud/cabin/error"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/goalie/internal/goalie"
	"github.com/quanxiang-cloud/goalie/pkg/code"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"github.com/quanxiang-cloud/goalie/pkg/header2"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//RoleAPI 角色api
type RoleAPI struct {
	role     goalie.Role
	roleFunc goalie.RoleFunc
}

//NewRoleAPI new
func NewRoleAPI(conf *config.Config, db *gorm.DB) *RoleAPI {
	return &RoleAPI{
		role:     goalie.NewRole(conf, db),
		roleFunc: goalie.NewRoleFunc(conf, db),
	}
}

// CreateRole 创建角色
func (g *RoleAPI) CreateRole(ctx *gin.Context) {
	r := new(goalie.CreateRoleRequest)
	if err := ctx.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams, err.Error())).Context(ctx)
		return
	}
	profile := header2.GetProfile(ctx)
	r.CreatedBy = profile.UserID
	r.TenantID = profile.TenantID
	resp.Format(g.role.CreateRole(ginheader.MutateContext(ctx), r)).Context(ctx)
}

// UpdateRole 更新角色信息
func (g *RoleAPI) UpdateRole(ctx *gin.Context) {

	r := new(goalie.UpdateRoleRequest)
	if err := ctx.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams, err.Error())).Context(ctx)
		return
	}
	profile := header2.GetProfile(ctx)
	r.UpdatedBy = profile.UserID

	resp.Format(g.role.UpdateRole(ginheader.MutateContext(ctx), r)).Context(ctx)
}

// GetRole 获取角色信息
func (g *RoleAPI) GetRole(ctx *gin.Context) {
	r := new(goalie.GetRoleRequest)
	if err := ctx.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams, err.Error())).Context(ctx)
		return
	}
	resp.Format(g.role.GetRole(ginheader.MutateContext(ctx), r)).Context(ctx)
}

// PageListRole 查询角色列表
func (g *RoleAPI) PageListRole(ctx *gin.Context) {
	r := new(goalie.PageListRoleRequest)
	if err := ctx.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams, err.Error())).Context(ctx)
		return
	}
	profile := header2.GetProfile(ctx)
	r.TenantID = profile.TenantID
	resp.Format(g.role.PageListRole(ginheader.MutateContext(ctx), r)).Context(ctx)
}

// DeleteRole 删除角色
func (g *RoleAPI) DeleteRole(ctx *gin.Context) {
	r := new(goalie.DeleteRoleRequest)
	if err := ctx.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams, err.Error())).Context(ctx)
		return
	}
	resp.Format(g.role.DeleteRole(ginheader.MutateContext(ctx), r)).Context(ctx)
}

//UpdateRoleFunc 绑定角色和功能
func (g *RoleAPI) UpdateRoleFunc(c *gin.Context) {
	r := new(goalie.UpdateRoleFuncRequest)
	if err := c.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(g.roleFunc.CreateInBatches(ginheader.MutateContext(c), r)).Context(c)
}

//ListRoleFunc 角色的功能
func (g *RoleAPI) ListRoleFunc(c *gin.Context) {
	r := new(goalie.ListRoleFuncRequest)
	if err := c.ShouldBind(r); err != nil {
		resp.Format(nil, error2.New(code.InvalidParams)).Context(c)
		return
	}
	resp.Format(g.roleFunc.ListRoleFunc(ginheader.MutateContext(c), r)).Context(c)
}
