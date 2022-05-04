package jwt

import (
    "os"
    "log"
    "encoding/json"
    //"strings"
    "errors"

	"github.com/gin-gonic/gin"
	jwtpackage "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

const JwtKeyname string = "access_token"
var secretKey string

func init () {
	err := godotenv.Load(".env")

	if err != nil {
        log.Panic(err)
    }

	secretKey = os.Getenv("JWT_SECRET_KEY")
}


func GenerateJWT(uid int) (string, error) {
	pl, err := GeneratePayload(uid)

	if err != nil {
		return "", errors.New("GenerateJWT error")
	}

    token := jwtpackage.NewWithClaims(jwtpackage.SigningMethodHS256, pl)

    return token.SignedString([]byte(secretKey))
}


func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pl, err := jwtAuth(c)

		if err != nil {
			c.Redirect(303, "/login")
			c.Abort()
			return
		}
		c.Set("payload", pl)
		c.Next()
	}
}


func JwtApiAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pl, err := jwtAuth(c)

		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("payload", pl)
		c.Next()
	}
}


func jwtAuth (c *gin.Context) (JwtPayload, error) {
	var pl JwtPayload

	tokenString, err := extractTokenString(c)
	if err != nil {
		return pl, err
	}
	token, err := toToken(tokenString)
	if err != nil {
		return pl, err
	}

	return extractPayload(token)
}

/*
func extractTokenString (c *gin.Context) (string, error) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		return "", errors.New("Hint: Authorization")
	}else if strings.Index(tokenString, "Bearer ") != 0 {
		return "", errors.New("Hint: Bearer")
	}

	return strings.TrimSpace(tokenString[7:]), nil
} 
*/

func extractTokenString (c *gin.Context) (string, error) {
	return c.Cookie(JwtKeyname)
} 


func toToken (tokenString string) (*jwtpackage.Token, error) {
	token, err := jwtpackage.Parse(tokenString, func(token *jwtpackage.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtpackage.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	return token, err
} 


func extractPayload (token *jwtpackage.Token) (JwtPayload, error) {
	var pl JwtPayload

	jsonString, err := json.Marshal(token.Claims.(jwtpackage.MapClaims))

    if err == nil {
        err = json.Unmarshal(jsonString, &pl)
    }

    return pl, err
} 
