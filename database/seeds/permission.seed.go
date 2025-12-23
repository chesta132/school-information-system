package seeds

import (
	"school-information-system/internal/libs/slicelib"
	"school-information-system/internal/models"

	"gorm.io/gorm"
)

var permissionSeeds = []models.Permission{
	{
		Name:        "role full manage",
		Resource:    models.ResourceRole,
		Description: "Full access to manage role of users",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
	{
		Name:        "permission full manage",
		Resource:    models.ResourcePermission,
		Description: "Full access to manage permissions of users with admin role",
		Actions:     []models.PermissionAction{models.ActionCreate, models.ActionRead, models.ActionUpdate, models.ActionDelete},
	},
}

func getNonExistingPermissionSeeds(db *gorm.DB) ([]models.Permission, error) {
	var existingNames []string
	err := db.Model(&models.Permission{}).
		Select("name").
		Where("name IN ?", slicelib.Map(permissionSeeds, func(idx int, perm models.Permission) string { return perm.Name })).
		Pluck("name", &existingNames).Error

	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]bool, len(existingNames))
	for _, name := range existingNames {
		existingMap[name] = true
	}

	return slicelib.Filter(permissionSeeds, func(idx int, perm models.Permission) bool {
		return !existingMap[perm.Name]
	}), nil
}

func PlantPermissions(db *gorm.DB) error {
	perms, err := getNonExistingPermissionSeeds(db)
	if err != nil {
		return err
	}
	if len(perms) == 0 {
		return nil
	}
	return db.Model(&models.Permission{}).Create(&perms).Error
}
