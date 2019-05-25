package main

import (
  "fmt"
  "net/http"
  "os/exec"
  "os"
  "hash/fnv"
  "strings"
)


type file struct {
  hash uint32
  content string
}

var (
  files map[string]file
)

func main() {
  files = make(map[string]file)
  http.HandleFunc("/", FileServer)
  http.ListenAndServe(":5555", nil)
}

func hash(s string) uint32 {
  h := fnv.New32a()
  h.Write([]byte(s))
  return h.Sum32()
}

func gethash(path string) string {
  f, ok := files[path]
  if ok {
    return fmt.Sprint(f.hash)
  } else {
    return "0"
  }
}

func getcontent(path string) string {
  f, ok := files[path]
  if ok {
    return f.content
  } else {
    return "No file loaded yet"
  }
}

func writeindices(w http.ResponseWriter) {
  for key, _ := range files {
      fmt.Fprintf(w, fmt.Sprintf("<a href=\"%s\">%s</a><br>", key, key))
  }
}

func FileServer(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    if r.URL.Path == "/" {
      writeindices(w)
      return
    }
    ts, ok := r.URL.Query()["ts"]
    if ok && len(ts) == 1 {
      fmt.Fprintf(w, gethash(r.URL.Path))
    } else {
      fmt.Fprintf(w, getcontent(r.URL.Path))
    }
  } else if r.Method == http.MethodPost {
    args := []string{"-t", "html", "-H", "reload.html", "-H", "style.html"}

    from, ok := r.URL.Query()["from"]
    if ok && len(from) == 1 && len(from[0]) > 0 {
      args = append(args, "-f", from[0])
    }

    cmd := exec.Command(os.Getenv("pandoc"), args...)
    cmd.Stdin = r.Body

    b, err := cmd.Output()
    if err != nil {
      panic(err)
    }
    content := string(b)
    ts := hash(content)
    f := file {ts, strings.Replace(content, "/*hash*/", fmt.Sprintf("\"%d\"", ts), 1)}

    filename, ok := r.URL.Query()["filename"]
    if ok && len(filename) == 1 {
      fmt.Printf("File added at /%s\n", filename[0])
      files[fmt.Sprintf("/%s", filename[0])] = f
    } else {
      files["unknown"] = f
    }
  }
}
