package dbmodel

import "gorm.io/gorm"

type Field struct {
	gorm.Model
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`

	User  User     `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Media []*Media `gorm:"many2many:media_field;constraint:OnDelete:CASCADE" json:"media"`
}

type FieldRepository interface {
	Create(field *Field) (*Field, error)
	FindAll() ([]*Field, error)
	FindById(id uint) (*Field, error)
	FindMediaByField(id uint) ([]Media, error)
	Update(field *Field) (*Field, error)
	Delete(id uint) error
}

type fieldRepository struct {
	db *gorm.DB
}

func NewFieldRepository(db *gorm.DB) FieldRepository {
	return &fieldRepository{db: db}
}

// Create the field
func (r *fieldRepository) Create(field *Field) (*Field, error) {
	if err := r.db.Create(field).Error; err != nil {
		return nil, err
	}

	return field, nil
}

// Find all field.
func (r *fieldRepository) FindAll() ([]*Field, error) {
	var fields []*Field
	if err := r.db.Find(&fields).Error; err != nil {
		return nil, err
	}

	return fields, nil
}

// Find a field by his id.
func (r *fieldRepository) FindById(id uint) (*Field, error) {
	var field Field
	if err := r.db.First(&field, id).Error; err != nil {
		return nil, err
	}

	return &field, nil
}

// Find a field by his id.
func (r *fieldRepository) FindMediaByField(id uint) ([]Media, error) {
	var field Field
	if err := r.db.Preload("Media").First(&field, id).Error; err != nil {
		return nil, err
	}
	media := make([]Media, len(field.Media))
	for i, u := range field.Media {
		media[i] = *u
	}

	return media, nil
}

// Update the given field.
func (r *fieldRepository) Update(field *Field) (*Field, error) {
	if err := r.db.Save(field).Error; err != nil {
		return nil, err
	}

	return field, nil
}

// Delete a field by his id.
func (r *fieldRepository) Delete(id uint) error {
	if err := r.db.Delete(Field{}, id).Error; err != nil {
		return err
	}

	return nil
}
