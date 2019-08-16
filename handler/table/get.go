package table

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/asynccnu/table_service_v2/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
)

// 获取课表
func Get(c *gin.Context) {
	log.Info("Get function called.")

	sid := c.MustGet("Sid").(string)
	password := c.MustGet("Password").(string)

	var tableList = make([]*model.TableItem, 0)

	tableFromXk, err := service.GetFromXk(c, sid, password)
	if err != nil {
		// 获取不到则查看数据库中是否有记录
		haveTable, err := model.HaveTable(sid)
		if err != nil {
			SendError(c, err, nil, err.Error())
			return

			// 没有记录则返回错误
		} else if !haveTable {
			SendError(c, errno.ErrDatabase, nil, "No table in database.")
			return
		}

		// 有记录就返回课表（教务和自加）
		tableList, err = model.GetTable(sid)
		if err != nil {
			SendError(c, err, nil, err.Error())
			return
		}

		SendResponse(c, nil, &tableList)
		return
	}

	// 将教务课表添加到数据库中
	if err = model.AddXKTable(sid, tableFromXk); err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	// 获取自定义课表
	tableList, err = model.GetSelfTable(sid)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	// 加入教务课表
	tableList = append(tableList, tableFromXk...)

	SendResponse(c, nil, &tableList)
	log.Info("Get table successfully.")
}
