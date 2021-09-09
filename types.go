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
	Requestid string                 `json:"requestid"`
	Status    string                 `json:"status"`
	Created   time.Time              `json:"created"`
	Pin       IpfsPin                `json:"pin"`
	Delegates []string               `json:"delegates"`
	Info      map[string]interface{} `json:"info"`
}
