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
	readonly = flag.Bool("readonly", true, "Set to true to restrict demo to reading data")
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

	log.Printf("Fetching information about content %s", vcid.String())
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

	log.Printf("Fetching deal information for %s", miner.String())
	minerdeals, err := c.PublicMinerDeals(miner).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(minerdeals) == 0 {
		log.Printf("  No deals found for miner")
	} else {
		log.Printf("  %d deals found for miner", len(minerdeals))
		log.Printf("  First deal proposal cid %s", minerdeals[0].PropCid)
		log.Printf("  First deal content cid %s", minerdeals[0].ContentCid)
	}

	log.Printf("Fetching failed deal information for %s", miner.String())
	minerfails, err := c.PublicMinerFailures(miner).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if len(minerfails) == 0 {
		log.Printf("  No failed deals found for miner")
	} else {
		log.Printf("  %d failed deals found for miner", len(minerfails))
		log.Printf("  First deal fail phase: %s", minerfails[0].Phase)
		log.Printf("  First deal fail message: %s", minerfails[0].Message)
	}

	log.Printf("Fetching storage ask information for %s", miner.String())
	minerask, err := c.PublicMinerStorageAsk(miner).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("  Price: %s", minerask.Price)
	log.Printf("  Verified price: %s", minerask.VerifiedPrice)
	log.Printf("  Piece size range: %d-%d", minerask.MinPieceSize, minerask.MaxPieceSize)

	//-------------------------------------------------------------
	// Authenticated demos
	//-------------------------------------------------------------

	if *token == "" {
		log.Printf("Skipping authenticated services demo, specify -token to enable")
		return
	}

	ac := c.WithToken(*token)

	if !*readonly {
		demoAddContent(ac)
	}

	log.Printf("Fetching list of pins")
	pinlist, err := ac.Pins.List().Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	if pinlist.Count == 0 {
		log.Printf("  No pins found")
	} else {
		log.Printf("  %d pins found", pinlist.Count)
		log.Printf("  First pin status: %s", pinlist.Results[0].Status)
		log.Printf("  First pin created: %s", pinlist.Results[0].Created)
		log.Printf("  First pin cid: %s", pinlist.Results[0].Pin.Cid)
	}
}

func demoAddContent(ac *creek.AuthedClient) {
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

	log.Printf("Adding pin")
	pinstatus, err := ac.Pins.Add(ci).Send()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("  Pin request id: %s", pinstatus.RequestId)
	log.Printf("  Pin status: %s", pinstatus.Status)
}
