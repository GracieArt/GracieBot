package like


import (
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"

  "github.com/enescakir/emoji"
)


type Like struct {
  extensionID string
  extensionInfo bubble.ExtensionInfo
  bot *bubble.Bot
  emoji *emoji.Emoji
}

func (l *Like) ExtensionID() string { return l.extensionID }
func (l *Like) ExtensionInfo() bubble.ExtensionInfo { return l.extensionInfo }


type Config struct {
  Emoji *emoji.Emoji
}


func New(cnf Config) *Like {
  l := &Like{
    extensionID: "graciebell.art.like",
    emoji: &emoji.YellowHeart,
    extensionInfo: bubble.ExtensionInfo {
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


// satisfy extension interface
func (l *Like) OnLifecycleEvent(bubble.LifecycleEvent) error { return nil }


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
