package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	. "github.com/asynccnu/table_service_v2/handler"
	"github.com/asynccnu/table_service_v2/model"
	"github.com/asynccnu/table_service_v2/pkg/errno"
	pb "github.com/asynccnu/table_service_v2/rpc"
	"github.com/asynccnu/table_service_v2/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 从服务器中获取教务课表
func GetFromXk(c *gin.Context, sid, password string) ([]*model.TableItem, error) {
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

	// 获取学年和学期
	xn, xqm := util.GetXnAndXq()

	table, err := client.GetUndergraduateTable(ctx, &pb.GradeRequest{
		Sid:      sid,
		Password: password,
		Xqm:      xqm,
		Xnm:      xn,
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
	for index, item := range table.Lists {
		t, err := model.Process(&model.TableRowItem{
			Kcmc: item.Kcmc,
			Zcd:  item.Zcd,
			Jcor: item.Jcor,
			Cdmc: item.Cdmc,
			Xm:   item.Xm,
			Xqj:  item.Xqj,
		})
		if err != nil {
			return tableList, err
		}
		t.Id = strconv.Itoa(index)
		t.Source = "xk"
		tableList = append(tableList, &t)
	}

	return tableList, nil
}
