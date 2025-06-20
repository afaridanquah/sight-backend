package otpbus

import (
	"crypto/rand"
	"math/big"

	"bitbucket.org/msafaridanquah/verifylab-service/business/domain/otpbus/valueobject"
)

func generateOTPCode() (string, valueobject.HashCode, error) {
	const charset = "0123456789"
	charsetLen := big.NewInt(int64(len(charset)))
	otp := make([]byte, 6)
	for i := range otp {
		randInt, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", valueobject.HashCode{}, err
		}
		otp[i] = charset[randInt.Int64()]
	}

	hashed, err := valueobject.ParseToHashCode(string(otp))
	if err != nil {
		return "", valueobject.HashCode{}, err
	}

	return string(otp), hashed, nil
}
