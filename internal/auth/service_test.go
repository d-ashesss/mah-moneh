package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"github.com/d-ashesss/mah-moneh/internal/auth"
	mocks "github.com/d-ashesss/mah-moneh/internal/mocks/auth"
	"github.com/d-ashesss/mah-moneh/internal/users"
	"github.com/gofrs/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

const TestKeyID = "test_key"

type AuthServiceTestSuite struct {
	suite.Suite

	privKey jwk.Key
	pubKey  jwk.Key

	users *mocks.UsersService
	srv   *auth.Service
}

func (ts *AuthServiceTestSuite) SetupSuite() {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		log.Fatalf("Failed to generate private key: %s", err)
	}

	privKey, err := jwk.FromRaw(rsaKey)
	if err != nil {
		log.Fatalf("Failed to create public key: %s", err)
	}
	if err := privKey.Set(jwk.KeyIDKey, TestKeyID); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	if err := privKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	ts.privKey = privKey

	pubKey, err := jwk.FromRaw(rsaKey.Public())
	if err != nil {
		log.Fatalf("Failed to create public key: %s", err)
	}
	if err = pubKey.Set(jwk.KeyIDKey, TestKeyID); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	if err = pubKey.Set(jwk.AlgorithmKey, jwa.RS256); err != nil {
		log.Fatalf("Failed to set key ID: %s", err)
	}
	ts.pubKey = pubKey
}

func (ts *AuthServiceTestSuite) SetupTest() {
	ts.users = mocks.NewUsersService(ts.T())
	cfg := &auth.Config{}
	ts.srv = auth.NewService(cfg, ts.users)
	if err := ts.srv.AddKey(ts.pubKey); err != nil {
		log.Fatalf("Failed to add test public key: %s", err)
	}
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_InvalidToken() {
	token := `token of gratitude`

	user, err := ts.srv.AuthenticateUser(context.Background(), token)
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_UnsignedToken() {
	ctx := context.Background()
	// JWT without signature part
	token := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiI3YzMxYTVlMS1jMjY1LTQ3NTUtOTEyZC0xMjBlNDA1YjQ2ZDMifQ."

	user, err := ts.srv.AuthenticateUser(ctx, token)
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_ExpiredToken() {
	ctx := context.Background()
	userUUID := uuid.Must(uuid.NewV4())

	tt := jwt.New()
	err := tt.Set(jwt.SubjectKey, userUUID.String())
	ts.Require().NoError(err)
	err = tt.Set(jwt.ExpirationKey, time.Now().Add(-10*time.Second))
	ts.Require().NoError(err)
	token, err := jwt.Sign(tt, jwt.WithKey(ts.privKey.Algorithm(), ts.privKey))
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(ctx, string(token))
	ts.Nil(user)
	ts.ErrorIs(err, jwt.ErrTokenExpired())
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_Valid_UserNotFound() {
	ctx := context.Background()
	userUUID := uuid.Must(uuid.NewV4())
	ts.users.On("GetUser", ctx, userUUID).Return(nil, errors.New("test error"))

	tt := jwt.New()
	err := tt.Set(jwt.SubjectKey, userUUID.String())
	ts.Require().NoError(err)
	token, err := jwt.Sign(tt, jwt.WithKey(ts.privKey.Algorithm(), ts.privKey))
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(ctx, string(token))
	ts.Nil(user)
	ts.Error(err)
}

func (ts *AuthServiceTestSuite) TestAuthenticateUser_Valid_UserFound() {
	ctx := context.Background()
	userUUID := uuid.Must(uuid.NewV4())
	ts.users.On("GetUser", ctx, userUUID).Return(&users.User{UUID: userUUID}, nil)

	tt := jwt.New()
	err := tt.Set(jwt.SubjectKey, userUUID.String())
	ts.Require().NoError(err)
	token, err := jwt.Sign(tt, jwt.WithKey(ts.privKey.Algorithm(), ts.privKey))
	ts.Require().NoError(err)

	user, err := ts.srv.AuthenticateUser(ctx, string(token))
	ts.Require().NoError(err)
	ts.Equal(userUUID, user.UUID)
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}
