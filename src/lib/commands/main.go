package commands

import (
  "math/rand"
  "time"
)

func init() {
  rand.Seed(int64(time.Now().Second()))
}
