package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	"github.com/aidilbaihaqi/prabodrive-be/internal/config"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/infrastructure/database"
	"github.com/aidilbaihaqi/prabodrive-be/internal/repository"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
)

type seedUser struct {
	name     string
	email    string
	password string
	role     string
}

var seeds = []seedUser{
	{
		name:     "Super Admin",
		email:    "admin@prabodrive.com",
		password: "Admin@prabodrive123",
		role:     constants.RoleAdmin,
	},
	{
		name:     "Sample User",
		email:    "user@prabodrive.com",
		password: "User@prabodrive123",
		role:     constants.RoleUser,
	},
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using system environment")
	}

	cfg := config.Load()
	ctx := context.Background()

	db, err := database.NewPool(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	fmt.Println("Seeding users...")
	fmt.Println()

	for _, s := range seeds {
		exists, err := userRepo.ExistsByEmail(ctx, s.email)
		if err != nil {
			log.Fatalf("check email %s: %v", s.email, err)
		}
		if exists {
			fmt.Printf("  ⏭  %s (%s) — already exists, skipped\n", s.email, s.role)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(s.password), 12)
		if err != nil {
			log.Fatalf("hash password for %s: %v", s.email, err)
		}

		now := time.Now()
		user := &domain.User{
			ID:           uuid.New().String(),
			Email:        s.email,
			PasswordHash: string(hash),
			Name:         s.name,
			Role:         s.role,
			QuotaUsed:    0,
			QuotaMax:     constants.DefaultQuotaMax,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("create user %s: %v", s.email, err)
		}

		fmt.Printf("  ✓  %s (%s) — created\n", s.email, s.role)
	}

	fmt.Println()
	fmt.Println("Done. Credentials:")
	fmt.Println()
	for _, s := range seeds {
		fmt.Printf("  Role     : %s\n", s.role)
		fmt.Printf("  Email    : %s\n", s.email)
		fmt.Printf("  Password : %s\n", s.password)
		fmt.Println()
	}
}
