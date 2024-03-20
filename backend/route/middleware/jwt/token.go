package jwt

import (
	"net/http"
	"strings"
	"time"

	"github.com/apicat/apicat/backend/model/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
)

// 定义过期时间
var TokenExpireDuration = time.Hour * 24 * 30

// 定义secret
var secretKey = []byte("这是一段用于生成token的密钥")

type MyClaims struct {
	UserID uint `json:"cid"`
	jwt.RegisteredClaims
}

type NotAbortPathList []*NotAbortPath

type NotAbortPath struct {
	Method []string
	Path   string
}

func Generate(userid uint) (string, error) {
	c := MyClaims{
		UserID: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			Issuer:    "apicat",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(secretKey)
}

const ctxKey = "selfuser"

func JwtUser(notAbortPathList NotAbortPathList) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token, _ := request.ParseFromRequest(
			ctx.Request,
			request.AuthorizationHeaderExtractor,
			func(t *jwt.Token) (interface{}, error) { return secretKey, nil },
		)

		if token != nil {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if _, exists := claims["cid"]; exists {
					usr := &user.User{ID: uint(claims["cid"].(float64))}
					if ok, _ := usr.Get(ctx); ok {
						ctx.Set(ctxKey, usr)
						ctx.Next()
						return
					}
				}
			}
		}

		for _, v := range notAbortPathList {
			methodMatch := false
			if len(v.Method) == 1 && strings.ToLower(v.Method[0]) == "all" {
				methodMatch = true
			} else if len(v.Method) > 0 {
				for _, m := range v.Method {
					if strings.ToUpper(m) == ctx.Request.Method {
						methodMatch = true
						break
					}
				}
			}

			if methodMatch {
				pathMatch := true
				requestPathSplit := strings.Split(ctx.Request.URL.Path, "/")
				notAbortPathSplit := strings.Split(v.Path, "/")
				if len(requestPathSplit) < len(notAbortPathSplit) {
					// 如果规则路径比请求路径长，肯定不匹配
					continue
				}
				for i, v := range requestPathSplit {
					if i >= len(notAbortPathSplit) {
						// 如果请求路由已经校验到比规则路径长还不匹配，就代表已经不匹配规则了
						pathMatch = false
						break
					}

					if notAbortPathSplit[i] == "*" {
						// 如果规则中有*，表示后面的都不用匹配了
						break
					}

					if strings.HasPrefix(notAbortPathSplit[i], ":") {
						// :开头的表示参数，不用匹配
						continue
					}

					if v != notAbortPathSplit[i] {
						pathMatch = false
						break
					}
				}
				if pathMatch {
					ctx.Next()
					return
				}
			}
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Login status has expired, please log in again.",
			"action":  "login",
		})
	}
}

func GetUser(ctx *gin.Context) *user.User {
	v, ok := ctx.Get(ctxKey)
	if ok && v != nil {
		return v.(*user.User)
	}
	return nil
}
