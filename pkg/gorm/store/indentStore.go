package store

import (
	"github.com/jordy2254/indoormaprestapi/pkg/model"
	"gorm.io/gorm"
)

type IndentStore struct {
	DB *gorm.DB
}

func NewIndentStore(DB *gorm.DB) IndentStore {
	return IndentStore{DB: DB}
}

func (is *IndentStore) CreateIndent(indent *model.Indent) {
	is.DB.Create(indent)
}

func (is *IndentStore) GetIndentById(id int) model.Indent {
	var indent model.Indent
	is.DB.Find(&indent, "indents.id=?", id)
	return indent
}

func (is *IndentStore) UpdateIndent(indent model.Indent) {
	is.DB.Model(&indent).UpdateColumns(&indent)
}

func (is *IndentStore) DeleteIndent(id int) {
	is.DB.Delete(&model.Indent{}, id)
}
