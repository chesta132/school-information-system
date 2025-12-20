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

	query := fmt.Sprintf(`
		DO $$
		BEGIN
				IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN
						CREATE TYPE %s AS ENUM (%s);
				END IF;
		END
		$$;
	`, name, name, enum)

	return db.Exec(query).Error
}
