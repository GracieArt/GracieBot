package cmd

import (
  "fmt"
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"

  "github.com/thoas/go-funk"
)


// returns default commands
func (man *CmdManager) defaultCommands() []*Command {
  return []*Command{
    man.cmd_help(),
    man.cmd_commands(),
    man.cmd_extensions(),
  }
}





// Gives helpful info on a specific command
func (man *CmdManager) cmd_help() *Command {
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


    Run: func (data Call) (*discordgo.MessageSend, error) {
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





// A big ol command list
func (man *CmdManager) cmd_commands() *Command {
  const cmdName = "commands"

  return &Command{
    Name: cmdName,
    Group: "utility",
    Description: "Displays a list of available commands.",
    Args: []Arg{
      Arg{
        Key: "group",
        Type: ArgType_String,
      },
      Arg{
        Key: "page_number",
        Type: ArgType_Int,
      },
    },


    Run: func (data Call) (*discordgo.MessageSend, error) {
      response := &discordgo.MessageSend{}


      // Show command groups if group argument is empty
      if data.Args["group"] == nil {
        embed := &discordgo.MessageEmbed{
          Title: "Commands",
          Description: fmt.Sprintf(
            "Use `%s` to view the commands in a specific group.",
            data.Prefix + " commands <group>",
          ),
        }

        // create the embed fields with group name and the number of commands in it
        for group, cmds := range man.cmdsByGroup {
          fSuffix := "command"
          if len(cmds) > 1 { fSuffix = fSuffix+"s" }

          embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: strings.Title(group),
            Value: fmt.Sprintf("%d %s", len(cmds), fSuffix),
            Inline: true,
          })
        }

        response.Embeds = []*discordgo.MessageEmbed{ embed }
        return response, nil
      }


      inputGroup := data.Args["group"].(string)

      // Show list of commands in group if group argument is valid
      if cmds, ok := man.cmdsByGroup[inputGroup]; ok {

        // paginate the list of commands
        pageSize := 10
        pages := funk.Chunk(cmds, pageSize).([][]*Command)
        lastPage := len(pages)
        page := 1

        if data.Args["page_number"] != nil {
          inputPageNum := data.Args["page_number"].(int)

          if inputPageNum < 1 {
            response.Content = "`<page_number>` must be a positive."
            return response, nil
          }
          page = funk.MinInt([]int{ inputPageNum, lastPage })
        }

        cmdsOnPage := pages[page - 1]

        // create the embed
        embed := &discordgo.MessageEmbed{
          Title: strings.Title(inputGroup) + " | Commands",
          Description: fmt.Sprintf(
            "To change pages, use `%s %s %s <page_number>`.\n"+
            "For more info on a particular command, use `%s help <command>`.",
            data.Prefix, cmdName, inputGroup, data.Prefix,
          ),
          Footer: &discordgo.MessageEmbedFooter{
            Text: fmt.Sprintf("Page %d of %d", page, lastPage),
          },
        }

        // create the embed fields with command name and description
        for _, cmd := range cmdsOnPage {
          embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: cmd.Name,
            Value: cmd.Description,
          })
        }

        response.Embeds = []*discordgo.MessageEmbed{ embed }


      // invalid group argument
      } else {
        response.Content = fmt.Sprintf(
          "Unrecognized group. Use `%s %s` for a list of command groups.",
          data.Prefix, cmdName,
        )
      }

      return response, nil
    },
  }
}







// Sends a list of all the registered extensions
func (man *CmdManager) cmd_extensions() *Command {
  const thisCmdName = "extensions"

  return &Command{
    Name: thisCmdName,
    Group: "utility",
    Description: "Lists all registered extensions.",
    // need to add arg for page number, but will do that later when pagination is seperated

    Run: func (data Call) (*discordgo.MessageSend, error) {
      response := &discordgo.MessageSend{
        Embeds: []*discordgo.MessageEmbed{
          &discordgo.MessageEmbed{
            Title: "Extensions",
            Description: "The following extensions are registered with this bot:",
            Fields: funk.Map(
              man.bot.ExtManager.Info(),
              func (i core.ExtensionInfo) *discordgo.MessageEmbedField {
                return &discordgo.MessageEmbedField{ i.Name, i.Description, false }
              },
            ).([]*discordgo.MessageEmbedField),
          },
        },
      }

      return response, nil
    },
  }
}
