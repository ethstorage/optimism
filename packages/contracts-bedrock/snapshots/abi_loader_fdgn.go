// +build faultdisputegamen

package snapshots

import (
	_ "embed"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

//go:embed abi/FaultDisputeGameN.json
var faultDisputeGameN []byte

func LoadFaultDisputeGameNABI() *abi.ABI {
	return loadABI(faultDisputeGameN)
}
