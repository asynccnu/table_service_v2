package table

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
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

// 获取课表
func Get(c *gin.Context) {
	var tableList = make([]*model.TableItem, 0)
	sid := c.MustGet("Sid").(string)

	tableFromXZ, err := GetFromXK(c)
	if err != nil {
		// 获取不到则查看数据库中是否有记录
		haveTable, err := model.HaveTable(sid)
		if err != nil {
			SendError(c, err, nil, err.Error())
			return

			// 没有记录则返回错误
		} else if !haveTable {
			//log.Error()
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
	if err = model.AddXKTable(sid, tableFromXZ); err != nil {
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
}

// 从服务器中获取教务课表
func GetFromXK(c *gin.Context) ([]*model.TableItem, error){
	var tableList = make([]*model.TableItem, 0)

	var r Reqeust
	if err := c.ShouldBindQuery(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return tableList, err
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(viper.GetString("data_service_url"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDataProviderClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	fmt.Println(c.MustGet("Sid").(string), c.MustGet("Password").(string))
	table, err := client.GetUndergraduateTable(ctx, &pb.GradeRequest{
		Sid:		c.MustGet("Sid").(string),
		Password: 	c.MustGet("Password").(string),
		Xqm:		r.XQM,
		Xnm:		r.XNM,
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
				return tableList, err
			}
		}
		SendError(c, err, nil, err.Error())
		return tableList, err
	}

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

	fmt.Println(tableList)

	return tableList, nil
}