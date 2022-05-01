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

    // formats option values into string map
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
      s.respond(i.Interaction, errorResponse(
        "Encountered an error while running the command.", err.Error()))
    }

    if res != nil {
      s.respond(i.Interaction, res)
    }
  }
}


func errorResponse(errMsg, details string) *discordgo.InteractionResponse {
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


func (s *Slash) respond(
  i *discordgo.Interaction,
  res *discordgo.InteractionResponse,
) {
  err := s.bot.Session.InteractionRespond(i, res)
  if err == nil { return }
  bubble.Log(bubble.Warning, s.toyID,
    fmt.Sprint("Failed to respond to interaction: ", err) )
}
