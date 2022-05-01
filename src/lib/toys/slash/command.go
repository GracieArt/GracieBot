package slash


import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"
)


type Command struct {
  appCmd *discordgo.ApplicationCommand
  category string
  Handle CmdHandler
}


func NewCommand(conf CmdConfig) *Command {
  cmd := &Command {
    appCmd: &discordgo.ApplicationCommand{
      Name: conf.Name,
      Description: conf.Description,
      Options: conf.Options,
    },
    category: conf.Category,
    Handle: conf.Handle,
  }
  return cmd
}

func (c Command) Category() string { return c.category }


type CmdConfig struct {
  Name, Description, Category string
  Options []*discordgo.ApplicationCommandOption
  Handle CmdHandler
}


type CmdData struct {
  Bot *bubble.Bot
  DM bool
  GuildID, ChannelID string
  Invoker CmdInvoker
  Options map[string]*discordgo.ApplicationCommandInteractionDataOption
  Interaction *discordgo.Interaction
}


type CmdInvoker struct {
  *discordgo.User
  Member *discordgo.Member
}


type CmdHandler func(CmdData) (*discordgo.InteractionResponse, error)


func NewInteractionResponse(
  t discordgo.InteractionResponseType,
) (
  *discordgo.InteractionResponse,
) {
  return &discordgo.InteractionResponse{
    Type: t,
    Data: &discordgo.InteractionResponseData{},
  }
}
