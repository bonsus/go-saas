package media

import (
	"context"
	"errors"
	"mime/multipart"
	"strings"

	"github.com/bonsus/go-saas/internal/utils/s3"
	"github.com/gofiber/fiber/v2"
)

type service struct {
	repo repository
}

func NewService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) Update(c context.Context, req updateRequest, id string) (media *Medias, errorsMap map[string][]string, err error) {
	errorsMap = map[string][]string{}
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
	req.Alt = strings.TrimSpace(req.Alt)
	req.Status = strings.TrimSpace(req.Status)

	if req.Status != "" && req.Status != "public" && req.Status != "private" {
		errorsMap["status"] = append(errorsMap["status"], "status is invalid")
	}

	if len(errorsMap) > 0 {
		return nil, errorsMap, errors.New("")
	}

	if check := s.repo.check(id); !check {
		return nil, nil, errors.New("data not found")
	}

	media, err = s.repo.Update(req, id)
	if err != nil {
		return nil, nil, err
	}

	return media, nil, nil
}

func (s *service) Upload(c *fiber.Ctx, file *multipart.FileHeader) (*Medias, error) {
	media, err := s3.UploadTemp(c, file)
	if err != nil {
		return nil, err
	}

	myMedia := Medias{
		Name: media.Name,
		Alt:  media.Name,
		Type: media.Type,
	}
	for _, me := range media.Sizes {
		err = s3.UploadTos3(me.File, me.FileS3, media.Type, "public-read")
		if err != nil && me.Id == "original" {
			s3.RemoveTemp(*media)
			return nil, errors.New("failed to upload file")
		} else if err == nil {
			myMedia.Files = append(myMedia.Files, MediaFile{
				Id:       me.Id,
				File:     me.FileS3,
				Width:    me.Width,
				Height:   me.Height,
				Filesize: me.Filesize,
			})
		}
	}

	result, err := s.repo.Create(myMedia)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Read(c context.Context, id string) (result *Medias, err error) {
	result, err = s.repo.Read(id)
	return
}

func (s *service) Index(c context.Context, param ParamIndex) (result *MediaIndex, err error) {
	result, err = s.repo.Index(param)
	return
}

func (s *service) Delete(c context.Context, Id string) (err error) {
	media, err := s.repo.Read(Id)
	if err != nil {
		return err
	}
	for _, md := range media.Files {
		err = s3.DeleteFileFromS3(md.File)
		if err != nil {
			return err
		}
	}
	err = s.repo.Delete(Id)
	if err != nil {
		return err
	}
	return
}
