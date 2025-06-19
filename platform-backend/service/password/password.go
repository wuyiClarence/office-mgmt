package password

import (
	"crypto/des"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func CheckPassWord(password string) error {
	num := `[0-9]{1}`
	a_z := `[a-z]{1}`
	A_Z := `[A-Z]{1}`

	err := errors.New("密码应该包含大小写字母、数字，且长度大于等于6位")
	if len(password) < 6 {
		return err
	}

	if b, err := regexp.MatchString(num, password); !b || err != nil {
		return err
	}

	if b, err := regexp.MatchString(a_z, password); !b || err != nil {
		return err
	}

	if b, err := regexp.MatchString(A_Z, password); !b || err != nil {
		return err
	}

	if strings.Contains(password, " ") {
		return errors.New("密码不应该有空格")
	}

	return nil
}

// Encrypt encrypts password using DES-DCB.
func Encrypt(key, password string) (string, error) {
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	mode := ecb.NewECBEncrypter(block)
	p := padding.NewPkcs7Padding(mode.BlockSize())

	pt := []byte(password)
	pt, err = p.Pad(pt)
	if err != nil {
		return "", err
	}

	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return base64.StdEncoding.EncodeToString(ct), nil
}

// Decrypt decrypts password using DES-DCB.
func Decrypt(key, ciperPassword string) (string, error) {
	block, err := des.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	mode := ecb.NewECBDecrypter(block)
	p := padding.NewPkcs7Padding(mode.BlockSize())

	ct, err := base64.StdEncoding.DecodeString(ciperPassword)
	if err != nil {
		return "", err
	}
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)

	pt, err = p.Unpad(pt)
	if err != nil {
		return "", err
	}
	return string(pt), err
}
