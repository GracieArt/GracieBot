package bellhop

import (
  "github.com/gracieart/bubblebot"
  _"github.com/bwmarrin/discordgo"
)


type Bellhop struct {
  toyID string
  toyInfo bubble.ToyInfo
  bot *bubble.Bot

  storage *bubble.Storage
}

func (b *Bellhop) ToyID() string { return b.toyID }
func (b *Bellhop) ToyInfo() bubble.ToyInfo { return b.toyInfo }


func New() *Bellhop {
  b := &Bellhop{
    toyID : "graciebell.art.bellhop",
    toyInfo: bubble.ToyInfo{
      Name: "Bellhop",
      Description: "Welcomes new members.",
    },
  }
  return b
}



func (b *Bellhop) Load(bot *bubble.Bot, s *bubble.Storage) error {
  b.bot = bot
  b.storage = s

  bot.Session.AddHandler(b.onJoin)

  return nil
}



func (b *Bellhop) OnLifecycleEvent(l bubble.LifecycleEvent) {}
