package commands

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/extensions/slash"

  "github.com/tmdvs/Go-Emoji-Utils"
)


var Poll = slash.NewCommand(slash.CmdConfig{
  Name: "poll",
  Category: "fun",
  Description: "Makes a quick reaction poll",
  Options: []*discordgo.ApplicationCommandOption{
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "content",
      Description: "The content of the poll",
      Required: true },
  },

  Handle: func (
    data slash.CmdData,
  ) (
    res *discordgo.InteractionResponse,
    err error,
  ) {
    res = slash.NewInteractionResponse(discordgo.InteractionResponseChannelMessageWithSource)

    content := data.Options["content"].StringValue()

    // find emojis
    emojis := emoji.FindAll(content)
    switch len(emojis) {
    case 0:
      res.Data.Content = "You must use emojis for the poll choices"
      res.Data.Flags = uint64(discordgo.MessageFlagsEphemeral)
      return
    case 1:
      res.Data.Content = "You must provide more than one emoji choice for a fair poll"
      res.Data.Flags = uint64(discordgo.MessageFlagsEphemeral)
      return
    }

    // Acknowledge the interaction
    res.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
    err = data.Bot.Session.InteractionRespond(data.Interaction, res)
    if err != nil { return }

    // create the embed
    embed := &discordgo.MessageEmbed {
      Title: "Poll",
      Description: content,
    }

    // send the poll
    pollMsg, err := data.Bot.Session.FollowupMessageCreate(
      data.Interaction,
      true,
      &discordgo.WebhookParams{ Embeds: []*discordgo.MessageEmbed{ embed } },
    )
    if err != nil { return }

    // add emoji reactions (the poll choices)
    for _, e := range emojis {
      err = data.Bot.Session.MessageReactionAdd(
        pollMsg.ChannelID,
        pollMsg.ID,
        e.Match.(emoji.Emoji).Value,
      )
      if err != nil { return }
    }

    // no response necessary (the poll kinda was the response)
    res = nil
    return
  },
})
