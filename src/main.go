package main

import (
  "fmt"
  "flag"
  "os"

  "github.com/gracieart/graciebot-core"
  "github.com/gracieart/graciebot/src/extensions/cmd"
  "github.com/gracieart/graciebot/src/extensions/like"
  "github.com/gracieart/graciebot/src/extensions/graciepost"
  "github.com/gracieart/graciebot/src/setup/commands"

  "github.com/joho/godotenv"
  "github.com/enescakir/emoji"
)


// Variables used for command line parameters
var (
	Token string
)


func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}



func main() {
  if err := godotenv.Load(); err != nil {
    panic(fmt.Errorf("error loading .env file: %s", err))
  }


  ext_cmd, err := cmd.New(cmd.Config{
    Prefix: "g!",
    Commands: []*cmd.Command{
      commands.Poll,
    },
  })
  if err != nil { panic(fmt.Errorf("error creating Classic Commands extension: %s", err)) }

  ext_like := like.New(like.Config{
    Emoji: &emoji.YellowHeart,
  })

  ext_graciepost := graciepost.New(graciepost.Config{
    LikeExtension: ext_like,
  })


	GracieBot, err := core.NewBot(core.Config{
    Token: os.Getenv("DISCORD_API_TOKEN"),
    Extensions: []core.Extension{
      ext_cmd,
      ext_graciepost,
      ext_like,
    },
  })


	if err != nil {
		fmt.Printf("error creating bot %s", err)
		return
	}

	GracieBot.Connect()
}
