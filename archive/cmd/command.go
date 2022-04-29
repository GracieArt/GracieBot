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


type CmdHandler func(CallData) (*discordgo.MessageSend, error)

type ArgVals map[string]any

type CallData struct {
  Prefix string
  Msg *discordgo.Message
  Args ArgVals
  Bot *core.Bot
}



type Arg struct {
  Key string
  Type ArgType
  Required bool
}



type Command struct {
  Name, Group, Description string
  Args []Arg
  Run CmdHandler
}
