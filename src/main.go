package main

import (
  "fmt"
  "os"
  "flag"

  "github.com/gracieart/bubblebot"
  "github.com/gracieart/graciebot/src/lib/toys"

  "github.com/joho/godotenv"
)


// Varaibles used for command line parameters
var (
  HideTimestamps bool
  DevMode bool
)


func init() {
  flag.BoolVar(&HideTimestamps, "hidets", false, "hide timestamps")
  flag.BoolVar(&DevMode, "dev", false, "developer mode")
  flag.Parse()
}


func main() {
  // Load environment variables
  if err := godotenv.Load(); err != nil {
    panic(fmt.Errorf("error loading .env file: %w", err))
  }

  // Initialize the bot
	GracieBot, err := bubble.NewBot(bubble.Config{
    Name: "GracieBot",
    Token: os.Getenv("DISCORD_API_TOKEN"),
    Toys: toys.Toys( toys.Config{
      DevMode: DevMode,
      GraciePostKey: os.Getenv("GRACIEPOST_KEY") } ),
    HideTimestamps: HideTimestamps,
  })
	if err != nil { panic(fmt.Errorf("error creating bot: %w", err)) }

  // Start the bot
	GracieBot.Start()
}
