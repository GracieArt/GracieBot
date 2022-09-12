package fun

import (
  "github.com/gracieart/graciebot/src/lib/toys/slash"
  "math/rand"
  "time"
)

func init() {
  rand.Seed(int64(time.Now().Second()))
}

func Commands() []*slash.Command {
  return []*slash.Command {
    EightBall,
    Minesweeper,
    Poll,
  }
}
