package commands

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/lib/toys/slash"

  "strings"
  "strconv"

  emojiUtils "github.com/tmdvs/Go-Emoji-Utils"
)


type Choice struct {
  emoji string
  label string
}


var Poll = slash.NewCommand(slash.CmdConfig{
  Name: "poll",
  Category: "fun",
  Description: "Makes a quick reaction poll",
  Options: []*discordgo.ApplicationCommandOption{
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "title",
      Description: "The title of the poll",
      Required: true },
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "choice1",
      Description: "First choice",
      Required: true },
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "choice2",
      Description: "Second choice",
      Required: true },
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "choice3",
      Description: "Third choice",
      Required: false },
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "choice4",
      Description: "Fourth choice",
      Required: false },
  },

  Handle: func (
    data slash.CmdData,
  ) (
    res *discordgo.InteractionResponse,
    err error,
  ) {
    res = slash.NewInteractionResponse(discordgo.InteractionResponseChannelMessageWithSource)


    //numberofChoices := len(data.Options)-1
    var choices []Choice

    // create the embed
    embed := &discordgo.MessageEmbed {
      Title: data.Options["title"].StringValue(),
    }

    // create list of choices
    i := 0
    for optName := range data.Options {
      if !strings.HasPrefix(optName, "choice") { continue }
      label := data.Options["choice"+strconv.Itoa(i+1)].StringValue()
      // use emoji contained in the input, or use a number emoji if none
      emojis := emojiUtils.FindAll(label)
      emoji := ""
      if len(emojis) > 0 {
        emoji = emojis[0].Match.(emojiUtils.Emoji).Value
      } else {
        switch i {
        case 0:
          emoji = "1️⃣"
        case 1:
          emoji = "2️⃣"
        case 2:
          emoji = "3️⃣"
        case 3:
          emoji = "4️⃣"
        }
        label = emoji + " " + label
      }
      embed.Description = embed.Description + label + "\n"
      choices = append(choices, Choice{
        emoji,
        label,
      })
      i++
    }

    // Acknowledge the interaction
    res.Type = discordgo.InteractionResponseDeferredChannelMessageWithSource
    err = data.Bot.Session.InteractionRespond(data.Interaction, res)
    if err != nil { return }

    // send the poll
    pollMsg, err := data.Bot.Session.FollowupMessageCreate(
      data.Interaction,
      true,
      &discordgo.WebhookParams{ Embeds: []*discordgo.MessageEmbed{ embed } },
    )
    if err != nil { return }

    // add emoji reactions for each poll choice
    for _, choice := range choices {
      err = data.Bot.Session.MessageReactionAdd(
        pollMsg.ChannelID,
        pollMsg.ID,
        choice.emoji,
      )
      if err != nil { return }
    }

    // no response necessary (the poll kinda was the response)
    res = nil
    return
  },
})
