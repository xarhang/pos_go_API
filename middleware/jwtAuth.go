package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type JwtAuth struct {
}
type MyCustomClaims struct {
	Username string `json:"username"`
	RuleID   int    `json:"rule_id"`
	jwt.StandardClaims
}

func Authorize(obj string, act string, adapter *gormadapter.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user/subject
		val, existed := c.Get("username")
		if !existed {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user hasn't logged in yet"})
			return
		}
		// Casbin enforces policy
		ok, err := enforce(val.(string), obj, act, adapter)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": "error occurred when authorizing user"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden", "OK": ok})
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string, adapter *gormadapter.Adapter) (bool, error) {
	// Load model configuration file and policy store adapter
	enforcer, err := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	if err != nil {
		return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}
	// Load policies from DB dynamically
	err = enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	// Verify
	ok, err := enforcer.Enforce(sub, obj, act)

	return ok, err
}
func AuthorizationJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.Request.Header.Get("Authorization")
		if auth == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Authorization Header Not Found"})
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
		ctx.Set("username", claims.Username)
		ctx.Set("rule_id", claims.RuleID)
		ctx.Next()

	}

}
