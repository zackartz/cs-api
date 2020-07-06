package models

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"html"
	"strings"
	"time"
)

type CodeSnippet struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Slug      string    `gorm:"not null" json:"slug"`
	Code      string    `gorm:"type:longtext;not null" json:"code"`
	Language  string    `gorm:"not null" json:"language"`
	Private   bool      `gorm:"not null" json:"private"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *CodeSnippet) Prepare() {
	guid := uuid.New()
	c.ID = 0
	c.Code = html.EscapeString(strings.TrimSpace(c.Code))
	c.Language = html.EscapeString(strings.TrimSpace(c.Language))
	c.Slug = guid.String()
	c.UpdatedAt = time.Now()
	c.CreatedAt = time.Now()
}

func (c *CodeSnippet) SaveCode(db *gorm.DB) (*CodeSnippet, error) {
	var err error
	err = db.Debug().Create(&c).Error
	if err != nil {
		return &CodeSnippet{}, err
	}
	return c, nil
}

func (c *CodeSnippet) FindAllCodeSnippets(db *gorm.DB) (*[]CodeSnippet, error) {
	var err error
	var snippets []CodeSnippet
	err = db.Debug().Model(&CodeSnippet{}).Where("private = ?", false).Limit(100).Find(&snippets).Error
	return &snippets, err
}

func (c *CodeSnippet) FindSnippetBySlug(db *gorm.DB, slug string) (*CodeSnippet, error) {
	var err error
	err = db.Debug().Model(&CodeSnippet{}).Where("slug = ?", slug).Take(&c).Error
	if err != nil {
		return &CodeSnippet{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &CodeSnippet{}, errors.New("snippet not found")
	}
	return c, err
}

func (c *CodeSnippet) FindSnippetByID(db *gorm.DB, cid uint32) (*CodeSnippet, error) {
	var err error
	err = db.Debug().Model(&CodeSnippet{}).Where("id = ?", cid).Take(&c).Error
	if err != nil {
		return &CodeSnippet{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &CodeSnippet{}, errors.New("snippet not found")
	}
	return c, err
}

func (c *CodeSnippet) Delete(db *gorm.DB, cid uint32) (int64, error) {
	db = db.Debug().Model(&CodeSnippet{}).Where("id = ?", cid).Take(&CodeSnippet{}).Delete(&CodeSnippet{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (c *CodeSnippet) Validate() error {
	if c.Code == "" {
		return errors.New("code is required")
	}
	return nil
}
