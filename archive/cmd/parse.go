package cmd


import (
  "fmt"
  "errors"
  "strings"
  "unicode"
  "strconv"
  "regexp"

  "github.com/bwmarrin/discordgo"

  ."github.com/kevinwallace/fieldsn"
)


type cmdParser struct {
  prefix string
  parseData parseData

  mentionRegexp regexp.Regexp
}

func newCmdParser(prefix string) *cmdParser {
  mentionRegexp, err := regexp.Compile(`^<(@&?!?|#)(\d+)>$`)
  if err != nil { panic(err) }

  return cmdParser {
    prefix: prefix,
    mentionRegexp: mentionRegexp,
  }
}



type parseData struct {
  // data used for each parse
  msg *discordgo.Message
  input string
  rawCmdName string
  cmdName string
  inputArgs []string
}



func (p *cmdParser) parse(
  msg *discordgo.Message,
) (
  command *Command,
  argVals ArgVals,
  errMsg string,
) {

  p.parseData = parseData{
    msg: msg,
    input: strings.TrimSpace(strings.TrimPrefix(msg.Content, p.prefix)),
  }

  words := strings.Fields(p.parseData.input)
  if len(words) == 0 {
    errMsg = fmt.Sprintf(
      "You must provide a command.\n",
      "For a list of commands, use `%s commands`.", p.prefix)
    return
  }

  p.parseData.rawCmdName = words[0]
  lowerCmdName = strings.ToLower(p.parseData.rawCmdName)
  p.parseData.command, exists := man.commands[lowerCmdName]
  if !exists {
    errMsg = fmt.Sprintf(
      "Unrecognized command.\n",
      "For a list of commands, use `%s commands`.", p.prefix)
    return
  }

  // get arguments
  argVals, errMsg = p.parseArgs()
  return
}


func (p *cmdParser) parseArgs() (argVals ArgVals, errMsg string) {
  argVals = make(map[string]any)
  if len(p.parseData.command.Args) == 0 { return }

  p.parseData.inputArgs = FieldsN(
    strings.TrimPrefix(input, p.parseData.rawCmdName),
    len(p.parseData.command.Args))

  for i, arg := range p.parseData.command.Args {
    if i >= len(p.parseData.inputArgs) { continue }
    inputArg := p.parseData.inputArgs[i]

    switch arg.Type {
    case ArgType_String:
      argVals[arg.Key] = inputArg

    case ArgType_Int:
      val, errMsg := p.parseIntArg(inputArg)
      if errMsg != "" { return }
      argVals[arg.Key] = val

    case ArgType_Member:
      p.parseMentionable(arg.Type)
      mentionMatches := p.mentionRegexp.FindStringSubmatch(inputArg)
      if len(mentionMatches) == 3 {
        id = mentionMatches[2]
        prefix = mentionMatches[1]
        if prefix != "@" && prefix != "@!" {
          errMsg = "sdsdsdgdfg"
        }
      }

    case ArgType_Channel:
      mentionTypeAsString := "member"
      if arg.Type == ArgType_Channel { mentionTypeAsString = "channel" }
      id := inputArg
      mentionMatches := p.mentionRegexp.FindStringSubmatch(inputArg)
      if len(mentionMatches) == 2 {
        id = mentionMatches[1]
        mentionType = p.parseMentionType(arg.Type)
      }
      res, errMsg := resolveID()
      if errMsg != "" { return }
      argVals[arg.Key] = member


    case ArgType_Channel:

      argVals[arg.Key] = channel
    }
  }
}



func (p *cmdParser) parseIntArg() (val int, errMsg string) {
  val, err := strconv.Atoi(p.parseData.inputArgs[i])
  if err != nil { errMsg = fmt.Sprintf("`<%s>` must be an integer.", arg.Key) }
  return
}



type MentionType int
const (
  MentionType_User = iota // <@ or <@!
  MentionType_Channel // <#
  MentionType_Role // <@&
)


func (p *cmdParser) parseIfMention(
  input string,
  mentionType MentionType,
) (
  id string,
) {
  if !strings.HasPrefix(input, "<") { return input }
  mentionMatches := p.mentionRegexp.FindStringSubmatch(input)
  if len(mentionMatches) != 3 { return }
  switch {
  case strings.HasPrefix("<@&"):
    t = MentionType_Role
  case strings.HasPrefix("<#"):
    t = MentionType_Channel
  case strings.HasPrefix("<@"):
    t = MentionType_User
  }
  return

  mentionPrefix := getMentionPrefix(t)
  id := strings.TrimPrefix(strings.TrimSuffix())

  channelID := strings.TrimPrefix(strings.TrimSuffix(p.parseData.inputArgs[i], ">"), "<#")
  for _, r := range channelID {
    if !unicode.IsDigit(r) {
      errMsg = fmt.Sprintf(
        "`<%s>` must be either a channel mention or a channel ID.", arg.Key)
      return
    }
  }
  channel, err := man.bot.Session.Channel(channelID)
  if err != nil || channel.GuildID != m.GuildID {
    errMsg = fmt.Sprintf(
      "The channel provided for `<%s>` could not be found.", arg.Key,
    )
    return
  }
}



// returns nil if the string is only comprised
// of alphanumeric characters and symbols
func validateKeyword(p string) error {
  if len(p) == 0 { return errors.New("got empty string") }
  for _, r := range p {
    isAlphanumeric := unicode.IsLetter(r) || unicode.IsDigit(r)
    isSymbol := unicode.IsSymbol(r) || unicode.IsPunct(r)

    if !isAlphanumeric && !isSymbol {
      return fmt.Errorf("expecting letter, number, or symbol, but got %q", r)
    }
  }
  return nil
}
