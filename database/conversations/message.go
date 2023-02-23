package conversations

import (
	"time"
	"unsafe"

	"github.com/golang-jwt/jwt/v5"
)

type Message struct {
	ID string `json:"id" gorm:"primaryKey"`

	Conversation uint   `json:"conversation" gorm:"not null"`
	Certificate  string `json:"certificate" gorm:"not null"`
	Creation     int64  `json:"creation" gorm:"autoUpdateTime:milli"` // Unix timestamp (ms)
	Data         string `json:"data" gorm:"not null"`                 // Encrypted data
}

var MessageOwnerKey = "KVA83skdBfR0z3RyrRSNfAk76T1APKl0ep10YzrmEcD8vcrVUS8WO39sHWTuxl6OJHiZ8gIbZcGOi0FinSomHDbTnjGIn7oBadlNX5XmSEsgxagCsAVD3oCgacPGpCI4OOefRynmp4abUwQTaCakQv9PkgjgDwyZ3w2APH47sy4sFsKhX80LbfYKapxpnC4U61PVZnPI"

func CheckSize(message string) bool {
	return unsafe.Sizeof(message) > 1000*24
}

type CertificateClaims struct {
	MID string `json:"mid"` // Message ID
	Sd  int64  `json:"sd"`  // Sender ID
	Ct  int64  `json:"ct"`  // Creation time
	jwt.RegisteredClaims
}

func GenerateCertificate(id string, sender int64) (string, error) {

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, CertificateClaims{
		MID: id,
		Sd:  sender,
		Ct:  time.Now().Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "chat-node",
		},
	})

	token, err := tk.SignedString([]byte(MessageOwnerKey))

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetCertificateClaims(certificate string) (*CertificateClaims, bool) {

	token, err := jwt.ParseWithClaims(certificate, &CertificateClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(MessageOwnerKey), nil
	}, jwt.WithLeeway(5*time.Minute))

	if err != nil {
		return &CertificateClaims{}, false
	}

	if claims, ok := token.Claims.(*CertificateClaims); ok && token.Valid {
		return claims, true
	}

	return &CertificateClaims{}, false
}

func (m *CertificateClaims) Valid(id string, sender int64) bool {
	return m.MID == id && m.Sd == sender
}
