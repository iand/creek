package creek

import "time"

type Error struct {
	Error string `json:"error"`
}

type Health struct {
	Status string `json:"status"`
}

type PublicStats struct {
	TotalStorage     int64 `json:"totalStorage"`
	TotalFilesStored int64 `json:"totalFiles"`
	DealsOnChain     int64 `json:"dealsOnChain"`
}

type PublicNodeInfo struct {
	PrimaryAddress string `json:"primaryAddress"`
}

type AddedContent struct {
	Cid       string   `json:"cid"`
	EstuaryId uint     `json:"estuaryId"`
	Providers []string `json:"providers"`
}

type IpfsPin struct {
	Cid     string                 `json:"cid"`
	Name    string                 `json:"name"`
	Origins []string               `json:"origins"`
	Meta    map[string]interface{} `json:"meta"`
}

type IpfsPinStatus struct {
	RequestId string                 `json:"requestid"`
	Status    string                 `json:"status"`
	Created   time.Time              `json:"created"`
	Pin       IpfsPin                `json:"pin"`
	Delegates []string               `json:"delegates"`
	Info      map[string]interface{} `json:"info"`
}

type ContentInfo struct {
	Content      Content       `json:"content"`
	AggregatedIn *Content      `json:"aggregatedIn,omitempty"`
	Deals        []ContentDeal `json:"deals"`
}

type Content struct {
	ID           uint   `json:"id"`
	Cid          string `json:"cid"`
	Name         string `json:"name"`
	UserID       uint   `json:"userId"`
	Description  string `json:"description"`
	Size         int64  `json:"size"`
	Active       bool   `json:"active"`
	Offloaded    bool   `json:"offloaded"`
	Replication  int    `json:"replication"`
	AggregatedIn uint   `json:"aggregatedIn"`
	Aggregate    bool   `json:"aggregate"`
	Pinning      bool   `json:"pinning"`
	PinMeta      string `json:"pinMeta"`
	Failed       bool   `json:"failed"`
	Location     string `json:"location"`
	DagSplit     bool   `json:"dagSplit"`
}

type ContentDeal struct {
	ID               uint      `json:"id"`
	Content          uint      `json:"content"`
	PropCid          string    `json:"propCid"`
	Miner            string    `json:"miner"`
	DealID           int64     `json:"dealId"`
	Failed           bool      `json:"failed"`
	Verified         bool      `json:"verified"`
	FailedAt         time.Time `json:"failedAt,omitempty"`
	DTChan           string    `json:"dtChan"`
	TransferStarted  time.Time `json:"transferStarted"`
	TransferFinished time.Time `json:"transferFinished"`
	OnChainAt        time.Time `json:"onChainAt"`
	SealedAt         time.Time `json:"sealedAt"`
}

type MinerStats struct {
	Miner           string          `json:"miner"`
	Name            string          `json:"name"`
	Version         string          `json:"version"`
	UsedByEstuary   bool            `json:"usedByEstuary"`
	DealCount       int64           `json:"dealCount"`
	ErrorCount      int64           `json:"errorCount"`
	Suspended       bool            `json:"suspended"`
	SuspendedReason string          `json:"suspendedReason"`
	ChainInfo       *MinerChainInfo `json:"chainInfo"`
}

type MinerChainInfo struct {
	PeerID    string   `json:"peerId"`
	Addresses []string `json:"addresses"`
	Owner     string   `json:"owner"`
	Worker    string   `json:"worker"`
}

type MinerDeal struct {
	ID               uint      `json:"id"`
	Content          uint      `json:"content"`
	PropCid          string    `json:"propCid"`
	Miner            string    `json:"miner"`
	DealID           int64     `json:"dealId"`
	Failed           bool      `json:"failed"`
	Verified         bool      `json:"verified"`
	FailedAt         time.Time `json:"failedAt,omitempty"`
	DTChan           string    `json:"dtChan"`
	TransferStarted  time.Time `json:"transferStarted"`
	TransferFinished time.Time `json:"transferFinished"`
	OnChainAt        time.Time `json:"onChainAt"`
	SealedAt         time.Time `json:"sealedAt"`
	ContentCid       string    `json:"contentCid"`
}

type MinerDealFailure struct {
	ID           uint   `json:"id"`
	Miner        string `json:"miner"`
	Phase        string `json:"phase"`
	Message      string `json:"message"`
	Content      uint   `json:"content"`
	MinerVersion string `json:"minerVersion"`
}

type MinerStorageAsk struct {
	Miner         string `json:"miner"`
	Price         string `json:"price"`
	VerifiedPrice string `json:"verifiedPrice"`
	MinPieceSize  uint64 `json:"minPieceSize"`
	MaxPieceSize  uint64 `json:"maxPieceSize"`
}

type PinList struct {
	Count   int             `json:"count"`
	Results []IpfsPinStatus `json:"results"`
}
