package table

import (
	"context"
	"encoding/base64"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
	"time"

	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	pb "github.com/asynccnu/table_service_v2/rpc"
	"github.com/lexkong/log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// 学期和学年
var (
	Xq = "1"
	Xn = "2019"
)

// 获取课表
func Get(c *gin.Context) {
	log.Info("Get function called.")

	sid := c.GetHeader("sid")

	bs, err := base64.StdEncoding.DecodeString(c.GetHeader("Authorization"))

	if err != nil {
		SendBadRequest(c, err, nil, "Base64 decode error.")
		return
	}

	arr := strings.Split(string(bs), ":")
	password := arr[1]
	//fmt.Println(sid, password)

	if password == "" {
		SendBadRequest(c, errno.ErrTokenInvalid, nil, "No password")
		return
	}

	var tableList = make([]*model.TableItem, 0)

	tableFromXk, err := GetFromXk(c, sid, password)
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

	// 获取课表（教务和自加）
	tableList, err = model.GetTable(sid)
	if err != nil {
		SendError(c, err, nil, err.Error())
		return
	}

	SendResponse(c, nil, &tableList)
	log.Info("Get table successfully.")
}

// 从服务器中获取教务课表
func GetFromXk(c *gin.Context, sid, password string) ([]*model.TableItem, error){
	var tableList = make([]*model.TableItem, 0)

	// Set up a connection to the server.
	conn, err := grpc.Dial(viper.GetString("data_service_url"), grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect", err)
		return tableList, err
	}
	defer conn.Close()

	client := pb.NewDataProviderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	table, err := client.GetUndergraduateTable(ctx, &pb.GradeRequest{
		Sid:		sid,
		Password:	password,
		Xqm:		Xq,
		Xnm:		Xn,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.Unauthenticated {
				c.JSON(http.StatusOK, Response{
					Code:    errno.ErrPasswordIncorrect.Code,
					Message: st.Message(),
					Data:    nil,
				})
				return nil, err
			}
		}
		SendError(c, err, nil, err.Error())
		return nil, err
	}

	// 获取加工后的课表
	for _, item := range table.Lists {
		t, err := model.Process(&model.TableRowItem{
			Kcmc: item.Kcmc,
			Zcd: item.Zcd,
			Jcor: item.Jcor,
			Cdmc: item.Cdmc,
			Xm: item.Xm,
			Xqj: item.Xqj,
		})
		if err != nil {
			return tableList, err
		}
		tableList = append(tableList, &t)
	}

	// 测试输出
	//for _, item := range tableList {
	//	fmt.Println(*item)
	//}

	return tableList, nil
}