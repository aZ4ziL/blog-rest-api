package models

import (
	"time"

	"gorm.io/gorm"
)

// Category
// this type will declare which categories are available in the blog application.
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:30;unique;index" json:"title"`
	Slug      string         `gorm:"size:100;unique;index" json:"slug"`
	Logo      string         `gorm:"size:255;null" json:"logo"`
	Approved  bool           `gorm:"default:0" json:"approved"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Articles  []Article      `gorm:"foreignKey:CategoryID" json:"articles"`
}

type categoryModel struct {
	db *gorm.DB
}

func NewCategoryModel() *categoryModel {
	return &categoryModel{db}
}

// CreateCategory
// create a new category
func (c *categoryModel) CreateCategory(category *Category) error {
	return c.db.Create(category).Error
}

// GetAllCategories
// returns all categories
func (c *categoryModel) GetAllCategories() []Category {
	var categories []Category
	c.db.Model(&Category{}).Preload("Articles").Find(&categories)
	return categories
}

// GetCategoryByID
// return category by passing the `id`
func (c *categoryModel) GetCategoryByID(id uint) (Category, error) {
	var category Category
	err := c.db.Model(&Category{}).Where("id = ?", id).Preload("Articles").First(&category).Error
	return category, err
}

// GetCategoryBySlug
// return category by passing the `slug`
func (c *categoryModel) GetCategoryBySlug(slug string) (Category, error) {
	var category Category
	err := c.db.Model(&Category{}).Where("slug = ?", slug).Preload("Articles").First(&category).Error
	return category, err
}
