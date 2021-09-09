package main

import (
	"flag"
	"log"
	"os"

	"github.com/iand/creek"
	"github.com/ipfs/go-cid"
)

var (
	token    = flag.String("token", "", "Your authentication token")
	filename = flag.String("file", "", "A file to upload to estuary")
)

func main() {
	flag.Parse()

	c := creek.NewDefault()

	h, err := c.Health().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("status: %s", h.Status)

	ps, err := c.PublicStats().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("PublicStats.TotalStorage: %v", ps.TotalStorage)
	log.Printf("PublicStats.TotalFilesStored: %v", ps.TotalFilesStored)
	log.Printf("PublicStats.DealsOnChain: %v", ps.DealsOnChain)

	pi, err := c.PublicNodeInfo().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("PublicNodeInfo.PrimaryAddress: %v", pi.PrimaryAddress)

	if *token == "" {
		log.Printf("use -token to demo authenticated services")
		return
	}

	ac := c.WithToken(*token)

	if *filename != "" {
		fi, err := os.Open(*filename)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		defer fi.Close()

		add, err := ac.ContentAdd("eh-whs.nt", fi).Send()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		log.Printf("ContentAdd: %+v", add)
	}

	ci, err := cid.Decode("bafybeiaivflo5qbiy6wmi3i7rcobsv45za2jswrely4xnmcuijf2g7lbca")
	if err != nil {
		log.Fatalf("decode cid: %v", err)
	}
	pin, err := ac.ContentAddFromIpfs(ci).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("ContentAddFromIpfs: %+v", pin)
}
