package mysql

import (
	"fmt"
	"github.com/shayamvlmna/lift/pkg/models"

	"gorm.io/gorm"
)

func OpenAdminDb() (*gorm.DB, error) {
	Db, err := openDB()
	if err != nil {
		return nil, err
	}
	admin := &models.Admin{}
	err = Db.AutoMigrate(&admin)
	if err != nil {
		return nil, err
	}
	fmt.Println("admin db opened")
	return Db, nil
}
func AddAdmin(admin *models.Admin) error {
	db, err := OpenAdminDb()
	if err != nil {
		return err
	}
	// defer closeUserdb(db)
	result := db.Create(admin)

	return result.Error
}

func GetAdmin(username string) (models.Admin, bool) {
	db, err := OpenAdminDb()
	if err != nil {
		fmt.Println(err)
	}
	admin := &models.Admin{}
	result := db.Where("username=?", username).First(&admin)
	if result.Error == gorm.ErrRecordNotFound {
		return *admin, false
	}
	return *admin, true
}
