package commands

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/extensions/cmd"

  "github.com/tmdvs/Go-Emoji-Utils"
)


var Poll = &cmd.Command{
  Name: "poll",
  Group: "fun",
  Description: "Makes a quick reaction poll.",
  Args: []cmd.Arg{
    cmd.Arg{
      Key: "content",
      Type: cmd.ArgType_String,
    },
  },

  Run: func (data cmd.Call) (*discordgo.MessageSend, error) {
    response := &discordgo.MessageSend{}

    if data.Args["content"] == nil {
      response.Content = "You must provide the poll message."
      return response, nil
    }


    content := data.Args["content"].(string)

    // find emojis
    emojis := emoji.FindAll(content)
    switch len(emojis) {
    case 0:
      response.Content = "You must use emojis for the poll choices"
      return response, nil
    case 1:
      response.Content = "You must provide more than one emoji choice for a fair poll"
      return response, nil
    }

    // delete the message that triggered the command
    err := data.Bot.Session.ChannelMessageDelete(
      data.Msg.ChannelID,
      data.Msg.ID,
    )
    if err != nil { return nil, err }

    // create the embed
    embed := &discordgo.MessageEmbed {
      Title: "Poll",
      Footer: &discordgo.MessageEmbedFooter{
        Text: "Sent by " + data.Msg.Author.Username,
      },
      Description: content,
    }

    // send the poll
    pollMsg, err := data.Bot.Session.ChannelMessageSendEmbed(
      data.Msg.ChannelID,
      embed,
    )
    if err != nil { return nil, err }

    // add emoji reactions (the poll choices)
    for _, e := range emojis {
      err = data.Bot.Session.MessageReactionAdd(
        pollMsg.ChannelID,
        pollMsg.ID,
        e.Match.(emoji.Emoji).Value,
      )
      if err != nil { return nil, err }
    }

    // no response necessary (the poll kinda was the response)
    return nil, nil
  },
}
