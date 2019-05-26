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
  "time"
)


type file struct {
  hash string
  content string
  la time.Time
}

var (
  files map[string]file
)

func main() {
  files = make(map[string]file)
  writeindices()
  go watchdelete()
  http.HandleFunc("/", FileServer)
  http.ListenAndServe(":5555", nil)
}

func watchdelete() {
  th := 10 * time.Minute
  for {
    for key, _ := range files {
      if key == "" {
        continue
      }
      if time.Now().Sub(files[key].la) > th {
        fmt.Printf("Removing file %s", key)
        delete(files, key)
        writeindices()
      }
    }
    time.Sleep(th)
  }
}

func hash(s string) string {
  h := fnv.New32a()
  h.Write([]byte(s))
  return fmt.Sprint(h.Sum32())
}

func gethash(path string) string {
  f, ok := files[path[1:]]
  if ok {
    f.la = time.Now()
    files[path[1:]] = f
    return f.hash
  }
  return ""
}

func getcontent(path string) string {
  f, ok := files[path[1:]]
  if ok {
    f.la = time.Now()
    files[path[1:]] = f
    return f.content
  }
  return ""
}

func writeindices() {
  str := "## Files:\n"
  for key, _ := range files {
    if key == "" {
      continue
    }
    str += fmt.Sprintf("* [%[1]s](/%[1]s)\n", key)
  }
  files[""] = pandoc(ioutil.NopCloser(bytes.NewReader([]byte(str))), "gfm")
}

func pandoc(in io.ReadCloser, from string) file {
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
  content := string(b)
  ts := hash(content)
  return file {ts, strings.Replace(content, "/*hash*/", fmt.Sprintf("\"%s\"", ts), 1), time.Now()}
}

func FileServer(w http.ResponseWriter, r *http.Request) {
  if r.Method == http.MethodGet {
    ts, ok := r.URL.Query()["ts"]
    if ok && len(ts) == 1 {
      fmt.Fprintf(w, gethash(r.URL.Path))
    } else {
      fmt.Fprintf(w, getcontent(r.URL.Path))
    }
  } else if r.Method == http.MethodPost {
    filename_q, ok := r.URL.Query()["filename"]
    if !ok || len(filename_q) != 1 {
      fmt.Println("No file name provided...")
      return
    }
    filename := filename_q[0]

    from_q, ok := r.URL.Query()["from"]
    from := ""
    if ok && len(from_q) == 1 && len(from_q[0]) > 0 {
      from = from_q[0]
    }

    f := pandoc(r.Body, from)

    fmt.Printf("File added at /%s\n", filename)
    if _, ok := files[filename]; !ok {
      defer writeindices()
    }
    files[filename] = f
  }
}
