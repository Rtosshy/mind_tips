package auth

import (
	"database/sql"
	"log"
	"mind_tips/internal/models"
	"mind_tips/internal/utils"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var IdentityKey = "user_name"

func NewJWTMiddleware(db *sql.DB) (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("your_secret_key_here"),
		Timeout:     time.Hour * 24,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]interface{}); ok {
				log.Println("PayloadFunc succeeded")
				ret := jwt.MapClaims{IdentityKey: v[IdentityKey]}
				log.Println(ret)
				return jwt.MapClaims{
					IdentityKey: v[IdentityKey],
				}
			}
			log.Println("PayloadFunc failed")
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return map[string]interface{}{
				IdentityKey: claims[IdentityKey],
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var user models.UserAuth
			if err := c.ShouldBindJSON(&user); err != nil {
				return nil, jwt.ErrMissingLoginValues
			}
			// 本来ここでデータベースからのユーザー認証を行います
			var storedUser models.UserAuth
			err := db.QueryRow("SELECT user_name, password FROM user WHERE user_name = ?", user.UserName).Scan(&storedUser.UserName, &storedUser.Password)
			if err != nil {
				if err == sql.ErrNoRows {
					utils.LogError(c, http.StatusUnauthorized, err, "Invalid name or password")
					return nil, jwt.ErrFailedAuthentication
				}
				utils.LogError(c, http.StatusInternalServerError, err, "Failed to retrieve user")
				return nil, jwt.ErrFailedAuthentication
			}
			// パスワードの比較
			err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
			if err != nil {
				utils.LogError(c, http.StatusUnauthorized, err, "Invalid name or password")
				return nil, jwt.ErrFailedAuthentication
			}
			log.Printf("User %s logged in successfully", storedUser.UserName)
			// 認証成功、ユーザー情報を返す
			var data = map[string]interface{}{
				IdentityKey: storedUser.UserName,
			}
			log.Printf("Authenticator returning data: %#v", data)
			return data, nil

			// 単純な実装
			// return storedUser.UserName, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Printf("Data received in Authorizator: %#v\n", data)
			// dataがmap[string]interface{}であるか確認
			v, ok := data.(map[string]interface{})
			log.Println(v[IdentityKey], ok)
			if v, ok := data.(map[string]interface{}); ok {
				userName, exists := v[IdentityKey]
				if !exists {
					log.Println("Authorization failed: user_name not found in claims")
					return false
				}
				log.Printf("Authorization succeeded for user: %v", userName)
				return true
			}
			log.Println("Authorization failed: invalid data format")
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			utils.LogError(c, code, nil, message)
		},
		TokenLookup:    "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:  "Bearer",
		SendCookie:     true,
		CookieName:     "jwt",
		CookieHTTPOnly: true,
		TimeFunc:       time.Now,
	})
}
