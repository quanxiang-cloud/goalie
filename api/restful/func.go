package restful

import (
	"github.com/gin-gonic/gin"
	ginheader "github.com/quanxiang-cloud/cabin/tailormade/header"
	"github.com/quanxiang-cloud/cabin/tailormade/resp"
	"github.com/quanxiang-cloud/goalie/internal/goalie"
	"github.com/quanxiang-cloud/goalie/pkg/config"
	"github.com/quanxiang-cloud/goalie/pkg/header2"
	"gorm.io/gorm"
	"net/http"
)

//FuncAPI 系统功能接口
type FuncAPI struct {
	funcs goalie.Func
}

//NewFuncAPI new
func NewFuncAPI(conf *config.Config, db *gorm.DB) *FuncAPI {
	return &FuncAPI{
		funcs: goalie.NewFunc(conf, db),
	}
}

//List 集合
func (f *FuncAPI) List(c *gin.Context) {
	r := &goalie.GetListFuncRequest{}
	resp.Format(f.funcs.List(ginheader.MutateContext(c), r)).Context(c)
}

// ListFuncTag 获取全量tag
func (f *FuncAPI) ListFuncTag(c *gin.Context) {
	r := &goalie.ListFuncTagRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	resp.Format(f.funcs.ListFuncTag(ginheader.MutateContext(c), r)).Context(c)
}

// ListUserFuncTag 获取用户funcTag
func (f *FuncAPI) ListUserFuncTag(c *gin.Context) {
	r := &goalie.ListUserFuncTagRequest{}
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
	resp.Format(f.funcs.ListUserFuncTag(ginheader.MutateContext(c), r)).Context(c)
}
