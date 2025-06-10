package database

import "gorm.io/gorm"

type Entity interface {
	Key() string
}

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any]() *Repository[T] {
	if db == nil {
		panic("database is not connected")
	}
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Insert(entity *T) error {
	tx := r.db.Create(entity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Repository[T]) Update(entity *T) error {
	tx := r.db.Save(entity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Repository[T]) Delete(entity *T) error {
	tx := r.db.Delete(entity)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Repository[T]) FindByID(id interface{}) (*T, error) {
	var entity T
	tx := r.db.First(&entity, "id = ?", id)
	if tx.Error != nil {
		return &entity, tx.Error
	}
	return &entity, nil
}

func (r *Repository[T]) FindAll(conditions ...interface{}) ([]T, error) {
	var entities []T
	tx := r.db.Find(&entities, conditions...)
	if tx.Error != nil {
		return entities, tx.Error
	}
	return entities, nil
}

func (r *Repository[T]) Count(entity *T, count *int64) error {
	tx := r.db.Model(entity).Count(count)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *Repository[T]) Where(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}
