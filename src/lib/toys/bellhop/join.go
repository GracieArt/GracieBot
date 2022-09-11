package bellhop

import (
  "github.com/bwmarrin/discordgo"
  _"log"
)


func (b *Bellhop) onJoin(_ *discordgo.Session, g *discordgo.GuildMemberAdd) {
  joinMsgEnabled, _ := b.bot.Option(b, g.Member.GuildID, "join_message_enabled")
  if joinMsgEnabled != nil && joinMsgEnabled.(bool) { b.sendJoinMsg(g.Member) }

  joinDMEnabled, _ := b.bot.Option(b, g.Member.GuildID, "join_dm_enabled")
  if joinDMEnabled != nil && joinDMEnabled.(bool) { b.sendJoinDM(g.Member) }
}


func (b *Bellhop) sendJoinMsg(m *discordgo.Member) {
  ch, _ := b.bot.Option(b, m.GuildID, "join_message_channel")
  if ch == nil { return }

  msg, _ := b.bot.Option(b, m.GuildID, "join_message")
  if msg == nil { return }
  content := msg.(string)

  shouldMention, _ := b.bot.Option(b, m.GuildID, "join_message_should_mention")
  if shouldMention != nil && shouldMention.(bool) {
    content = m.Mention() + " " + content
  }

  b.bot.Session.ChannelMessageSend(ch.(string), content)
}


func (b *Bellhop) sendJoinDM(m *discordgo.Member) {
  msg, _ := b.bot.Option(b, m.GuildID, "join_dm")
  if msg == nil { return }

  ch, _ := b.bot.Session.UserChannelCreate(m.User.ID)

  b.bot.Session.ChannelMessageSend(ch.ID, msg.(string))
}
