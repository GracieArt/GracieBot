package bellhop

import (
  "github.com/gracieart/bubblebot"
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/lib/toys/slash"
  "log"
  "fmt"
)

var _ = log.Print // debugging

func (b *Bellhop) Commands() []*slash.Command {
  return []*slash.Command{
    b.Cmd_welcome(),
  }
}

func (b *Bellhop) Cmd_welcome() *slash.Command {
  return slash.NewCommand(slash.CmdConfig{
    Name: "welcome",
    Category: "automation",
    Description: "Welcome message for new members",
    AdminOnly: true,
    GuildOnly: true,
    Options: []*discordgo.ApplicationCommandOption{
      {
        Type: discordgo.ApplicationCommandOptionSubCommand,
        Name: "create",
        Description: "Set up a welcome message for new members",
        Options: []*discordgo.ApplicationCommandOption{
          {
            Type: discordgo.ApplicationCommandOptionChannel,
            Name: "channel",
            Description: "The channel that the welcome message will be sent in",
            Required: true,
            ChannelTypes: []discordgo.ChannelType{ discordgo.ChannelTypeGuildText } },
          {
            Type: discordgo.ApplicationCommandOptionBoolean,
            Name: "should-mention",
            Description: "Whether the message should ping the new member or not",
            Required: true },
          {
            Type: discordgo.ApplicationCommandOptionString,
            Name: "message",
            Description: "The message that will be sent",
            Required: true },
        },
      },

      {
        Type: discordgo.ApplicationCommandOptionSubCommand,
        Name: "test",
        Description: "Test the welcome message as if you were joining the server",
      },

      {
        Type: discordgo.ApplicationCommandOptionSubCommandGroup,
        Name: "options",
        Description: "Change welcome message related options",
        Options: []*discordgo.ApplicationCommandOption{
          {
            Type: discordgo.ApplicationCommandOptionSubCommand,
            Name: "enabled",
            Description: "Enable/disable the welcome message",
            Options: []*discordgo.ApplicationCommandOption{
              {
                Type: discordgo.ApplicationCommandOptionBoolean,
                Name: "value",
                Description: "Enable/disable the welcome message",
                Required: true },
            },
          },
          {
            Type: discordgo.ApplicationCommandOptionSubCommand,
            Name: "should-mention",
            Description: "Pings the new member in the welcome message",
            Options: []*discordgo.ApplicationCommandOption{
              {
                Type: discordgo.ApplicationCommandOptionBoolean,
                Name: "value",
                Description: "Pings the new member in the welcome message",
                Required: true },
            },
          },
          {
            Type: discordgo.ApplicationCommandOptionSubCommand,
            Name: "channel",
            Description: "The channel that the welcome message will be sent in",
            Options: []*discordgo.ApplicationCommandOption{
              {
                Type: discordgo.ApplicationCommandOptionChannel,
                Name: "value",
                Description: "The channel that the welcome message will be sent in",
                ChannelTypes: []discordgo.ChannelType{ discordgo.ChannelTypeGuildText },
                Required: true },
            },
          },
          {
            Type: discordgo.ApplicationCommandOptionSubCommand,
            Name: "message",
            Description: "The message that will be sent",
            Options: []*discordgo.ApplicationCommandOption{
              {
                Type: discordgo.ApplicationCommandOptionString,
                Name: "value",
                Description: "The message that will be sent",
                Required: true },
            },
          },
        },
      },
    },

    Handle: func (
      data slash.CmdData,
    ) (
      res *discordgo.InteractionResponse,
      err error,
    ) {
      res = slash.NewInteractionResponse(discordgo.InteractionResponseChannelMessageWithSource)

      // get guild options, or make them if they dont exist yet
      var guildOptions bubble.Entry
      if !b.storage.HasGuild(data.GuildID) {
        guild, err := b.bot.Session.Guild(data.GuildID)
        if err != nil { return res, err }
        guildOptions, err = b.storage.InsertOne(guild)
        if err != nil { return res, err }
      } else {
        guildOptions, err = b.storage.Guild(data.GuildID)
        if err != nil { return res, err }
      }

      switch data.SubcommandName {
      case "create":
        joinMsg := data.Options["message"].StringValue()
        ch := data.Options["channel"].ChannelValue(b.bot.Session)
        shouldMention := data.Options["should-mention"].BoolValue()

        guildOptions.Set("join_message_enabled", true)
        guildOptions.Set("join_message", joinMsg)
        guildOptions.Set("join_message_channel", ch.ID)
        guildOptions.Set("join_message_should_mention", shouldMention)

        b.storage.Save(guildOptions)

        if shouldMention { joinMsg = "\\@User " + joinMsg }

        res.Data.Content = fmt.Sprint(
          "Welcome message set up! ",
          "Test it out with `/welcome test`",
        )


      case "test":
        enabled := guildOptions.Get("join_message_enabled")
        if enabled == nil || enabled.(bool) == false {
          res.Data.Flags = discordgo.MessageFlagsEphemeral
          res.Data.Content = fmt.Sprint(
            "You need to enable welcome messages with ",
            "`/welcome options enabled true`",
          )
          break
        }

        ch := guildOptions.Get("join_message_channel")
        if ch == nil {
          res.Data.Flags = discordgo.MessageFlagsEphemeral
          res.Data.Content = fmt.Sprint(
            "Looks like you haven't set up a welcome message yet. ",
            "Do that with `/welcome create`",
          )
          break
        }

        data.Invoker.Member.GuildID = data.GuildID
        b.sendJoinMsg(data.Invoker.Member, guildOptions)

        res.Data.Content = "Welcome message sent"


      case "options":
        switch data.NestedSubcommandName {
        case "enabled":
          val := data.Options["value"].BoolValue()
          guildOptions.Set("join_message_enabled", val)
          if val {
            res.Data.Content = "Enabled welcome message"
          } else {
            res.Data.Content = "Disabled welcome message"
          }

        case "should-mention":
          val := data.Options["value"].BoolValue()
          guildOptions.Set("join_message_should_mention", val)
          if val {
            res.Data.Content = "Welcome message will ping new users"
          } else {
            res.Data.Content = "Welcome message won't ping new users"
          }

        case "channel":
          ch := data.Options["value"].ChannelValue(b.bot.Session)
          guildOptions.Set("channel", ch.ID)
          res.Data.Content = "Welcome messages will be sent in " + ch.Mention()

        case "message":
          val := data.Options["value"].StringValue()
          guildOptions.Set("join_message", val)
          res.Data.Content = "Welcome message set to:\n> " + val
        }

        b.storage.Save(guildOptions)
      }


      return res, nil
    },
  })
}
