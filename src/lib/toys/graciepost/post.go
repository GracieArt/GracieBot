package graciepost


import (
  "strings"
  "regexp"
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"
)


var collapseWhitespace *regexp.Regexp


func init() {
  collapseWhitespace = regexp.MustCompile(`(\r\s?|\s){2,}`)
}


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


// creates and sends the message to the appropriate channel
func (g *GraciePost) createMsg(meta PostMeta) *discordgo.MessageSend {
  // just return the link if its set to override the embed
  if meta.OverrideEmbed {
    return &discordgo.MessageSend{ Content: meta.PostLink }
  }

  // Truncate the description (if any)
  if len(meta.Desc) > 0 {
    // collapse consecutive whitespace characters
    meta.Desc = collapseWhitespace.ReplaceAllString(meta.Desc, "$1")

    // truncate string if more than 2 newlines
    newlines := strings.Count(meta.Desc, "\n")
    if newlines > 2 {
      meta.Desc = strings.Join(strings.SplitN(meta.Desc, "\n", 2), "\n")
    }

    // truncate the string if more than the char limit
    if len(meta.Desc) > g.charLimit {
      meta.Desc = meta.Desc[:g.charLimit] + " ..."
    }
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
