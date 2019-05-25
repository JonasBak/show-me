package main

import (
  "fmt"
  "net/http"
  "os/exec"
  "os"
  "hash/fnv"
  "strings"
  "io"
  "io/ioutil"
  "bytes"
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
  str := "## Files:\n"
  for key, _ := range files {
    str += fmt.Sprintf("* [%s](%s)\n", key, key)
  }
  str = pandoc(ioutil.NopCloser(bytes.NewReader([]byte(str))), "")
  // TODO use amount of files as ts for / path
  fmt.Fprintf(w, strings.Replace(str, "/*hash*/", "\"0\"", 1))
}

func pandoc(in io.ReadCloser, from string) string {
  args := []string{"-t", "html", "-H", "reload.html", "-H", "style.html"}
  if len(from) > 0 {
    args = append(args, "-f", from)
  }

  cmd := exec.Command(os.Getenv("pandoc"), args...)
  cmd.Stdin = in

  b, err := cmd.Output()
  if err != nil {
    panic(err)
  }
  return string(b)
}

func FileServer(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    ts, ok := r.URL.Query()["ts"]
    if ok && len(ts) == 1 {
      fmt.Fprintf(w, gethash(r.URL.Path))
    } else {
      if r.URL.Path == "/" {
        writeindices(w)
        return
      }
      fmt.Fprintf(w, getcontent(r.URL.Path))
    }
  } else if r.Method == http.MethodPost {

    from_q, ok := r.URL.Query()["from"]
    from := ""
    if !ok && len(from_q) == 1 && len(from_q[0]) > 0 {
      from = from_q[0]
    }

    content := pandoc(r.Body, from)
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
