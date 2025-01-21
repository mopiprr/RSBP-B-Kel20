package migrations

import (
	"github.com/mci-its/backend-service/data-layer/migrations/seeds"
	"gorm.io/gorm"
)

func Seeder(db *gorm.DB) error {
	if err := seeds.SeedAll(db); err != nil {
		return err
	}

	return nil
}
