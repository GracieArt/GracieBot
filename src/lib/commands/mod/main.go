package mod

import (
  "github.com/gracieart/graciebot/src/lib/toys/slash"
)

func Commands() []*slash.Command {
  return []*slash.Command {
    Prune,
  }
}
