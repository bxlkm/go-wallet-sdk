// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package nft_candy_machine_v2

import (
	"errors"

	"github.com/bxlkm/go-wallet-sdk/coins/solana/base"
)

// WithdrawFunds is the `withdrawFunds` instruction.
type WithdrawFunds struct {

	// [0] = [WRITE] candyMachine
	//
	// [1] = [SIGNER] authority
	base.AccountMetaSlice `bin:"-"`
}

// NewWithdrawFundsInstructionBuilder creates a new `WithdrawFunds` instruction builder.
func NewWithdrawFundsInstructionBuilder() *WithdrawFunds {
	nd := &WithdrawFunds{
		AccountMetaSlice: make(base.AccountMetaSlice, 2),
	}
	return nd
}

// SetCandyMachineAccount sets the "candyMachine" account.
func (inst *WithdrawFunds) SetCandyMachineAccount(candyMachine base.PublicKey) *WithdrawFunds {
	inst.AccountMetaSlice[0] = base.Meta(candyMachine).WRITE()
	return inst
}

// GetCandyMachineAccount gets the "candyMachine" account.
func (inst *WithdrawFunds) GetCandyMachineAccount() *base.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetAuthorityAccount sets the "authority" account.
func (inst *WithdrawFunds) SetAuthorityAccount(authority base.PublicKey) *WithdrawFunds {
	inst.AccountMetaSlice[1] = base.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *WithdrawFunds) GetAuthorityAccount() *base.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

func (inst WithdrawFunds) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_WithdrawFunds,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst WithdrawFunds) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *WithdrawFunds) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.CandyMachine is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
	}
	return nil
}

func (obj WithdrawFunds) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	return nil
}

// NewWithdrawFundsInstruction declares a new WithdrawFunds instruction with the provided parameters and accounts.
func NewWithdrawFundsInstruction(
	// Accounts:
	candyMachine base.PublicKey,
	authority base.PublicKey) *WithdrawFunds {
	return NewWithdrawFundsInstructionBuilder().
		SetCandyMachineAccount(candyMachine).
		SetAuthorityAccount(authority)
}
