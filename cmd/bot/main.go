package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/akionka/akionkabot/internal/cache"
	"github.com/akionka/akionkabot/internal/d2pt"
	"github.com/akionka/akionkabot/internal/postgres"
	"github.com/akionka/akionkabot/internal/s3"
	"github.com/akionka/akionkabot/internal/service"
	"github.com/akionka/akionkabot/internal/stratz"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mymmrac/telego"

	gocache "github.com/patrickmn/go-cache"

	_ "net/http/pprof"

	vault "github.com/hashicorp/vault/api"
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
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = "http://127.0.0.1:8200"

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	vaultClient.SetToken(os.Getenv("VAULT_TOKEN"))

	secret, err := vaultClient.KVv2("akionka-bot").Get(ctx, "config")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	dbURL := secret.Data["DATABASE_URL"].(string)
	botToken := secret.Data["BOT_TOKEN"].(string)
	stratzToken := secret.Data["STRATZ_TOKEN"].(string)

	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	bot, err := telego.NewBot(botToken, telego.WithLogger(logger))
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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

	c := gocache.New(time.Minute, time.Minute*5)

	questionRepo := postgres.NewQuestionRepository(pool, logger.Logger)
	heroRepo := postgres.NewHeroRepository(pool, logger.Logger)
	itemRepo := postgres.NewItemRepository(pool, logger.Logger)
	matchRepo := postgres.NewMatchRepository(pool, logger.Logger)
	steamAccountRepo := postgres.NewSteamAccountRepository(pool, logger.Logger)
	userRepo := cache.NewUserRepository(postgres.NewUserRepository(pool, logger.Logger), c)
	heroImageFetcher := s3.NewHeroImageFetcher(minioClient, c)
	itemImageFetcher := s3.NewItemImageFetcher(minioClient, c)

	d2ptClient := d2pt.NewClient(&http.Client{})
	stratzClient := cache.NewCachedStratzClient(stratz.NewClient(&http.Client{}, stratzToken), c)

	questionService := service.NewQuestionService(questionRepo, matchRepo, steamAccountRepo, d2ptClient, stratzClient, heroRepo, itemRepo, stratzClient, heroImageFetcher)
	matchService := service.NewMatchService(matchRepo, itemImageFetcher)
	userService := service.NewUserService(userRepo)
	playerService := service.NewSteamAccountService(stratzClient)

	collager := NewDefaultCollager(c)

	akionkaBot := NewBot(
		bot,
		logger,
		collager,
		questionService,
		matchService,
		userService,
		playerService,
	)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	logger.Info("starting bot", slog.String("bot_token", botToken))

	akionkaBot.Start(ctx)
}
