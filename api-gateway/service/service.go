package service

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type authService struct {
	secretKey     string
	expiryMinutes int
}

func NewAuthService(secretKey string, expiryMinutes int) AuthService {
	return &authService{
		secretKey:     secretKey,
		expiryMinutes: expiryMinutes,
	}
}

func (s *authService) GenerateToken(userID string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.expiryMinutes) * time.Minute)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

type ProxyResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

type ProxyService interface {
	ProxyRequest(method, path string, body io.Reader, headers map[string]string) (*ProxyResponse, error)
}

type proxyService struct {
	baseURL    string
	httpClient *http.Client
}

func NewProxyService(baseURL string) ProxyService {
	return &proxyService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *proxyService) ProxyRequest(method, path string, body io.Reader, headers map[string]string) (*ProxyResponse, error) {
	relativePath := path
	if strings.HasPrefix(path, "/api/v1") {
		relativePath = strings.TrimPrefix(path, "/api/v1")
	}

	targetURL := s.baseURL + relativePath

	req, err := http.NewRequest(method, targetURL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respHeaders := make(map[string]string)
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	return &ProxyResponse{
		StatusCode: resp.StatusCode,
		Headers:    respHeaders,
		Body:       respBody,
	}, nil
}
