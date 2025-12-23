package seeds

import (
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	if err := PlantPermissions(db); err != nil {
		log.Fatal("[SEEDS-MIGRATE] failed to migrate permission seeds", err.Error())
	}
}
