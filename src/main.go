package main

import (
  "fmt"
  "os"

  "github.com/gracieart/bubblebot"
  "github.com/gracieart/graciebot/src/lib/toys"

  "github.com/joho/godotenv"
)


func main() {
  // Load environment variables
  if err := godotenv.Load(); err != nil {
    panic(fmt.Errorf("error loading .env file: %w", err))
  }

  // Initialize the bot
	GracieBot, err := bubble.NewBot(bubble.Config{
    Name: "GracieBot",
    Token: os.Getenv("DISCORD_API_TOKEN"),
    Toys: toys.Toys,
  })
	if err != nil { panic(fmt.Errorf("error creating bot: %w", err)) }

  // Start the bot
	GracieBot.Start()
}
