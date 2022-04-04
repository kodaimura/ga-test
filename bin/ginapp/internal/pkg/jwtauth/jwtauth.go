package jwtauth

import (
    "os"
    "log"
    "encoding/json"
    "time"
    //"strings"
    "errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)


const JwtKeyName string = "token"
const JwtExpires time.Duration = 108000
var secretKey string

func init () {
	err := godotenv.Load(".env")

	if err != nil {
        log.Panic(err)
    }

	secretKey = os.Getenv("JWT_SECRET_KEY")
}


type JwtPayload struct {
	UId int `json:"uid"`
    UserName string `json:"username"`
    jwt.StandardClaims
}


func ExtractUId (c *gin.Context) (int, error) {
	payload := c.Keys["payload"]
	if payload == nil {
		return -1, errors.New("ExtractUId error")
	} else {
		return payload.(JwtPayload).UId, nil
	}
	
}


func ExtractUserName (c *gin.Context) (string, error) {
	payload := c.Keys["payload"]
	if payload == nil {
		return "", errors.New("ExtractUserName error")
	} else {
		return payload.(JwtPayload).UserName, nil
	}
}


func GenerateJWT(uid int, userName string) (string, error) {
	payload := JwtPayload{
        UId: uid,
        UserName: userName,
    }

    payload.IssuedAt =  time.Now().Unix()
    payload.ExpiresAt = time.Now().Add(time.Second * JwtExpires).Unix()

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

    return token.SignedString([]byte(secretKey))
}


func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := jwtAuth(c)

		if err != nil {
			c.Redirect(303, "/login")
			c.Abort()
			return
		}
		c.Set("payload", payload)
		c.Next()
	}
}


func JwtApiAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := jwtAuth(c)

		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("payload", payload)
		c.Next()
	}
}


func jwtAuth (c *gin.Context) (JwtPayload, error) {
	var payload JwtPayload

	tokenString, err := extractTokenString(c)
	if err != nil {
		return payload, err
	}
	token, err := toToken(tokenString)
	if err != nil {
		return payload, err
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
	return c.Cookie(JwtKeyName)
} 


func toToken (tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	return token, err
} 


func extractPayload (token *jwt.Token) (JwtPayload, error) {
	var payload JwtPayload

	jsonString, err := json.Marshal(token.Claims.(jwt.MapClaims))

    if err == nil {
        err = json.Unmarshal(jsonString, &payload)
    }

    return payload, err
} 
