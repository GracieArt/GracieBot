package graciepost

import (
  "encoding/json"
  "net/http"
  "time"

  "github.com/gracieart/bubblebot"
)


func (g *GraciePost) setHandler() { http.HandleFunc("/", g.HandleRequest) }


func (g *GraciePost) HandleRequest(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Access-Control-Allow-Origin", "*")

  switch r.Method {
  case "GET":
    w.Header().Set("Content-Type", "application/json")
    w.Write(g.GetChannels())


  case "POST":
    decoder := json.NewDecoder(r.Body)
    var meta PostMeta
    err := decoder.Decode(&meta)
    if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    if meta.Key != g.key {
      http.Error(w, "401 unauthorized.", http.StatusUnauthorized)
      return
    }
    g.Post(meta)


  default:
    http.Error(w, "501 method not implmemented.", http.StatusNotImplemented)
  }
}


func (g *GraciePost) listen() {
  for {
    bubble.Log(bubble.Info, g.toyID, "Starting server on port " + g.port)

    err := http.ListenAndServe(":"+g.port, nil)

    bubble.Log(bubble.Error, g.toyID,
      "The server had to stop due to an error: " + err.Error() )
    bubble.Log(bubble.Info, g.toyID, "Restarting in 10 seconds")
    time.Sleep(10*time.Second)
  }
}
