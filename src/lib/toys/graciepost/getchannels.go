package graciepost


import (
  "encoding/json"

  "github.com/bwmarrin/discordgo"

  "github.com/thoas/go-funk"
)


type Menu struct {
  Title string `json:"title"`
  ID string `json:"id"`
  ParentID string `json:"parentId"`
}

type MenuLevel struct {
  Name string `json:"name"`
  Items []Menu `json:"items"`
}

// get the object with all the menus to send back to the GraciePost toy
func (g *GraciePost) GetChannels() []byte {
  // get guilds. (have to do it this way kuz the guilds in State.Ready
  // are not populated with the channels)
  tempGuilds := g.bot.Session.State.Ready.Guilds
  guilds := funk.Map(tempGuilds,
    func(tempGuild *discordgo.Guild) *discordgo.Guild {
      guild, _ := g.bot.Session.Guild(tempGuild.ID)
      return guild
    },
  ).([]*discordgo.Guild)

  // we create this map to filter out the empty categories later on
  categorySizes := make(map[string]int)

  // get all channels
  textChannels := []*discordgo.Channel{}
  categories := []*discordgo.Channel{}
  for _, guild := range guilds {
    guildChannels, _ := g.bot.Session.GuildChannels(guild.ID)

    // separate text and category channels
    for _, ch := range guildChannels {
      switch ch.Type {
      case discordgo.ChannelTypeGuildText:
        textChannels = append(textChannels, ch)
        if ch.ParentID != "" { categorySizes[ch.ParentID]++ }
      case discordgo.ChannelTypeGuildCategory:
        categories = append(categories, ch)
      }
    }
  }

  // filter out categories with no text channels
  categories = funk.Filter(categories, func(cat *discordgo.Channel) bool {
    if categorySizes[cat.ID] > 0 { return true } else { return false }
  }).([]*discordgo.Channel)

  // create the struct to be sent as json
  menus := []MenuLevel{
    MenuLevel{
      Name: "guilds",
      Items: funk.Map(guilds, func(guild *discordgo.Guild) Menu {
        return Menu{
          Title: guild.Name,
          ID: guild.ID,
        }
      }).([]Menu),
    },

    MenuLevel{
      Name: "categories",
      Items: funk.Map(categories, func(cat *discordgo.Channel) Menu {
        return Menu{
          Title: cat.Name,
          ID: cat.ID,
          ParentID: cat.GuildID,
        }
      }).([]Menu),
    },

    MenuLevel{
      Name: "channels",
      Items: funk.Map(textChannels, func(ch *discordgo.Channel) Menu {
        m := Menu{
          Title: ch.Name,
          ID: ch.ID,
          ParentID: ch.ParentID,
        }
        if m.ParentID == "" { m.ParentID = ch.GuildID }
        return m
      }).([]Menu),
    },
  }

  // marshal and return
  res, _ := json.Marshal(menus)
  return res
}
