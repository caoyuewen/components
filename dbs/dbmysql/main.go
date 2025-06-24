package dbmysql

import (
	"errors"
	"fmt"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrate(tables []interface{}) {
	if err := dbc.AutoMigrate(tables...); err != nil {
		log.Errorf("AutoMigrate err: %v", err)
	}
}

type BaseRepository[T any] struct {
	db       *gorm.DB
	pkColumn string // 默认主键字段名
}

func NewBaseRepository[T any](db *gorm.DB, pkColumn string) BaseRepository[T] {
	return BaseRepository[T]{
		db:       db,
		pkColumn: pkColumn,
	}
}

func (r *BaseRepository[T]) Insert(obj *T) error {
	err := r.db.Create(obj).Error
	if err != nil {
		log.Errorf("Insert err: %s", err.Error())
		return err
	}
	return nil
}

func (r *BaseRepository[T]) Update(obj *T) error {
	err := r.db.Save(obj).Error
	if err != nil {
		log.Errorf("Update err: %s", err.Error())
		return err
	}
	return nil
}

func (r *BaseRepository[T]) Delete(id any) error {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("Delete err, id is not int64: %s", err.Error())
		return err
	}
	err = r.db.Delete(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error
	if err != nil {
		log.Errorf("Delete err : %s", err.Error())
		return err
	}
	return nil
}

func (r *BaseRepository[T]) FindByID(id any) (*T, error) {
	var t *T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("FindByID err: %s", err.Error())
		return t, err
	}
	err = r.db.First(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("FindByID err : %s", err.Error())
		}
		return t, err
	}
	return t, err
}

// FindAll 开发禁用
func (r *BaseRepository[T]) FindAll() ([]*T, error) {
	var list []*T
	err := r.db.Find(&list).Error
	if err != nil {
		log.Errorf("FindAll err: %s", err.Error())
		return list, err
	}
	return list, err
}

func (r *BaseRepository[T]) FindOne(query any, args ...any) (*T, error) {
	var t T
	err := r.db.Where(query, args...).First(&t).Error
	if err != nil {
		log.Errorf("FindOne err: %s", err.Error())
		return &t, err
	}
	return &t, err
}

func (r *BaseRepository[T]) Count(query any, args ...any) (int64, error) {
	var count int64
	err := r.db.Model((*T)(nil)).Where(query, args...).Count(&count).Error
	if err != nil {
		log.Errorf("Count err : %s", err.Error())
		return count, err
	}
	return count, err
}

func (r *BaseRepository[T]) FindPage(offset, limit int, order string, query any, args ...any) ([]*T, int64, error) {
	var (
		list  []*T
		count int64
		db    = r.db.Model((*T)(nil)).Where(query, args...)
	)

	// 加上排序、分页、查找
	if order != "" {
		db = db.Order(order) // 比如 "id desc" 或 "created_at asc"
	}

	if err := db.Count(&count).Error; err != nil {
		log.Errorf("FindPage Count err : %s", err.Error())
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("FindPage err : %s", err.Error())
		}
		return list, 0, err
	}

	return list, count, nil
}

func (r *BaseRepository[T]) Find(order string, query any, args ...any) ([]*T, error) {
	var (
		list []*T

		db = r.db.Model((*T)(nil)).Where(query, args...)
	)

	if order != "" {
		db = db.Order(order) // 比如 "id desc" 或 "created_at asc"
	}

	err := db.Find(&list).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Errorf("Find err : %s", err.Error())
		}
		return list, err
	}

	return list, nil
}
