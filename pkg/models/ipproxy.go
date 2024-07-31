package models

import (
	"time"

	"gorm.io/gorm/clause"
)

type ProxyType int8

const (
	TransparentProxy ProxyType = iota + 1
	AnonymousProxy
	HighAnonymityProxy
)

type ProtocolType int8

const (
	HTTP ProtocolType = iota + 1
	HTTPS
	Both
)

type IPProxy struct {
	ID        uint64       `gorm:"primaryKey;autoIncrement"`
	IP        string       `gorm:"type:varchar(63);not null;uniqueIndex:idx_ip_port"`
	Port      int          `gorm:"not null;uniqueIndex:idx_ip_port"`
	Type      ProxyType    `gorm:"not null"`
	Protocol  ProtocolType `gorm:"not null"`
	Country   string       `gorm:"type:varchar(63)"`
	Region    string       `gorm:"type:varchar(63)"`
	ISP       string       `gorm:"type:varchar(63)"`
	Source    string       `gorm:"type:varchar(63);not null"`
	UpdatedAt int64        `gorm:"not null"`
}

func UpsertIPProxy(proxy *IPProxy) error {
	if proxy.UpdatedAt <= 0 {
		proxy.UpdatedAt = int64(time.Now().Unix())
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ip"}, {Name: "port"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "protocol", "country", "region", "isp", "source", "updated_at"}),
	}).Create(proxy).Error
}

type ListIPProxyRequest struct {
	Type     ProxyType    `form:"type" json:"type"`
	Protocol ProtocolType `form:"protocol" json:"protocol"`
	Limit    int64        `form:"limit" json:"limit"`
}

func ListIPProxy(req *ListIPProxyRequest) []*IPProxy {
	q := db.Table("ip_proxies")
	if req != nil {
		if req.Protocol != 0 {
			q = q.Where("protocol = ?", req.Protocol)
		}
		if req.Type != 0 {
			q = q.Where("type = ?", req.Type)
		}
		if req.Limit > 0 {
			q = q.Limit(int(req.Limit))
		}
	}

	var resp []*IPProxy
	q.Scan(&resp)
	return resp
}

func GetIPProxyCount() int64 {
	var cnt int64
	db.Table("ip_proxies").Count(&cnt)
	return cnt
}

type DeleteIPProxyRequest struct {
	IP   string `form:"ip" json:"ip"`
	Port int    `form:"port" json:"port"`
}

func DeleteIPProxy(req *DeleteIPProxyRequest) error {
	return db.Where("ip = ? and port = ?", req.IP, req.Port).Delete(&IPProxy{}).Error
}
