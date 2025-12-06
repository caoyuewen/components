package dbmysql

import (
	"context"
	"errors"
	"fmt"

	"github.com/caoyuewen/components/util"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ErrRecordNotFound 记录不存在错误
var ErrRecordNotFound = gorm.ErrRecordNotFound

// BaseRepository 泛型基础仓库
type BaseRepository[T any] struct {
	pkColumn string
}

// NewBaseRepository 创建基础仓库实例
func NewBaseRepository[T any](pkColumn string) BaseRepository[T] {
	if pkColumn == "" {
		pkColumn = "id"
	}
	return BaseRepository[T]{pkColumn: pkColumn}
}

// WithTx 返回一个带事务的仓库实例
func (r *BaseRepository[T]) WithTx(tx *gorm.DB) *TxRepository[T] {
	return &TxRepository[T]{
		pkColumn: r.pkColumn,
		tx:       tx,
	}
}

// ==================== 基础 CRUD ====================

// Insert 插入单条记录（支持值或指针）
func (r *BaseRepository[T]) Insert(obj T) error {
	err := Client().Create(&obj).Error
	if err != nil {
		log.Errorf("Insert err: %v", err)
	}
	return err
}

// InsertBatch 批量插入
func (r *BaseRepository[T]) InsertBatch(objs []T) error {
	if len(objs) == 0 {
		return nil
	}
	err := Client().Create(&objs).Error
	if err != nil {
		log.Errorf("InsertBatch err: %v", err)
	}
	return err
}

// InsertBatchSize 批量插入（指定批次大小）
func (r *BaseRepository[T]) InsertBatchSize(objs []T, batchSize int) error {
	if len(objs) == 0 {
		return nil
	}
	err := Client().CreateInBatches(&objs, batchSize).Error
	if err != nil {
		log.Errorf("InsertBatchSize err: %v", err)
	}
	return err
}

// Save 保存记录（存在则更新，不存在则插入）
func (r *BaseRepository[T]) Save(obj *T) error {
	err := Client().Save(obj).Error
	if err != nil {
		log.Errorf("Save err: %v", err)
	}
	return err
}

// Update 更新记录（根据主键更新）
func (r *BaseRepository[T]) Update(obj T) error {
	err := Client().Save(&obj).Error
	if err != nil {
		log.Errorf("Update err: %v", err)
	}
	return err
}

// ==================== 更新操作 ====================

// UpdateByID 根据 ID 更新
func (r *BaseRepository[T]) UpdateByID(id any, updates map[string]any) (int64, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("UpdateByID err, invalid id: %v", err)
		return 0, err
	}

	result := Client().Model(&t).Where(r.pkColumn+" = ?", idI64).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateByID err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateById 根据 ID 更新（兼容旧代码）
func (r *BaseRepository[T]) UpdateById(id any, updates map[string]any) (int64, error) {
	return r.UpdateByID(id, updates)
}

// UpdateByIDs 根据多个 ID 批量更新
func (r *BaseRepository[T]) UpdateByIDs(ids []int64, updates map[string]any) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	var t T
	result := Client().Model(&t).Where(r.pkColumn+" IN ?", ids).Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateByIDs err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateWhere 条件更新
func (r *BaseRepository[T]) UpdateWhere(updates map[string]any, conds ...any) (int64, error) {
	var t T
	db := Client().Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	result := db.Updates(updates)
	if result.Error != nil {
		log.Errorf("UpdateWhere err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// ==================== 删除操作 ====================

// DeleteByID 根据 ID 删除
func (r *BaseRepository[T]) DeleteByID(id any) error {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("DeleteByID err, invalid id: %v", err)
		return err
	}
	err = Client().Delete(&t, r.pkColumn+" = ?", idI64).Error
	if err != nil {
		log.Errorf("DeleteByID err: %v", err)
	}
	return err
}

// Delete 根据 ID 删除（兼容旧代码）
func (r *BaseRepository[T]) Delete(id any) error {
	return r.DeleteByID(id)
}

// DeleteByIDs 根据多个 ID 批量删除
func (r *BaseRepository[T]) DeleteByIDs(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	var t T
	err := Client().Delete(&t, r.pkColumn+" IN ?", ids).Error
	if err != nil {
		log.Errorf("DeleteByIDs err: %v", err)
	}
	return err
}

// DeleteWhere 条件删除
func (r *BaseRepository[T]) DeleteWhere(conds ...any) error {
	var t T
	db := Client().Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.Delete(&t).Error
	if err != nil {
		log.Errorf("DeleteWhere err: %v", err)
	}
	return err
}

// ==================== 查询操作 ====================

// FindByID 根据 ID 查询
func (r *BaseRepository[T]) FindByID(id any) (T, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("FindByID err, invalid id: %v", err)
		return t, err
	}
	err = Client().First(&t, r.pkColumn+" = ?", idI64).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("FindByID err: %v", err)
	}
	return t, err
}

// FindByIDs 根据多个 ID 查询
func (r *BaseRepository[T]) FindByIDs(ids []int64) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}
	var list []T
	err := Client().Where(r.pkColumn+" IN ?", ids).Find(&list).Error
	if err != nil {
		log.Errorf("FindByIDs err: %v", err)
	}
	return list, err
}

// FindOne 查询单条记录
func (r *BaseRepository[T]) FindOne(conds ...any) (T, error) {
	var t T
	db := Client().Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.First(&t).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("FindOne err: %v", err)
	}
	return t, err
}

// FindAll 查询所有记录
func (r *BaseRepository[T]) FindAll() ([]T, error) {
	var list []T
	err := Client().Find(&list).Error
	if err != nil {
		log.Errorf("FindAll err: %v", err)
	}
	return list, err
}

// Find 条件查询
func (r *BaseRepository[T]) Find(order string, conds ...any) ([]T, error) {
	var list []T
	db := Client().Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		log.Errorf("Find err: %v", err)
	}
	return list, err
}

// FindSelect 条件查询（指定字段）
func (r *BaseRepository[T]) FindSelect(fields []string, order string, conds ...any) ([]T, error) {
	var list []T
	db := Client().Model((*T)(nil))
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	for _, cond := range conds {
		db = db.Where(cond)
	}
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		log.Errorf("FindSelect err: %v", err)
	}
	return list, err
}

// ==================== 分页查询 ====================

// PageResult 分页结果
type PageResult[T any] struct {
	List     []T   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

// FindPage 分页查询（使用 offset/limit）
func (r *BaseRepository[T]) FindPage(offset, limit int, order string, conds ...any) ([]T, int64, error) {
	var (
		list  []T
		count int64
		db    = Client().Model((*T)(nil))
	)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	if err := db.Count(&count).Error; err != nil {
		log.Errorf("FindPage count err: %v", err)
		return nil, 0, err
	}
	if order != "" {
		db = db.Order(order)
	}
	err := db.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		log.Errorf("FindPage find err: %v", err)
		return nil, 0, err
	}
	return list, count, nil
}

// Paginate 分页查询（使用 page/pageSize）
func (r *BaseRepository[T]) Paginate(page, pageSize int, order string, conds ...any) (*PageResult[T], error) {
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
	return &PageResult[T]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// ==================== 聚合查询 ====================

// Count 统计数量
func (r *BaseRepository[T]) Count(conds ...any) (int64, error) {
	var count int64
	db := Client().Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.Count(&count).Error
	if err != nil {
		log.Errorf("Count err: %v", err)
	}
	return count, err
}

// Exists 判断是否存在
func (r *BaseRepository[T]) Exists(conds ...any) (bool, error) {
	count, err := r.Count(conds...)
	return count > 0, err
}

// ExistsByID 根据 ID 判断是否存在
func (r *BaseRepository[T]) ExistsByID(id any) (bool, error) {
	idI64, err := util.ToInt64E(id)
	if err != nil {
		return false, err
	}
	return r.Exists(fmt.Sprintf("%s = ?", r.pkColumn), idI64)
}

// ==================== 原生查询 ====================

// Raw 执行原生 SQL 查询
func (r *BaseRepository[T]) Raw(sql string, args ...any) ([]T, error) {
	var list []T
	err := Client().Raw(sql, args...).Scan(&list).Error
	if err != nil {
		log.Errorf("Raw query err: %v", err)
	}
	return list, err
}

// Exec 执行原生 SQL
func (r *BaseRepository[T]) Exec(sql string, args ...any) (int64, error) {
	result := Client().Exec(sql, args...)
	if result.Error != nil {
		log.Errorf("Exec err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// ==================== 辅助方法 ====================

// IsNotFound 判断是否为记录不存在错误
func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// ==================== 事务仓库 ====================

// TxRepository 事务仓库（通过 WithTx 获取）
type TxRepository[T any] struct {
	pkColumn string
	tx       *gorm.DB
}

// Insert 事务内插入
func (r *TxRepository[T]) Insert(obj T) error {
	err := r.tx.Create(&obj).Error
	if err != nil {
		log.Errorf("[TX] Insert err: %v", err)
	}
	return err
}

// InsertBatch 事务内批量插入
func (r *TxRepository[T]) InsertBatch(objs []T) error {
	if len(objs) == 0 {
		return nil
	}
	err := r.tx.Create(&objs).Error
	if err != nil {
		log.Errorf("[TX] InsertBatch err: %v", err)
	}
	return err
}

// Save 事务内保存
func (r *TxRepository[T]) Save(obj *T) error {
	err := r.tx.Save(obj).Error
	if err != nil {
		log.Errorf("[TX] Save err: %v", err)
	}
	return err
}

// Update 事务内更新记录
func (r *TxRepository[T]) Update(obj T) error {
	err := r.tx.Save(&obj).Error
	if err != nil {
		log.Errorf("[TX] Update err: %v", err)
	}
	return err
}

// UpdateByID 事务内根据 ID 更新
func (r *TxRepository[T]) UpdateByID(id any, updates map[string]any) (int64, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("[TX] UpdateByID err, invalid id: %v", err)
		return 0, err
	}

	result := r.tx.Model(&t).Where(r.pkColumn+" = ?", idI64).Updates(updates)
	if result.Error != nil {
		log.Errorf("[TX] UpdateByID err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateById 事务内根据 ID 更新（兼容旧代码）
func (r *TxRepository[T]) UpdateById(id any, updates map[string]any) (int64, error) {
	return r.UpdateByID(id, updates)
}

// UpdateByIDs 事务内根据多个 ID 批量更新
func (r *TxRepository[T]) UpdateByIDs(ids []int64, updates map[string]any) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	var t T
	result := r.tx.Model(&t).Where(r.pkColumn+" IN ?", ids).Updates(updates)
	if result.Error != nil {
		log.Errorf("[TX] UpdateByIDs err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// UpdateWhere 事务内条件更新
func (r *TxRepository[T]) UpdateWhere(updates map[string]any, conds ...any) (int64, error) {
	var t T
	db := r.tx.Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	result := db.Updates(updates)
	if result.Error != nil {
		log.Errorf("[TX] UpdateWhere err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// DeleteByID 事务内根据 ID 删除
func (r *TxRepository[T]) DeleteByID(id any) error {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("[TX] DeleteByID err, invalid id: %v", err)
		return err
	}
	err = r.tx.Delete(&t, r.pkColumn+" = ?", idI64).Error
	if err != nil {
		log.Errorf("[TX] DeleteByID err: %v", err)
	}
	return err
}

// Delete 事务内根据 ID 删除（兼容旧代码）
func (r *TxRepository[T]) Delete(id any) error {
	return r.DeleteByID(id)
}

// DeleteByIDs 事务内根据多个 ID 批量删除
func (r *TxRepository[T]) DeleteByIDs(ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	var t T
	err := r.tx.Delete(&t, r.pkColumn+" IN ?", ids).Error
	if err != nil {
		log.Errorf("[TX] DeleteByIDs err: %v", err)
	}
	return err
}

// DeleteWhere 事务内条件删除
func (r *TxRepository[T]) DeleteWhere(conds ...any) error {
	var t T
	db := r.tx.Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.Delete(&t).Error
	if err != nil {
		log.Errorf("[TX] DeleteWhere err: %v", err)
	}
	return err
}

// FindByID 事务内根据 ID 查询
func (r *TxRepository[T]) FindByID(id any) (T, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("[TX] FindByID err, invalid id: %v", err)
		return t, err
	}
	err = r.tx.First(&t, r.pkColumn+" = ?", idI64).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("[TX] FindByID err: %v", err)
	}
	return t, err
}

// FindByIDForUpdate 事务内根据 ID 查询并锁定
func (r *TxRepository[T]) FindByIDForUpdate(id any) (T, error) {
	var t T
	idI64, err := util.ToInt64E(id)
	if err != nil {
		log.Errorf("[TX] FindByIDForUpdate err, invalid id: %v", err)
		return t, err
	}
	err = r.tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&t, r.pkColumn+" = ?", idI64).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("[TX] FindByIDForUpdate err: %v", err)
	}
	return t, err
}

// FindByIDs 事务内根据多个 ID 查询
func (r *TxRepository[T]) FindByIDs(ids []int64) ([]T, error) {
	if len(ids) == 0 {
		return []T{}, nil
	}
	var list []T
	err := r.tx.Where(r.pkColumn+" IN ?", ids).Find(&list).Error
	if err != nil {
		log.Errorf("[TX] FindByIDs err: %v", err)
	}
	return list, err
}

// FindOne 事务内查询单条记录
func (r *TxRepository[T]) FindOne(conds ...any) (T, error) {
	var t T
	db := r.tx.Model(&t)
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.First(&t).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("[TX] FindOne err: %v", err)
	}
	return t, err
}

// FindOneForUpdate 事务内查询单条并锁定
func (r *TxRepository[T]) FindOneForUpdate(conds ...any) (T, error) {
	var t T
	db := r.tx.Model(&t).Clauses(clause.Locking{Strength: "UPDATE"})
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.First(&t).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("[TX] FindOneForUpdate err: %v", err)
	}
	return t, err
}

// Find 事务内条件查询
func (r *TxRepository[T]) Find(order string, conds ...any) ([]T, error) {
	var list []T
	db := r.tx.Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	if order != "" {
		db = db.Order(order)
	}
	err := db.Find(&list).Error
	if err != nil {
		log.Errorf("[TX] Find err: %v", err)
	}
	return list, err
}

// Count 事务内统计数量
func (r *TxRepository[T]) Count(conds ...any) (int64, error) {
	var count int64
	db := r.tx.Model((*T)(nil))
	for _, cond := range conds {
		db = db.Where(cond)
	}
	err := db.Count(&count).Error
	if err != nil {
		log.Errorf("[TX] Count err: %v", err)
	}
	return count, err
}

// Exists 事务内判断是否存在
func (r *TxRepository[T]) Exists(conds ...any) (bool, error) {
	count, err := r.Count(conds...)
	return count > 0, err
}

// Exec 事务内执行原生 SQL
func (r *TxRepository[T]) Exec(sql string, args ...any) (int64, error) {
	result := r.tx.Exec(sql, args...)
	if result.Error != nil {
		log.Errorf("[TX] Exec err: %v", result.Error)
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// ==================== 事务便捷方法 ====================

// Transaction 执行事务（自动提交/回滚）
func Transaction(fn func(tx *gorm.DB) error) error {
	return TransactionCtx(context.Background(), fn)
}

// TransactionCtx 执行事务（带 Context）
func TransactionCtx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return Client().WithContext(ctx).Transaction(fn)
}

