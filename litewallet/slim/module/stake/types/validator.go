package types

import (
	"github.com/QOSGroup/litewallet/litewallet/slim/codec"
	"github.com/QOSGroup/litewallet/litewallet/slim/tendermint/crypto"
	"time"

	btypes "github.com/QOSGroup/litewallet/litewallet/slim/base/types"
	abci "github.com/tendermint/tendermint/abci/types"
	//tmtypes "github.com/tendermint/tendermint/types"
)

type InactiveCode int8

// 验证节点绑定QOS与共识权重换算， 1 voting power = 1 bond tokens
var PowerReduction = btypes.OneInt()

const (
	//Active 可获得挖矿奖励状态
	Active int8 = iota

	//Inactive
	Inactive

	//Inactive Code
	Revoke        InactiveCode = iota // 2
	MissVoteBlock                     // 3
	MaxValidator                      // 4
	DoubleSign                        // 5
)

// Description - description fields for a validator
type Description struct {
	Moniker string `json:"moniker"` // name
	Logo    string `json:"logo"`    // optional logo link
	Website string `json:"website"` // optional website link
	Details string `json:"details"` // optional details
}

type Validator struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	Owner           btypes.AccAddress `json:"owner"`
	ConsPubKey      crypto.PubKey     `json:"consensus_pubkey"`
	BondTokens      btypes.BigInt     `json:"bond_tokens"` //不能超过int64最大值
	Description     Description       `json:"description"`
	Commission      Commission        `json:"commission"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactive_code"`
	InactiveTime   time.Time    `json:"inactive_time"`
	InactiveHeight int64        `json:"inactive_height"`

	MinPeriod  int64 `json:"min_period"`
	BondHeight int64 `json:"bond_height"`
}

type jsonifyValidator struct {
	OperatorAddress btypes.ValAddress `json:"validator_address"`
	Owner           btypes.AccAddress `json:"owner"`
	ConsPubKey      string            `json:"consensus_pubkey"`
	BondTokens      btypes.BigInt     `json:"bond_tokens"`
	Description     Description       `json:"description"`

	Status         int8         `json:"status"`
	InactiveCode   InactiveCode `json:"inactive_code"`
	InactiveTime   time.Time    `json:"inactive_time"`
	InactiveHeight int64        `json:"inactive_height"`

	MinPeriod  int64 `json:"min_period"`
	BondHeight int64 `json:"bond_height"`
}

func (val Validator) ConsAddress() btypes.ConsAddress {
	return btypes.ConsAddress(val.ConsPubKey.Address())
}

func (val Validator) ConsensusPower() int64 {
	return val.BondTokens.Div(PowerReduction).Int64()
}

func (val Validator) ToABCIValidator() (abciVal abci.Validator) {
	abciVal.Power = val.ConsensusPower()
	abciVal.Address = val.ConsAddress()
	return
}

//func (val Validator) ToABCIValidatorUpdate(isRemoved bool) (abciVal abci.ValidatorUpdate) {
//	abciVal.PubKey = tmtypes.TM2PB.PubKey(val.ConsPubKey)
//	if isRemoved {
//		abciVal.Power = int64(0)
//	} else {
//		abciVal.Power = val.ConsensusPower()
//	}
//	return
//}

func (val Validator) IsActive() bool {
	return val.Status == Active
}

func (val Validator) GetBondTokens() btypes.BigInt {
	return val.BondTokens
}

func (val Validator) GetConsensusPubKey() crypto.PubKey {
	return val.ConsPubKey
}

func (val Validator) GetOwner() btypes.AccAddress {
	return val.Owner
}

func (val Validator) MarshalJSON() ([]byte, error) {
	bechPubKey, err := btypes.ConsensusPubKeyString(val.ConsPubKey)
	if err != nil {
		return nil, err
	}

	return codec.Cdc.MarshalJSON(jsonifyValidator{
		OperatorAddress: val.OperatorAddress,
		Owner:           val.Owner,
		ConsPubKey:      bechPubKey,
		BondTokens:      val.BondTokens,
		Description:     val.Description,

		Status:         val.Status,
		InactiveCode:   val.InactiveCode,
		InactiveTime:   val.InactiveTime,
		InactiveHeight: val.InactiveHeight,

		MinPeriod:  val.MinPeriod,
		BondHeight: val.BondHeight,
	})
}

func (val *Validator) UnmarshalJSON(data []byte) error {

	jv := &jsonifyValidator{}
	if err := codec.Cdc.UnmarshalJSON(data, jv); err != nil {
		return err
	}

	consPubKey, err := btypes.GetConsensusPubKeyBech32(jv.ConsPubKey)
	if err != nil {
		return err
	}

	*val = Validator{
		OperatorAddress: jv.OperatorAddress,
		Owner:           jv.Owner,
		ConsPubKey:      consPubKey,
		BondTokens:      jv.BondTokens,
		Description:     jv.Description,

		Status:         jv.Status,
		InactiveCode:   jv.InactiveCode,
		InactiveTime:   jv.InactiveTime,
		InactiveHeight: jv.InactiveHeight,

		MinPeriod:  jv.MinPeriod,
		BondHeight: jv.BondHeight,
	}

	return nil
}

func (val Validator) GetValidatorAddress() btypes.ValAddress {
	return val.OperatorAddress
}
