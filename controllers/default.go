package controllers

import (
	"JWTLearning/models"
	"encoding/json"
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type MainController struct {
	beego.Controller
}

var (
	jwtKey = []byte("secret_key")

	users = map[string]string{
		"user1": "password1",
		"user2": "password2",
	}
)

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

func (c *MainController) Home() {
	tokenStr, existTokenStr := c.GetSecureCookie("secret_key", "token")
	CheckExists(c.Controller, existTokenStr, "500")
	claims := &models.Claims{}
	token, getTokenError := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	CheckCustomError(c.Controller, getTokenError, jwt.ErrSignatureInvalid, "Unauthorized Status", "500")
	CheckExists(c.Controller, token.Valid, "500")
	c.Ctx.ResponseWriter.Write([]byte(fmt.Sprintf("Hello %s", claims.Username)))
}

func (c *MainController) Login() {
	var credentials models.Credentials
	unmarshalError := json.Unmarshal(c.Ctx.Input.RequestBody, &credentials)
	CheckError(c.Controller, unmarshalError, "500")
	expectedPassword, ok := users[credentials.Username]
	if !ok || expectedPassword != credentials.Password {
		credentialsFail := errors.New("credentials not matched")
		CheckError(c.Controller, credentialsFail, "500")
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &models.Claims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: expirationTime,
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, generateTokenStringError := token.SignedString(jwtKey)
	CheckError(c.Controller, generateTokenStringError, "500")
	c.SetSecureCookie("secret_key", "token", tokenString)
}
func (c *MainController) Refresh() {
	tokenStr, existTokenStr := c.GetSecureCookie("secret_key", "token")
	CheckExists(c.Controller, existTokenStr, "500")
	claims := &models.Claims{}
	token, getTokenError := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	CheckCustomError(c.Controller, getTokenError, jwt.ErrSignatureInvalid, "Unauthorized Status", "500")
	CheckExists(c.Controller, token.Valid, "500")
	var timeCheckError error
	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) < 30*time.Second {
		timeCheckError = errors.New("time expiry less than 30 seconds")
	}
	CheckError(c.Controller, timeCheckError, "500")
	expirationTime := time.Now().Add(time.Minute * 5)
	claims.ExpiresAt.Time = expirationTime
	tokens := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, generateTokenStringError := tokens.SignedString(jwtKey)
	CheckError(c.Controller, generateTokenStringError, "500")
	c.SetSecureCookie("secret_key", "token", tokenString)
}
