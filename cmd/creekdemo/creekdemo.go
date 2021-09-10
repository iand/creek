package main

import (
	"flag"
	"log"
	"os"

	"github.com/filecoin-project/go-address"
	"github.com/iand/creek"
	"github.com/ipfs/go-cid"
)

var (
	token    = flag.String("token", "", "Your authentication token")
	filename = flag.String("file", "", "A file to upload to estuary")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	c := creek.NewDefault()

	log.Printf("Fetching health of Estuary node")
	h, err := c.Health().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("  Status: %s", h.Status)

	log.Printf("Fetching stats about Estuary node")
	ps, err := c.PublicStats().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("  Total storage: %v", ps.TotalStorage)
	log.Printf("  Total files stored: %v", ps.TotalFilesStored)
	log.Printf("  Deals on chain: %v", ps.DealsOnChain)

	log.Printf("Fetching information about Estuary node")
	pi, err := c.PublicNodeInfo().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("  Primary address: %v", pi.PrimaryAddress)

	vcid, err := cid.Decode("QmVrrF7DTnbqKvWR7P7ihJKp4N5fKmBX29m5CHbW9WLep9")
	if err != nil {
		log.Fatalf("decode cid: %v", err)
	}

	log.Printf("Fetching information about content ", vcid.String())
	infos, err := c.PublicContentByCid(vcid).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for _, info := range infos {
		log.Printf("  Content cid: %s", info.Content.Cid)
		log.Printf("  Content name: %s", info.Content.Name)
		log.Printf("  Content size: %d", info.Content.Size)
	}

	miner, err := address.NewFromString("f02620")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Printf("Fetching statistics for %s", miner.String())
	minerstats, err := c.PublicMinerStats(miner).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("  Miner name: %s", minerstats.Name)
	log.Printf("  Miner version: %s", minerstats.Version)
	log.Printf("  Miner deal count: %v", minerstats.DealCount)

	if *token == "" {
		log.Printf("use -token to demo authenticated services")
		return
	}

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
