package toys

import (
  "os"

  "github.com/gracieart/bubblebot"
  "github.com/gracieart/graciebot/src/lib/toys/slash"
  "github.com/gracieart/graciebot/src/lib/toys/graciepost"
  "github.com/gracieart/graciebot/src/lib/toys/like"
  "github.com/gracieart/graciebot/src/lib/commands"

  "github.com/enescakir/emoji"
)


var Toys []bubble.Toy


func init() {
  Slash := slash.New( commands.Poll )

  GraciePost := graciepost.New(graciepost.Config{
    Key: os.Getenv("GRACIEPOST_KEY"),
  })

  Like := like.New(like.Config{
    Emoji: &emoji.YellowHeart,
  })

  Toys = []bubble.Toy{
    Slash,
    GraciePost,
    Like,
  }
}
