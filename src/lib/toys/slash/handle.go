package slash

import (
  "fmt"

  "github.com/gracieart/bubblebot"

  "github.com/bwmarrin/discordgo"
)

const (
  subcommand = discordgo.ApplicationCommandOptionSubCommand
  subcommandGroup = discordgo.ApplicationCommandOptionSubCommandGroup
)



// takes the data from an interaction, formats and sends it to the appropriate
// command handler, and ensures the interaction gets a response
func (s *Slash) handleCommand(
  sesh *discordgo.Session,
  i *discordgo.InteractionCreate,
) {

  appCmdData := i.ApplicationCommandData()
  c, exists := s.commands[appCmdData.Name]

  if !exists {
    // Idk a situation where this could actually happen, so if it does happen
    // I need to know about it kuz that means something's very wrong lmao
    res := errorResponse("Unknown command")
    _ = s.bot.Session.InteractionRespond(i.Interaction, res)
    return
  }


  // create the struct of data formatted for the command handler
  data := CmdData {
    Bot: s.bot,
    GuildID: i.GuildID,
    ChannelID: i.ChannelID,
    Interaction: i.Interaction,
  }


  // format the option data, if any
  if len(appCmdData.Options) > 0 {
    var options []*discordgo.ApplicationCommandInteractionDataOption
    t := appCmdData.Options[0].Type
    isSubcommand := (t == subcommand || t == subcommandGroup)

    if !isSubcommand {
      // set options
      options = appCmdData.Options

    } else {
      // set subcommand names
      data.SubcommandName = appCmdData.Options[0].Name
      if t == subcommandGroup {
        data.NestedSubcommandName = appCmdData.Options[0].Options[0].Name
      }

      // set subcommand options
      switch t {
      case subcommand:
        options = appCmdData.Options[0].Options
      case subcommandGroup:
        options = appCmdData.Options[0].Options[0].Options
      }
    }

    // format options into string map
    data.Options = make(
      map[string]*discordgo.ApplicationCommandInteractionDataOption,
      len(options),
    )
    for _, opt := range options { data.Options[opt.Name] = opt }
  }


  // format data about who sent the command so that a member object can
  // be included if it was sent from a guild
  if i.Member == nil {
    data.DM = true
    data.Invoker = CmdInvoker{ i.User, nil }
  } else {
    data.DM = false
    data.Invoker = CmdInvoker{ i.Member.User, i.Member }
  }


  // run the handler function and get back an InteractionResponse
  res, err := c.Handle(data)

  // respond to the interaction
  var response *discordgo.InteractionResponse
  if err != nil {
    response = errorResponse(err.Error())
  } else if res != nil {
    response = res
  } else {
    response = errorResponse("Got nil pointer for handler response")
  }

  err = s.bot.Session.InteractionRespond(i.Interaction, response)
  if err != nil {
    bubble.Log(
      bubble.Warning,
      s.toyID,
      fmt.Sprint("Failed to respond to interaction: ", err),
    )
  }
}




// returns an InteractionResponse to be used when an error occured while
// executing the command's handler funcction
func errorResponse(details string) *discordgo.InteractionResponse {
  embed := &discordgo.MessageEmbed{
    Color: 15483205,
    Title: "Error",
    Description: "Encountered an error while running the command.",
    Footer: &discordgo.MessageEmbedFooter{
      Text: "Show this to the bot's owner if you believe there is a problem."},
  }

  if details != "" {
    embed.Description = embed.Description + "```" + details + "```"
  }

  return &discordgo.InteractionResponse{
    Type: discordgo.InteractionResponseChannelMessageWithSource,
    Data: &discordgo.InteractionResponseData{
      Flags: discordgo.MessageFlagsEphemeral,
      Embeds: []*discordgo.MessageEmbed{ embed },
    },
  }
}
