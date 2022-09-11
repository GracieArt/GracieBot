package commands

import (
  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot/src/lib/toys/slash"

  "math/rand"
  "strconv"
  "strings"
)


var (
  gameWidth = 9
  gameHeight = 9
  mines = 9
)

var Minesweeper = slash.NewCommand(slash.CmdConfig{
  Name: "minesweeper",
  Category: "fun",
  Description: "Generates a minesweeper puzzle",

  Handle: func (
    data slash.CmdData,
  ) (
    res *discordgo.InteractionResponse,
    err error,
  ) {
    res =  slash.NewInteractionResponse(discordgo.InteractionResponseChannelMessageWithSource)

    // create and zero the array
    minefield := make([][]int, gameHeight)
    for i := range minefield {
      minefield[i] = make([]int, gameWidth)
    }

    // place the mines and make the clues
    for i := 0; i < mines; i++ {
      var mineX, mineY int

      // loop until u find an empty space to put a mine and place it
      for {
        mineX = rand.Intn(gameWidth)
        mineY = rand.Intn(gameHeight)
        if minefield[mineY][mineX] != 9 {
          minefield[mineY][mineX] = 9
          break
        }
      }

      // update the clues around the mine
      for y := mineY-1; y <= mineY+1; y++ {
        if y < 0 || y >= gameHeight { continue }
        for x := mineX-1; x <= mineX+1; x++ {
          if x < 0 || x >= gameWidth { continue }
          if minefield[y][x] == 9 { continue }
          minefield[y][x]++
        }
      }
    }


    // pick a space that has no clues to be already uncovered
    var freeSpaceX, freeSpaceY int
    foundFreeSpace := false
    for !foundFreeSpace {
      freeSpaceX = rand.Intn(gameWidth-2)+1
      freeSpaceY = rand.Intn(gameHeight-2)+1
      if minefield[freeSpaceY][freeSpaceX] == 0 { foundFreeSpace = true }
    }


    // fake minesweeper HUD
    res.Data.Content = (
      strconv.Itoa(mines/10) + "\uFE0F" + "\u20E3" +
      strconv.Itoa(mines % 10) + "\uFE0F" + "\u20E3" +
      strings.Repeat(":white_large_square:", 2) +
      ":open_mouth:" +
      strings.Repeat(":white_large_square:", 2) +
      strings.Repeat(":zero:", 2) + "\n\n" )



    // translate the minefield values into the string that'll be posted
    for y := 0; y < gameHeight; y++ {
      row := ""

      for x := 0; x < gameWidth; x++ {
        val := minefield[y][x]
        symbol := ""
        if val == 0 {
          symbol = ":white_large_square:"
        } else if val > 0 && val < 9 {
          symbol = strconv.Itoa(val) + "\uFE0F" + "\u20E3"
        } else {
          symbol = ":boom:"
        }

        inFreeSpaceY := y >= freeSpaceY-1 && y <= freeSpaceY+1
        inFreeSpaceX := x >= freeSpaceX-1 && x <= freeSpaceX+1
        if inFreeSpaceY && inFreeSpaceX {
          row += symbol
        } else {
          row += "||" + symbol + "||"
        }
      }

      res.Data.Content += row + "\n"
    }

    return res, nil
  },
})
