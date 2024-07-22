package models

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
