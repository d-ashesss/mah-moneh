package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/auth"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/auth"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type AuthServiceTestSuite struct {
	suite.Suite

	privKey *rsa.PrivateKey
	pubKey  []byte

	users *mocks.UsersService
	srv   *auth.Service
}

func (ts *AuthServiceTestSuite) SetupSuite() {
	privKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Failed to generate private key: %s", err)
	}
	ts.privKey = privKey
	der, err := x509.MarshalPKIXPublicKey(privKey.Public())
	if err != nil {
		log.Fatalf("Failed to marshal public key: %s", err)
	}
	ts.pubKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: der})
}

func (ts *AuthServiceTestSuite) SetupTest() {
	ts.users = mocks.NewUsersService(ts.T())
	cfg := &auth.Config{PublicKey: string(ts.pubKey)}
	ts.srv = auth.NewService(cfg, ts.users)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_InvalidToken() {
	token := `token of gratitude`

	user, err := ts.srv.AuthenticateUser(context.Background(), token)
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_UnsignedToken() {
	t := jwt.NewWithClaims(jwt.SigningMethodNone, &jwt.RegisteredClaims{Subject: uuid.Must(uuid.NewV4()).String()})
	token, err := t.SignedString(jwt.UnsafeAllowNoneSignatureType)
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(context.Background(), token)
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_ExpiredToken() {
	ctx := context.Background()

	exp := jwt.NewNumericDate(time.Now().Add(-10 * time.Second))
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.RegisteredClaims{Subject: uuid.Must(uuid.NewV4()).String(), ExpiresAt: exp})
	token, err := t.SignedString(ts.privKey)
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(ctx, token)
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_Valid_UserNotFound() {
	ctx := context.Background()
	userUUID := uuid.Must(uuid.NewV4())
	ts.users.On("GetUser", ctx, userUUID).Return(nil, errors.New("test error"))

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.RegisteredClaims{Subject: userUUID.String()})
	token, err := t.SignedString(ts.privKey)
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(ctx, token)
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_Valid_UserFound() {
	ctx := context.Background()
	userUUID := uuid.Must(uuid.NewV4())
	ts.users.On("GetUser", ctx, userUUID).Return(&users.User{UUID: userUUID}, nil)

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.RegisteredClaims{Subject: userUUID.String()})
	token, err := t.SignedString(ts.privKey)
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(ctx, token)
	ts.Require().NoError(err)
	ts.Equal(userUUID, user.UUID)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
