package cmd

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"

  "github.com/thoas/go-funk"
)


// Sends a list of all the registered extensions
func (man *CmdManager) stdCmd_extensions() *Command {
  const thisCmdName = "extensions"

  return &Command{
    Name: thisCmdName,
    Group: "information",
    Description: "Lists all registered extensions.",
    // need to add arg for page number, but will do that later when pagination is seperated

    Run: func (data CallData) (*discordgo.MessageSend, error) {
      response := &discordgo.MessageSend{
        Embeds: []*discordgo.MessageEmbed{
          &discordgo.MessageEmbed{
            Title: "Extensions",
            Description: "The following extensions are registered with this bot:",
            Fields: funk.Map(
              man.bot.Extensions(),
              func (e core.Extension) *discordgo.MessageEmbedField {
                info := e.Info()
                return &discordgo.MessageEmbedField{
                  info.Name, info.Description, false,
                }
              },
            ).([]*discordgo.MessageEmbedField),
          },
        },
      }

      return response, nil
    },
  }
}
