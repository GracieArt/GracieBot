package slash

import (
  "fmt"

  "github.com/gracieart/bubblebot"
)


func (s *Slash) registerCommands() error {
  bubble.Log(bubble.Info, s.toyID, "Registering commands")

  err := s.RegisterCommands(
    s.stdlib_help(),
    s.stdlib_commands(),
    s.stdlib_toys())
  if err != nil { return err }

  return s.RegisterCommands(s.cmdsToRegister...)
}


func (s *Slash) removeAllCommands() error {
  bubble.Log(bubble.Info, s.toyID, "Removing commands")

  for _, c := range s.commands {
    err := s.bot.Session.ApplicationCommandDelete(s.bot.UserID(), "", c.appCmd.ID)
    if err != nil { return fmt.Errorf(
      "couldn't delete %q command: %w", c.appCmd.Name, err)}
  }
  return nil
}


func (s *Slash) RegisterCommands(commands ...*Command) error {
  for _, c := range commands {
    if _, exists := s.commands[c.appCmd.Name]; exists {
      return fmt.Errorf("tried to register duplicate command %q", c.appCmd.Name)
    }

    cmd, err := s.bot.Session.ApplicationCommandCreate(
      s.bot.UserID(), "", c.appCmd)
    if err != nil { return fmt.Errorf(
      "couldn't create %q command: %w", c.appCmd.Name, err)}

    c.appCmd = cmd
    s.commands[c.appCmd.Name] = c
    s.cmdsByCat[c.category] = append(s.cmdsByCat[c.category], c)
  }
  return nil
}
