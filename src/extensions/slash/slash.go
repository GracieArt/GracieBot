package slash


import (
  "github.com/gracieart/graciebot-core"
)



type Slash struct {
  extensionID string
  extensionInfo core.ExtensionInfo
  bot *core.Bot
  cmdsToRegister []*Command
  commands map[string]*Command
  cmdsByCat map[string][]*Command
}

func (s *Slash) ExtensionID() string { return s.extensionID }
func (s *Slash) ExtensionInfo() core.ExtensionInfo { return s.extensionInfo }


func New(cmds ...*Command) (*Slash) {
  return &Slash{
    extensionID: "graciebell.art.slash",
    extensionInfo: core.ExtensionInfo{
      Name: "Slash Commands",
      Description: "Add interactive chat commands to your bot in a way " +
        "that's integrated with the Discord UI."},
    cmdsToRegister: cmds,
    commands: make(map[string]*Command),
    cmdsByCat: make(map[string][]*Command),
  }
}


func (s *Slash) Load(b *core.Bot) error {
  s.bot = b
  s.bot.Session.AddHandler(s.handleCommand)
  return nil
}



func (s *Slash) OnLifecycleEvent(
  l core.LifecycleEvent,
) (
  err error,
) {
  switch l {
  case core.LE_Connect:
    err = s.registerCommands()

  case core.LE_Close:
    err = s.removeAllCommands()
  }
  return
}
