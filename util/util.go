package util

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

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

	//fmt.Println(year, xq)
	return strconv.Itoa(year), xqn
}
