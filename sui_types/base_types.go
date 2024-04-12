package sui_types

import (
	"github.com/thorli9527/sui-wallet-sdk/lib"
	"github.com/thorli9527/sui-wallet-sdk/move_types"
)

type SuiAddress = move_types.AccountAddress

type SequenceNumber = uint64

type ObjectID = move_types.AccountAddress

func NewAddressFromHex(str string) (*SuiAddress, error) {
	return move_types.NewAccountAddressHex(str)
}

func NewObjectIdFromHex(str string) (*ObjectID, error) {
	return move_types.NewAccountAddressHex(str)
}

// ObjectRef for BCS, need to keep this order
type ObjectRef struct {
	ObjectId ObjectID       `json:"objectId"`
	Version  SequenceNumber `json:"version"`
	Digest   ObjectDigest   `json:"digest"`
}

type MoveObjectType_ struct {
	Other     *move_types.StructTag
	GasCoin   *lib.EmptyEnum
	StakedSui *lib.EmptyEnum
	Coin      *move_types.TypeTag
}

func (o MoveObjectType_) IsBcsEnum() {

}
