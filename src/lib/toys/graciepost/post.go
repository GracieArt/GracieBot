package graciepost


import (
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"

  "github.com/thoas/go-funk"
)


func (g *GraciePost) Post(meta PostMeta) {
  // send the message and save the struct to add a like button to it later
  msg, err := g.bot.Session.ChannelMessageSendComplex(
    meta.Channel, g.createMsg(meta) )

  if err != nil {
    bubble.Log(bubble.Error, g.toyID,
      "Encountered an error when trying to send message: " + err.Error() )
    return
  }

  // add like button if the plugin is connected
  if g.like != nil { g.like.AddLike(msg) }
}


func (g *GraciePost) createMsg(meta PostMeta) *discordgo.MessageSend {
  // just return the link if its set to override embed
  if meta.OverrideEmbed {
    return &discordgo.MessageSend{ Content: meta.PostLink }
  }

  // Truncate the description (if any)
  if len(meta.Desc) > 0 {
    cutoff := 0
    // set cutoff point at the first double linebreak (if any)
    if i := strings.Index(meta.Desc, "\n\n"); i != -1 { cutoff = i }
    // set cutoff point to the CharLimit if it exceeds it
    cutoff = funk.MinInt( []int{ cutoff, g.charLimit } )
    // truncate the string if cutoff was set
    if cutoff > 0 { meta.Desc = meta.Desc[:cutoff] + " ..." }
  }

  // create & send the embed
  msg := &discordgo.MessageSend{
    Embeds: []*discordgo.MessageEmbed{
      &discordgo.MessageEmbed{
        Title: meta.Title,
        Author: &discordgo.MessageEmbedAuthor{
          Name: meta.Artist,
          IconURL: meta.Pfp,
        },
        Description: meta.Desc,
        Fields: []*discordgo.MessageEmbedField{
          &discordgo.MessageEmbedField{
            Name: "Post link:",
            Value: meta.PostLink,
          },
        },
        Image: &discordgo.MessageEmbedImage{
          URL: meta.ImageLink,
        },
        Footer: &discordgo.MessageEmbedFooter{
          Text: "Retrieved from " + meta.SiteName + " using GraciePost",
        },
      },
    },
  }

  //log.Print(msg.Embeds[0])
  return msg
}
