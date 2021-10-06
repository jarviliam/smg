package main

import (
	"log"
	"os"

	"github.com/jarviliam/smg/internal/extractor"
	"github.com/jarviliam/smg/internal/mapper"
	"github.com/jarviliam/smg/internal/prioritizer"
	"github.com/jarviliam/smg/internal/requester"
	"github.com/jarviliam/smg/internal/spider"
	"github.com/jarviliam/smg/internal/storage"
	"github.com/jarviliam/smg/internal/target"
)

func main() {
	if len(os.Args[1]) == 0 {
		return
	}
	t := target.NewTarget(os.Args[1])
	t.Priority = 1
	//TODO: Move into SMG struct
	prior := prioritizer.NewPrioritizer()
	e, err := extractor.NewExtractor(t.BaseURL)
	if err != nil {
		return
	}
	e.SetUA("Go/Sitemap")
	err = spider.NewSpider().Run(
		t,
		storage.NewStorage(),
		requester.NewRequester(),
		prior,
		mapper.NewMapper(),
		e,
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("shutdown")
}
