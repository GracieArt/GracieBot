package cmd


import (
  "fmt"
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"
)



type CmdManager struct {
  id string
  bot *core.Bot
  prefix string
  commands map[string]*Command
  cmdsByGroup map[string][]*Command
  info core.ExtensionInfo
}

func (man *CmdManager) ID() string { return man.id }
func (man *CmdManager) Info() core.ExtensionInfo { return man.info }



type Config struct {
  Prefix string
  Commands []*Command
}



func New(cnf Config) (*CmdManager, error) {
  // check if prefix is valid
  if err := validateKeyword(cnf.Prefix); err != nil {
    return nil, fmt.Errorf("invalid prefix %q: %s", err)
  }

  man := &CmdManager{
    id: "graciebell.art.cmd",
    prefix: cnf.Prefix,
    commands: make(map[string]*Command),
    cmdsByGroup: make(map[string][]*Command),
    info: core.ExtensionInfo{
      Name: "Classic Commands",
      Description: "Interact with your bot using a custom prefix.",
    },
  }

  // add the default commands to the list
  cmds := append(
    cnf.Commands,
    man.stdCmd_help(),
    man.stdCmd_commands(),
    man.stdCmd_extensions(),
  )

  // validate and register each command
  for _, c := range cmds {
    err := man.registerCommand(c)
    if err != nil { return nil, err }
  }

  return man, nil
}



func (man *CmdManager) registerCommand(c *Command) error {
  // validate command name
  if err := validateKeyword(c.Name); err != nil {
    return errors.New("command name can only contain letters, digits, and symbols")
  }

  // check for command name conflict
  if _, exists := man.commands[c.Name]; exists {
    return fmt.Errorf("command with name %q already exists", c.Name)
  }

  // validate argument keys
  hasKeyBeenUsed := make(map[string]bool)
  for _, a := range c.Args {
    if err := validateKeyword(a.Key); err != nil {
      return errors.New("argument key can only contain letters, digits, and symbols")
    }
    if hasKeyBeenUsed[a.Key] {
      return fmt.Errorf("multiple arguments with key %q", a.Key)
    }
    hasKeyBeenUsed[a.Key] = true
  }

  // register the command
  man.commands[c.Name] = c
  man.cmdsByGroup[c.Group] = append(man.cmdsByGroup[c.Group], c)

  return nil
}




// register the message handler on load
func (man *CmdManager) Load(b *core.Bot) error {
  b.AddMsgHandler( func (m *discordgo.MessageCreate) bool {
    if !strings.HasPrefix(m.Content, man.prefix) { return false }
    errMsg := man.Eval(m.Message)
    if errMsg != "" {
      man.bot.Session.ChannelMessageSendReply(
  			m.ChannelID, errMsg, m.Reference(),
  		)
    }
    return true
  })
  man.bot = b
  return nil
}

func (man *CmdManager) OnConnect() {} // satisfies extension interface


// parse the message and try to execute the command
func (man *CmdManager) Eval(m *discordgo.Message) (errMsg string) {
  command, argVals, errMsg := man.parseCommand(m)
  if errMsg != "" { return }


	// run the command
	res, err := command.Run(CallData{
    Prefix: man.prefix,
    Msg: m,
    Args: argVals,
    Bot: man.bot,
  })

	// report error if there was one
	if err != nil {
    return fmt.Sprintf(
      "An error occurred while running the command:```%s```\n"+
      "Please show this message to the bot owner so the error can be fixed.",
      err.Error(),
    )
	}

	// send a response if there was one
	if res != nil {
    res.Reference = m.Reference() // mandatory reply
    man.bot.Session.ChannelMessageSendComplex(m.ChannelID, res)
  }

  return ""
}
