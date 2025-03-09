package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/akionka/akionkabot/internal/d2pt"
	"github.com/akionka/akionkabot/internal/postgres"
	"github.com/akionka/akionkabot/internal/s3"
	"github.com/akionka/akionkabot/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mymmrac/telego"
	"github.com/patrickmn/go-cache"

	_ "net/http/pprof"
)

func main() {
	baseHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "bot_token" {
				return slog.String("bot_token", "")
			}
			return a
		},
	})
	customHandler := &TelegoContextHandler{Handler: baseHandler}

	logger := &TelegoLogger{
		Logger:    slog.New(customHandler),
		LogErrors: true,
		LogDebug:  false,
	}

	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	botToken := os.Getenv("BOT_TOKEN")

	bot, err := telego.NewBot(botToken, telego.WithLogger(logger))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	d2ptClient := d2pt.NewClient(&http.Client{})

	endpoint := "localhost:9000"
	accessKeyID := "akionka"
	secretAccessKey := "password"

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
	})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	c := cache.New(time.Minute, time.Minute*5)

	questionRepo := postgres.NewQuestionRepository(pool, logger.Logger)
	heroRepo := postgres.NewHeroRepository(pool, logger.Logger)
	itemRepo := postgres.NewItemRepository(pool, logger.Logger)
	userRepo := postgres.NewUserRepository(pool, logger.Logger)
	heroImageFetcher := s3.NewHeroImageFetcher(minioClient, c)
	itemImageFetcher := s3.NewItemImageFetcher(minioClient, c)

	questionService := service.NewQuestionService(questionRepo, d2ptClient, heroRepo, itemRepo, heroImageFetcher, itemImageFetcher)
	userService := service.NewUserService(userRepo, c)

	collager := NewDefaultCollager(c)

	akionkaBot := NewBot(
		bot,
		logger,
		collager,
		questionService,
		userService,
	)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	logger.Info("starting bot", slog.String("bot_token", botToken))

	akionkaBot.Start(ctx)

}
