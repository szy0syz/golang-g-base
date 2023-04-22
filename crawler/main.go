package main

import (
	"fmt"
	"g-base/crawler/fetcher"
	"regexp"
)

func main() {
	all, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	if err != nil {
		panic(err)
	}

	printCityList(all)
}

func printCityList(contexts []byte) {
	re := regexp.MustCompile(`<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`)
	matches := re.FindAllSubmatch(contexts, -1)
	for _, m := range matches {
		fmt.Printf("City: %s, URL:%s\n", m[2], m[1])
	}
}
