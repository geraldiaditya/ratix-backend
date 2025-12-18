package main

import (
	"log"
	"time"

	"github.com/geraldiaditya/ratix-backend/internal/config"
	"github.com/geraldiaditya/ratix-backend/internal/infrastructure"
	cinemaDomain "github.com/geraldiaditya/ratix-backend/internal/modules/cinema/domain"
	cinemaHandler "github.com/geraldiaditya/ratix-backend/internal/modules/cinema/handler"
	cinemaRepository "github.com/geraldiaditya/ratix-backend/internal/modules/cinema/repository"
	cinemaService "github.com/geraldiaditya/ratix-backend/internal/modules/cinema/service"
	movieDomain "github.com/geraldiaditya/ratix-backend/internal/modules/movie/domain"
	movieHandler "github.com/geraldiaditya/ratix-backend/internal/modules/movie/handler"
	movieRepository "github.com/geraldiaditya/ratix-backend/internal/modules/movie/repository"
	movieService "github.com/geraldiaditya/ratix-backend/internal/modules/movie/service"
	ticketDomain "github.com/geraldiaditya/ratix-backend/internal/modules/ticket/domain"
	ticketHandler "github.com/geraldiaditya/ratix-backend/internal/modules/ticket/handler"
	ticketRepository "github.com/geraldiaditya/ratix-backend/internal/modules/ticket/repository"
	ticketService "github.com/geraldiaditya/ratix-backend/internal/modules/ticket/service"
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
		// 4. Auto Migrate Phase 1: Independent Tables
		// We migrate Cinema and Theater first because Showtime depends on Cinema.
		// We migrate Genre first because Movie depends on Genre.
		if err := db.AutoMigrate(
			&userDomain.User{},
			&cinemaDomain.Cinema{},
			&cinemaDomain.Theater{},
			&movieDomain.Genre{},
		); err != nil {
			log.Printf("Warning: Failed to auto migrate phase 1: %v", err)
		}

		// 5. Seed Phase 1: Cinemas & Genres
		// Ensure Cinemas exist before Showtimes are migrated (Showtime has FK to Cinema)
		var cinemaCount int64
		db.Model(&cinemaDomain.Cinema{}).Count(&cinemaCount)
		if cinemaCount == 0 {
			log.Println("Seeding dummy cinema data...")
			jakartaCinema := cinemaDomain.Cinema{Name: "Cinema XXI, Grand Indonesia", City: "Jakarta", Address: "Jl. M.H. Thamrin No.1", BasePrice: 50000}
			bandungCinema := cinemaDomain.Cinema{Name: "CGV, Paris Van Java", City: "Bandung", Address: "Jl. Sukajadi No.131-139", BasePrice: 35000}
			db.Create(&jakartaCinema) // ID likely 1
			db.Create(&bandungCinema) // ID likely 2
		}

		var genreCount int64
		db.Model(&movieDomain.Genre{}).Count(&genreCount)
		if genreCount == 0 {
			log.Println("Seeding dummy genre data...")
			action := movieDomain.Genre{Name: "Action"}
			fantasy := movieDomain.Genre{Name: "Fantasy"}
			db.Create(&action)
			db.Create(&fantasy)
		}

		// 6. Auto Migrate Phase 2: Dependent Tables
		// Now safe to migrate Showtime (FK to Cinema) and Movie (FK to Genre, etc)
		if err := db.AutoMigrate(
			&movieDomain.Movie{},
			&movieDomain.CastMember{},
			&movieDomain.Showtime{},
			&ticketDomain.Ticket{},
		); err != nil {
			log.Printf("Warning: Failed to auto migrate phase 2: %v", err)
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

		// Initialize TicketRepo first as CinemaService needs it
		ticketRepo := ticketRepository.NewPostgresTicketRepository(db)
		ticketService := ticketService.NewTicketService(ticketRepo)
		ticketHandler := ticketHandler.NewTicketHandler(ticketService)

		cinemaRepo := cinemaRepository.NewPostgresCinemaRepository(db)
		cinemaService := cinemaService.NewCinemaService(cinemaRepo, ticketRepo)
		cinemaHandler := cinemaHandler.NewCinemaHandler(cinemaService)

		// Seeder Phase 2: Movies & Tickets
		var count int64
		db.Model(&movieDomain.Movie{}).Count(&count)
		if count == 0 {
			log.Println("Seeding dummy movie data...")

			// Fetch Genres (assuming they exist from Phase 1)
			var action, fantasy movieDomain.Genre
			db.Where("name = ?", "Action").First(&action)
			db.Where("name = ?", "Fantasy").First(&fantasy)

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
					{StartTime: time.Now().Add(1 * time.Hour)},
					{StartTime: time.Now().Add(4 * time.Hour)},
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

		// Seed Ticket
		var ticketCount int64
		db.Model(&ticketDomain.Ticket{}).Count(&ticketCount)
		if ticketCount == 0 {
			var m1 movieDomain.Movie
			if err := db.First(&m1, 1).Error; err == nil {
				log.Println("Seeding dummy ticket data...")
				ticket1 := ticketDomain.Ticket{
					UserID:      1,
					MovieID:     m1.ID,
					ShowtimeID:  1,
					BookingCode: "BOOK-12345",
					Seats:       "G14, G15",
					CinemaName:  "Cinema XXI, Mall Grand Indonesia",
					TheaterName: "Studio 1",
					Price:       75000,
					Status:      "active",
					CreatedAt:   time.Now(),
				}
				db.Create(&ticket1)

				ticket2 := ticketDomain.Ticket{
					UserID:      1,
					MovieID:     m1.ID,
					ShowtimeID:  1,
					BookingCode: "BOOK-67890",
					Seats:       "A1",
					CinemaName:  "CGV, Central Park",
					TheaterName: "Velvet Class",
					Price:       150000,
					Status:      "history",
					CreatedAt:   time.Now().AddDate(0, -1, 0),
				}
				db.Create(&ticket2)
				log.Println("Ticket seeding complete.")
			}
		}

		// 4. Setup Fiber App
		app := fiber.New()
		userHandler.RegisterRoutes(app)
		movieHandler.RegisterRoutes(app)
		ticketHandler.RegisterRoutes(app)
		cinemaHandler.RegisterRoutes(app)

		// 5. Start Server
		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	}
}
