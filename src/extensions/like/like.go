package like


import (
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"

  "github.com/enescakir/emoji"
)


type Like struct {
  extensionID string
  extensionInfo core.ExtensionInfo
  bot *core.Bot
  emoji *emoji.Emoji
}

func (l *Like) ExtensionID() string { return l.extensionID }
func (l *Like) ExtensionInfo() core.ExtensionInfo { return l.extensionInfo }


type Config struct {
  Emoji *emoji.Emoji
}


func New(cnf Config) *Like {
  l := &Like{
    extensionID: "graciebell.art.like",
    emoji: &emoji.YellowHeart,
    extensionInfo: core.ExtensionInfo {
      Name: "Like",
      Description: "Adds a like button to every message that contains media.",
    },
  }
  if cnf.Emoji != nil { l.emoji = cnf.Emoji }
  return l
}

func (l *Like) Load(b *core.Bot) error {
  b.AddMsgHandler(func (m *discordgo.MessageCreate) bool {
    if !containsMedia(m) || m.Author.Bot { return false }
    l.AddLike(m.Message)
    return true
  })
  l.bot = b
  return nil
}


// satisfy extension interface
func (l *Like) OnLifecycleEvent(core.LifecycleEvent) error { return nil }


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
