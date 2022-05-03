package toys

import (
  "github.com/gracieart/bubblebot"
  "github.com/gracieart/graciebot/src/lib/toys/slash"
  "github.com/gracieart/graciebot/src/lib/toys/graciepost"
  "github.com/gracieart/graciebot/src/lib/toys/like"
  "github.com/gracieart/graciebot/src/lib/commands"

  "github.com/enescakir/emoji"
)

type Config struct {
  DevMode bool
  GraciePostKey string
}


func Toys(conf Config) []bubble.Toy {
  Slash := slash.New( slash.Config{
    Commands: []*slash.Command{ commands.Poll },
    DevMode: conf.DevMode,
  } )

  GraciePost := graciepost.New(graciepost.Config{
    Key: conf.GraciePostKey,
  })

  Like := like.New(like.Config{
    Emoji: &emoji.YellowHeart,
  })

  return []bubble.Toy{
    Slash,
    GraciePost,
    Like,
  }
}
