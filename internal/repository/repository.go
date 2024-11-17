package repository

import "gorm.io/gorm"

type Repository interface {
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
	Offset(offset int) *gorm.DB
	Limit(limit int) *gorm.DB
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(value interface{}) *gorm.DB {
	return r.db.Create(value)
}

func (r *GormRepository) Save(value interface{}) *gorm.DB {
	return r.db.Save(value)
}

func (r *GormRepository) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return r.db.First(dest, conds...)
}

func (r *GormRepository) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return r.db.Find(dest, conds...)
}

func (r *GormRepository) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return r.db.Delete(value, conds...)
}

func (r *GormRepository) Where(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

func (r *GormRepository) Offset(offset int) *gorm.DB {
	return r.db.Offset(offset)
}

func (r *GormRepository) Limit(limit int) *gorm.DB {
	return r.db.Limit(limit)
}
