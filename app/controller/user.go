package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/rudirahardian/go_env/app/service"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserForm struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
	Name string `form:"name" json:"name"`
	Foto gin.Context `form:"foto" json:"foto"`
}

type Result struct {
	Token string `json:"token"`
	Data service.Claims `json:"data"`
}

func V1UserLogin(c *gin.Context) {
	var user Credential
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "can't bind struct",
		})
	}

	users, err := service.FindUser(user.Username,user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "wrong username or password",
		})
		return
	}
	if len(users) == 0{
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "wrong username or password",
		})
		return
	}else{
		data := users[0]

		expirationTime := time.Now().Add(5 * time.Minute)
		claims := &service.Claims{
			Username: data.Username,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"token": tokenString,
		})
	}
}

func V1UserRegister(c *gin.Context) {
	request := c
	user, err := service.InsertUser(request);

	if err != nil{
		c.JSON(400, gin.H{"message": err.Error(), "data": user})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"message": "success", "data": user})
}

func V1UserGet(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	data,_ := service.ExtractClaims(tokenString)
	value := *data
	result := Result {
        Token : tokenString,
        Data : value,
    }

	c.JSON(http.StatusCreated, result)
}