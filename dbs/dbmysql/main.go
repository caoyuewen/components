package dbmysql

import (
	"fmt"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
)

type BaseRepository[T any] struct {
	pkColumn string
}

func NewBaseRepository[T any](pkColumn string) BaseRepository[T] {
	return BaseRepository[T]{pkColumn: pkColumn}
}

func (r *BaseRepository[T]) Insert(obj T) error {
	err := dbc.Create(&obj).Error
	if err != nil {
		log.Errorf("Insert err: %s", err.Error())
	}
	return err
}

func (r *BaseRepository[T]) Update(obj T) error {
	err := dbc.Save(&obj).Error
	if err != nil {
		log.Errorf("Update err: %s", err.Error())
	}
	return err
}

func (r *BaseRepository[T]) UpdateById(id string, updates map[string]interface{}) (int64, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("UpdateById err, id is not int64: %s", err.Error())
		return 0, err
	}

	result := dbc.Model(&t).Where("id = ?", idI64).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateById err: %s", result.Error)
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		log.Warnf("UpdateById: no rows affected for id = %v", id)
	}
	return result.RowsAffected, nil
}

func (r *BaseRepository[T]) UpdateByIDs(ids []int64, updates map[string]interface{}) (int64, error) {
	var t T
	if len(ids) == 0 {
		log.Warn("UpdateByIDs: empty ID list")
		return 0, nil
	}

	result := dbc.Model(&t).Where("id IN ?", ids).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateByIDs err: %s", result.Error)
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		log.Warnf("UpdateByIDs: no rows affected for ids = %v", ids)
	}
	return result.RowsAffected, nil
}

func (r *BaseRepository[T]) UpdateMany(conds []interface{}, updates map[string]interface{}) (int64, error) {
	var t T
	db := dbc.Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	result := db.Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateMany err: %s", result.Error)
		return 0, result.Error
	}

	if result.RowsAffected == 0 {
		log.Warn("UpdateMany: no rows affected")
	}
	return result.RowsAffected, nil
}

func (r *BaseRepository[T]) UpdateWhereRaw(whereSQL string, args []any, updates map[string]interface{}) (int64, error) {
	var t T
	result := dbc.Model(&t).Where(whereSQL, args...).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateWhereRaw err: %s", result.Error)
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		log.Warnf("UpdateWhereRaw: no rows affected for condition = %s", whereSQL)
	}
	return result.RowsAffected, nil
}

func (r *BaseRepository[T]) Delete(id any) error {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("Delete err, id is not int64: %s", err.Error())
		return err
	}
	err = dbc.Delete(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error
	if err != nil {
		log.Errorf("Delete err: %s", err.Error())
	}
	return err
}

func (r *BaseRepository[T]) DeleteWhere(conds ...interface{}) error {
	var t T
	db := dbc.Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.Delete(&t).Error
	if err != nil {
		log.Errorf("DeleteWhere err: %s", err.Error())
		return err
	}
	return nil
}

func (r *BaseRepository[T]) FindByID(id any) (T, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("FindByID err (invalid id): %s", err.Error())
		return t, err
	}
	err = dbc.First(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error
	if err != nil {
		log.Errorf("FindByID err: %s", err.Error())
		return t, err
	}
	return t, nil
}

func (r *BaseRepository[T]) FindOne(conds ...interface{}) (T, error) {
	var t T
	db := dbc.Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.First(&t).Error
	if err != nil {
		log.Errorf("FindOne err: %s", err.Error())
		return t, err
	}
	return t, nil
}

func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var list []T
	err := dbc.Find(&list).Error
	if err != nil {
		log.Errorf("FindAll err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

func (r *BaseRepository[T]) Find(order string, conds ...interface{}) ([]T, error) {
	var list []T
	db := dbc.Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		log.Errorf("Find err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

func (r *BaseRepository[T]) Count(conds ...interface{}) (int64, error) {
	var count int64
	db := dbc.Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.Count(&count).Error
	if err != nil {
		log.Errorf("Count err: %s", err.Error())
		return 0, err
	}
	return count, nil
}

func (r *BaseRepository[T]) FindPage(offset, limit int, order string, conds ...interface{}) ([]T, int64, error) {
	var (
		list  []T
		count int64
		db    = dbc.Model((*T)(nil))
	)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	if order != "" {
		db = db.Order(order)
	}
	if err := db.Count(&count).Error; err != nil {
		log.Errorf("FindPage count err: %s", err.Error())
		return nil, 0, err
	}
	err := db.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		log.Errorf("FindPage find err: %s", err.Error())
		return nil, 0, err
	}
	return list, count, nil
}
