package table

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

type AddResponse struct {
	Id string
}

// 自定义添加课程
func Add(c *gin.Context) {
	log.Infof("Add function called.")

	var table model.TableItem

	if err := c.BindJSON(&table); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	sid := c.MustGet("Sid").(string)

	id, err := model.AddSelfTable(sid, &table)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	SendResponse(c, nil, AddResponse{Id: id})
	log.Info("Add table successfully.")
}
