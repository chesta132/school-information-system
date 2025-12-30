package seeds

import (
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

const (
	PermRoleName       = "role full manage"
	PermPermissionName = "permission full manage"

	PermAdminName   = "teacher full manage"
	PermTeacherName = "teacher full manage"
	PermStudentName = "student full manage"

	PermSubjectName = "subject full manage"
	PermClassName   = "class full manage"
)

var PermissionSeeds = []models.Permission{
	{
		Name:        PermRoleName,
		Resource:    models.ResourceRole,
		Description: "Full access to manage role of users",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
	{
		Name:        PermPermissionName,
		Resource:    models.ResourcePermission,
		Description: "Full access to manage permissions of users with admin role",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},

	{
		Name:        PermTeacherName,
		Resource:    models.ResourceAdmin,
		Description: "Full access to manage admins",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
	{
		Name:        PermAdminName,
		Resource:    models.ResourceTeacher,
		Description: "Full access to manage teachers",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
	{
		Name:        PermAdminName,
		Resource:    models.ResourceStudent,
		Description: "Full access to manage students",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},

	{
		Name:        PermSubjectName,
		Resource:    models.ResourceSubject,
		Description: "Full access to manage subjects",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
	{
		Name:        PermClassName,
		Resource:    models.ResourceClass,
		Description: "Full access to manage class",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
}

var PermissionSeedIDs = make(map[string]struct{}, len(PermissionSeeds))

func getNonExistingPermissionSeeds(db *gorm.DB) ([]models.Permission, error) {
	var existingNames []string
	err := db.Model(&models.Permission{}).
		Select("name").
		Where("name IN ?", slicelib.Map(PermissionSeeds, func(idx int, perm models.Permission) string { return perm.Name })).
		Pluck("name", &existingNames).Error

	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool, len(existingNames))
	for _, name := range existingNames {
		existingMap[name] = true
	}

	return slicelib.Filter(PermissionSeeds, func(idx int, perm models.Permission) bool {
		return !existingMap[perm.Name]
	}), nil
}

func populateSeedIDs(db *gorm.DB) error {
	err := db.
		Model(&models.Permission{}).
		Where(
			"name IN ?",
			slicelib.Map(PermissionSeeds, func(idx int, perm models.Permission) string { return perm.Name }),
		).
		Find(&PermissionSeeds).Error
	if err != nil {
		return err
	}

	nameToIndex := make(map[string]int, len(PermissionSeeds))
	for i, seed := range PermissionSeeds {
		nameToIndex[seed.Name] = i
	}

	for i, perm := range PermissionSeeds {
		PermissionSeeds[i] = perm
		PermissionSeedIDs[perm.ID] = struct{}{}
	}

	return nil
}

func PlantPermissions(db *gorm.DB) error {
	perms, err := getNonExistingPermissionSeeds(db)
	if err != nil {
		return err
	}
	if len(perms) == 0 {
		return populateSeedIDs(db)
	}

	err = db.Model(&models.Permission{}).Create(&perms).Error
	if err != nil {
		return err
	}

	return populateSeedIDs(db)
}

func IsPermissionSeed(permissions ...*models.Permission) bool {
	for _, perm := range permissions {
		if _, ok := PermissionSeedIDs[perm.ID]; ok {
			return true
		}
	}
	return false
}
