package main

import (
  "fmt"
  "os"

  "github.com/gracieart/graciebot-core"
  "github.com/gracieart/graciebot/src/extensions/cmd"
  "github.com/gracieart/graciebot/src/extensions/like"
  "github.com/gracieart/graciebot/src/extensions/graciepost"
  "github.com/gracieart/graciebot/src/setup/commands"

  "github.com/joho/godotenv"
  "github.com/enescakir/emoji"
)


func main() {
  if err := godotenv.Load(); err != nil {
    panic(fmt.Errorf("error loading .env file: %s", err))
  }


  ClassicCommands, err := cmd.New(cmd.Config{
    Prefix: "g!",
    Commands: []*cmd.Command{
      commands.Poll,
    },
  })
  if err != nil {
    panic(fmt.Errorf("error creating Classic Commands extension: [%w]", err))
  }

  Like := like.New(like.Config{
    Emoji: &emoji.YellowHeart,
  })

  GraciePost := graciepost.New(graciepost.Config{
    Key: os.Getenv("GRACIEPOST_KEY"),
  })


	GracieBot, err := core.NewBot(core.Config{
    Token: os.Getenv("DISCORD_API_TOKEN"),
    Extensions: []core.Extension{
      ClassicCommands,
      GraciePost,
      Like,
    },
  })
	if err != nil { panic(fmt.Errorf("error creating bot: [%w]", err)) }

	GracieBot.Connect()
}
