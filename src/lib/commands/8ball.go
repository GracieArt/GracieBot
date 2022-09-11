package commands

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/lib/toys/slash"

  "math/rand"
)

var eightBallAnswers = []string{
  "It is certain.",
  "It is decidedly so.",
  "Without a doubt.",
  "Yes definitely.",
  "You may rely on it.",
  "As I see it,  yes.",
  "Most likely.",
  "Outlook good.",
  "Yes.",
  "Signs point to yes.",

  "Reply hazy, try again.",
  "Ask again later.",
  "Better not tell you now.",
  "Cannot predict now.",
  "Concentrate and ask again.",

  "Don't count on it.",
  "My reply is no.",
  "My sources say no.",
  "Outlook not soo good.",
  "Very doubtful.",
}


var EightBall = slash.NewCommand(slash.CmdConfig{
  Name: "8ball",
  Category: "fun",
  Description: "Ask a yes/no question and the 8ball will answer",
  Options: []*discordgo.ApplicationCommandOption{
    {
      Type: discordgo.ApplicationCommandOptionString,
      Name: "question",
      Description: "Yes/no question",
      Required: true },
  },

  Handle: func (
    data slash.CmdData,
  ) (
    res *discordgo.InteractionResponse,
    err error,
  ) {
    res =  slash.NewInteractionResponse(discordgo.InteractionResponseChannelMessageWithSource)

    answer := eightBallAnswers[rand.Intn(len(eightBallAnswers))]

    res.Data.Content = (
      ":speech_balloon: *" + data.Options["question"].StringValue() + "*\n" +
      ":8ball: ||" + answer + "||" )

    return res, nil
  },
})
