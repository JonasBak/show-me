package main

import (
  "fmt"
  "net/http"
  "os/exec"
  "hash/fnv"
  "strings"
)


type file struct {
  hash uint32
  content string
}

var (
	loaded file
)

func main() {
  loaded = file {0, "No file loaded yet"}
  http.HandleFunc("/", HelloServer)
  http.HandleFunc("/ts", TsServer)
  http.ListenAndServe(":5555", nil)
}

func hash(s string) uint32 {
  h := fnv.New32a()
  h.Write([]byte(s))
  return h.Sum32()
}

func TsServer(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, fmt.Sprint(loaded.hash))
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    fmt.Fprintf(w, loaded.content)
  } else if r.Method == http.MethodPost {
    args := []string{"-t", "html", "-H", "reload.html"}

    from, ok := r.URL.Query()["from"]
    if ok && len(from) == 1 {
      args = append(args, "-f", from[0])
    }

    cmd := exec.Command("pandoc", args...)
    cmd.Stdin = r.Body

    b, err := cmd.Output()
    if err != nil {
      panic(err)
    }
    content := string(b)
    ts := hash(content)
    loaded = file {ts, strings.Replace(content, "/*hash*/", fmt.Sprintf("\"%d\"", ts), 1)}
  }
}
