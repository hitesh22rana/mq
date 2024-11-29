// pkg/utils/generator.go

package utils

import (
	"time"
	
	"github.com/rs/xid"
)

var (
	messageIDPrefix    = "msg"
	subscriberIDPrefix = "sub"
)

type Generator interface {
	GetUniqueMessageID() string
	GetUniqueSubscriberID() string
	GetCurrentTimestamp() int64
}

type generator struct{}

func NewGenerator() Generator {
	return &generator{}
}

func (g *generator) GetUniqueMessageID() string {
	return messageIDPrefix + xid.New().String()
}

func (g *generator) GetUniqueSubscriberID() string {
	return subscriberIDPrefix + xid.New().String()
}

func (g *generator) GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}
