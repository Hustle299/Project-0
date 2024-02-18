package models

import (
	"github.com/jinzhu/gorm"
)

// tao services cho tung loai
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: NewGalleryServices(db),
		db:      db,
		Image:   NewImageService(),
	}, nil
}

type Services struct {
	Gallery GalleryServices
	User    UserServices
	db      *gorm.DB
	Image   ImageService
}

func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}, &Gallery{}).Error
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}, &Gallery{}).Error
	if err != nil {
		return err
	}
	return s.AutoMigrate()
}
