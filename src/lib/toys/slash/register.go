package slash

import (
  "fmt"

  "github.com/gracieart/bubblebot"
)


func (s *Slash) registerCommands() {
  bubble.Log(bubble.Info, s.toyID, "Registering commands")

  stdlib := []*Command{
    s.stdlib_help(),
    s.stdlib_commands(),
    s.stdlib_toys(),
    //s.stdlib_about(),
  }

  s.cmdsToRegister = append(s.cmdsToRegister, stdlib...)

  registeredAll := s.RegisterCommands(s.cmdsToRegister...)
  if !registeredAll {
    bubble.Log(bubble.Warning, s.toyID,
      "One or more commands could not be registered")
  }
}


func (s *Slash) removeAllCommands() {
  bubble.Log(bubble.Info, s.toyID, "Removing commands")

  for _, c := range s.commands {
    err := s.bot.Session.ApplicationCommandDelete(
      s.bot.UserID(), s.guild, c.appCmd.ID)
    if err != nil {
      bubble.Log(bubble.Warning, s.toyID,
        fmt.Sprintf("Failed to delete %q command: %v", c.appCmd.Name, err) )
    }
  }
  return
}


func (s *Slash) RegisterCommands(commands ...*Command) (registeredAll bool) {
  registeredAll = true

  for _, c := range commands {
    if _, exists := s.commands[c.appCmd.Name]; exists {
      bubble.Log(bubble.Warning, s.toyID,
        fmt.Sprintf("Tried to register duplicate command %q", c.appCmd.Name) )
      registeredAll = false
      continue
    }

    cmd, err := s.bot.Session.ApplicationCommandCreate(
      s.bot.UserID(), s.guild, c.appCmd)
    if err != nil {
      bubble.Log(bubble.Warning, s.toyID,
        fmt.Sprintf("Failed to create %q command: %v", c.appCmd.Name, err) )
      registeredAll = false
      continue
    }

    c.appCmd = cmd
    s.commands[c.appCmd.Name] = c
    s.cmdsByCat[c.category] = append(s.cmdsByCat[c.category], c)
  }

  return
}
