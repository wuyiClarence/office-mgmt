package repository

import (
	"errors"
	"platform-backend/dto"
	"platform-backend/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{BaseRepository: NewBaseRepository[models.User](db)}
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("删除条件为空！")
	}
	user := &models.User{}
	if err := r.db.WithContext(ctx).First(user, id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Delete(user, id).Error
}

func (r *UserRepository) GetChildUserIds(ctx context.Context, userID int64, isSuperAdmin bool) ([]int64, error) {
	var allIDs []int64
	if isSuperAdmin {
		err := r.db.WithContext(ctx).Model(models.User{}).Where("id != ?", userID).Select("id").Find(&allIDs).Error
		if err != nil {
			return nil, err
		}
		return allIDs, nil
	} else {
		if err := r.db.WithContext(ctx).Model(models.User{}).Where("create_user = ?", userID).Select("id").Find(&allIDs).Error; err != nil {
			return nil, err
		}

		result := []int64{}
		for _, id := range allIDs {
			result = append(result, id)
			// 递归查询下级的下级
			subIDs, err := r.GetChildUserIds(ctx, id, false)
			if err != nil {
				return nil, err
			}
			result = append(result, subIDs...)
		}
		return result, nil
	}
}

func (r *UserRepository) GetUserByUserIds(ctx context.Context, userIds []int64, page int, pageSize int) ([]*dto.UserListItem, int64, error) {
	var results []*dto.UserListItem
	var total int64

	db := r.db.WithContext(ctx).Model(models.User{})
	db = db.Joins("left JOIN user b ON user.create_user = b.id").
		Where("user.id IN (?)", userIds).Where("user.status = 0")

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		db = db.Limit(pageSize).Offset(offset)
	}
	db = db.Order("user.created_at desc")

	err := db.Select("user.*,user.id as user_id,b.user_name as create_user_name,b.user_display_name as create_user_display_name").Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
