package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrHandle = errors.New("Invalid Token")

type Claims struct {
	UserID               uint   `json:"user_id"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        //lib jwt untuk register struct default from jwt
}

type Manager struct {
	Secret    []byte `json:"uid"` //byte itu angka dengan maksimal 255
	AccessTTL time.Duration
}

func newManager(secret string, accessTTL time.Duration) *Manager {
	return &Manager{Secret: []byte(secret), AccessTTL: time.Duration(accessTTL)}
}

// GenerateAccess membuat jwt berumur pendek.
func (m *Manager) GenerateAccess(userID uint, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10), // merubah userid menjadi string karena jwt butuh string dan Go tidak bisa otomatis converse
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.AccessTTL)),
		},
	}
	// 	jwt.NewWithClaims(HS256, claims) membuat token berisi claims, ditandatangani dengan algoritma HS256.
	// .SignedString(m.secret) menandatangani memakai secret (JWT_SECRET) dan menghasilkan string xxxxx.yyyyy.zzzzz.
	return jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(m.Secret)
}

// ParseAccess -> verifikasi tanda tangan dan masa berlaku token
func (m *Manager) ParseAccess(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	t, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(t *jwt.Token) (any, error) { return m.Secret, nil },
		jwt.WithValidMethods([]string{"HS256"}), //hasnya terima HS256
	)

	if err != nil || !t.Valid {
		return nil, ErrHandle
	}

	return claims, nil
}

// NewRefresh menghasilkan token acak asli (raw) dan versi hash-nya untuk database.
func NewRefresh() (raw string, hash string, err error) {
	// 1. Buat wadah kosong berupa potongan byte (angka) sebanyak 32 slot.
	// Contoh value saat ini: b = [0, 0, 0, 0, 0, ...] (32 buah angka nol)
	b := make([]byte, 32)

	// 2. Isi wadah 'b' dengan angka acak menggunakan library 'crypto/rand'.
	// Kita membuang nilai pertama (_) karena kita hanya butuh mengecek error-nya.
	// Contoh value setelah diisi: b = [214, 55, 12, 99, 255, ...]
	if _, err = rand.Read(b); err != nil {
		return "", "", err // Jika gagal baca angka acak, kembalikan error
	}

	// 3. Ubah deretan angka acak tadi menjadi teks Base64 yang aman untuk URL/HTTP.
	// RawURLEncoding menghilangkan karakter seperti '+' atau '/' agar tidak merusak URL.
	// Contoh value: raw = "1aB2cD3eF4gH5iJ6kL7mN8oP9qR0sT1uV2wX3yZ4"
	raw = base64.RawURLEncoding.EncodeToString(b)

	// 4. Kembalikan teks raw, hasil hash dari teks raw, dan nil (tanpa error).
	return raw, HashRefresh(raw), nil
}

// HashRefresh mengubah token asli menjadi teks acakan (hash) SHA-256.
func HashRefresh(raw string) string {
	// 1. Ubah string raw menjadi byte, lalu proses menggunakan algoritma SHA-256.
	// Output 'sum' berupa array berukuran pasti 32 byte.
	// Contoh value: sum = [43, 211, 99, 2, 105, ...]
	sum := sha256.Sum256([]byte(raw))

	// 2. Ubah array byte tadi menjadi format Hexadecimal (kombinasi angka 0-9 dan huruf a-f).
	// sum[:] digunakan untuk memotong array menjadi slice agar diterima oleh EncodeToString.
	// Contoh value akhir: "2bd3630269..." (Teks inilah yang disimpan ke database)
	return hex.EncodeToString(sum[:])
}
