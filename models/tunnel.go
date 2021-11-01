package models

type Tunnel struct {
	Hash       string `json:"hash" gorm:"primary_key"`
	Type       string `gorm:"column:type" json:"type"`
	DstHost    string `json:"dstHost"`
	DstPort    int    `json:"dstPort"`
	DstTLS     bool   `json:"dstTLS"`
	Autostart  bool   `json:"autoStart"`
	Connected  bool   `json:"connected"`
	PublicAddr string `json:"publicAddr"`
}
