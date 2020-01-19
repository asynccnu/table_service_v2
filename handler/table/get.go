package table

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/asynccnu/table_service_v2/service"
	"strconv"

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
// 1. Start During 改为了int32  ok
// 2. day 改为了星期一 --- 星期日
// 3. 增加了Source属性  ok
// 4. 增加了Color属性  color可取0 1 2 3 ok
func adaptTableItemList(items []*model.TableItem) ([]*model.TableAdaptItem, error) {
	respItems := make([]*model.TableAdaptItem, 0)

	color := 0
	daysMap := map[string]string {
		"1":"星期一",
		"2":"星期二",
		"3":"星期三",
		"4":"星期四",
		"5":"星期五",
		"6":"星期六",
		"7":"星期日",
	}

	for _, item := range items {
		startInt32, err := strconv.Atoi(item.Start)
		if err != nil{
			return nil, err
		}

		DuringInt32, err := strconv.Atoi(item.During)
		if err != nil{
			return nil, err
		}

		dayInWeek := ""
		if v, ok := daysMap[item.Day]; ok {
			dayInWeek = v
		} else {
			return nil, errno.ErrWeekConvert
		}

		adaptItem := model.TableAdaptItem{
			Id:      item.Id,
			Course:  item.Course,
			Teacher: item.Teacher,
			Place:   item.Place,
			Start:   int32(startInt32),
			During:  int32(DuringInt32),
			Day:     dayInWeek,
			Source:  item.Source,
			Weeks:   item.Weeks,
			Remind:  item.Remind,
			Color:   int32(color),
		}
		// color 在0 1 2 3之内取值
		color = (color+1)%4
		respItems = append(respItems, &adaptItem)
	}

	return respItems, nil
}