package util

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/teris-io/shortid"
)

var envXn string
var envXq string

func init() {
	// xn xq 当代码不准确时使用环境变量用来手动调整学年学期
	// export CCNUBOX_TABLE_XN=2018
	// export CCNUBOX_TABLE_XQ=3
	// 当xn或xq为 空字符串 "" 时忽略环境变量 即不设置该环境变量时忽略
	envXn = viper.GetString("table.xn")
	envXq = viper.GetString("table.xq")
}

func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestID, ok := v.(string); ok {
		return requestID
	}
	return ""
}

// 获取学年和学期
func GetXnAndXq() (string, string) {
	var xqn = "3"

	now := time.Now()
	year := now.Year()
	month := now.Month()

	if month <= 7 {
		year--
	}

	// 第三学期的时间还有点问题
	if month >= 2 && month < 7 {
		xqn = "12"
	} else if month == 7 {
		xqn = "16"
	}

	xn := strconv.Itoa(year)
	xq := xqn

	// 使用环境变量对自动获取的学年名学期名进行兜底
	if envXn != "" {
		xn = envXn
	}
	if envXq != "" {
		xq = envXq
	}

	return xn, xq
}
