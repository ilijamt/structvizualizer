package structvizualizer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

type ObjectConnection struct {
	Base     string
	Embedded string
	Label    string
}

func NewObjectConnection(base string, embedded string, label string) ObjectConnection {
	return ObjectConnection{
		Base:     base,
		Embedded: embedded,
		Label:    label,
	}
}

func (o ObjectConnection) Hash() string {
	b, _ := json.Marshal(o)
	hasher := md5.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
}
