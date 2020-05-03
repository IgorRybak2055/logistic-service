// Package services contains the basic logic of a Ragger.
package services

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/IgorRybak2055/logistic-service/internal/models"
	"github.com/IgorRybak2055/logistic-service/internal/repository"
	"github.com/IgorRybak2055/logistic-service/pkg/email"
)

type accountService struct {
	accountRepo repository.Account
	log         *logrus.Logger
}

// NewAccountService will create new accountService object representation of Account interface
func NewAccountService(ar repository.Account, logger *logrus.Logger) Account {
	return &accountService{
		accountRepo: ar,
		log:         logger,
	}
}

func encryptPassword(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", errors.Wrap(err, "encrypting password")
	}

	return string(b), nil
}

// CreateAccount allows to create a new user in the system,
// checks the incoming of this account according to the requirements.
func (s accountService) CreateAccount(ctx context.Context, account models.Account) (models.Account, error) {
	var err error

	if err = account.Validate(); err != nil {
		s.log.Debugf("create account: failed to validate account: %s", err)
		return models.Account{}, err
	}

	account.Password, err = encryptPassword(account.Password)
	if err != nil {
		s.log.Debugf("create account: failed to encrypting password: %s", err)
		return models.Account{}, err
	}

	var createTime = time.Now().UTC()

	account.CreatedAt, account.UpdatedAt = createTime, createTime

	s.log.Debug("create account: successfully created")

	return s.accountRepo.CreateAccount(ctx, account)
}

// Login provides access rights to protected resources and generates access tokens,
// checking the user for incoming  email and password.
func (s accountService) Login(ctx context.Context, email, password string) (models.Account, error) {
	var account, err = s.accountRepo.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Debug("failed login in account: invalid email")
			return models.Account{}, errors.New("invalid credentials")
		}

		s.log.Debugf("login: failed to login: %s", err)

		return models.Account{}, errors.Wrap(err, "logging in ragger")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			s.log.Debug("failed login in account: invalid password")

			return models.Account{}, errors.New("invalid password")
		}

		s.log.Debugf("login: failed to compare password: %s", err)

		return models.Account{}, errors.Wrap(err, "comparing passwords")
	}

	tokens, err := generateTokenPair(account.ID, account.CompanyId, account.UpdatedAt)
	if err != nil {
		s.log.Debugf("login: failed to generate token: %s", err)

		return models.Account{}, errors.Wrap(err, "generating tokens")
	}

	account.Token = tokens
	account.Password = ""

	s.log.Debugf("login: user %v successfully", account.ID)

	return account, nil
}

// validateRefreshTokenClaims checks that the refresh token is valid.
func (s accountService) validateRefreshTokenClaims(ctx context.Context,
	claims map[string]interface{}) (models.Account, bool) {
	var accountID = int64(claims["id"].(float64))

	var account, err = s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		s.log.Debugf("validateRefreshTokenClaims: failed to get user by id: %s", err)

		return models.Account{}, false
	}

	ct, err := time.Parse(time.RFC3339, claims["updated_at"].(string))
	if err != nil {
		s.log.Debugf("validateRefreshTokenClaims: failed to parsing time: %s", err)

		return models.Account{}, false
	}

	return account, account.UpdatedAt.Equal(ct)
}

// GenerateToken allows using a refresh token to generate a new pair of tokens.
func (s accountService) GenerateToken(ctx context.Context, refreshToken string) (map[string]string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Debugf("generate token: failed to check signing method")

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKEN_PASSWORD")), nil
	})
	if err != nil {
		s.log.Debugf("generate token: failed to parsing token: %s", err)

		return nil, err
	}

	var (
		claims jwt.MapClaims
		ok     bool
	)

	if claims, ok = token.Claims.(jwt.MapClaims); !(ok && token.Valid) {
		s.log.Debug("generate token: failed to getting claims")

		return nil, err
	}

	var account models.Account

	if account, ok = s.validateRefreshTokenClaims(ctx, claims); !ok {
		s.log.Debug("generate token: failed to validate refresh token claims: invalid refresh token")

		return nil, errors.New("invalid refresh token")
	}

	newTokenPair, err := generateTokenPair(account.ID, account.CompanyId, account.UpdatedAt)
	if err != nil {
		s.log.Debugf("generate token: failed to generate news tokens: %s", err)

		return nil, err
	}

	return newTokenPair, nil
}

func generateTokenPair(accountID, companyID int64, accountUpdateTime time.Time) (map[string]string, error) {
	var (
		token  = jwt.New(jwt.SigningMethodHS256)
		claims = token.Claims.(jwt.MapClaims)
	)

	claims["id"] = accountID
	claims["company_id"] = companyID
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()

	var t, err = token.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	if err != nil {
		return nil, err
	}

	var (
		refreshToken = jwt.New(jwt.SigningMethodHS256)
		rtClaims     = refreshToken.Claims.(jwt.MapClaims)
	)

	rtClaims["id"] = accountID
	rtClaims["updated_at"] = accountUpdateTime
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	var rt string

	rt, err = refreshToken.SignedString([]byte(os.Getenv("TOKEN_PASSWORD")))
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}

// CheckAccountByEmail checks if users with received email exist.
func (s accountService) RestorePassword(ctx context.Context, ch chan email.MessageData, emailAddress string) error {
	var account, err = s.accountRepo.GetByEmail(ctx, emailAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			s.log.Debug("failed login in account: invalid email")
			return errors.New("invalid credentials")
		}

		s.log.Debugf("login: failed to login: %s", err)

		return errors.Wrap(err, "logging in ragger")
	}

	s.log.Info("checked account: ", account)

	tokens, err := generateTokenPair(account.ID, account.CompanyId, time.Now())
	if err != nil {
		s.log.Debugf("CheckAccountByEmail: failed to generate news tokens: %s", err)

		return err
	}

	ch <- email.MessageData{
		RecvEmail: emailAddress,
		UserID:    string(account.ID),
		UserToken: tokens["access_token"],
	}

	return nil
}

func (s accountService) SetNewPassword(ctx context.Context, newPassword string) error {
	var err error

	newPassword, err = encryptPassword(newPassword)
	if err != nil {
		s.log.Debugf("updating password: failed to encrypting password: %s", err)

		return err
	}

	err = s.accountRepo.SetNewPassword(ctx, newPassword)
	if err != nil {
		s.log.Debugf("updating password: failed to set new password: %s", err)

		return err
	}

	return nil
}
