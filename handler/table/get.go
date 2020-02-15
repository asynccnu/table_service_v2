package table

import (
	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	"github.com/asynccnu/table_service_v2/service"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 获取课表
func Get(c *gin.Context) {
	log.Info("Get function called.")

	sid := c.MustGet("Sid").(string)
	password := c.MustGet("Password").(string)

	var tableList = make([]*model.TableItem, 0)

	tableFromCache := false
	tableFromXk, err := service.GetFromXk(c, sid, password)
	if err != nil {
		// 获取不到则查看数据库中是否有记录
		log.Warn("Can't get table from Xk")

		// 先判断是否是因为密码错误，密码错误则返回Password错误
		st, ok := status.FromError(err)
		if !ok {
			SendError(c, err, nil, err.Error())
			return
		}
		if st.Code() == codes.Unauthenticated {
			SendUnauthorized(c, errno.ErrPasswordIncorrect, nil)
			return
		}

		// 不是密码错误，尝试从缓存中获取
		haveTable, err := model.HaveTable(sid)
		if err != nil {
			// 缓存获取失败
			SendError(c, err, nil, err.Error())
			return
			// 数据库中没有则返回错误
		} else if !haveTable {
			log.Warn("Can't get table form cache.")
			SendError(c, errno.ErrNoTable, nil, "No table in database and can't get table from internet.")
			return
		}

		// 不是密码错误，数据库中也有相应记录, 则尝试从数据库中获取选课课表
		tableList, err = model.GetXkTable(sid)
		if err != nil {
			SendError(c, err, nil, err.Error())
			return
		}

		tableFromXk = tableList
		tableFromCache = true
	}

	// 将教务课表添加到数据库中
	if !tableFromCache {
		if err = model.AddXKTable(sid, tableFromXk); err != nil {
			SendError(c, err, nil, err.Error())
			return
		}
	}

	// 获取自定义课表
	userTableList, err := model.GetSelfTable(sid)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	tableList = make([]*model.TableItem, 0)
	// 加入教务课表
	tableList = append(tableList, tableFromXk...)
	// 加入用户自己的课表
	tableList = append(tableList, userTableList...)

	// table adapt
	adaptedTableList, err := adaptTableItemList(tableList)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	SendResponse(c, nil, &adaptedTableList)
	log.Info("Get table successfully.")
	return
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
		color = (color + 1) % 4
		respItems = append(respItems, &adaptItem)
	}

	response := model.TableAdaptListObject{
		Table: respItems,
	}

	return &response, nil
}
