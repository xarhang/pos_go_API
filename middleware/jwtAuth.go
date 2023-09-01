package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JwtAuth struct {
}
type MyCustomClaims struct {
	// Username string `json:"username"`
	RuleID int `json:"rule_id,omitempty"`
	jwt.StandardClaims
}

func AuthorizationJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"}) //Authorization Header Not Found"})
			return
		}
		// fmt.Println(auth)
		if !strings.HasPrefix(auth, "Bearer") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Authorization Header is incorrect"})
			return
		}
		splitToken := strings.Split(auth, "Bearer ")
		auth = splitToken[1]
		token, err := jwt.ParseWithClaims(auth, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRETKEY")), nil
		})
		// if token.Method.Alg() != "HS512" {
		// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Incorect ALG"})
		// 	return
		// }
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": err.Error()})
			return
		}

		claims, ok := token.Claims.(*MyCustomClaims)
		if ok && token.Valid {
			// fmt.Println(claims.Username, claims.ExpiresAt)
			if claims.ExpiresAt <= time.Now().Unix() {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Token Expired"})
				return
			}
		}
		ctx.Set("aud", claims.Audience)
		ctx.Set("rule_id", claims.RuleID)
		ctx.Next()

	}

}
