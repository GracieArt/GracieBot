package slash

import (
  "fmt"

  "github.com/gracieart/bubblebot"

  "github.com/bwmarrin/discordgo"
)



func (s *Slash) handleCommand(
  sesh *discordgo.Session,
  i *discordgo.InteractionCreate,
) {
  appCmdData := i.ApplicationCommandData()

  if c, ok := s.commands[appCmdData.Name]; ok {
    data := CmdData {
      Bot: s.bot,
      GuildID: i.GuildID,
      ChannelID: i.ChannelID,
      Interaction: i.Interaction,
    }

    // format option values into string map
    data.Options = make(
      map[string]*discordgo.ApplicationCommandInteractionDataOption,
      len(appCmdData.Options) )
    for _, opt := range appCmdData.Options {
      data.Options[opt.Name] = opt
    }

    if i.Member == nil {
      data.DM = true
      data.Invoker = CmdInvoker{ i.User, nil }
    } else {
      data.DM = false
      data.Invoker = CmdInvoker{ i.Member.User, i.Member }
    }

    res, err := c.Handle(data)

    if err != nil {
      err := sesh.InteractionRespond(i.Interaction, ErrorResponse(
        "Encountered an error while running the command.", err.Error()))
      if err != nil {
        bubble.Log(bubble.Error, s.toyID, fmt.Sprint(
          "Error responding to interaction: ", err))
      }
    }

    if res != nil {
      err := sesh.InteractionRespond(i.Interaction, res)
      if err != nil {
        bubble.Log(bubble.Error, s.toyID, fmt.Sprint(
          "Error responding to interaction: ", err))
      }
    }
  }
}


func ErrorResponse(errMsg, details string) *discordgo.InteractionResponse {
  embed := &discordgo.MessageEmbed{
    Color: 15483205,
    Title: "Error",
    Description: errMsg,
    Footer: &discordgo.MessageEmbedFooter{
      Text: "Show this to the bot's owner if you believe there is a problem."},
  }

  if details != "" {
    embed.Description = embed.Description + "```" + details + "```"
  }

  return &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Flags: uint64(discordgo.MessageFlagsEphemeral),
      Embeds: []*discordgo.MessageEmbed{ embed },
    },
  }
}
