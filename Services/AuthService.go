package Services

import (
	"ct-backend/Model"
	"ct-backend/Model/Dto"
	"ct-backend/Repository"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type (
	IAuthService interface {
		Login(email string, password string) (user *Model.User, token string, err error)
		Register(request *Dto.RegisterRequest) (err error)
		RequestOtp(email string) (err error)
		VerifyOtp(email string, otp string) (err error)
	}

	AuthService struct {
		repo Repository.IAuthRepository
	}
)

func AuthServiceProvider(repo Repository.IAuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (h *AuthService) Login(email string, password string) (user *Model.User, token string, err error) {

	// get user information
	if user, err = h.repo.GetUserInformation(email); err != nil {
		return nil, "", err
	} else if user == nil {
		return nil, "", errors.New("user not found")
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid password")
	}

	// generate token
	expirationTime := time.Now().Add(3 * 30 * 24 * time.Hour)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"sub": user.ID,
		"exp": expirationTime.Unix(),
	})

	if token, err = jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return nil, "", err
	}

	return user, token, err
}

func (h *AuthService) Register(request *Dto.RegisterRequest) (err error) {
	var (
		hashedPassword []byte
	)

	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	request.Password = string(hashedPassword)

	if err = h.repo.InsertUserInformation(request); err != nil {
		return err
	}

	return err
}

func (h *AuthService) RequestOtp(email string) (err error) {
	randomNumber := rand.Intn(9000) + 1000
	if err = h.repo.SetOtpCode(email, strconv.Itoa(randomNumber)); err != nil {
		return err
	}

	// send otp to email

	return err
}

func (h *AuthService) VerifyOtp(email string, otp string) (err error) {
	var (
		user *Model.User
	)

	if user, err = h.repo.GetUserInformation(email); err != nil {
		return err
	} else if user == nil {
		return errors.New("user not found")
	}

	if user.OtpCode != otp {
		return errors.New("invalid otp")
	} else {

		if user.UpdatedAt.Before(time.Now().Add(-2 * time.Minute)) {
			return errors.New("otp expired")
		}

		if err = h.repo.SetVerificationStatus(email); err != nil {
			return err
		}
	}

	return
}
