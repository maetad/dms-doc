package doc

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Get() ([]Doc, error)
	Create(d *Doc) error
	Find(id string) (Doc, error)
	Update(id string, d *Doc) error
	Delete(id string) error
}

type repo struct {
	DB *gorm.DB
}

type NotFoundError struct {
	Id string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Doc id %s is not found", e.Id)
}

func (r *repo) Get() ([]Doc, error) {
	var docs []Doc

	if err := r.DB.Preload("User").Find(&docs).Error; err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *repo) Create(d *Doc) error {
	return r.DB.Create(d).Error
}

func (r *repo) Find(id string) (Doc, error) {
	var doc Doc

	if err := r.DB.Preload("User").First(&doc, id).Error; err != nil {
		return Doc{}, &NotFoundError{Id: id}
	}

	return doc, nil
}

func (r *repo) Update(id string, d *Doc) error {
	doc, err := r.Find(id)

	if err != nil {
		return &NotFoundError{Id: id}
	}

	d.ID = doc.ID

	r.DB.Model(&Doc{}).Where("id = ?", id).Updates(d)

	return nil
}

func (r *repo) Delete(id string) error {
	doc, err := r.Find(id)

	if err != nil {
		return &NotFoundError{Id: id}
	}

	return r.DB.Delete(&doc).Error
}

func CreateRepository(db *gorm.DB) Repository {
	return &repo{
		DB: db,
	}
}
