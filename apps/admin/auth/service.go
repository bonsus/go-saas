package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/bonsus/go-saas/internal/config"
	myredis "github.com/bonsus/go-saas/internal/redis"
	token "github.com/bonsus/go-saas/internal/utils/jwt"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Login(ctx context.Context, req Request) (res Response, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	err = validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch field {
			case "Username":
				errorsMap["usename"] = append(errorsMap["usename"], "usename is required")
			case "Password":
				errorsMap["password"] = append(errorsMap["password"], "password is required")
			}
		}
	}

	if len(errorsMap) > 0 {
		return res, errorsMap, errors.New("")
	}

	user, err := s.repo.FindUserByUsername(ctx, req)
	if err != nil {
		return res, errorsMap, errors.New("username or password does not match")
	} else if err := user.ValidatePassword(req.Password); err != nil {
		return res, errorsMap, errors.New("username or password does not match")
	} else if user.Status != "active" {
		return res, errorsMap, errors.New("account is inactive")
	}

	cfg := config.GetConfig()
	claims := token.Claims{
		Id:         user.Id,
		Type:       "admin",
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
		Admin: *user,
	}
	return res, nil, nil
}

func (s *service) Register(req RegisterRequest) (*Admin, map[string][]string, error) {
	errorsMap := map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch field {
			case "Name":
				errorsMap["name"] = append(errorsMap["name"], "name is required")
			case "Email":
				errorsMap["email"] = append(errorsMap["email"], "email is required")
			case "Username":
				errorsMap["username"] = append(errorsMap["username"], "username is required")
			case "Password":
				errorsMap["password"] = append(errorsMap["password"], "password is required")
			case "RoleId":
				errorsMap["role_id"] = append(errorsMap["role_id"], "role is required")
			}
		}
	}

	if err := s.repo.FindByEmail(req.Email); err {
		errorsMap["email"] = append(errorsMap["email"], "email is already exists")
	}
	if err := s.repo.FindByUsername(req.Username); err {
		errorsMap["username"] = append(errorsMap["username"], "username is already exists")
	}
	if err := s.repo.FindRoleByID(req.RoleId); !err {
		errorsMap["role_id"] = append(errorsMap["role_id"], "role is invalid")
	}
	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("validation failed")
	}

	savedUser, err := s.repo.Register(req)
	if err != nil {
		return nil, nil, err
	}

	return savedUser, nil, nil
}

func (s *service) Me(adminId string) (*Admin, error) {
	var admin *Admin
	cacheId := "user:" + adminId
	err := myredis.GetData(cacheId, &admin)
	if err != nil {
		admin, err = s.repo.FindUserById(adminId)
		if err != nil {
			return nil, err
		}
		myredis.SetData(cacheId, admin, 30*time.Minute)
	}

	return admin, nil
}

func (s *service) Update(req RegisterRequest, Id string) (*Admin, map[string][]string, error) {
	errorsMap := map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Email == "" {
		errorsMap["email"] = append(errorsMap["email"], "email is required")
	}
	if req.Username == "" {
		errorsMap["username"] = append(errorsMap["username"], "username is required")
	}

	if err := s.repo.FindByEmail(req.Email, Id); err {
		errorsMap["email"] = append(errorsMap["email"], "email is already exists")
	}
	if err := s.repo.FindByUsername(req.Username, Id); err {
		errorsMap["username"] = append(errorsMap["username"], "username is already exists")
	}
	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("validation failed")
	}

	savedUser, err := s.repo.Update(req, Id)
	if err != nil {
		return nil, nil, err
	}

	myredis.RemoveData("user:" + Id)

	return savedUser, nil, nil
}

func (s *service) UpdatePassword(req RegisterRequest, Id string) (*Admin, map[string][]string, error) {
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

	user, err := s.repo.FindUserPasswordById(Id)
	if err != nil {
		errorsMap["password"] = append(errorsMap["password"], "current password does not match")
	} else if err := user.ValidatePassword(req.Password); err != nil {
		errorsMap["password"] = append(errorsMap["password"], "current password does not matchs")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("validation failed")
	}

	savedUser, err := s.repo.UpdatePassword(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return savedUser, nil, nil
}
