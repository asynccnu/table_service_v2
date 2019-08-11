package table

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
)

type AddResponse struct {
	Id 	int64
}

// 自定义添加课程
func Add(c *gin.Context) {
	var table model.TableItem

	if err := c.BindJSON(&table); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	sid := c.Request.Header.Get("sid")
	if sid == "" {
		SendBadRequest(c, errno.ErrBind, nil, "No sid")
		return
	}

	id, err := model.AddSelfTable(sid, &table)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	SendResponse(c, nil, AddResponse{Id: id})
}