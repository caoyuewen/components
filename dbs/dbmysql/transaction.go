package dbmysql

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// txKeyType 事务上下文键类型
type txKeyType struct{}

var txKey = txKeyType{}

// ==================== 基础事务操作 ====================

// Transaction 执行事务，自动处理提交和回滚
func Transaction(fn func(tx *gorm.DB) error) error {
	return Client().Transaction(fn)
}

// TransactionWithContext 带上下文的事务执行
func TransactionWithContext(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) error {
	return Client().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey, tx)
		return fn(txCtx, tx)
	})
}

// ==================== 手动事务管理 ====================

// BeginTx 开始一个事务，返回事务实例
func BeginTx() *gorm.DB {
	tx := Client().Begin()
	if tx.Error != nil {
		log.Errorf("[MYSQL] Begin transaction failed: %v", tx.Error)
		return nil
	}
	return tx
}

// BeginTxWithContext 带上下文开始事务
func BeginTxWithContext(ctx context.Context) *gorm.DB {
	tx := Client().WithContext(ctx).Begin()
	if tx.Error != nil {
		log.Errorf("[MYSQL] Begin transaction with context failed: %v", tx.Error)
		return nil
	}
	return tx
}

// CommitTx 提交事务
func CommitTx(tx *gorm.DB) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	if err := tx.Commit().Error; err != nil {
		log.Errorf("[MYSQL] Commit transaction failed: %v", err)
		return err
	}
	return nil
}

// RollbackTx 回滚事务
func RollbackTx(tx *gorm.DB) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	if err := tx.Rollback().Error; err != nil {
		log.Errorf("[MYSQL] Rollback transaction failed: %v", err)
		return err
	}
	return nil
}

// ==================== 上下文事务管理 ====================

// GetTxFromContext 从上下文获取事务
func GetTxFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return nil
}

// GetDBOrTx 获取数据库连接或事务（如果上下文中有事务则使用事务）
func GetDBOrTx(ctx context.Context) *gorm.DB {
	if tx := GetTxFromContext(ctx); tx != nil {
		return tx
	}
	return Client()
}

// WithTx 在事务中执行操作，支持嵌套调用
func WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	// 检查是否已在事务中
	if tx := GetTxFromContext(ctx); tx != nil {
		return fn(ctx)
	}

	// 开启新事务
	return TransactionWithContext(ctx, func(txCtx context.Context, tx *gorm.DB) error {
		return fn(txCtx)
	})
}

// ==================== 安全事务操作 ====================

// SafeTransaction 安全事务，带有 panic 恢复
func SafeTransaction(fn func(tx *gorm.DB) error) (err error) {
	tx := BeginTx()
	if tx == nil {
		return fmt.Errorf("failed to begin transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			RollbackTx(tx)
			err = fmt.Errorf("transaction panic: %v", r)
			log.Errorf("[MYSQL] Transaction panic recovered: %v", r)
		}
	}()

	if err = fn(tx); err != nil {
		RollbackTx(tx)
		return err
	}

	return CommitTx(tx)
}

// SafeTransactionWithContext 带上下文的安全事务
func SafeTransactionWithContext(ctx context.Context, fn func(ctx context.Context, tx *gorm.DB) error) (err error) {
	tx := BeginTxWithContext(ctx)
	if tx == nil {
		return fmt.Errorf("failed to begin transaction")
	}

	txCtx := context.WithValue(ctx, txKey, tx)

	defer func() {
		if r := recover(); r != nil {
			RollbackTx(tx)
			err = fmt.Errorf("transaction panic: %v", r)
			log.Errorf("[MYSQL] Transaction panic recovered: %v", r)
		}
	}()

	if err = fn(txCtx, tx); err != nil {
		RollbackTx(tx)
		return err
	}

	return CommitTx(tx)
}

// ==================== 事务辅助函数 ====================

// TxDo 简化事务操作
func TxDo(operations ...func(tx *gorm.DB) error) error {
	return Transaction(func(tx *gorm.DB) error {
		for _, op := range operations {
			if err := op(tx); err != nil {
				return err
			}
		}
		return nil
	})
}

// TxDoWithResult 事务操作并返回结果
func TxDoWithResult[T any](fn func(tx *gorm.DB) (T, error)) (result T, err error) {
	err = Transaction(func(tx *gorm.DB) error {
		result, err = fn(tx)
		return err
	})
	return
}
