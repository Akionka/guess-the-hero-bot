package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/akionka/akionkabot/d2pt"
	"github.com/akionka/akionkabot/postgres"
	"github.com/akionka/akionkabot/s3"
	"github.com/akionka/akionkabot/service"
	"github.com/patrickmn/go-cache"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mymmrac/telego"

	_ "net/http/pprof"
)

func main() {
	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := telego.NewBot(botToken)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	d2ptClient := d2pt.NewClient(&http.Client{})
	endpoint := "localhost:9000"
	accessKeyID := "akionka"
	secretAccessKey := "password"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	questionRepo := postgres.NewQuestionRepository(pool)
	heroRepo := postgres.NewHeroRepository(pool)
	itemRepo := postgres.NewItemRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	heroImageFetcher := s3.NewHeroImageFetcher(minioClient)
	itemImageFetcher := s3.NewItemImageFetcher(minioClient)

	c := cache.New(time.Minute, time.Minute*5)

	questionService := service.NewQuestionService(questionRepo, d2ptClient, heroRepo, itemRepo, heroImageFetcher, itemImageFetcher)
	userService := service.NewUserService(userRepo, c)

	akionkaBot := NewBot(
		bot,
		NewDefaultCollager(),
		questionService,
		userService,
	)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	akionkaBot.Start(ctx)

}
