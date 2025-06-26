package dbmysql

// Package dbmysql 提供通用的数据库操作封装。
//
// base_repository.go 封装了对任意模型（Model）的基础增删改查操作，适用于 GORM 的泛型支持。
// 它简化了常规的数据库操作逻辑，使各业务仓储层可通过组合 BaseRepository 来复用常用方法。
//
// 使用方式：
//   repo := NewBaseRepository[YourModel](db, "id")
//   err := repo.Insert(&obj)
//   data, err := repo.FindByID(123)
//   list, count, err := repo.FindPage(0, 10, "created_at desc", "status = ?", 1)
//
// 设计规范：
// - Insert / Update / Delete / Count 操作记录错误日志并返回错误
// - Find / FindAll / FindPage 查询不到数据时返回空 slice，不视为错误
// - FindByID / FindOne 查询不到数据时返回 gorm.ErrRecordNotFound，便于上层判断“未找到”
// - 默认主键字段为 BIGINT UNSIGNED，建议业务 ID 使用 string 表示（全局唯一）
//
// 注意事项：
// - 泛型 T 必须是结构体类型，不能是指针类型
// - 使用 ToInt64E 将 string 或 interface{} 类型转换为 int64 用于主键查找

import (
	"fmt"
	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
)

type BaseRepository[T any] struct {
	pkColumn string
}

func NewBaseRepository[T any](pkColumn string) BaseRepository[T] {
	return BaseRepository[T]{
		pkColumn: pkColumn,
	}
}

func (r *BaseRepository[T]) Insert(obj *T) error {
	err := dbc.Create(obj).Error
	if err != nil {
		log.Errorf("Insert err: %s", err.Error())
	}
	return err
}

func (r *BaseRepository[T]) Update(obj *T) error {
	err := dbc.Save(obj).Error
	if err != nil {
		log.Errorf("Update err: %s", err.Error())
	}
	return err
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

func (r *BaseRepository[T]) DeleteWhere(query any, args ...any) error {
	var t T
	err := dbc.Where(query, args...).Delete(&t).Error
	if err != nil {
		log.Errorf("DeleteWhere err: %s", err.Error())
		return err
	}
	return nil
}

func (r *BaseRepository[T]) FindByID(id any) (*T, error) {
	var t *T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("FindByID err (invalid id): %s", err.Error())
		return nil, err
	}
	err = dbc.First(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error
	if err != nil {
		log.Errorf("FindByID err: %s", err.Error())
		return nil, err
	}
	return t, nil
}

func (r *BaseRepository[T]) FindOne(query any, args ...any) (*T, error) {
	var t T
	err := dbc.Where(query, args...).First(&t).Error
	if err != nil {
		log.Errorf("FindOne err: %s", err.Error())
		return nil, err
	}
	return &t, nil
}

func (r *BaseRepository[T]) FindAll() ([]*T, error) {
	var list []*T
	err := dbc.Find(&list).Error
	if err != nil {
		log.Errorf("FindAll err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

func (r *BaseRepository[T]) Find(order string, query any, args ...any) ([]*T, error) {
	var list []*T
	db := dbc.Model((*T)(nil)).Where(query, args...)

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

func (r *BaseRepository[T]) Count(query any, args ...any) (int64, error) {
	var count int64
	err := dbc.Model((*T)(nil)).Where(query, args...).Count(&count).Error
	if err != nil {
		log.Errorf("Count err: %s", err.Error())
		return 0, err
	}
	return count, nil
}

func (r *BaseRepository[T]) FindPage(offset, limit int, order string, query any, args ...any) ([]*T, int64, error) {
	var (
		list  []*T
		count int64
		db    = dbc.Model((*T)(nil)).Where(query, args...)
	)

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
