package auth

import (
	"context"
	"errors"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	myredis "github.com/bonsus/go-saas/internal/redis"
	token "github.com/bonsus/go-saas/internal/utils/jwt"
)

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Register(req RegisterRequest) (*User, map[string][]string, error) {
	errorsMap := map[string][]string{}
	var err error
	req.Company = strings.TrimSpace(req.Company)
	req.Name = strings.TrimSpace(req.Name)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Company == "" {
		errorsMap["company"] = append(errorsMap["company"], "company is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Phone == "" {
		errorsMap["phone"] = append(errorsMap["phone"], "phone is required")
	} else {
		req.Phone, err = validatePhone(req.Phone)
		if err != nil {
			errorsMap["phone"] = append(errorsMap["phone"], "phone is invalid")
		}
		if check := s.repo.CheckPhone(req.Phone); check {
			errorsMap["phone"] = append(errorsMap["phone"], "phone is already used")
		}
	}
	if req.Email == "" {
		errorsMap["email"] = append(errorsMap["email"], "email is required")
	} else if check := s.repo.CheckEmail(req.Email); check {
		errorsMap["email"] = append(errorsMap["email"], "email is already registered")
	} else {
		if check := validateEmail(req.Email); !check {
			errorsMap["email"] = append(errorsMap["email"], "email is invalid")
		}
	}
	if req.Username != "" {
		if check := s.repo.CheckUsername(req.Username); check {
			errorsMap["username"] = append(errorsMap["username"], "username is already used")
		}
	}
	if req.Password == "" {
		errorsMap["password"] = append(errorsMap["password"], "password is required")
	}
	if len(req.Password) < 8 {
		errorsMap["password"] = append(errorsMap["password"], "password must have at least 8 characters")
	}
	if req.PasswordConfirmation != req.Password {
		errorsMap["password_confirmation"] = append(errorsMap["password_confirmation"], "password confirmation does not match")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("validation failed")
	}

	user, err := s.repo.Register(req)
	if err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}
func (s *service) Login(ctx context.Context, req Request) (res Response, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	if len(errorsMap) > 0 {
		return res, errorsMap, errors.New("")
	}
	user, err := s.repo.Login(req)

	if err != nil {
		return res, errorsMap, errors.New("login failed. Invalid credentials")
	} else if err := user.ValidatePassword(req.Password); err != nil {
		return res, errorsMap, errors.New("login failed. Invalid credentials")
	} else if user.Status == "unverify" {
		return res, errorsMap, errors.New("account is not active yet, please verify via email")
	} else if user.Status != "active" {
		return res, errorsMap, errors.New("account is not active")
	}
	cfg := config.GetConfig()
	claims := token.Claims{
		Id:         user.Id,
		ExpireTime: time.Duration(time.Now().Unix() + int64(cfg.JWT.ExpireTime)),
	}
	token, err := user.GenerateToken(claims, cfg.JWT.Key)
	if err != nil {
		slog.ErrorContext(ctx, "[login] error when create token", slog.Any("error", err.Error()))
		return res, nil, err
	}
	user.Password = ""
	res = Response{
		Token: token,
		User:  *user,
	}
	return res, nil, nil
}

func (s *service) Me(userId string) (*User, error) {
	var user *User
	cacheId := "user:" + userId
	err := myredis.GetData(cacheId, &user)
	if err != nil {
		user, err = s.repo.ReadDetail(userId)
		if err != nil {
			return nil, err
		}
		myredis.SetData(cacheId, user, 30*time.Minute)
	}

	return user, nil
}

func (s *service) Update(req Request, Id string) (*User, map[string][]string, error) {
	_, err := s.repo.Read(Id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	errorsMap := map[string][]string{}
	req.Company = strings.TrimSpace(req.Company)
	req.Name = strings.TrimSpace(req.Name)

	if req.Company == "" {
		errorsMap["company"] = append(errorsMap["company"], "company is required")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("validation failed")
	}

	user, err := s.repo.Update(req, Id)
	if err != nil {
		return nil, nil, err
	}

	myredis.RemoveData("user:" + Id)

	return user, nil, nil
}

func (s *service) UpdatePassword(req Request, Id string) (*User, map[string][]string, error) {
	_, err := s.repo.Read(Id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	errorsMap := map[string][]string{}

	if req.Password == "" {
		errorsMap["password"] = append(errorsMap["password"], "current password is required")
	}
	if req.NewPassword == "" {
		errorsMap["new_password"] = append(errorsMap["new_password"], "new password is required")
	}
	if len(req.NewPassword) < 8 {
		errorsMap["new_password"] = append(errorsMap["new_password"], "password must have at least 8 characters")
	}
	if req.NewPasswordConfirmation == "" {
		errorsMap["new_password_confirmation"] = append(errorsMap["new_password_confirmation"], "password confirmation is required")
	}
	if req.NewPasswordConfirmation != "" && req.NewPassword != req.NewPasswordConfirmation {
		errorsMap["new_password_confirmation"] = append(errorsMap["new_password_confirmation"], "password confirmation does not match")
	}

	user, err := s.repo.LoginById(Id)
	if err != nil {
		errorsMap["password"] = append(errorsMap["password"], "current password does not match")
	} else if err := user.ValidatePassword(req.Password); err != nil {
		errorsMap["password"] = append(errorsMap["password"], "current password does not matchs")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("validation failed")
	}

	user, err = s.repo.UpdatePassword(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return user, nil, nil
}

func validatePhone(input string) (string, error) {
	var res string
	switch {
	case strings.HasPrefix(input, "08"):
		res = "62" + input[1:]
	case strings.HasPrefix(input, "8"):
		res = "62" + input
	case strings.HasPrefix(input, "62"):
		res = input
	default:
		return input, errors.New("phone is invalid")
	}
	if len(res) < 10 {
		return res, errors.New("phone is invalid")
	}
	return res, nil
}
func validateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
