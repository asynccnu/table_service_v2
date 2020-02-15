package table

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteItem struct {
	Sid string `json:"sid" bson:"sid"`
	Id  string `json:"id" bson:"id" binding:"required"`
}

func Delete(c *gin.Context) {
	log.Info("Delete function called.")

	sid := c.MustGet("Sid").(string)
	id := c.Query("id")

	if id == "" {
		SendBadRequest(c, errno.ErrBind, nil, "No id.")
		return
	}

	// 教务处的课表不可删除，只能删除自定义课程
	// 验证id是否属于自加的课程
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		SendBadRequest(c, errno.ErrDeleteXKTable, nil, "id error")
		return
	}

	if delCount, err := model.DeleteTable(sid, id, objId); err != nil {
		log.Error("DeleteTable function error", err)
		SendError(c, errno.ErrDeleteTable, nil, err.Error())
		return
	} else if delCount == 0 {
		SendBadRequest(c, errno.ErrBind, nil, "This table does not exist.")
		return
	}

	SendResponse(c, nil, nil)
	log.Info("Delete table successfully.")
}
