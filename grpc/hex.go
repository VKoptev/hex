package grpc

import (
	"github.com/VKoptev/hex"
)

// ToEntity casts Hex to hex.Hex.
func (x *Hex) ToEntity() hex.Hex {
	if x == nil {
		return hex.ZE
	}

	return hex.New(int(x.GetQ()), int(x.GetR()))
}
