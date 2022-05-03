package slash


import (
  "runtime/debug"

  "github.com/gracieart/bubblebot"
)



type Slash struct {
  toyID string
  toyInfo bubble.ToyInfo
  bot *bubble.Bot
  cmdsToRegister []*Command
  commands map[string]*Command
  cmdsByCat map[string][]*Command

  devMode bool
  guild string

  // used for the "about" command
  bubbleBotVersion string
  bubbleBotPath string
}

func (s *Slash) ToyID() string { return s.toyID }
func (s *Slash) ToyInfo() bubble.ToyInfo { return s.toyInfo }


type Config struct {
  Commands []*Command
  DevMode bool
}

func New(conf Config) (*Slash) {
  s := &Slash{
    toyID: "graciebell.art.slash",
    toyInfo: bubble.ToyInfo{
      Name: "Slash Commands",
      Description: "Add interactive chat commands to your bot in a way " +
        "that's integrated with the Discord UI."},
    cmdsToRegister: conf.Commands,
    commands: make(map[string]*Command),
    cmdsByCat: make(map[string][]*Command),
    devMode: conf.DevMode,
  }

  if s.devMode { s.guild = "652282433911128074" }

  // get the version of the bubblebot package, for the "about" command
  buildInfo, _ := debug.ReadBuildInfo()
  for _, d := range buildInfo.Deps {
    if d.Path != "github.com/gracieart/bubblebot" { continue }
    s.bubbleBotVersion = d.Version
    s.bubbleBotPath = d.Path
    break
  }

  return s
}


func (s *Slash) Load(b *bubble.Bot) error {
  //s.bubbleBotVersion = debug
  b.Session.AddHandler(s.handleCommand)
  s.bot = b
  return nil
}



func (s *Slash) OnLifecycleEvent(l bubble.LifecycleEvent) {
  switch l {
  case bubble.Connect:
    s.registerCommands()

  case bubble.Close:
    s.removeAllCommands()
  }
  return
}
