package toys

import (
  "github.com/gracieart/bubblebot"
  "github.com/gracieart/graciebot/src/lib/toys/slash"
  "github.com/gracieart/graciebot/src/lib/toys/graciepost"
  "github.com/gracieart/graciebot/src/lib/toys/like"
  "github.com/gracieart/graciebot/src/lib/toys/bellhop"

  "github.com/gracieart/graciebot/src/lib/commands/fun"

  "github.com/enescakir/emoji"
)

type Config struct {
  DevMode bool
  GraciePostKey string
}


func Toys(conf Config) []bubble.Toy {
  Bellhop := bellhop.New()

  c := []*slash.Command{}
  c = append(c, fun.Commands()...)
  c = append(c, Bellhop.Commands()...)

  Slash := slash.New( slash.Config{
    Commands: c,
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
    Bellhop,
  }
}
