package slash

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"

  "github.com/thoas/go-funk"
)


// Sends a list of all the registered extensions
func (s *Slash) stdlib_toys() *Command {
  return NewCommand(CmdConfig{
    Name: "toys",
    Category: "information",
    Description: "List all registered toys.",
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
        Title: "Toys",
        Description: "The following toys are registered with this bot:",
        Fields: funk.Map(
          s.bot.Toys(),
          func (e bubble.Toy) *discordgo.MessageEmbedField {
            info := e.ToyInfo()
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
