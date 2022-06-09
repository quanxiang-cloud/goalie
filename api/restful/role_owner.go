package restful

import (
	"fmt"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/goalie/internal/goalie"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"github.com/quanxiang-cloud/goalie/pkg/header2"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"net/http"
)

//RoleOwnerAPI 角色和管理者api
type RoleOwnerAPI struct {
	roleOwner goalie.RoleOwner
}

//NewRoleOwnerAPI new
func NewRoleOwnerAPI(conf *config.Config, db *gorm.DB) *RoleOwnerAPI {
	return &RoleOwnerAPI{
		roleOwner: goalie.NewRoleOwner(conf, db),
	}
}

// ListRoleOwner 获取角色拥有者列表
func (g *RoleOwnerAPI) ListRoleOwner(c *gin.Context) {
	r := &goalie.ListRoleOwnerRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(g.roleOwner.ListRoleOwner(ginheader.MutateContext(c), c.Request.Header.Clone(), r)).Context(c)
}

// UpdateRoleOwner 修改角色拥有者列表
func (g *RoleOwnerAPI) UpdateRoleOwner(c *gin.Context) {
	r := &goalie.UpdateRoleOwnerRequest{}
	if err := c.ShouldBind(r); err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(g.roleOwner.UpdateRoleOwner(ginheader.MutateContext(c), r)).Context(c)
}

// ListUserRoleTags token获取用户角色tag列表
func (g *RoleOwnerAPI) ListUserRoleTags(c *gin.Context) {
	r := &goalie.ListUserRoleRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	profile := header2.GetProfile(c)
	r.UserID = profile.UserID
	departments := header2.GetDepartments(c)
	if len(departments) > 0 {
		for k := range departments {
			r.DepartmentID = append(r.DepartmentID, departments[k]...)
		}
	}
	result, err := g.roleOwner.ListUserRole(ginheader.MutateContext(c), r)
	role := make([]string, 0, len(result.Roles))
	roleID := make([]string, 0, len(result.Roles))
	for _, r := range result.Roles {
		role = append(role, r.Tag)
		roleID = append(roleID, r.RoleID)
	}
	header2.SetRole(c, role...)
	resp.Format(role, err).Context(c)
}

// ListUserRole token获取用户角色列表
func (g *RoleOwnerAPI) ListUserRole(c *gin.Context) {
	r := &goalie.ListUserRoleRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	profile := header2.GetProfile(c)
	r.UserID = profile.UserID
	departments := header2.GetDepartments(c)
	if len(departments) > 0 {
		for k := range departments {
			r.DepartmentID = append(r.DepartmentID, departments[k]...)
		}
	}
	result, err := g.roleOwner.ListUserRole(ginheader.MutateContext(c), r)

	resp.Format(result, err).Context(c)
}

// ListOwnerRole 获取拥有者角色列表
func (g *RoleOwnerAPI) ListOwnerRole(c *gin.Context) {
	r := &goalie.ListOwnerRoleRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(g.roleOwner.ListOwnerRole(ginheader.MutateContext(c), r)).Context(c)
}

// DeleteOwnerRole 删除拥有者角色
func (g *RoleOwnerAPI) DeleteOwnerRole(c *gin.Context) {
	r := &goalie.DeleteOwnerRoleRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(g.roleOwner.DeleteOwnerRole(ginheader.MutateContext(c), r)).Context(c)
}

// UpdateOwnerRole 修改拥有者角色
func (g *RoleOwnerAPI) UpdateOwnerRole(c *gin.Context) {
	r := &goalie.UpdateOwnerRoleRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(g.roleOwner.UpdateOwnerRole(ginheader.MutateContext(c), r)).Context(c)
}

// OthUserRoleList 其它服务获取用户角色列表
func (g *RoleOwnerAPI) OthUserRoleList(c *gin.Context) {
	r := &goalie.OthOneUserRoleListRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	result, err := g.roleOwner.OthUserRoleList(ginheader.MutateContext(c), r)
	resp.Format(result, err).Context(c)
}

// OthRoleListByIDs 其它服务根据ids获取用户角色列表
func (g *RoleOwnerAPI) OthRoleListByIDs(c *gin.Context) {
	r := &goalie.OthRoleListByIDsRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	result, err := g.roleOwner.OthRoleListByIDs(ginheader.MutateContext(c), r)
	resp.Format(result, err).Context(c)
}

// OthDelByOwner 其它服务通过用户ID删除权限
func (g *RoleOwnerAPI) OthDelByOwner(c *gin.Context) {
	r := &goalie.OthDelRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resp.Format(g.roleOwner.OthDelUserRoleList(ginheader.MutateContext(c), r)).Context(c)
}
