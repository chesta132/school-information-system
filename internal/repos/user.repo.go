package repos

import (
	"context"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
	create[models.User]
	read[models.User]
	update[models.User]
	archivable[models.User]
}

func NewUser(db *gorm.DB) *User {
	return &User{db, create[models.User]{db}, read[models.User]{db}, update[models.User]{db}, archivable[models.User]{db}}
}

func (r *User) WithTx(tx *gorm.DB) *User {
	return NewUser(tx)
}

func (r *User) DB() *gorm.DB {
	return r.db
}

// remove this func if not needed later
func (r *User) DeleteUser(ctx context.Context, where any, args ...any) error {
	// delete parents if there is no student related to parents
	return r.db.Transaction(func(tx *gorm.DB) error {
		// user.StudentProfile only have ID and UserID
		user, err := gorm.G[models.User](tx.Unscoped()).
			Preload("StudentProfile", func(db gorm.PreloadBuilder) error {
				db.Select("id, user_id")
				return nil
			}).
			Where(where, args...).
			First(ctx)

		if err != nil {
			return err
		}

		if user.Role != models.RoleStudent || user.StudentProfile == nil {
			_, err := gorm.G[models.User](tx.Unscoped()).Where(where, args...).Delete(ctx)
			return err
		}

		var orphanParentIDs []string
		err = tx.Raw(`
			SELECT sp.parent_id 
			FROM student_parents sp
			WHERE sp.student_profile_id = ?
			AND NOT EXISTS (
				SELECT 1 FROM student_parents sp2 
				WHERE sp2.parent_id = sp.parent_id 
				AND sp2.student_profile_id != ?
			)
		`, user.StudentProfile.ID, user.StudentProfile.ID).Scan(&orphanParentIDs).Error

		if err != nil {
			return err
		}

		if len(orphanParentIDs) > 0 {
			if _, err := gorm.G[models.Parent](tx).Where("id IN ?", orphanParentIDs).Delete(ctx); err != nil {
				return err
			}
		}

		_, err = gorm.G[models.User](tx.Unscoped()).Where(where, args...).Delete(ctx)
		return err
	})
}

func (r *User) DeleteAllUser(ctx context.Context, where any, args ...any) (rowsAffected int, err error) {
	// delete parents if there is no student related to parents
	err = r.db.Transaction(func(tx *gorm.DB) error {
		// each user.StudentProfile only have ID and UserID
		users, err := gorm.G[models.User](tx.Unscoped()).
			Preload("StudentProfile", func(db gorm.PreloadBuilder) error {
				db.Select("id, user_id")
				return nil
			}).
			Where(where, args...).
			Find(ctx)

		if err != nil {
			return err
		}

		if len(users) == 0 {
			return nil
		}

		var studentProfileIDs []string
		for _, user := range users {
			if user.Role == models.RoleStudent && user.StudentProfile != nil {
				studentProfileIDs = append(studentProfileIDs, user.StudentProfile.ID)
			}
		}

		if len(studentProfileIDs) > 0 {
			var orphanParentIDs []string
			err = tx.Raw(`
				SELECT sp.parent_id 
				FROM student_parents sp
				WHERE sp.student_profile_id IN ?
				AND NOT EXISTS (
					SELECT 1 FROM student_parents sp2 
					WHERE sp2.parent_id = sp.parent_id 
					AND sp2.student_profile_id NOT IN ?
				)
				GROUP BY sp.parent_id
			`, studentProfileIDs, studentProfileIDs).Scan(&orphanParentIDs).Error

			if err != nil {
				return err
			}

			if len(orphanParentIDs) > 0 {
				if _, err := gorm.G[models.Parent](tx).Where("id IN ?", orphanParentIDs).Delete(ctx); err != nil {
					return err
				}
			}
		}

		rowsAffected, err = gorm.G[models.User](tx.Unscoped()).Where(where, args...).Delete(ctx)
		return err
	})
	return
}
