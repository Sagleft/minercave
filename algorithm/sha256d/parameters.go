/*
Author: Mathieu Mailhos
Filename: parameters.go
Description: Define mining parameters specific to Bitcoin BP023
*/

package sha256d

import (
	"fmt"
	"math/big"
	"minercave/utils"
)

//Gettarget TODO Calculate the right target depending on the difficulty. This one is totally made up for testing purpose.
func Gettarget(difficulty utils.Difficulty) string {
	diffAsBig := new(big.Float).SetUint64(uint64(difficulty))
	bigEndian := "0x00000000ffff0000000000000000000000000000000000000000000000000000"
	targetAsBigInt, _ := new(big.Int).SetString(bigEndian, 0)
	targetAsBigFloat := new(big.Float)
	targetAsBigFloat.SetInt(targetAsBigInt).Quo(targetAsBigFloat, diffAsBig)
	target, _ := targetAsBigFloat.Int(nil)
	return fmt.Sprintf("%064x", target)
}
