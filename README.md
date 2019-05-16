# show-me
> Tool for displaying markdown (and some other formats) quickly as a html page in your browser

Requires docker to be installed. The docker container runs a simple web server that uses pandoc to convert a file to a html page that is served locally.

## Installation
Run the following commands to install, the file can be placed in any directory in your `PATH`.
```
curl -L https://raw.githubusercontent.com/JonasBak/show-me/master/bin/show-me -o show-me
chmod +x ./show-me
sudo mv ./show-me /usr/local/bin/show-me
```

## Instructions
Run with `show-me [-w] FILE`. The command will print out the url that shows your file. The page will reload if the command is run again. If the `-w` flag is used, the page will update automatically when the file changes.

Stop the server with `show-me kill`.

## Todos
 * `show-me --help`
 * Nicer styling
 * Run in the background with `-w` flag
