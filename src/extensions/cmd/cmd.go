package cmd


import (
  "fmt"
  "strings"
  "errors"

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
  cmds := append(cnf.Commands, man.defaultCommands()...)

  // validate and register each command
  for _, c := range cmds {

    // validate command name
    if err := validateKeyword(c.Name); err != nil {
      return nil, errors.New("command name can only contain letters, digits, and symbols")
    }

    // check for command name conflict
    if _, exists := man.commands[c.Name]; exists {
      return nil, fmt.Errorf("command with name %q already exists", c.Name)
    }

    // validate argument keys
    for _, a := range c.Args {
      if err := validateKeyword(a.Key); err != nil {
        return nil, errors.New("argument key can only contain letters, digits, and symbols")
      }
    }

    // register the command
    man.commands[c.Name] = c
    man.cmdsByGroup[c.Group] = append(man.cmdsByGroup[c.Group], c)
  }

  return man, nil
}




// register the message handler on load
func (man *CmdManager) Load(b *core.Bot) {
  b.MsgManager.AddHandler( func (m *discordgo.MessageCreate) bool {
    if !strings.HasPrefix(m.Content, man.prefix) { return false }
    man.Eval(m.Message)
    return true
  })
  man.bot = b
}

func (man *CmdManager) OnConnect() {} // satisfies extension interface



// parse the message and try to execute the command
func (man *CmdManager) Eval(m *discordgo.Message) {

  // separate command name from arguments
	cmdName, argsStr := parseCmd(m.Content, man.prefix)

  // check if command exists
	c, exists := man.commands[cmdName];
	if !exists {
		man.bot.Session.ChannelMessageSendReply(
			m.ChannelID,
      fmt.Sprintf(
        "Unrecognized command.\nFor a list of commands, use `%s commands`.",
        man.prefix,
      ),
			m.Reference(),
		)
		return
	}


  callData := Call{
    Prefix: man.prefix,
    Msg: m,
    Args: make(map[string]interface{}),
    Bot: man.bot,
  }


	// parse and type check arguments
  if len(c.Args) > 0 {
    inputs := splitArgs(argsStr, c)

    for i, input := range inputs {
      arg := c.Args[i]
      parsedVal, err := parseArgType(input, arg.Type)

      if err != nil {
        man.bot.Session.ChannelMessageSendReply(
          m.ChannelID,
          fmt.Sprintf("`<%s>` %s", arg.Key, err),
          m.Reference(),
        )
        return
      }

      callData.Args[c.Args[i].Key] = parsedVal
    }
  }


	// run the command
	res, err := c.Run(callData)

	// report error if there was one
	if err != nil {
		man.bot.Session.ChannelMessageSendReply(
			m.ChannelID,
			fmt.Sprintf(
        "An error occurred while running the command:```%s```\n"+
        "Please show this message to the bot owner so the error can be fixed.",
        err.Error(),
      ),
      m.Reference(),
		)
	}

	// send a response if there was one
	if res != nil {
    res.Reference = m.Reference() // mandatory reply
    man.bot.Session.ChannelMessageSendComplex(m.ChannelID, res)
  }
}
