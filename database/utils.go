package database

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func CreateEnum(db *gorm.DB, name string, enums []string) error {
	for i := range enums {
		enums[i] = "'" + enums[i] + "'"
	}
	enum := strings.Join(enums, ",")

	// check existing value of enum from name
	var existingValues []string
	err := db.Raw(`
		SELECT e.enumlabel 
		FROM pg_type t 
		JOIN pg_enum e ON t.oid = e.enumtypid 
		WHERE t.typname = ?
		ORDER BY e.enumsortorder
	`, name).Scan(&existingValues).Error

	if err == nil && len(existingValues) > 0 {
		if len(existingValues) == len(enums) {
			allSame := true
			for i, val := range existingValues {
				if "'"+val+"'" != enums[i] {
					allSame = false
					break
				}
			}
			if allSame {
				// skip same
				return nil
			}
		}

		// delete to update
		if err := db.Exec(fmt.Sprintf("DROP TYPE IF EXISTS %s CASCADE", name)).Error; err != nil {
			return err
		}
	}

	// create enum
	query := fmt.Sprintf("CREATE TYPE %s AS ENUM (%s)", name, enum)
	return db.Exec(query).Error
}
