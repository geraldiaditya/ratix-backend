package main

import (
	"log"
	"time"

	"github.com/geraldiaditya/ratix-backend/internal/config"
	"github.com/geraldiaditya/ratix-backend/internal/infrastructure"
	movieDomain "github.com/geraldiaditya/ratix-backend/internal/modules/movie/domain"
	movieHandler "github.com/geraldiaditya/ratix-backend/internal/modules/movie/handler"
	movieRepository "github.com/geraldiaditya/ratix-backend/internal/modules/movie/repository"
	movieService "github.com/geraldiaditya/ratix-backend/internal/modules/movie/service"
	userDomain "github.com/geraldiaditya/ratix-backend/internal/modules/user/domain"
	"github.com/geraldiaditya/ratix-backend/internal/modules/user/handler"
	"github.com/geraldiaditya/ratix-backend/internal/modules/user/repository"
	"github.com/geraldiaditya/ratix-backend/internal/modules/user/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Config
	cfg := config.Load()

	// 2. Initialize Infrastructure
	db, err := infrastructure.NewPostgresDB(cfg.Database.DSN)
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v. Continuing without DB for demo purposes.", err)
		// In production, we might want to panic or exit here
	} else {
		// Auto Migrate
		// In a real app, you might want to do this in a separate command or use a migration tool like golang-migrate
		// But GORM AutoMigrate is great for rapid development
		if err := db.AutoMigrate(
			&userDomain.User{},
			&movieDomain.Movie{},
			&movieDomain.Genre{},
			&movieDomain.CastMember{},
			&movieDomain.Showtime{},
		); err != nil {
			log.Printf("Warning: Failed to auto migrate: %v", err)
		}
	}

	// 3. Initialize Modules
	// User Module
	validate := validator.New()
	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userService, validate)

	movieRepo := movieRepository.NewPostgresMovieRepository(db)
	movieService := movieService.NewMovieService(movieRepo)
	movieHandler := movieHandler.NewMovieHandler(movieService)

	// Seeder (Simple check & insert)
	var count int64
	db.Model(&movieDomain.Movie{}).Count(&count)
	if count == 0 {
		log.Println("Seeding dummy movie data...")
		// Genres
		action := movieDomain.Genre{Name: "Action"}
		fantasy := movieDomain.Genre{Name: "Fantasy"}
		db.Create(&action)
		db.Create(&fantasy)

		// Movie 1: The Crimson Blade
		movie1 := movieDomain.Movie{
			Title:       "The Crimson Blade",
			Description: "A legendary warrior awakens to defend his kingdom from an ancient evil.",
			Duration:    135,
			Rating:      8.9,
			PosterURL:   "https://example.com/poster1.jpg",
			ReleaseDate: time.Now(),
			Status:      "now_showing",
			Genres:      []movieDomain.Genre{action, fantasy},
			Cast: []movieDomain.CastMember{
				{Name: "John Smith", Role: "Actor", CharacterName: "Blade", PhotoURL: "https://example.com/pro1.jpg"},
				{Name: "Alan Smithee", Role: "Director"},
			},
			Showtimes: []movieDomain.Showtime{
				{StartTime: time.Now().Add(1 * time.Hour), Price: 15.00},
				{StartTime: time.Now().Add(4 * time.Hour), Price: 15.00},
			},
		}
		db.Create(&movie1)

		// Movie 2: Echoes of Tomorrow
		movie2 := movieDomain.Movie{
			Title:       "Echoes of Tomorrow",
			Description: "A sci-fi thriller about time travel.",
			Duration:    120,
			Rating:      9.1,
			PosterURL:   "https://example.com/poster2.jpg",
			ReleaseDate: time.Now().AddDate(0, 0, 7),
			Status:      "coming_soon",
			Genres:      []movieDomain.Genre{fantasy},
		}
		db.Create(&movie2)
		log.Println("Seeding complete.")
	}

	// 4. Setup Fiber App
	app := fiber.New()
	userHandler.RegisterRoutes(app)
	movieHandler.RegisterRoutes(app)

	// 5. Start Server
	log.Printf("Starting server on port %s", cfg.ServerPort)
	if err := app.Listen(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
