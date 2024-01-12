package model

import (
	"context"

	"gorm.io/gorm"
)

type TestModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *TestModelRepository {
	return &TestModelRepository{
		db: db,
	}
}

func (r *TestModelRepository) FindAll(ctx context.Context) ([]TestModel, error) {

	var models []TestModel

	res := r.db.WithContext(ctx).Find(&models)
	if res.Error != nil {
		return nil, res.Error
	}

	return models, nil
}

func (r *TestModelRepository) Create(ctx context.Context, model *TestModel) error {
	res := r.db.WithContext(ctx).Create(model)

	return res.Error
}
