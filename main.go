package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("dev.env"); err != nil {
		log.Fatalln("error loading .env file:", err)
	}

	discordToken := os.Getenv("DISCORD_TOKEN")
	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalln("error loading discord bot:", err)
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	session.AddHandler(onMsgCreate)
	if err := session.Open(); err != nil {
		log.Fatalln("error opening connection:", err)
	}
	defer session.Close()
	log.Println("Bot is now running. Press CTRL-C to exit.")

	shutdownSigs := make(chan os.Signal, 1)
	signal.Notify(shutdownSigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-shutdownSigs
}
