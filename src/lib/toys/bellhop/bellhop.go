package bellhop

import (
  "github.com/gracieart/bubblebot"
  _"github.com/bwmarrin/discordgo"
  "github.com/ostafen/clover"
)


type Bellhop struct {
  toyID string
  toyInfo bubble.ToyInfo
  bot *bubble.Bot

  storage *clover.DB
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



func (b *Bellhop) Load(bot *bubble.Bot) error {
  b.bot = bot

  bot.Session.AddHandler(b.onJoin)

  return nil
}



func (b *Bellhop) OnLifecycleEvent(l bubble.LifecycleEvent) {}
