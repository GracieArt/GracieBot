package slash


import (
  "github.com/gracieart/bubblebot"
)



type Slash struct {
  extensionID string
  extensionInfo bubble.ExtensionInfo
  bot *bubble.Bot
  cmdsToRegister []*Command
  commands map[string]*Command
  cmdsByCat map[string][]*Command
}

func (s *Slash) ExtensionID() string { return s.extensionID }
func (s *Slash) ExtensionInfo() bubble.ExtensionInfo { return s.extensionInfo }


func New(cmds ...*Command) (*Slash) {
  return &Slash{
    extensionID: "graciebell.art.slash",
    extensionInfo: bubble.ExtensionInfo{
      Name: "Slash Commands",
      Description: "Add interactive chat commands to your bot in a way " +
        "that's integrated with the Discord UI."},
    cmdsToRegister: cmds,
    commands: make(map[string]*Command),
    cmdsByCat: make(map[string][]*Command),
  }
}


func (s *Slash) Load(b *bubble.Bot) error {
  s.bot = b
  s.bot.Session.AddHandler(s.handleCommand)
  return nil
}



func (s *Slash) OnLifecycleEvent(
  l bubble.LifecycleEvent,
) (
  err error,
) {
  switch l {
  case bubble.LE_Connect:
    err = s.registerCommands()

  case bubble.LE_Close:
    err = s.removeAllCommands()
  }
  return
}
