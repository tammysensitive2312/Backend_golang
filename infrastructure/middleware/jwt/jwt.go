package jwt

import (
	"Backend_golang_project/infrastructure/config"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"time"
)

/**  khái niệm
jwt : json web token
cấu trúc của jwt có 3 phần :
+ header : có 2 phần nhỏ
  + typ: loại token (thường là jwt)
  + alg (algorithm): thuật toán mã hóa sử dụng để ký token,ví dụ như HMAC SHA256 (HS256)
+ payload : chứa dữ liệu (claims) mà response muốn truyền tải:
  + Registered Claims: Các claims chuẩn đã được định nghĩa sẵn trong JWT, bao gồm:
    iss (issuer): Người phát hành token.
    sub (subject): Chủ đề của token (thường là ID của người dùng).
    aud (audience): Đối tượng mà token hướng tới.
    exp (expiration): Thời gian hết hạn của token (UNIX timestamp).
    nbf (not before): Thời điểm token có hiệu lực (UNIX timestamp).
    iat (issued at): Thời điểm token được tạo ra (UNIX timestamp).
    jti (JWT ID): ID duy nhất của token.

  + Public Claims: Các claims do người dùng tự định nghĩa, ví dụ như thông tin về người dùng.

  + Private Claims: Các claims riêng giữa các bên trao đổi.

+ signature : là phần bảo mật của JWT, được tạo ra bằng cách kết hợp Header và Payload, sau đó ký bằng
một thuật toán mã hóa đã được chỉ định trong Header và một secret key
*/

/** ứng dụng
+ stateless : ko cần lưu trữ trên server (khác với session)
+ confidentiality : dữ liệu trong jwt đều được ký và mã hóa để đảm bảo C I (A)
*/

/**
một mục đích sử dụng nữa của JWT là cơ chế Refresh token (đã học qua trong java)
+ nếu access token được set thời gian sống ngắn từ vài phút đến vài giờ (để 10p)
thì refresh token sẽ được set thời gian sống lâu hơn thường là từ ngày đến tuần (để 1 ngày)
response sẽ dùng token này để gửi tới một endpoint xin cấp phát access token mới
*/

/** workflow
khi mà client login server sẽ xác thực thông tin dựa trên nhưng gì client gửi
-> tạo ra token pair (access & refresh) & và gửi lại cho client (js)
-> client lưu trữ token : (nhiệm vụ phía fe)
         + at : bộ nhớ tạm || localStorage || sessionStorage
         + rt : httpOnly cookie
-> client sử dụng at để truy cập vào các endpoint (tài nguyên hệ thống) (inject trong bearer token)
    + access token oke và chưa hết hạn cho truy cập tiếp
    + else -> lỗi 401 cho client -> client gửi refresh token tới endpoint cụ thể xin cấp lại at
    + khi rt hết hạn thì yêu cầu đăng nhập lại
*/

type JWTClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

// GenerateJwtToken tạo cặp token mới
func GenerateJwtToken(config *config.Config, ID int) (string, string, error) {
	atExpTime := time.Now().Add(time.Duration(config.JwtConfig.AccessTokenExp) * time.Minute)
	rtExpTime := time.Now().Add(time.Duration(config.JwtConfig.RefreshTokenExp) * time.Minute)

	accessClaims := &JWTClaims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(atExpTime),
		},
	}

	refreshClaims := &RefreshTokenClaims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(rtExpTime),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(config.JwtConfig.SecretKey))
	if err != nil {
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(config.JwtConfig.SecretKey))
	if err != nil {
		return "", "", err
	}

	encryptedRefreshToken, err := encryptRefreshToken(refreshTokenString, config.JwtConfig.SecretKey)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, encryptedRefreshToken, nil
}

// ClaimToken xác thực token
func ClaimToken(tokenString string, config *config.Config) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtConfig.SecretKey), nil
		},
	)

	if err != nil {
		return nil, handleTokenError(err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func handleTokenError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		return errors.New("malformed token")
	case errors.Is(err, jwt.ErrTokenExpired):
		return errors.New("token is expired")
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return errors.New("token not active yet")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return errors.New("invalid token signature")
	default:
		return fmt.Errorf("couldn't handle this token: %w", err)
	}
}

// ClaimRefreshToken em nghĩ có thể tái sử dụng đoạn mã bên trên theo một cách nào đó,
// vì ở hàm này chỉ đơn giản là copy bên trên và thêm phần giải mã
func ClaimRefreshToken(encryptedToken string, config *config.Config) (*RefreshTokenClaims, error) {
	refreshToken, err := decryptRefreshToken(encryptedToken, config.JwtConfig.SecretKey)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		refreshToken,
		&RefreshTokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtConfig.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// giải mã Refresh token từ thuật toán AES-GCM
func decryptRefreshToken(encryptedToken string, secretKey string) (string, error) {
	// cách tạo salt em đọc được ở trên diễn đàn nào đó, nó là chuỗi ngẫy nhiên được băm
	// từ secret key bằng SHA1
	sha1Hasher := sha1.New()
	io.WriteString(sha1Hasher, secretKey)
	salt := string(sha1Hasher.Sum(nil))[0:16]

	// mã hóa salt
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// giải mã từ base64 về dạng byte
	data, err := base64.URLEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", err
	}

	// tách nonce và bản mã
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	//giải mã
	decrypted, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

// ngược lại với hàm giải mã
func encryptRefreshToken(refreshToken string, secretKey string) (string, error) {
	sha1Hasher := sha1.New()
	io.WriteString(sha1Hasher, secretKey)
	salt := string(sha1Hasher.Sum(nil))[0:16]

	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(refreshToken), nil)
	return base64.URLEncoding.EncodeToString(encrypted), nil
}
