package dbmysql

import (
	"context"
	"fmt"

	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// BaseRepository 通用 Repository 基类，支持泛型
type BaseRepository[T any] struct {
	pkColumn string
}

// NewBaseRepository 创建新的 BaseRepository 实例
func NewBaseRepository[T any](pkColumn string) BaseRepository[T] {
	return BaseRepository[T]{pkColumn: pkColumn}
}

// ==================== 创建操作 ====================

// Insert 插入单条记录
func (r *BaseRepository[T]) Insert(obj T) error {
	return r.InsertWithDB(Client(), obj)
}

// InsertWithDB 使用指定 DB 插入（支持事务）
func (r *BaseRepository[T]) InsertWithDB(db *gorm.DB, obj T) error {
	if err := db.Create(&obj).Error; err != nil {
		log.Errorf("Insert err: %s", err.Error())
		return err
	}
	return nil
}

// InsertBatch 批量插入
func (r *BaseRepository[T]) InsertBatch(objs []T) error {
	return r.InsertBatchWithDB(Client(), objs)
}

// InsertBatchWithDB 使用指定 DB 批量插入
func (r *BaseRepository[T]) InsertBatchWithDB(db *gorm.DB, objs []T) error {
	if len(objs) == 0 {
		return nil
	}
	if err := db.CreateInBatches(objs, 100).Error; err != nil {
		log.Errorf("InsertBatch err: %s", err.Error())
		return err
	}
	return nil
}

// InsertOrUpdate 插入或更新（Upsert）
func (r *BaseRepository[T]) InsertOrUpdate(obj T, updateColumns []string) error {
	return r.InsertOrUpdateWithDB(Client(), obj, updateColumns)
}

// InsertOrUpdateWithDB 使用指定 DB 插入或更新
func (r *BaseRepository[T]) InsertOrUpdateWithDB(db *gorm.DB, obj T, updateColumns []string) error {
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: r.pkColumn}},
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(&obj).Error; err != nil {
		log.Errorf("InsertOrUpdate err: %s", err.Error())
		return err
	}
	return nil
}

// ==================== 更新操作 ====================

// Update 根据主键更新整个对象
func (r *BaseRepository[T]) Update(obj T) error {
	return r.UpdateWithDB(Client(), obj)
}

// UpdateWithDB 使用指定 DB 更新
func (r *BaseRepository[T]) UpdateWithDB(db *gorm.DB, obj T) error {
	if err := db.Save(&obj).Error; err != nil {
		log.Errorf("Update err: %s", err.Error())
		return err
	}
	return nil
}

// UpdateByID 根据 ID 更新指定字段
func (r *BaseRepository[T]) UpdateByID(id any, updates map[string]interface{}) (int64, error) {
	return r.UpdateByIDWithDB(Client(), id, updates)
}

// UpdateByIDWithDB 使用指定 DB 根据 ID 更新
func (r *BaseRepository[T]) UpdateByIDWithDB(db *gorm.DB, id any, updates map[string]interface{}) (int64, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("UpdateByID err, invalid id: %s", err.Error())
		return 0, err
	}

	result := db.Model(&t).Where("id = ?", idI64).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateByID err: %s", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateByIDs 根据多个 ID 批量更新
func (r *BaseRepository[T]) UpdateByIDs(ids []int64, updates map[string]interface{}) (int64, error) {
	return r.UpdateByIDsWithDB(Client(), ids, updates)
}

// UpdateByIDsWithDB 使用指定 DB 批量更新
func (r *BaseRepository[T]) UpdateByIDsWithDB(db *gorm.DB, ids []int64, updates map[string]interface{}) (int64, error) {
	var t T
	if len(ids) == 0 {
		return 0, nil
	}

	result := db.Model(&t).Where("id IN ?", ids).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateByIDs err: %s", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateWhere 根据条件更新
func (r *BaseRepository[T]) UpdateWhere(updates map[string]interface{}, conds ...interface{}) (int64, error) {
	return r.UpdateWhereWithDB(Client(), updates, conds...)
}

// UpdateWhereWithDB 使用指定 DB 根据条件更新
func (r *BaseRepository[T]) UpdateWhereWithDB(db *gorm.DB, updates map[string]interface{}, conds ...interface{}) (int64, error) {
	var t T
	query := db.Model(&t)
	for _, cond := range conds {
		query = query.Where(cond)
	}
	result := query.Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateWhere err: %s", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateMany 根据多个条件更新（旧方法保持兼容）
func (r *BaseRepository[T]) UpdateMany(conds []interface{}, updates map[string]interface{}) (int64, error) {
	return r.UpdateWhere(updates, conds...)
}

// UpdateWhereRaw 使用原始 SQL 条件更新
func (r *BaseRepository[T]) UpdateWhereRaw(whereSQL string, args []any, updates map[string]interface{}) (int64, error) {
	return r.UpdateWhereRawWithDB(Client(), whereSQL, args, updates)
}

// UpdateWhereRawWithDB 使用指定 DB 原始 SQL 更新
func (r *BaseRepository[T]) UpdateWhereRawWithDB(db *gorm.DB, whereSQL string, args []any, updates map[string]interface{}) (int64, error) {
	var t T
	result := db.Model(&t).Where(whereSQL, args...).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateWhereRaw err: %s", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// ==================== 删除操作 ====================

// Delete 根据 ID 删除
func (r *BaseRepository[T]) Delete(id any) error {
	return r.DeleteWithDB(Client(), id)
}

// DeleteWithDB 使用指定 DB 删除
func (r *BaseRepository[T]) DeleteWithDB(db *gorm.DB, id any) error {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("Delete err, invalid id: %s", err.Error())
		return err
	}
	if err := db.Delete(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error; err != nil {
		log.Errorf("Delete err: %s", err.Error())
		return err
	}
	return nil
}

// DeleteByIDs 根据多个 ID 批量删除
func (r *BaseRepository[T]) DeleteByIDs(ids []int64) error {
	return r.DeleteByIDsWithDB(Client(), ids)
}

// DeleteByIDsWithDB 使用指定 DB 批量删除
func (r *BaseRepository[T]) DeleteByIDsWithDB(db *gorm.DB, ids []int64) error {
	var t T
	if len(ids) == 0 {
		return nil
	}
	if err := db.Where("id IN ?", ids).Delete(&t).Error; err != nil {
		log.Errorf("DeleteByIDs err: %s", err.Error())
		return err
	}
	return nil
}

// DeleteWhere 根据条件删除
func (r *BaseRepository[T]) DeleteWhere(conds ...interface{}) error {
	return r.DeleteWhereWithDB(Client(), conds...)
}

// DeleteWhereWithDB 使用指定 DB 条件删除
func (r *BaseRepository[T]) DeleteWhereWithDB(db *gorm.DB, conds ...interface{}) error {
	var t T
	query := db.Model(&t)
	for _, cond := range conds {
		query = query.Where(cond)
	}
	if err := query.Delete(&t).Error; err != nil {
		log.Errorf("DeleteWhere err: %s", err.Error())
		return err
	}
	return nil
}

// ==================== 查询操作 ====================

// FindByID 根据 ID 查询单条记录
func (r *BaseRepository[T]) FindByID(id any) (T, error) {
	return r.FindByIDWithDB(Client(), id)
}

// FindByIDWithDB 使用指定 DB 查询
func (r *BaseRepository[T]) FindByIDWithDB(db *gorm.DB, id any) (T, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("FindByID err (invalid id): %s", err.Error())
		return t, err
	}
	if err := db.First(&t, fmt.Sprintf("%s = ?", r.pkColumn), idI64).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Errorf("FindByID err: %s", err.Error())
		}
		return t, err
	}
	return t, nil
}

// FindByIDs 根据多个 ID 查询
func (r *BaseRepository[T]) FindByIDs(ids []int64) ([]T, error) {
	return r.FindByIDsWithDB(Client(), ids)
}

// FindByIDsWithDB 使用指定 DB 批量查询
func (r *BaseRepository[T]) FindByIDsWithDB(db *gorm.DB, ids []int64) ([]T, error) {
	var list []T
	if len(ids) == 0 {
		return list, nil
	}
	if err := db.Where(fmt.Sprintf("%s IN ?", r.pkColumn), ids).Find(&list).Error; err != nil {
		log.Errorf("FindByIDs err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

// FindOne 根据条件查询单条记录
func (r *BaseRepository[T]) FindOne(conds ...interface{}) (T, error) {
	return r.FindOneWithDB(Client(), conds...)
}

// FindOneWithDB 使用指定 DB 查询单条
func (r *BaseRepository[T]) FindOneWithDB(db *gorm.DB, conds ...interface{}) (T, error) {
	var t T
	query := db.Model((*T)(nil))
	
	// 处理条件参数
	if len(conds) > 0 {
		// 如果第一个参数是字符串且包含 ?，则作为单个 Where 条件处理
		if str, ok := conds[0].(string); ok && len(conds) > 1 {
			// 检查是否包含占位符
			hasPlaceholder := false
			for _, char := range str {
				if char == '?' {
					hasPlaceholder = true
					break
				}
			}
			if hasPlaceholder {
				// 将第一个参数作为条件，其余作为值
				query = query.Where(str, conds[1:]...)
			} else {
				// 没有占位符，按原来的方式处理
				for _, cond := range conds {
					query = query.Where(cond)
				}
			}
		} else {
			// 其他情况，按原来的方式处理
			for _, cond := range conds {
				query = query.Where(cond)
			}
		}
	}
	
	if err := query.First(&t).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Errorf("FindOne err: %s", err.Error())
		}
		return t, err
	}
	return t, nil
}

// FindAll 查询所有记录
func (r *BaseRepository[T]) FindAll() ([]T, error) {
	return r.FindAllWithDB(Client())
}

// FindAllWithDB 使用指定 DB 查询所有
func (r *BaseRepository[T]) FindAllWithDB(db *gorm.DB) ([]T, error) {
	var list []T
	if err := db.Find(&list).Error; err != nil {
		log.Errorf("FindAll err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

// Find 根据条件查询多条记录
func (r *BaseRepository[T]) Find(order string, conds ...interface{}) ([]T, error) {
	return r.FindWithDB(Client(), order, conds...)
}

// FindWithDB 使用指定 DB 条件查询
func (r *BaseRepository[T]) FindWithDB(db *gorm.DB, order string, conds ...interface{}) ([]T, error) {
	var list []T
	query := db.Model((*T)(nil))
	for _, cond := range conds {
		query = query.Where(cond)
	}
	if order != "" {
		query = query.Order(order)
	}
	if err := query.Find(&list).Error; err != nil {
		log.Errorf("Find err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

// ==================== 统计查询 ====================

// Count 统计符合条件的记录数
func (r *BaseRepository[T]) Count(conds ...interface{}) (int64, error) {
	return r.CountWithDB(Client(), conds...)
}

// CountWithDB 使用指定 DB 统计
func (r *BaseRepository[T]) CountWithDB(db *gorm.DB, conds ...interface{}) (int64, error) {
	var count int64
	query := db.Model((*T)(nil))
	for _, cond := range conds {
		query = query.Where(cond)
	}
	if err := query.Count(&count).Error; err != nil {
		log.Errorf("Count err: %s", err.Error())
		return 0, err
	}
	return count, nil
}

// Exists 检查是否存在符合条件的记录
func (r *BaseRepository[T]) Exists(conds ...interface{}) (bool, error) {
	count, err := r.Count(conds...)
	return count > 0, err
}

// ExistsWithDB 使用指定 DB 检查存在
func (r *BaseRepository[T]) ExistsWithDB(db *gorm.DB, conds ...interface{}) (bool, error) {
	count, err := r.CountWithDB(db, conds...)
	return count > 0, err
}

// ==================== 分页查询 ====================

// FindPage 分页查询
func (r *BaseRepository[T]) FindPage(offset, limit int, order string, conds ...interface{}) ([]T, int64, error) {
	return r.FindPageWithDB(Client(), offset, limit, order, conds...)
}

// FindPageWithDB 使用指定 DB 分页查询
func (r *BaseRepository[T]) FindPageWithDB(db *gorm.DB, offset, limit int, order string, conds ...interface{}) ([]T, int64, error) {
	var (
		list  []T
		count int64
	)

	query := db.Model((*T)(nil))
	for _, cond := range conds {
		query = query.Where(cond)
	}

	// 先统计总数
	if err := query.Count(&count).Error; err != nil {
		log.Errorf("FindPage count err: %s", err.Error())
		return nil, 0, err
	}

	// 再查询数据
	if order != "" {
		query = query.Order(order)
	}
	if err := query.Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		log.Errorf("FindPage find err: %s", err.Error())
		return nil, 0, err
	}

	return list, count, nil
}

// PageResult 分页结果
type PageResult[T any] struct {
	List       []T   `json:"list"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// FindPageResult 分页查询，返回 PageResult
func (r *BaseRepository[T]) FindPageResult(page, pageSize int, order string, conds ...interface{}) (*PageResult[T], error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	list, total, err := r.FindPage(offset, pageSize, order, conds...)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &PageResult[T]{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// ==================== 事务操作 ====================

// Transaction 在事务中执行操作
func (r *BaseRepository[T]) Transaction(fn func(tx *gorm.DB) error) error {
	return Client().Transaction(fn)
}

// TransactionWithContext 带上下文的事务
func (r *BaseRepository[T]) TransactionWithContext(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return Client().WithContext(ctx).Transaction(fn)
}

// ==================== 原始 SQL ====================

// Raw 执行原始 SQL 查询
func (r *BaseRepository[T]) Raw(sql string, args ...interface{}) ([]T, error) {
	var list []T
	if err := Client().Raw(sql, args...).Scan(&list).Error; err != nil {
		log.Errorf("Raw err: %s", err.Error())
		return nil, err
	}
	return list, nil
}

// Exec 执行原始 SQL（INSERT/UPDATE/DELETE）
func (r *BaseRepository[T]) Exec(sql string, args ...interface{}) (int64, error) {
	result := Client().Exec(sql, args...)
	if result.Error != nil {
		log.Errorf("Exec err: %s", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}
