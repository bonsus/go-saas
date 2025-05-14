package admins

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	myredis "github.com/bonsus/go-saas/internal/redis"
)

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(c context.Context, req Request) (admin *Admin, errorsMap map[string][]string, err error) {
	req, errorsMap = s.validateRequest(req, "")

	if req.Password == "" {
		errorsMap["password"] = append(errorsMap["password"], "password is required")
	}
	if len(req.Password) < 8 {
		errorsMap["password"] = append(errorsMap["password"], "password must have at least 8 characters")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Create(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) Update(c context.Context, req Request, id string) (admin *Admin, errorsMap map[string][]string, err error) {
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	req, errorsMap = s.validateRequest(req, id)
	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.Update(req, id)
	if err != nil {
		return nil, nil, err
	}

	myredis.RemoveData("user:" + id)

	return result, nil, nil
}

func (s *service) UpdateStatus(c context.Context, req Request, id string) (admin *Admin, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.UpdateStatus(req, id)
	if err != nil {
		return nil, nil, err
	}
	myredis.RemoveData("user:" + id)

	return result, nil, nil
}

func (s *service) UpdatePassword(c context.Context, req Request, id string) (admin *Admin, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	_, err = s.repo.Read(id)
	if err != nil {
		return nil, nil, errors.New("data not found")
	}
	if req.Password == "" {
		errorsMap["password"] = append(errorsMap["password"], "new password is required")
	}
	if len(req.Password) < 8 {
		errorsMap["password"] = append(errorsMap["password"], "new password must have at least 8 characters")
	}
	if req.PasswordConfirmation != req.Password {
		errorsMap["password_confirmation"] = append(errorsMap["password_confirmation"], "password confirmation does not match")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.UpdatePassword(req, id)
	if err != nil {
		return nil, nil, err
	}
	myredis.RemoveData("user:" + id)

	return result, nil, nil
}

func (s *service) Index(c context.Context, param ParamIndex) (result *AdminIndex, err error) {
	result, err = s.repo.Index(param)
	return
}

func (s *service) Read(c context.Context, id string) (result *Admin, err error) {
	result, err = s.repo.Read(id)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (s *service) Delete(c context.Context, id string) (err error) {
	err = s.repo.Delete(id)
	return
}

func (s *service) validateRequest(req Request, Id string) (request Request, errorsMap map[string][]string) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if req.Email == "" {
		errorsMap["email"] = append(errorsMap["email"], "email is required")
	}
	if req.Name == "" {
		errorsMap["username"] = append(errorsMap["username"], "username is required")
	}
	if req.Status != "active" && req.Status != "inactive" && req.Status != "deleted" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}
	if req.RoleId == "" {
		errorsMap["role_id"] = append(errorsMap["role_id"], "role is required")
	}
	if err := s.repo.FindRoleByID(req.RoleId); !err {
		errorsMap["role_id"] = append(errorsMap["role_id"], "role is invalid")
	}

	if err := s.repo.FindByEmail(req.Email, Id); err {
		errorsMap["email"] = append(errorsMap["email"], "email is already exists")
	}
	if err := s.repo.FindByUsername(req.Username, Id); err {
		errorsMap["username"] = append(errorsMap["username"], "username is already exists")
	}
	return req, errorsMap
}
func IsValidEmail(email string) error {
	if email != "" {
		regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		re := regexp.MustCompile(regex)
		if !re.MatchString(email) {
			return errors.New("invalid email format")
		}
	}
	return nil
}
func ConvertPhoneNumber(phone string) (string, error) {
	if phone != "" {
		// Hapus semua spasi dan karakter selain angka
		re := regexp.MustCompile(`[^0-9]`)
		phone = re.ReplaceAllString(phone, "")

		// Jika nomor diawali dengan "620", ubah menjadi "62"
		if strings.HasPrefix(phone, "620") {
			phone = "62" + phone[3:]
		}

		// Jika nomor diawali dengan "0", ubah menjadi "62"
		if strings.HasPrefix(phone, "0") {
			phone = "62" + phone[1:]
		}

		// Pastikan nomor diawali dengan "62"
		if !strings.HasPrefix(phone, "62") {
			return "", errors.New("phone number must start with '62'")
		}

		// Validasi panjang nomor (minimal 10 digit setelah konversi)
		if len(phone) < 10 || len(phone) > 15 {
			return "", errors.New("invalid phone number length")
		}
	}
	return phone, nil
}
func ValidateBirthdate(birthdate string) (string, error) {
	if birthdate != "" {
		const layout = "2006-01-02" // Format YYYY-MM-DD
		parsedDate, err := time.Parse(layout, birthdate)
		if err != nil {
			return "", errors.New("invalid date format, expected YYYY-MM-DD")
		}

		// Cek rentang usia
		today := time.Now()
		age := today.Year() - parsedDate.Year()

		// Periksa apakah ulang tahun sudah terjadi tahun ini atau belum
		if today.YearDay() < parsedDate.YearDay() {
			age--
		}

		if age < 1 || age > 120 {
			return "", errors.New("invalid birthdate: age must be between 1 and 120 years")
		}

		// Return tanggal dalam format YYYY-MM-DD
		return parsedDate.Format("2006-01-02"), nil
	}
	return birthdate, nil
}

func (s *service) RoleIndex(c context.Context) (result []AdminRole, err error) {
	result, err = s.repo.RoleIndex()
	return
}

func (s *service) RoleCreate(c context.Context, req RoleRequest) (role *AdminRole, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)

	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.RoleCreate(req)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}
func (s *service) RoleUpdate(c context.Context, req RoleRequest, Id string) (role *AdminRole, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)

	if err := s.repo.FindRoleByID(Id); !err {
		return nil, nil, errors.New("data not found")
	}
	if req.Name == "" {
		errorsMap["name"] = append(errorsMap["name"], "name is required")
	}
	if err := s.repo.FindRoleByName(req.Name, Id); err {
		errorsMap["name"] = append(errorsMap["name"], "name is already exists")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}
	result, err := s.repo.RoleUpdate(req, Id)
	if err != nil {
		return nil, nil, err
	}

	return result, nil, nil
}

func (s *service) RoleDelete(c context.Context, id string) (err error) {
	err = s.repo.RoleDelete(id)
	return
}
