package user

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/asynccnu/table_service_v2/service"

	"github.com/gin-gonic/gin"
)

// List list the users in the database.
func List(c *gin.Context) {
	var r ListRequest
	if err := c.Bind(&r); err != nil {
		SendBadRequest(c, errno.ErrBind, nil, err.Error())
		return
	}

	infos, count, err := service.ListUser(r.Username, r.Offset, r.Limit)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	SendResponse(c, nil, ListResponse{
		TotalCount: count,
		UserList:   infos,
	})
}
