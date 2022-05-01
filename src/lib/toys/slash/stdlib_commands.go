package slash

import (
  "fmt"
  "strings"

  "github.com/bwmarrin/discordgo"

  "github.com/thoas/go-funk"
)


// A big ol command list
func (s *Slash) stdlib_commands() *Command {
  return NewCommand(CmdConfig{
    Name: "commands",
    Category: "utility",
    Description: "List available commands by category.",
    Options: []*discordgo.ApplicationCommandOption{
      {
        Type: discordgo.ApplicationCommandOptionString,
        Name: "category",
        Description: "The category that you'd like to see the commands from.",
        Required: false },
      {
        Type: discordgo.ApplicationCommandOptionInteger,
        Name: "page-number",
        Description: "View different pages of the command list.",
        Required: false },
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

      // Show command categories if category argument is not set
      if _, ok := data.Options["category"]; !ok {
        embed := &discordgo.MessageEmbed{
          Title: "Commands",
          Description: "Use `/commands <category>`" +
            "to view the commands in a specific category.",
        }

        // create the embed fields with group name and the number of commands
        for cat, cmds := range s.cmdsByCat {
          fSuffix := "command"
          if len(cmds) > 1 { fSuffix = fSuffix+"s" }

          embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: strings.Title(cat),
            Value: fmt.Sprintf("%d %s", len(cmds), fSuffix),
            Inline: true,
          })
        }

        res.Data.Embeds = []*discordgo.MessageEmbed{ embed }
        return
      }


      cat := data.Options["category"].StringValue()

      // Show list of commands in category if argument is valid
      if cmds, ok := s.cmdsByCat[cat]; ok {

        // paginate the list of commands
        pageSize := 10
        pages := funk.Chunk(cmds, pageSize).([][]*Command)
        lastPage := len(pages)
        page := 1

        if opt, ok := data.Options["page-number"]; ok {
          inputPageNum := int(opt.IntValue())

          if inputPageNum < 1 {
            res.Data.Content = "`<page-number>` must be positive."
            return
          }
          page = funk.MinInt([]int{ inputPageNum, lastPage })
        }

        cmdsOnPage := pages[page - 1]

        // create the embed
        embed := &discordgo.MessageEmbed{
          Title: strings.Title(cat) + " | Commands",
          Description: "To change pages, use `/commands <page_number>`.\n" +
            "For more info on a particular command, use `/help <command>`.",
          Footer: &discordgo.MessageEmbedFooter{
            Text: fmt.Sprintf("Page %d of %d", page, lastPage) },
        }

        // create the embed fields with command name and description
        for _, cmd := range cmdsOnPage {
          embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: cmd.appCmd.Name,
            Value: cmd.appCmd.Description })
        }

        res.Data.Embeds = []*discordgo.MessageEmbed{ embed }


      // invalid group argument
      } else {
        res.Data.Content = "Unrecognized group." +
          "Use `/commands` for a list of command groups."
      }

      return
    },
  })
}
