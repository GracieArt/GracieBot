package slash

import (
  "fmt"
  "strings"

  "github.com/bwmarrin/discordgo"
)


// Gives helpful info on a specific command
func (s *Slash) stdlib_help() *Command {
  return NewCommand(CmdConfig{
    Name: "help",
    Category: "utility",
    Description: "Get info on a specific command.",
    Options: []*discordgo.ApplicationCommandOption{
      {
        Type: discordgo.ApplicationCommandOptionString,
        Name: "command",
        Description: "The command you need help with.",
        Required: true },
    },


    Handle: func (
      data CmdData,
    ) (
      res *discordgo.InteractionResponse,
      err error,
    ) {
      res = NewInteractionResponse(
        discordgo.InteractionResponseChannelMessageWithSource,
      )

      // get lowercase command name argument
      cmdName := strings.ToLower(data.Options["command"].StringValue())

      // check that command exists
      cmd, exists := s.commands[cmdName];
      if !exists {
        res.Data.Content = fmt.Sprintf(
          "`%s` is not a recognized command.\nFor a list of commands, use `/commands`.",
          cmdName,
        )
        return
      }

      // Show info on command
      // generate the line that illustrates the form of the command
      syntaxElems := []string{ cmd.appCmd.Name }
      for _, opt := range cmd.appCmd.Options {
        syntaxElems = append(syntaxElems, "<" + opt.Name + ">")
      }
      syntax := strings.Join(syntaxElems, " ")

      // create the embed
      embed := &discordgo.MessageEmbed {
        Title: strings.Title(cmd.appCmd.Name) + " | Help",
        Description: fmt.Sprintf(
          "%s\n\n**Syntax**```%s```",
          cmd.appCmd.Description, syntax,
        ),
      }

      res.Data.Embeds = []*discordgo.MessageEmbed{ embed }
      return
    },
  })
}
