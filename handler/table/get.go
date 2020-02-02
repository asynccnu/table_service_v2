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
	usertableList, err := model.GetSelfTable(sid)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}


	tableList = make([]*model.TableItem, 0)
	// 加入教务课表
	tableList = append(tableList, tableFromXk...)
	// 加入用户自己的课表
	tableList = append(tableList, usertableList...)

	// table adapt

	adaptedTableList, err := adaptTableItemList(tableList)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	SendResponse(c, nil, &adaptedTableList)
	log.Info("Get table successfully.")
}


// 与原服务格式统一
// 与TableItem的不同之处:
// 1. 增加了Source属性  ok
// 2. 增加了Color属性  color可取0 1 2 3 ok
func adaptTableItemList(items []*model.TableItem) (*model.TableAdaptListObject, error) {
	respItems := make([]*model.TableAdaptItem, 0)

	color := 0

	for _, item := range items {
		adaptItem := model.TableAdaptItem{
			Id:      item.Id,
			Course:  item.Course,
			Teacher: item.Teacher,
			Place:   item.Place,
			Start:   item.Start,
			During:  item.During,
			Day:     item.Day,
			Source:  item.Source,
			Weeks:   item.Weeks,
			Remind:  item.Remind,
			Color:   int32(color),
		}
		// color 在0 1 2 3之内取值
		color = (color+1)%4
		respItems = append(respItems, &adaptItem)
	}

	response := model.TableAdaptListObject{
		Table: respItems,
	}

	return &response, nil
}