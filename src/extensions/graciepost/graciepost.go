package graciepost

import (
  "net/http"
  "log"
  "encoding/json"
  "strings"

  "github.com/bwmarrin/discordgo"
  "github.com/gracieart/graciebot-core"

  "github.com/gracieart/graciebot/src/extensions/like"
  "github.com/thoas/go-funk"
)


type PostMeta struct {
  Key           string
  Title         string
  Artist        string
  Pfp           string
  Desc          string
  PostLink      string
  ImageLink     string
  SiteName      string
  Channel       string
  OverrideEmbed bool
}



type GraciePost struct {
  bot *core.Bot
  port string
  charLimit int
  ext_like *like.Like
  likeExtensionEnabled bool
  info core.ExtensionInfo
}

func (g *GraciePost) Info() core.ExtensionInfo { return g.info }


type Config struct {
  CharLimit int
  LikeExtension *like.Like
  Port string
}


func New(cnf Config) *GraciePost {
  gp := &GraciePost{
    port : "30034",
    ext_like : cnf.LikeExtension,
    charLimit : 180,
    info: core.ExtensionInfo{
      Name: "GraciePost",
      Description: "Post images from your browser using the GraciePost Firefox extension.",
    },
  }
  if cnf.Port != "" { gp.port = cnf.Port }
  if cnf.CharLimit > 0 { gp.charLimit = cnf.CharLimit }
  if cnf.LikeExtension != nil { gp.likeExtensionEnabled = true }
  return gp
}


func (g *GraciePost) Load(b *core.Bot) {
  g.bot = b

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
      http.Error(w, "404 not found.", http.StatusNotFound)
      return
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")

    switch r.Method {
    case "GET":
      w.Header().Set("Content-Type", "application/json")
      w.Write(g.GetChannels())
    case "POST":
      decoder := json.NewDecoder(r.Body)
      var meta PostMeta
      err := decoder.Decode(&meta)
      if err != nil { http.Error(w, err.Error(), http.StatusBadRequest) }
      g.Post(meta)
    default:
      http.Error(w, "501 method not implmemented.", http.StatusNotImplemented)
    }
  })
}



func (g *GraciePost) OnConnect() {
  log.Fatal(http.ListenAndServe(":" + g.port, nil))
}



func (g *GraciePost) Post(meta PostMeta) {
  msg, err := g.bot.Session.ChannelMessageSendComplex(meta.Channel, g.createMsg(meta))
  log.Print(err) // need better error handling overall

  // add like button if the plugin is connected
  if g.likeExtensionEnabled { g.ext_like.AddLike(msg) }
}


func (g *GraciePost) createMsg(meta PostMeta) *discordgo.MessageSend {
  // just return the link if its set to override embed
  if meta.OverrideEmbed {
    return &discordgo.MessageSend{ Content: meta.PostLink }
  }

  // Truncate the description (if any)
  if len(meta.Desc) > 0 {
    cutoff := 0
    // set cutoff point at the first double linebreak (if any)
    if i := strings.Index(meta.Desc, "\n\n"); i != -1 { cutoff = i }
    // set cutoff point to the CharLimit if it exceeds it
    cutoff = funk.MinInt( []int{ cutoff, g.charLimit } )
    // truncate the string if cutoff was set
    if cutoff > 0 { meta.Desc = meta.Desc[:cutoff] + " ..." }
  }

  // create & send the embed
  msg := &discordgo.MessageSend{
    Embeds: []*discordgo.MessageEmbed{
      &discordgo.MessageEmbed{
        Title: meta.Title,
        Author: &discordgo.MessageEmbedAuthor{
          Name: meta.Artist,
          IconURL: meta.Pfp,
        },
        Description: meta.Desc,
        Fields: []*discordgo.MessageEmbedField{
          &discordgo.MessageEmbedField{
            Name: "Post link:",
            Value: meta.PostLink,
          },
        },
        Image: &discordgo.MessageEmbedImage{
          URL: meta.ImageLink,
        },
        Footer: &discordgo.MessageEmbedFooter{
          Text: "Retrieved from " + meta.SiteName + " using GraciePost",
        },
      },
    },
  }

  //log.Print(msg.Embeds[0])
  return msg
}



type Menu struct {
  Title string `json:"title"`
  ID string `json:"id"`
  ParentID string `json:"parentId"`
}

type MenuLevel struct {
  Name string `json:"name"`
  Items []Menu `json:"items"`
}

func (g *GraciePost) GetChannels() []byte {
  // get guilds
  guilds := funk.Map(g.bot.Session.State.Ready.Guilds,
    func(tempGuild *discordgo.Guild) *discordgo.Guild {
      guild, _ := g.bot.Session.Guild(tempGuild.ID)
      return guild
    },
  ).([]*discordgo.Guild)

  // get all channels
  channels := []*discordgo.Channel{}
  for _, guild := range guilds {
    guildChannels, _ := g.bot.Session.GuildChannels(guild.ID)
    channels = append(channels, guildChannels...)
  }

  // filter just the categories
  categories := funk.Filter(channels, func(ch *discordgo.Channel) bool {
    return ch.Type == discordgo.ChannelTypeGuildCategory
  }).([]*discordgo.Channel)

  // filter just the text channels
  textChannels := funk.Filter(channels, func(ch *discordgo.Channel) bool {
    return ch.Type == discordgo.ChannelTypeGuildText
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
        return Menu{
          Title: ch.Name,
          ID: ch.ID,
          ParentID: ch.ParentID,
        }
      }).([]Menu),
    },
  }

  // marshal and return
  res, _ := json.Marshal(menus)
  return res
}
