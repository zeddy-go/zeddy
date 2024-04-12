package ginx

import (
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/zeddy-go/zeddy/errx"
	"github.com/zeddy-go/zeddy/tool/jwt"
	"net/http"
	"strings"
)

func NewJwtAuthMiddlewareBuilder() *JwtAuthMiddlewareBuilder {
	return &JwtAuthMiddlewareBuilder{
		header: "Authorization",
		query:  "token",
		idKey:  jwt.JwtSubject,
	}
}

type JwtAuthMiddlewareBuilder struct {
	header string                           //头字段
	query  string                           //querystring字段
	salt   string                           //盐
	idKey  string                           //用户唯一标识
	check  func(claims jwt2.MapClaims) bool //判断是否通过
}

func (j *JwtAuthMiddlewareBuilder) SetHeader(header string) *JwtAuthMiddlewareBuilder {
	j.header = header
	return j
}

func (j *JwtAuthMiddlewareBuilder) SetQuery(queryField string) *JwtAuthMiddlewareBuilder {
	j.query = queryField
	return j
}

func (j *JwtAuthMiddlewareBuilder) SetSalt(salt string) *JwtAuthMiddlewareBuilder {
	j.salt = salt
	return j
}

func (j *JwtAuthMiddlewareBuilder) SetIdKey(idKey string) *JwtAuthMiddlewareBuilder {
	j.idKey = idKey
	return j
}

func (j *JwtAuthMiddlewareBuilder) SetUserIdentifier(f func(claims jwt2.MapClaims) bool) *JwtAuthMiddlewareBuilder {
	j.check = f
	return j
}

func (j *JwtAuthMiddlewareBuilder) Build() func(*gin.Context) error {
	return func(c *gin.Context) (err error) {
		var token string
		if t := c.Request.Header.Get(j.header); t != "" {
			arr := strings.Split(t, " ")
			if len(arr) < 2 {
				return errx.New("token格式错误", errx.WithCode(http.StatusUnauthorized), errx.WithAbort())
			}
			token = arr[1]
		} else if t := c.Query(j.query); t != "" {
			token = t
		} else {
			return errx.New("无有效token", errx.WithCode(http.StatusUnauthorized), errx.WithAbort())
		}

		t, err := jwt.Parse(token, j.salt)
		if err != nil {
			return errx.New(err.Error(), errx.WithCode(http.StatusUnauthorized), errx.WithAbort())
		}
		if !t.Valid {
			return errx.New("无效的token", errx.WithCode(http.StatusUnauthorized), errx.WithAbort())
		}

		claims, ok := t.Claims.(jwt2.MapClaims)
		if !ok {
			return errx.New("解析token失败", errx.WithCode(http.StatusUnauthorized), errx.WithAbort())
		}

		if j.check != nil && !j.check(claims) {
			return errx.New("未找到用户", errx.WithCode(http.StatusUnauthorized), errx.WithAbort())
		}

		c.Set("claims", claims)

		return nil
	}
}
