package conversations

import (
	"os"
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
)

type Message struct {
	ID string `json:"id" gorm:"primaryKey"`

	Conversation string `json:"conversation" gorm:"not null"`
	Certificate  string `json:"certificate" gorm:"not null"`
	Creation     int64  `json:"creation"`               // Unix timestamp (SET BY THE CLIENT, EXTREMELY IMPORTANT FOR SIGNATURES)
	Data         string `json:"data" gorm:"not null"`   // Encrypted data
	Edited       bool   `json:"edited" gorm:"not null"` // Edited flag
	Sender       string `json:"sender" gorm:"not null"` // Sender ID (of conversation token)
}

func CheckSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*12
}

type CertificateClaims struct {
	MID string `json:"mid"` // Message ID
	Sd  string `json:"sd"`  // Sender ID
	Ct  int64  `json:"ct"`  // Creation time
	jwt.RegisteredClaims
}

func GenerateCertificate(id string, sender string) (string, error) {

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, CertificateClaims{
		MID: id,
		Sd:  sender,
		Ct:  time.Now().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "chat-node",
		},
	})

	token, err := tk.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetCertificateClaims(certificate string) (*CertificateClaims, bool) {

	token, err := jwt.ParseWithClaims(certificate, &CertificateClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithLeeway(5*time.Minute))

	if err != nil {
		return &CertificateClaims{}, false
	}

	if claims, ok := token.Claims.(*CertificateClaims); ok && token.Valid {
		return claims, true
	}

	return &CertificateClaims{}, false
}

func (m *CertificateClaims) Valid(id string, sender string) bool {
	return m.MID == id && m.Sd == sender
}
