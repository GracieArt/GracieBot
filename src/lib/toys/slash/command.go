package slash


import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/bubblebot"
)

var adminPermission int64 = discordgo.PermissionAdministrator

type Command struct {
  appCmd *discordgo.ApplicationCommand
  category string
  Handle CmdHandler
  //allowedInDMs *bool
}


func NewCommand(conf CmdConfig) *Command {
  cmd := &Command {
    appCmd: &discordgo.ApplicationCommand{
      Name: conf.Name,
      Description: conf.Description,
      Options: conf.Options,
      DMPermission: new(bool),
    },
    category: conf.Category,
    Handle: conf.Handle,
  }

  *cmd.appCmd.DMPermission = !conf.GuildOnly

  if conf.AdminOnly == true {
    cmd.appCmd.DefaultMemberPermissions = &adminPermission
  }

  return cmd
}

func (c Command) Category() string { return c.category }


type CmdConfig struct {
  Name, Description, Category string
  Options []*discordgo.ApplicationCommandOption
  Handle CmdHandler
  AdminOnly bool
  GuildOnly bool
}


type CmdData struct {
  Bot *bubble.Bot
  DM bool
  GuildID, ChannelID string
  Invoker CmdInvoker
  Options map[string]*discordgo.ApplicationCommandInteractionDataOption
  Interaction *discordgo.Interaction
  SubcommandName, NestedSubcommandName string
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
