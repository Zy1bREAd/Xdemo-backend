package middleware

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const SECRETKEY string = "oceanwang"

type UserClaim struct {
	OrginalClaims jwt.RegisteredClaims
	UserID        uint   `json:"user_id"`
	Account       string `json:"account"`
	// Name          string `json:"username"`
	// Team          string `json:"team,omitempty"`
	Status string `json:"status"`
}

func (uc UserClaim) GetAudience() (jwt.ClaimStrings, error) {
	return uc.OrginalClaims.Audience, nil
}

func (uc UserClaim) GetExpirationTime() (*jwt.NumericDate, error) {
	return uc.OrginalClaims.ExpiresAt, nil
}

func (uc UserClaim) GetNotBefore() (*jwt.NumericDate, error) {
	return uc.OrginalClaims.NotBefore, nil
}

func (uc UserClaim) GetIssuedAt() (*jwt.NumericDate, error) {
	return uc.OrginalClaims.IssuedAt, nil
}

func (uc UserClaim) GetIssuer() (string, error) {
	return uc.OrginalClaims.Issuer, nil
}

func (uc UserClaim) GetSubject() (string, error) {
	return uc.OrginalClaims.Subject, nil
}

// 生成JWT
func GenerateJWT(id uint, account string, status string) (string, error) {
	myClaim := UserClaim{
		UserID:  id,
		Account: account,
		// Name:    userObj.Name,
		// Team:    userObj.Team,
		Status: status,
		OrginalClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 60 * time.Minute)), // 一天后过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "XDemo_Ocean",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaim)
	tokenStr, err := token.SignedString([]byte(SECRETKEY)) // 通常要根据加密签名的算法决定传入参数的类型，但一般都是[]byte
	if err != nil {
		// 这里应该组装好错误信息，
		log.Fatal("JWT签名加密时发生错误 ", err)
		return "", err
	}
	return tokenStr, nil
}

// 解析并验证JWT
func ParseJWT(tokenStr string, secret string, claims jwt.Claims) (*UserClaim, error) {
	// 解析
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRETKEY), nil
	})
	if err != nil {
		log.Println("解析JWT时发生错误", err)
		return nil, err
	}
	// 验证jwt的声明；通过转换jwt声明转换成我们自定义声明是否正确，并且判断token是否有效。
	if claims, ok := token.Claims.(*UserClaim); ok && token.Valid {
		return claims, nil
	} else {
		log.Println("验证JWT Token失败,", err)
		return nil, err
	}
}

// 验证用户Token是否合法
func IsUserTokenVaild(tokenStr string) error {
	_, err := ParseJWT(tokenStr, SECRETKEY, UserClaim{})
	return err
}
