package app

import (
	"github.com/jinzhu/gorm"

)

// DocumentDAO handles database operations for the Document model
type DocumentDAO struct {
	db *gorm.DB
}



// NewDocumentDAO creates a new instance of DocumentDAO
func NewDocumentDAO(db *gorm.DB) *DocumentDAO {
	return &DocumentDAO{db: db}
}

// Create inserts a new Document into the database
func (dao *DocumentDAO) Create(document *Document) error {
	return dao.db.Create(document).Error
}

// GetByID retrieves a Document by its ID
func (dao *DocumentDAO) GetByID(id uint) (*Document, error) {
	var document Document
	if err := dao.db.First(&document, id).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

// GetAll retrieves all Documents from the database
func (dao *DocumentDAO) GetAll() ([]Document, error) {
	var documents []Document
	if err := dao.db.Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}

// Update modifies an existing Document
func (dao *DocumentDAO) Update(document *Document) error {
	return dao.db.Save(document).Error
}

// Delete removes a Document from the database
func (dao *DocumentDAO) Delete(id uint) error {
	return dao.db.Delete(&Document{}, id).Error
}
