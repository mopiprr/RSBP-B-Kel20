package cmd

import (
	"log"
	"os"

	"github.com/mci-its/backend-service/data-layer/migrations"
	"gorm.io/gorm"
)

func Commands(db *gorm.DB) {
	migrate := false
	seed := false
	fresh := false

	for _, arg := range os.Args[1:] {
		if arg == "--migrate" {
			migrate = true
		}

		if arg == "--fresh" {
			fresh = true
		}

		if arg == "--seed" {
			seed = true
		}

	}

	if migrate {
		if err := migrations.Migrate(db); err != nil {
			log.Fatalf("error migration: %v", err)
		}
		log.Println("migration completed successfully")
	}

	if fresh {
		if err := migrations.Fresh(db); err != nil {
			log.Fatalf("Failed to perform fresh migration: %v", err)
		}
	}

	if seed {
		if err := migrations.Seeder(db); err != nil {
			log.Fatalf("error migration seeder: %v", err)
		}
		log.Println("seeder completed successfully")
	}

}
