package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"platform-backend/utils/log"
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

// Create 通用的创建方法，接受具体的 model 类型
func (b *BaseRepository[T]) Create(ctx context.Context, model *T) error {
	return b.db.WithContext(ctx).Create(model).Error
}

// BatchInsert 批量插入方法
func (r *BaseRepository[T]) BatchInsert(ctx context.Context, models []T) error {
	if len(models) == 0 {
		return nil // 如果没有数据，直接返回
	}

	return r.db.WithContext(ctx).Create(&models).Error
}

// Update 通用的更新方法
func (b *BaseRepository[T]) Update(ctx context.Context, model *T) error {
	return b.db.WithContext(ctx).Save(model).Error
}

func (b *BaseRepository[T]) UpdateByCondition(ctx context.Context, condition map[string]interface{}, model *T) error {
	return b.db.WithContext(ctx).Where(condition).Updates(model).Error
}

func (b *BaseRepository[T]) DeleteByCondition(ctx context.Context, condition map[string]interface{}) error {
	return b.db.WithContext(ctx).Where(condition).Delete(new(T)).Error
}

func (b *BaseRepository[T]) FindOne(ctx context.Context, condition map[string]interface{}) (*T, error) {
	var result T
	err := b.db.WithContext(ctx).Where(condition).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindOneForUpdate 实现带 FOR UPDATE 的单条记录查询
func (b *BaseRepository[T]) FindOneForUpdate(ctx context.Context, condition map[string]interface{}) (*T, error) {
	var result T
	// 加入 FOR UPDATE 锁
	err := b.db.WithContext(ctx).Where(condition).
		Clauses(clause.Locking{Strength: "UPDATE"}). // 添加 FOR UPDATE 锁
		First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// FindAll 实现带 FOR UPDATE 的单条记录查询
func (b *BaseRepository[T]) FindAll(ctx context.Context, condition map[string]interface{}, withLock bool) ([]T, error) {
	var result []T
	// 加入 FOR UPDATE 锁
	db := b.db.WithContext(ctx).Where(condition)
	if withLock {
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	err := db.Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (b *BaseRepository[T]) FindWithPagination(
	ctx context.Context,
	condition map[string]interface{},
	page, pageSize int,
	orderBy string,
) ([]T, int64, error) {
	var results []T
	var total int64

	// 获取总数
	if err := b.db.WithContext(ctx).Model(new(T)).Where(condition).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，带排序
	offset := (page - 1) * pageSize
	query := b.db.WithContext(ctx).Where(condition).Limit(pageSize).Offset(offset)

	if orderBy != "" {
		query = query.Order(orderBy) // 排序条件
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// IsUniqueConstraintError 判断是否是唯一索引冲突错误
func IsUniqueConstraintError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
}

type InTransaction func(ctx context.Context, tx *gorm.DB) error

func DoInTx(ctx context.Context, db *gorm.DB, fn InTransaction) (err error) {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(log.SQLLogger, "recovery panic: %s", debug.Stack())
			switch x := r.(type) {
			case string:
				err = fmt.Errorf("%s", x)
			case error:
				err = x
			default:
				err = fmt.Errorf("UnKnown panic: %+v", x)
			}
		}
		if err != nil {
			if e := tx.Rollback().Error; e != nil {
				err = fmt.Errorf("tx Rollback err is '%s', caused by other err '%s'", e.Error(), err.Error())
			}
		} else {
			err = tx.Commit().Error
		}
	}()
	return fn(ctx, tx)
}

func toInterfaceSlice(slice any) ([]interface{}, bool) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		return nil, false
	}

	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result, true
}
