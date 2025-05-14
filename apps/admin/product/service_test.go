package product

import (
	"github.com/bonsus/go-saas/internal/config"
	"github.com/bonsus/go-saas/internal/database"
)

var svc service

func init() {
	filename := "../../config.yaml"
	err := config.LoadConfig(filename)
	if err != nil {
		panic(err)
	}
	db, err := database.ConnectDB(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := NewRepository(db)
	svc = *NewService(repo)
}

// func TestAccount(t *testing.T) {
// 	t.Run("CreateSuccess", func(t *testing.T) {
// 		req := Request{
// 			PrivateId:     "",
// 			No:            uuid.NewString(),
// 			Name:          "Test Create Nama Akun ",
// 			NameAlias:     "Indonesia",
// 			Type:          "primary",
// 			Activity:      "aktifitas",
// 			Class:         "kelas",
// 			Normal:        "debit",
// 			AcceptPayment: "yes",
// 			Note:          "ini note ",
// 		}
// 		_, errors, err := svc.Create(context.Background(), req)
// 		assert.Nil(t, err)
// 		// fmt.Println(id)
// 		fmt.Println(errors)
// 		fmt.Println(err)
// 	})
// 	t.Run("UpdateSuccess", func(t *testing.T) {
// 		req := Request{
// 			PrivateId:     "",
// 			No:            uuid.NewString(),
// 			Name:          "Nama  Akun Sebelum Diedit",
// 			NameAlias:     "Nama Akun Sebelum Diedit",
// 			Type:          "primary",
// 			Activity:      "aktifitas ",
// 			Class:         "kelas",
// 			Normal:        "debit",
// 			AcceptPayment: "yes",
// 			Note:          "ini note",
// 		}
// 		result, _, err := svc.Create(context.Background(), req)
// 		reqEdit := Request{
// 			PrivateId:     "",
// 			No:            uuid.NewString(),
// 			Name:          "Edit Nama Akunx",
// 			NameAlias:     "Edit Nama Akun Sebelum Diedit",
// 			Type:          "edit primary",
// 			Activity:      "edit aktifitas",
// 			Class:         "edit kelas",
// 			Normal:        "edit debit",
// 			AcceptPayment: "edit yes",
// 			Note:          "edit ini note ",
// 		}
// 		_, errors, err := svc.Update(context.Background(), reqEdit, result.Id)
// 		assert.Nil(t, err)
// 		fmt.Println(errors)
// 		fmt.Println(err)
// 	})
// 	t.Run("ReadSuccess", func(t *testing.T) {
// 		req := Request{
// 			PrivateId:     "",
// 			No:            uuid.NewString(),
// 			Name:          "Test Akun Read",
// 			NameAlias:     "Alias Akun Read ",
// 			Type:          "primary",
// 			Activity:      "aktifitas ",
// 			Class:         "kelas",
// 			Normal:        "debit",
// 			AcceptPayment: "yes",
// 			Note:          "ini note",
// 		}
// 		result, _, err := svc.Create(context.Background(), req)
// 		_, err = svc.Read(context.Background(), result.Id)
// 		assert.Nil(t, err)
// 		// fmt.Println(id)
// 		fmt.Println(err)
// 	})
// 	t.Run("IndexSuccess", func(t *testing.T) {
// 		req := ParamIndex{
// 			Page:    1,
// 			PerPage: 3,
// 		}
// 		_, err := svc.Index(context.Background(), req)
// 		assert.Nil(t, err)
// 		fmt.Println(err)
// 		// fmt.Println(result)
// 	})
// }
