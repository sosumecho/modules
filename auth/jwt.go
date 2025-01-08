package auth

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/sosumecho/modules/exception"
)

const (
	tokenKey   = "auth.jwt.token"
	RefreshKey = "auth.jwt.refresh"
)

var (
	invalidTokenErr = exception.NewAuthError(errors.New("invalid token"))
	expiredErr      = exception.NewAuthError(errors.New("token expired"))
)

type JwtConf struct {
	Key     string `mapstructure:"key"`
	Refresh string `mapstructure:"refresh"`
	Expired string `mapstructure:"expired"`
}

// Jwt jwt
type Jwt struct {
	conf          *JwtConf
	claims        jwt.Claims
	refreshClaims jwt.Claims
}

func (j *Jwt) SetClaims(claims jwt.Claims) *Jwt {
	j.claims = claims
	return j
}

func (j *Jwt) SetRefreshClaims(claims jwt.Claims) *Jwt {
	j.refreshClaims = claims
	return j
}

func (j *Jwt) generateToken(claims jwt.Claims) (string, exception.Exception) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(j.conf.Key))
	if err != nil {
		return "", exception.NewBuildTokenError(err)
	}
	return accessToken, nil
}

// Build  生成token
func (j *Jwt) Build() (string, string, error) {
	accessToken, err := j.generateToken(j.claims)
	if err != nil {
		return "", "", exception.NewBuildTokenError(err)
	}
	refreshToken, refreshErr := j.generateToken(j.refreshClaims)
	if refreshErr != nil {
		return "", "", exception.NewBuildTokenError(err)
	}
	return accessToken, refreshToken, nil
}

// Builder token建造者
func Builder(conf *JwtConf) *Jwt {
	return &Jwt{
		conf: conf,
	}
}

// JwtAuth  验证器
type JwtAuth struct {
	// 在整个gin.Context 上线文中的 Get 操作的key名,可以获得 AuthEntity
	ContextKey string
	// IsContinue 是否失败后继续向下执行
	IsContinue    bool
	conf          *JwtConf
	token         *jwt.Token
	claims        jwt.Claims
	refreshClaims jwt.Claims
}

// SetContextKey  设置上下文中的 key
func (j *JwtAuth) SetContextKey(contextKey string) *JwtAuth {
	j.ContextKey = contextKey
	return j
}

// SetContinue  设置失败后是否继续执行
func (j *JwtAuth) SetContinue(isContinue bool) *JwtAuth {
	j.IsContinue = isContinue
	return j
}

// parse  得到 token
func (j *JwtAuth) parse(c *gin.Context) exception.Exception {
	var err error
	j.token, err = request.ParseFromRequest(c.Request,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (i interface{}, e error) {
			return []byte(j.conf.Key), nil
		},
		request.WithClaims(j.claims))
tag:
	if !j.IsContinue && err != nil {
		if errors.Is(err, request.ErrNoTokenInRequest) {
			// 从cookie中取
			var tokenStr = ""
			if tokenStr, err = c.Cookie(fmt.Sprintf("%s_token", j.ContextKey)); err == nil {
				j.token, err = jwt.ParseWithClaims(tokenStr, j.claims, func(token *jwt.Token) (i interface{}, err error) {
					return []byte(j.conf.Key), nil
				})
				goto tag
			}
		}

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				if j.token != nil && j.IsContinue {
					return nil
				}
				return expiredErr
			} else {
				return exception.NewAuthError(err)
			}
		}
	}
	return nil
}

// ParseToken  得到 token
func (j *JwtAuth) ParseToken(token string) exception.Exception {
	var err error
	j.token, err = jwt.NewParser().ParseWithClaims(token, j.claims, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(j.conf.Key), nil
	})
	if !j.IsContinue && err != nil {
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				if j.token != nil && j.IsContinue {
					return nil
				}
				return expiredErr
			} else {
				return exception.NewAuthError(err)
			}
		}
	}
	return nil
}

// Verify  校验 token
func (j *JwtAuth) Verify(c *gin.Context) exception.Exception {
	err := j.parse(c)
	if err != nil {
		return err
	}
	if j.token != nil && j.token.Valid {
		claims := j.token.Claims
		c.Set(j.ContextKey, claims) // 向下设置用户信息,控制器可直接获取
		return nil
	}
	// 如果失败了并且token存在
	if !j.IsContinue || j.token != nil {
		c.Abort()
		return invalidTokenErr
	}
	return nil
}

func (j *JwtAuth) VerifyString(token string) (jwt.Claims, exception.Exception) {
	err := j.ParseToken(token)
	if err != nil {
		return nil, err
	}
	if j.token != nil && j.token.Valid {
		claims := j.token.Claims
		return claims, nil
	}
	return nil, invalidTokenErr
}

func (j *JwtAuth) SetClaims(claims jwt.Claims) *JwtAuth {
	j.claims = claims
	return j
}

func (j *JwtAuth) SetRefreshClaims(claims jwt.Claims) *JwtAuth {
	j.refreshClaims = claims
	return j
}

// Parser 解析 token
func Parser(conf *JwtConf) *JwtAuth {
	return &JwtAuth{
		conf: conf,
	}
}
