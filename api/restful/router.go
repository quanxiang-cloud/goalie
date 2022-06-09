package restful

import (
	"context"
	"github.com/quanxiang-cloud/cabin/logger"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"github.com/quanxiang-cloud/goalie/pkg/probe"
	"github.com/quanxiang-cloud/goalie/pkg/util"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	ginlogger "github.com/quanxiang-cloud/cabin/tailormade/gin"
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
)

// Router 路由
type Router struct {
	c *config.Config

	engine *gin.Engine

	goalie *Goalie
}

// NewRouter 开启路由
func NewRouter(ctx context.Context, c *config.Config, log logger.AdaptedLogger, db *gorm.DB) (*Router, error) {
	engine, err := newRouter(c)
	if err != nil {
		return nil, err
	}

	v1 := engine.Group("/api/v1/goalie")

	goalie, err := NewGoalie(c, db)
	if err != nil {
		return nil, err
	}
	roleAPI := NewRoleAPI(c, db)
	funcAPI := NewFuncAPI(c, db)
	roleOwnerAPI := NewRoleOwnerAPI(c, db)

	k := v1.Group("/role")
	{
		k.POST("/add", roleAPI.CreateRole)
		k.PUT("/update", roleAPI.UpdateRole)
		k.DELETE("/del", roleAPI.DeleteRole)
		//k.POST("/getRole", goalie.GetRole)
		k.GET("/get", roleAPI.GetRole)

		//k.POST("/listRole", goalie.ListRole)
		k.GET("/list", roleAPI.PageListRole)

		//k.POST("/listRoleOwner", goalie.ListRoleOwner)
		k.GET("/owner/list", roleOwnerAPI.ListRoleOwner)

		//k.POST("/updateRoleOwner", goalie.UpdateRoleOwner)
		k.POST("/update/owner", roleOwnerAPI.UpdateRoleOwner)

		k.POST("/transferRoleSuper", goalie.TransferRoleSuper)

		//k.POST("/listUserRole", goalie.ListUserRole)
		k.GET("/user/list", roleOwnerAPI.ListUserRoleTags)

		k.GET("/now/list", roleOwnerAPI.ListUserRole)

		//k.POST("/listRoleFunc", goalie.ListRoleFunc)
		k.GET("/func/role/list", roleAPI.ListRoleFunc)

		//k.POST("/updateRoleFunc", goalie.UpdateRoleFunc)
		k.POST("/func/update", roleAPI.UpdateRoleFunc)

		//k.POST("/deleteOwnerRole", goalie.DeleteOwnerRole)
		k.DELETE("/owner/delete", roleOwnerAPI.DeleteOwnerRole)

		//k.POST("/updateOwnerRole", goalie.UpdateOwnerRole)
		k.PUT("/update/owner/role", roleOwnerAPI.UpdateOwnerRole)

		//k.POST("/listFuncTag", goalie.ListFuncTag)
		k.GET("/func/all/list", funcAPI.ListFuncTag)

		//k.POST("/listUserFuncTag", goalie.ListUserFuncTag)
		k.POST("/func/user/list", funcAPI.ListUserFuncTag)

		//k.POST("/listOwnerRole", goalie.ListOwnerRole)
		k.POST("/owner/role/list", roleOwnerAPI.ListOwnerRole)

		//k.GET("/role", goalie.ListUserRoles)

		k.POST("/del/owner", roleOwnerAPI.OthDelByOwner)
	}

	funcs := v1.Group("/func")
	{
		funcs.GET("/list", funcAPI.List) //可能待完善
	}

	{
		probe := probe.New(util.LoggerFromContext(ctx))
		engine.GET("liveness", func(c *gin.Context) {
			probe.LivenessProbe(c.Writer, c.Request)
		})

		engine.Any("readiness", func(c *gin.Context) {
			probe.ReadinessProbe(c.Writer, c.Request)
		})

	}

	return &Router{
		c:      c,
		engine: engine,
		goalie: goalie,
	}, nil
}

func newRouter(c *config.Config) (*gin.Engine, error) {
	if c.Model == "" || (c.Model != ReleaseMode && c.Model != DebugMode) {
		c.Model = ReleaseMode
	}
	gin.SetMode(c.Model)
	engine := gin.New()

	engine.Use(ginlogger.LoggerFunc(), ginlogger.LoggerFunc())

	return engine, nil
}

// Run 启动服务
func (r *Router) Run() {
	r.engine.Run(r.c.Port)
}

// Close 关闭服务
func (r *Router) Close() {
}
