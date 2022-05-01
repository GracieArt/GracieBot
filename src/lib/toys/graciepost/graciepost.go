package graciepost

import (
  "github.com/gracieart/bubblebot"
  "github.com/gracieart/graciebot/src/lib/toys/like"
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
  toyID string
  toyInfo bubble.ToyInfo
  bot *bubble.Bot
  port string
  charLimit int
  like *like.Like
  key string
}

func (g *GraciePost) ToyID() string { return g.toyID }
func (g *GraciePost) ToyInfo() bubble.ToyInfo { return g.toyInfo }


type Config struct {
  CharLimit int
  Port string
  Key string
}


func New(cnf Config) *GraciePost {
  gp := &GraciePost{
    toyID : "graciebell.art.graciepost",
    port : "30034",
    charLimit : 180,
    key : cnf.Key,
    toyInfo: bubble.ToyInfo{
      Name: "GraciePost",
      Description: "Post images from your browser using the GraciePost Firefox extension.",
    },
  }
  if cnf.Port != "" { gp.port = cnf.Port }
  if cnf.CharLimit > 0 { gp.charLimit = cnf.CharLimit }
  return gp
}



func (g *GraciePost) Load(b *bubble.Bot) error {
  g.bot = b

  if l, ok := g.bot.FindToy("graciebell.art.like"); ok {
    g.like = l.(*like.Like)
  }

  g.setHandler()

  return nil
}



func (g *GraciePost) OnLifecycleEvent(l bubble.LifecycleEvent) {
  switch l {
  case bubble.Connect:
    go g.listen()
  }
}
