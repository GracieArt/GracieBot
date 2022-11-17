package mod

import (
  "fmt"
  "strings"
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/lib/toys/slash"
)

var Prune = slash.NewCommand(slash.CmdConfig{
  Name: "prune",
  Category: "moderation",
  Description: "Bulk delete messages, automatically skipping pinned ones",
  AdminOnly: true,
  GuildOnly: true,
  Options: []*discordgo.ApplicationCommandOption{
    {
      Type: discordgo.ApplicationCommandOptionInteger,
      Name: "limit",
      Description: "The number of messages that will be scanned for deletion",
      Required: true,
      MinValue: new(float64), // no idea why this needs to be a pointer
      MaxValue: 30 },
    {
      Type: discordgo.ApplicationCommandOptionUser,
      Name: "filter-user",
      Description: "Must be from a specific user" },
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "filter-words",
      Description: "Must contain a specific word or phrase (case insensitive)" },
  },

  Handle: func (
    data slash.CmdData,
  ) (
    res *discordgo.InteractionResponse,
    err error,
  ) {
    res = slash.NewInteractionResponse(discordgo.InteractionResponseChannelMessageWithSource)

    // set the user filter, if one was provided
    var filterUser string
    if user, ok := data.Options["filter-user"]; ok {
      filterUser = user.UserValue(data.Bot.Session).ID
    }

    // set the words filter, if one was provided
    var filterWords string
    if words, ok := data.Options["filter-words"]; ok {
      filterWords = words.StringValue()
    }

    // get back list of messages from discord up to the limit provided
    messages, err := data.Bot.Session.ChannelMessages(
      data.ChannelID,
      int(data.Options["limit"].IntValue()),
      "", "", "" )
    if err != nil { return }

    // start checking messages that fit the filter
    var msgsToDelete []string
    for _, m := range messages {
      if m.Pinned { continue }
      if filterUser != "" && filterUser != m.Author.ID { continue }
      if filterWords != "" {
        s := strings.ToLower(m.Content)
        substr := strings.ToLower(filterWords)
        if !strings.Contains(s, substr) { continue }
      }
      msgsToDelete = append(msgsToDelete, m.ID)
    }

    // delete all that fit the filter
    err = data.Bot.Session.ChannelMessagesBulkDelete(data.ChannelID, msgsToDelete)
    if err != nil {
      return
    }

    // set the response
    res.Data.Content = fmt.Sprintf("Pruned %d message(s)", len(msgsToDelete))
    res.Data.Flags = discordgo.MessageFlagsEphemeral

    return
  },
})
