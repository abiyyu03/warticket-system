package token

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Role yang dibawa di klaim JWT.
const (
	RoleAuthor   = "author"
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

const (
	typeAccess  = "access"
	typeRefresh = "refresh"
)

var ErrInvalidToken = errors.New("token tidak valid")

// Claims: identitas + role + jenis token (access/refresh).
type Claims struct {
	UserID int64  `json:"uid"`
	Role   string `json:"role"`
	Type   string `json:"typ"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// New membangun Manager dari env (dengan default aman untuk dev/test).
func New() *Manager {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  time.Duration(intEnv("JWT_ACCESS_TTL_MIN", 15)) * time.Minute,
		refreshTTL: time.Duration(intEnv("JWT_REFRESH_TTL_HOURS", 168)) * time.Hour,
	}
}

func intEnv(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func (m *Manager) generate(userID int64, role, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
}

// GeneratePair menerbitkan access (short-lived) + refresh (long-lived) stateless.
func (m *Manager) GeneratePair(userID int64, role string) (access, refresh string, expiresIn int64, err error) {
	access, err = m.generate(userID, role, typeAccess, m.accessTTL)
	if err != nil {
		return
	}
	refresh, err = m.generate(userID, role, typeRefresh, m.refreshTTL)
	if err != nil {
		return
	}
	expiresIn = int64(m.accessTTL.Seconds())
	return
}

func (m *Manager) parse(tokenStr, wantType string) (Claims, error) {
	var claims Claims
	tok, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil || !tok.Valid || claims.Type != wantType {
		return Claims{}, ErrInvalidToken
	}
	return claims, nil
}

// ParseAccess memvalidasi access token dan mengembalikan klaimnya.
func (m *Manager) ParseAccess(tokenStr string) (Claims, error) {
	return m.parse(tokenStr, typeAccess)
}

// ParseRefresh memvalidasi refresh token dan mengembalikan klaimnya.
func (m *Manager) ParseRefresh(tokenStr string) (Claims, error) {
	return m.parse(tokenStr, typeRefresh)
}
