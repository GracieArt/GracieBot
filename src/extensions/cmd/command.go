package cmd


import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"
)


type ArgType int
const (
  ArgType_Int = iota
  ArgType_String
  ArgType_Member
  ArgType_Channel
)


type CmdHandler func(Call) (*discordgo.MessageSend, error)

type Call struct {
  Prefix string
  Msg *discordgo.Message
  Args map[string]interface{}
  Bot *core.Bot
}



type Arg struct {
  Key string
  Type ArgType
}



type Command struct {
  Name, Group, Description string
  Args []Arg
  Run CmdHandler
}
