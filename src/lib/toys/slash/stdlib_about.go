package slash

import (
  "github.com/bwmarrin/discordgo"
  _"github.com/gracieart/bubblebot"
)


// Sends a list of all the registered extensions
func (s *Slash) stdlib_about() *Command {
  return NewCommand(CmdConfig{
    Name: "about",
    Category: "information",
    Description: "List general information about this bot.",
    Handle: func (
      data CmdData,
    ) (
      res *discordgo.InteractionResponse,
      err error,
    ) {
      res = NewInteractionResponse(
        discordgo.InteractionResponseChannelMessageWithSource,
      )

      embed := &discordgo.MessageEmbed{
        Title: s.bot.Name(),
        Description: "Running [BubbleBot](https://"+s.bubbleBotPath+") " +
          "version `"+s.bubbleBotVersion+"`.",
        // Fields: []*discordgo.MessageEmbedField{
        //   {
        //     Name: "Owner:",
        //     Value: s.bot.Owner,
        //   },
        // },
      }

      res.Data.Embeds = []*discordgo.MessageEmbed{ embed }
      return
    },
  })
}
