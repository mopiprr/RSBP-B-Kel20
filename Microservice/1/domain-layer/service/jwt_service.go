package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTService interface {
	GenerateToken(userPkid string, rolePkid string, officePkid string) string
	ValidateToken(token string) (*jwt.Token, error)
	GetUserPkidByToken(token string) (int64, error)
	GenerateTokenRefresh(userPkid string, rolePkid string, officePkid string) string
}

type jwtCustomClaim struct {
	UserPkid   string `json:"user_pkid"`
	RolePkid   string `json:"role_pkid"`
	OfficePkid string `json:"office_pkid"`
	jwt.RegisteredClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "Template",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	return secretKey
}

func (j *jwtService) GenerateToken(userPkid string, rolePkid string, officePkid string) string {
	claims := jwtCustomClaim{
		userPkid,
		rolePkid,
		officePkid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 120)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return tx
}

func (j *jwtService) parseToken(t_ *jwt.Token) (any, error) {
	if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
	}
	return []byte(j.secretKey), nil
}

func (j *jwtService) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, j.parseToken)
}

func (j *jwtService) GetUserPkidByToken(token string) (int64, error) {
	t_Token, err := j.ValidateToken(token)
	if err != nil {
		return -1, err
	}

	claims := t_Token.Claims.(jwt.MapClaims)
	idStr := fmt.Sprintf("%v", claims["user_pkid"])
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (j *jwtService) GenerateTokenRefresh(userPkid string, rolePkid string, officePkid string) string {
	claims := jwtCustomClaim{
		userPkid,
		rolePkid,
		officePkid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 7)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tx, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		log.Println(err)
	}
	return tx
}
