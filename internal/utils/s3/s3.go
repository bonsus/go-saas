package s3

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/bonsus/go-saas/internal/config"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/image/webp"
)

type Media struct {
	Name        string
	Description string
	Alt         string
	Type        string
	Status      string
	Sizes       []Size
}
type Size struct {
	Id       string
	File     string
	FileS3   string
	Width    int
	Height   int
	Filesize string
}

func RemoveTemp(media Media) error {
	for _, sz := range media.Sizes {
		os.Remove(sz.File)
		dirPath := filepath.Dir(sz.File)
		os.Remove(dirPath)
		dirPath = filepath.Dir(dirPath)
		os.Remove(dirPath)
		dirPath = filepath.Dir(dirPath)
		os.Remove(dirPath)
	}
	return nil
}

func DeleteFileFromS3(fileS3 string) error {
	cfg := config.GetConfig()
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.S3.Region),
		Credentials:      credentials.NewStaticCredentials(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Endpoint:         aws.String(cfg.S3.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(cfg.S3.BucketName),
		Key:    aws.String(fileS3),
	}

	_, err = s3Client.DeleteObject(input)
	if err != nil {
		return err
	}
	return nil
}
func UploadTos3(fileTemp string, files3 string, fileType string, status string) (err error) {
	cfg := config.GetConfig()
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.S3.Region),
		Credentials:      credentials.NewStaticCredentials(cfg.S3.AccessKey, cfg.S3.SecretKey, ""),
		Endpoint:         aws.String(cfg.S3.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		return errors.New("failed create session s3")
	}

	s3Client := s3.New(sess)

	fileData, err := os.ReadFile(fileTemp)
	if err != nil {
		return errors.New("failed to read file")
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(cfg.S3.BucketName),
		Key:         aws.String(files3),
		Body:        bytes.NewReader(fileData),
		ACL:         aws.String(status),
		ContentType: aws.String(fileType),
	}

	_, err = s3Client.PutObject(input)
	if err != nil {
		return errors.New("failed to uplaod file")
	}

	return nil
}
func UploadTemp(c *fiber.Ctx, file *multipart.FileHeader) (*Media, error) {
	var sizes []Size
	cfg := config.GetConfig()
	t := time.Now()
	year := fmt.Sprintf("%d", t.Year())
	month := fmt.Sprintf("%02d", t.Month())
	s3Path := "/" + cfg.S3.Dir + "/" + year + "/" + month
	storageDir := filepath.Join(cfg.File.TempDir, s3Path)
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return nil, err
	}
	fileName := renameFile(file.Filename)
	fileTempPath := filepath.Join(storageDir, fileName)
	Name := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	Name = strings.ReplaceAll(Name, "-", " ")

	if err := c.SaveFile(file, fileTempPath); err != nil {
		return nil, err
	}

	webpFileName := fileName[:len(fileName)-len(filepath.Ext(fileName))] + ".webp"
	webpFilePath := filepath.Join(storageDir, webpFileName)
	fileS3Path := filepath.Join(s3Path, webpFileName)
	sizeOriginal, err := createOriginalWebp(fileTempPath, webpFilePath, fileS3Path)
	if err != nil {
		return nil, err
	}
	sizes = append(sizes, *sizeOriginal)

	for _, sz := range []string{"large", "medium", "small"} {
		sizeLarge, err := createWebp(sizeOriginal.File, sz, fileS3Path)
		if err != nil {
			return nil, err
		}
		sizes = append(sizes, *sizeLarge)
	}

	result := Media{
		Name:  Name,
		Alt:   webpFileName,
		Type:  "image/webp",
		Sizes: sizes,
	}
	return &result, nil
}

func createOriginalWebp(filePath string, webpFilePath string, fileS3Path string) (size *Size, err error) {

	if err := convertToWebP(filePath, webpFilePath); err != nil {
		return nil, err
	}
	width, height, err := getWebPDimensions(webpFilePath)
	if err != nil {
		return nil, err
	}
	fileSize, err := getFileSize(webpFilePath)
	if err != nil {
		return nil, err
	}

	return &Size{
		Id:       "original",
		File:     webpFilePath,
		FileS3:   fileS3Path,
		Width:    width,
		Height:   height,
		Filesize: formatSize(fileSize),
	}, err
}
func createWebp(filePath string, size string, fileS3Path string) (*Size, error) {
	newFilePath := strings.ReplaceAll(filePath, ".webp", "-"+size+".webp")
	fileS3Path = strings.ReplaceAll(fileS3Path, ".webp", "-"+size+".webp")
	switch size {
	case "large":
		err := exec.Command("convert", filePath, "-resize", "900x", "-filter", "Lanczos", "-sharpen", "0x1", newFilePath).Run()
		if err != nil {
			return nil, err
		}
	case "medium":
		err := exec.Command("convert", filePath, "-resize", "300x", "-filter", "Lanczos", "-sharpen", "0x1", newFilePath).Run()
		if err != nil {
			return nil, err
		}
	case "small":
		err := exec.Command("convert", filePath, "-resize", "100x100^", "-gravity", "center", "-crop", "100x100+0+0", "-filter", "Lanczos", "-sharpen", "0x1", newFilePath).Run()
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("size is invalid")
	}

	width, height, err := getWebPDimensions(newFilePath)
	if err != nil {
		return nil, err
	}
	fileSize, err := getFileSize(newFilePath)
	if err != nil {
		return nil, err
	}

	return &Size{
		Id:       size,
		File:     newFilePath,
		FileS3:   fileS3Path,
		Width:    width,
		Height:   height,
		Filesize: formatSize(fileSize),
	}, nil
}

func convertToWebP(inputPath, outputPath string) error {
	cmd := exec.Command("cwebp", "-q", "80", inputPath, "-o", outputPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	os.Remove(inputPath)
	return nil
}

func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

func getWebPDimensions(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, err := webp.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Width, img.Height, nil
}

func formatSize(size int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return strconv.FormatInt(size, 10) + " Byte"
	}
}

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func renameFile(filename string) string {
	ext := filepath.Ext(filename)             // Ekstensi file (.jpg)
	name := strings.TrimSuffix(filename, ext) // Nama file tanpa ekstensi
	randomStr := generateRandomString(10)     // String acak 8 karakter
	return fmt.Sprintf("%s-%s%s", name, randomStr, ext)
}
