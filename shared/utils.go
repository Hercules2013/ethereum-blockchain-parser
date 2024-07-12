package shared

import (
	"fmt"
)

// CurrentBlockToHex converts a block number to its hexadecimal representation
func CurrentBlockToHex(currentBlock int64) string {
	return fmt.Sprintf("0x%x", currentBlock)
}
