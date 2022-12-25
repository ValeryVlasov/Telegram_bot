package main

import (
	"github.com/ValeryVlasov/Telegram_bot/pkg/config"
	"github.com/ValeryVlasov/Telegram_bot/pkg/repository"
	"github.com/ValeryVlasov/Telegram_bot/pkg/repository/boltdb"
	"github.com/ValeryVlasov/Telegram_bot/pkg/server"
	"github.com/ValeryVlasov/Telegram_bot/pkg/telegram"
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	bot, err := tgbotapi.NewBotAPI("MyTelegramBotAPI")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	pocketClient, err := pocket.NewClient("PocketClientKey")
	if err != nil {
		log.Fatal(err)
	}

	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tokenRepository := boltdb.NewTokenRepository(db)

	telegramBot := telegram.NewBot(bot, pocketClient, tokenRepository, "http://localhost/")

	authorizationServer := server.NewAuthorizationServer(pocketClient, tokenRepository, "https://t.me/my_first_pocket_bot")

	go func() {
		if err := telegramBot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	if err = authorizationServer.Start(); err != nil {
		log.Fatal(err)
	}
}

func initDB() (*bolt.DB, error) {
	db, err := bolt.Open("bot.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
