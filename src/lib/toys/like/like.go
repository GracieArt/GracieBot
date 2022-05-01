package like


import (
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"

  "github.com/enescakir/emoji"
)


type Like struct {
  toyID string
  toyInfo bubble.ToyInfo
  bot *bubble.Bot
  emoji *emoji.Emoji
}

func (l *Like) ToyID() string { return l.toyID }
func (l *Like) ToyInfo() bubble.ToyInfo { return l.toyInfo }


type Config struct {
  Emoji *emoji.Emoji
}


func New(cnf Config) *Like {
  l := &Like{
    toyID: "graciebell.art.like",
    emoji: &emoji.YellowHeart,
    toyInfo: bubble.ToyInfo {
      Name: "Like",
      Description: "Adds a like button to every message that contains media.",
    },
  }
  if cnf.Emoji != nil { l.emoji = cnf.Emoji }
  return l
}

func (l *Like) Load(b *bubble.Bot) error {
  b.AddMsgHandler(func (m *discordgo.MessageCreate) bool {
    if !containsMedia(m) || m.Author.Bot { return false }
    l.AddLike(m.Message)
    return true
  })
  l.bot = b
  return nil
}


// satisfy toy interface
func (l *Like) OnLifecycleEvent(bubble.LifecycleEvent) {}


func (l *Like) AddLike(m *discordgo.Message) {
  l.bot.Session.MessageReactionAdd(m.ChannelID, m.ID, l.emoji.String())
}


func containsMedia(m *discordgo.MessageCreate) bool {
  hasLink := false
  for _, match := range []string{ "https://", "http://" } {
    if !strings.Contains(m.Content, match) { continue }
    hasLink = true
    break
  }

  hasAttachment := len(m.Attachments) > 0

  hasEmbed := len(m.Embeds) > 0

  return hasLink || hasAttachment || hasEmbed
}
