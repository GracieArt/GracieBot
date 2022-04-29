package slash

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"

  "github.com/thoas/go-funk"
)


// Sends a list of all the registered extensions
func (s *Slash) stdlib_extensions() *Command {
  return NewCommand(CmdConfig{
    Name: "extensions",
    Category: "information",
    Description: "List all registered extensions.",
    // need to add arg for page number, but will do that later when pagination is seperated

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
        Title: "Extensions",
        Description: "The following extensions are registered with this bot:",
        Fields: funk.Map(
          s.bot.Extensions(),
          func (e core.Extension) *discordgo.MessageEmbedField {
            info := e.ExtensionInfo()
            return &discordgo.MessageEmbedField{
              info.Name, info.Description, false,
            }
          },
        ).([]*discordgo.MessageEmbedField),
      }

      res.Data.Embeds = []*discordgo.MessageEmbed{ embed }
      return
    },
  })
}
