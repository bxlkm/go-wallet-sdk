package txnbuild

import (
	"github.com/bxlkm/go-wallet-sdk/coins/stellar/support/errors"
	"github.com/bxlkm/go-wallet-sdk/coins/stellar/xdr"
)

// AccountMerge represents the Stellar merge account operation. See
// https://developers.stellar.org/docs/start/list-of-operations/
type AccountMerge struct {
	Destination   string
	SourceAccount string
}

// BuildXDR for AccountMerge returns a fully configured XDR Operation.
func (am *AccountMerge) BuildXDR() (xdr.Operation, error) {
	var xdrOp xdr.MuxedAccount
	err := xdrOp.SetAddress(am.Destination)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to set destination address")
	}

	opType := xdr.OperationTypeAccountMerge
	body, err := xdr.NewOperationBody(opType, xdrOp)
	if err != nil {
		return xdr.Operation{}, errors.Wrap(err, "failed to build XDR OperationBody")
	}
	op := xdr.Operation{Body: body}
	SetOpSourceAccount(&op, am.SourceAccount)
	return op, nil
}

// FromXDR for AccountMerge initialises the txnbuild struct from the corresponding xdr Operation.
func (am *AccountMerge) FromXDR(xdrOp xdr.Operation) error {
	if xdrOp.Body.Type != xdr.OperationTypeAccountMerge {
		return errors.New("error parsing account_merge operation from xdr")
	}

	am.SourceAccount = accountFromXDR(xdrOp.SourceAccount)
	if xdrOp.Body.Destination != nil {
		am.Destination = xdrOp.Body.Destination.Address()
	}

	return nil
}

// Validate for AccountMerge validates the required struct fields. It returns an error if any of the fields are
// invalid. Otherwise, it returns nil.
func (am *AccountMerge) Validate() error {
	var err error
	_, err = xdr.AddressToMuxedAccount(am.Destination)
	if err != nil {
		return NewValidationError("Destination", err.Error())
	}
	return nil
}

// GetSourceAccount returns the source account of the operation, or the empty string if not
// set.
func (am *AccountMerge) GetSourceAccount() string {
	return am.SourceAccount
}
