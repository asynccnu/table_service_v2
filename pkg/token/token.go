package token

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	// ErrMissingHeader means the `Authorization` header was empty.
	ErrMissingHeader = errors.New("The length of the `Authorization` header is zero.")
)

// Context is the context of the JSON web token.
type Context struct {
	Sid      string
	Password string
}

// ParseRequest gets the token from the header and
// pass it to the Parse function to parses the token.
func ParseRequest(c *gin.Context) error {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		return ErrMissingHeader
	}

	var t string
	// Parse the header to get the token part.
	_, _ = fmt.Sscanf(header, "Basic %s", &t)

	sDec, _ := base64.StdEncoding.DecodeString(t)

	str := string(sDec)
	i := strings.Index(string(sDec), ":")
	if i < 0 {
		return ErrMissingHeader
	}

	c.Set("Sid", str[:i])
	c.Set("Password", str[i+1:])

	return nil
}
