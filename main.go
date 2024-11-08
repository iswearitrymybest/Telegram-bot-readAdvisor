package main

import (
	//"context"
	"context"
	"flag"
	"log"

	tgClient "TG_bot/clients/telegram"
	event_consumer "TG_bot/consumer/event-consumer"
	"TG_bot/events/telegram"
	"TG_bot/storage/sqlite"
)

const (
	//storagePath = "data/storage/files"
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {

	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"insert token for access to TG bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token not inserted")
	}
	return *token
}
