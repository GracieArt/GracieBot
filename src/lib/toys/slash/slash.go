package slash


import (
  "github.com/gracieart/bubblebot"
)



type Slash struct {
  toyID string
  toyInfo bubble.ToyInfo
  bot *bubble.Bot
  cmdsToRegister []*Command
  commands map[string]*Command
  cmdsByCat map[string][]*Command
}

func (s *Slash) ToyID() string { return s.toyID }
func (s *Slash) ToyInfo() bubble.ToyInfo { return s.toyInfo }


func New(cmds ...*Command) (*Slash) {
  return &Slash{
    toyID: "graciebell.art.slash",
    toyInfo: bubble.ToyInfo{
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
  case bubble.Connect:
    err = s.registerCommands()

  case bubble.Close:
    err = s.removeAllCommands()
  }
  return
}
