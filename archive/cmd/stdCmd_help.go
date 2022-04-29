package cmd

import (
  "fmt"
  "strings"

  "github.com/bwmarrin/discordgo"
)


// Gives helpful info on a specific command
func (man *CmdManager) stdCmd_help() *Command {
  const thisCmdName = "help"

  return &Command{
    Name: thisCmdName,
    Group: "utility",
    Description: "Gives info on a specific command.",
    Args: []Arg{
      Arg{
        Key: "command",
        Type: ArgType_String,
      },
    },


    Run: func (data CallData) (*discordgo.MessageSend, error) {
      response := &discordgo.MessageSend{}

      // no name provided
      if data.Args["command"] == nil {
        response.Content = fmt.Sprintf(
          "You must specify the command that you need help with.\n"+
          "For a list of commands, use `%s commands`.",
          data.Prefix,
        )
        return response, nil
      }

      // get lowercase command name argument
      cmdName := strings.ToLower(data.Args["command"].(string))

      // check that command exists
      cmd, exists := man.commands[cmdName];
      if !exists {
        response.Content = fmt.Sprintf(
          "`%s` is not a recognized command.\nFor a list of commands, use `%s commands`.",
          cmdName, data.Prefix,
        )
        return response, nil
      }

      // Show info on command
      // generate the line that illustrates the form of the command
      syntaxElems := []string{ data.Prefix, cmd.Name }
      for _, arg := range cmd.Args {
        syntaxElems = append(syntaxElems, "<" + arg.Key + ">")
      }
      syntax := strings.Join(syntaxElems, " ")

      // create the embed
      embed := &discordgo.MessageEmbed {
        Title: strings.Title(cmd.Name) + " | Help",
        Description: fmt.Sprintf(
          "%s\n\n**Syntax**```%s```",
          cmd.Description, syntax,
        ),
      }

      response.Embeds = []*discordgo.MessageEmbed{ embed }
      return response, nil
    },
  }
}
