package totp

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"net/url"
	"time"
)

func GenerateSecret() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate totp secret: %w", err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b), nil
}

func Validate(secret string, code string) bool {
	if len(code) != 6 {
		return false
	}

	secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return false
	}

	counter := time.Now().Unix() / 30

	// Check current and adjacent windows
	for _, c := range []int64{counter - 1, counter, counter + 1} {
		if computeTOTP(secretBytes, c) == code {
			return true
		}
	}

	return false
}

func computeTOTP(secret []byte, counter int64) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(counter))

	mac := hmac.New(sha1.New, secret)
	mac.Write(buf)
	hash := mac.Sum(nil)

	offset := hash[len(hash)-1] & 0x0f
	binary := binary.BigEndian.Uint32(hash[offset : offset+4])
	otp := (binary & 0x7fffffff) % 1000000

	return fmt.Sprintf("%06d", otp)
}

func QRURL(secret, username, issuer string) string {
	u := url.URL{
		Scheme: "otpauth",
		Host:   "totp",
		Path:   fmt.Sprintf("/%s:%s", issuer, username),
		RawQuery: url.Values{
			"secret": {secret},
			"issuer": {issuer},
			"algorithm": {"SHA1"},
			"digits":    {"6"},
			"period":    {"30"},
		}.Encode(),
	}
	return u.String()
}

func FormatSecretForDisplay(secret string) string {
	s, _ := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	return fmt.Sprintf("%x", s)
}

func MaskSecret(secret string) string {
	if len(secret) < 8 {
		return "****"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// FormatRecoveryCodes generates human-readable recovery codes.
func GenerateRecoveryCodes() ([]string, error) {
	codes := make([]string, 8)
	for i := range codes {
		b := make([]byte, 5)
		if _, err := rand.Read(b); err != nil {
			return nil, err
		}
		n := binary.BigEndian.Uint64(append([]byte{0, 0, 0}, b...))
		codes[i] = fmt.Sprintf("%010d", n%uint64(math.Pow10(10)))
	}
	return codes, nil
}
