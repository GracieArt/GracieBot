package bellhop

import (
  "github.com/gracieart/bubblebot"
  "github.com/bwmarrin/discordgo"
  _"log"
)


func (b *Bellhop) onJoin(_ *discordgo.Session, g *discordgo.GuildMemberAdd) {

  // get guild options
  if !b.storage.HasGuild(g.GuildID) { return }
  guildOptions, err := b.storage.Guild(g.GuildID)
  if err != nil { panic(err) }

  joinMsgEnabled := guildOptions.Get("join_message_enabled")
  if joinMsgEnabled != nil && joinMsgEnabled.(bool) { b.sendJoinMsg(g.Member, guildOptions) }

  // joinDMEnabled := guildOptions.Get("join_dm_enabled")
  // if joinDMEnabled != nil && joinDMEnabled.(bool) { b.sendJoinDM(g.Member) }
}


func (b *Bellhop) sendJoinMsg(m *discordgo.Member, guildOptions bubble.Entry) {
  ch := guildOptions.Get("join_message_channel")
  if ch == nil { return }

  msg := guildOptions.Get("join_message")
  if msg == nil { return }
  content := msg.(string)

  shouldMention := guildOptions.Get("join_message_should_mention")
  if shouldMention != nil && shouldMention.(bool) {
    content = m.Mention() + " " + content
  }

  b.bot.Session.ChannelMessageSend(ch.(string), content)
}


// func (b *Bellhop) sendJoinDM(m *discordgo.Member) {
//   msg, _ := b.bot.Option(b, m.GuildID, "join_dm")
//   if msg == nil { return }
//
//   ch, _ := b.bot.Session.UserChannelCreate(m.User.ID)
//
//   b.bot.Session.ChannelMessageSend(ch.ID, msg.(string))
// }
