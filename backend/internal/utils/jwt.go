   // utils/jwt.go
   package utils

   import (
       "fmt"
       "os"
       "time"
       "github.com/golang-jwt/jwt/v4"
   )

   

   var JWTSecret = []byte(getJWTSecret())

   type Claims struct {
       UserID uint `json:"user_id"` 
       jwt.RegisteredClaims
   }

   func (c *Claims) Valid() error {
       return c.RegisteredClaims.Valid()
   }

   func getJWTSecret() string {
       secret := os.Getenv("JWT_SECRET")
       if secret == "" {
           return "your-secret-key"
       }
       return secret
   }

   func GenerateToken(userID uint) (string, error) {
       claims := &Claims{
           UserID: userID,
           RegisteredClaims: jwt.RegisteredClaims{
               ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
               IssuedAt:  jwt.NewNumericDate(time.Now()),
               NotBefore: jwt.NewNumericDate(time.Now()),
           },
       }

       token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
       return token.SignedString(JWTSecret)
   }

   func ValidateToken(tokenString string, secret []byte, method jwt.SigningMethod) (*Claims, error) {
       token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
           if token.Method != method {
               return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
           }
           return secret, nil
       })

       if err != nil {
           return nil, err
       }

       if claims, ok := token.Claims.(*Claims); ok && token.Valid {
           return claims, nil
       }

       return nil, fmt.Errorf("invalid token")
   }
