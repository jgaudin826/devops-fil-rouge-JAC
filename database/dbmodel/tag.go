package dbmodel

import "gorm.io/gorm"

type Tag struct {
	gorm.Model
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`

	Media []*Media `gorm:"many2many:media_tag;constraint:OnDelete:CASCADE" json:"media"`
}

type TagRepository interface {
	Create(tag *Tag) (*Tag, error)
	FindAll() ([]*Tag, error)
	FindById(id uint) (*Tag, error)
	FindMediaByTag(id uint) ([]Media, error)
	Update(tag *Tag) (*Tag, error)
	Delete(id uint) error
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{db: db}
}

// Create the tag
func (r *tagRepository) Create(tag *Tag) (*Tag, error) {
	if err := r.db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// Find all tag.
func (r *tagRepository) FindAll() ([]*Tag, error) {
	var tags []*Tag
	if err := r.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	return tags, nil
}

// Find a tag by his id.
func (r *tagRepository) FindById(id uint) (*Tag, error) {
	var tag Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}

// Find media with tag.
func (r *tagRepository) FindMediaByTag(id uint) ([]Media, error) {
	var tag Tag
	if err := r.db.Preload("Media").First(&tag, id).Error; err != nil {
		return nil, err
	}
	media := make([]Media, len(tag.Media))
	for i, u := range tag.Media {
		media[i] = *u
	}

	return media, nil
}

// Update the given tag.
func (r *tagRepository) Update(tag *Tag) (*Tag, error) {
	if err := r.db.Save(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete a tag by is id.
func (r *tagRepository) Delete(id uint) error {
	if err := r.db.Delete(Tag{}, id).Error; err != nil {
		return err
	}

	return nil
}
