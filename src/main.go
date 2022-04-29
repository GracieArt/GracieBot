package main

import (
  "fmt"
  "os"

  "github.com/gracieart/graciebot-core"
  "github.com/gracieart/graciebot/src/extensions/slash"
  "github.com/gracieart/graciebot/src/extensions/like"
  "github.com/gracieart/graciebot/src/extensions/graciepost"
  "github.com/gracieart/graciebot/src/setup/commands"

  "github.com/joho/godotenv"
  "github.com/enescakir/emoji"
)


func main() {
  if err := godotenv.Load(); err != nil {
    panic(fmt.Errorf("error loading .env file: %w", err))
  }


  Slash := slash.New( commands.Poll )

  Like := like.New(like.Config{
    Emoji: &emoji.YellowHeart,
  })

  GraciePost := graciepost.New(graciepost.Config{
    Key: os.Getenv("GRACIEPOST_KEY"),
  })


	GracieBot, err := core.NewBot(core.Config{
    Token: os.Getenv("DISCORD_API_TOKEN"),
    Extensions: []core.Extension{
      Slash,
      GraciePost,
      Like,
    },
  })
	if err != nil { panic(fmt.Errorf("error creating bot: %w", err)) }

	GracieBot.Connect()
}
