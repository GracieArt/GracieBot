package like


import (
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"

  "github.com/enescakir/emoji"
)


type Like struct {
  id string
  bot *core.Bot
  emoji *emoji.Emoji
  info core.ExtensionInfo
}

func (l *Like) ID() string { return l.id }
func (l *Like) Info() core.ExtensionInfo { return l.info }


type Config struct {
  Emoji *emoji.Emoji
}


func New(cnf Config) *Like {
  l := &Like{
    id: "graciebell.art.like",
    emoji: &emoji.YellowHeart,
    info: core.ExtensionInfo {
      Name: "Like",
      Description: "Adds a like button to every message that contains media.",
    },
  }
  if cnf.Emoji != nil { l.emoji = cnf.Emoji }
  return l
}

func (l *Like) Load(b *core.Bot) {
  b.MsgManager.AddHandler(l.OnMessage)
  l.bot = b
}

func (l *Like) OnConnect() {} // satisfy extension interface


// the message handler
func (l *Like) OnMessage(m *discordgo.MessageCreate) bool {
  if !containsMedia(m) || m.Author.Bot { return false }
  l.AddLike(m.Message)
  return true
}


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
