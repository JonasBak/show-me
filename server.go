package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os/exec"
)

var (
	file string
)

func main() {
  file = "yeet"
  http.HandleFunc("/", HelloServer)
  http.ListenAndServe(":5555", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    fmt.Fprintf(w, file)
  } else if r.Method == http.MethodPost {
    args := []string{"-t", "html", "-H", "reload.html"}

    from, ok := r.URL.Query()["from"]
    if ok && len(from) == 1 {
      args = append(args, "-f", from[0])
    }

    cmd := exec.Command("pandoc", args...)
    cmd.Stdin = r.Body

    stdout, err := cmd.StdoutPipe()
    if err != nil {
      panic(err)
    }
    go func() {
      b, err := ioutil.ReadAll(stdout)
      if err != nil {
        panic(err)
      }
      // TODO replace hash
      file = string(b)
    }()

    cmd.Start()
    cmd.Wait()
  }
}
