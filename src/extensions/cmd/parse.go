package cmd


import (
  "fmt"
  "errors"
  "strings"
  "unicode"
  "strconv"
)


// Splits the message into the command name and the arguments as two strings
func parseCmd(s, p string) (string, string) {
	s = strings.TrimPrefix(s, p)
  s = strings.TrimSpace(s)

  fields := strings.SplitAfterN(s, " ", 2)
  cmdName := strings.TrimSpace(fields[0])

  args := ""
  if len(fields) > 1 { args = strings.TrimSpace(fields[1]) }

  return cmdName, args
}



// Splits arguments into an array of strings
func splitArgs(argsStr string, c *Command) []string {
  args := strings.SplitAfterN(argsStr, " ", len(c.Args))
  for i, v := range args {
    args[i] = strings.TrimSpace(v)
  }
  return args
}



// type checks the argument and returns
func parseArgType(input string, argType ArgType) (interface{}, error) {

  // return nil if the argument wasnt provided
  if input == "" { return nil, nil }

  var err error
  var parsedVal interface{}

  // type checkinggggggg
  switch argType {
  case ArgType_String:
    parsedVal = input


  case ArgType_Int:
    intVal, e := strconv.Atoi(input)
    if e != nil { err = errors.New("must be an integer") }
    parsedVal = intVal


  case ArgType_Member, ArgType_Channel:
    prefix := "@"
    if argType == ArgType_Channel { prefix = "#" }

    parsedVal := strings.TrimPrefix(strings.TrimSuffix(input, ">"), "<"+prefix)
    for _, r := range parsedVal {
      if unicode.IsDigit(r) { continue }

      typeWithArticle := "a member"
      if argType == ArgType_Channel { typeWithArticle = "a channel" }
      err = fmt.Errorf("must be %s", typeWithArticle)
    }
  }

  return parsedVal, err
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
