package main

import "fmt"

const (
	bannerText = `
 __  __     __     ______
/\ \_\ \   /\ \   /\__  _\
\ \  __ \  \ \ \  \/_/\ \/
 \ \_\ \_\  \ \_\    \ \_\
  \/_/\/_/   \/_/     \/_/
`
	usageText = `
Usage:
  -url
       HTTP server URL to make requests (required)
  -n
       Number of requests to make
  -c
       Concurrency level`
)

func banner() string { return bannerText[1:] }
func usage() string  { return usageText[1:] }

func main() {
	fmt.Println(banner())
	fmt.Println(usage())
}
