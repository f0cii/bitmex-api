package bitmex

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Wallet wallet
// swagger:model Wallet
type Wallet struct {

	// account
	// Required: true
	Account *int64 `json:"account"`

	// addr
	Addr string `json:"addr,omitempty"`

	// amount
	Amount int64 `json:"amount,omitempty"`

	// confirmed debit
	ConfirmedDebit int64 `json:"confirmedDebit,omitempty"`

	// currency
	// Required: true
	Currency *string `json:"currency"`

	// delta amount
	DeltaAmount int64 `json:"deltaAmount,omitempty"`

	// delta deposited
	DeltaDeposited int64 `json:"deltaDeposited,omitempty"`

	// delta transfer in
	DeltaTransferIn int64 `json:"deltaTransferIn,omitempty"`

	// delta transfer out
	DeltaTransferOut int64 `json:"deltaTransferOut,omitempty"`

	// delta withdrawn
	DeltaWithdrawn int64 `json:"deltaWithdrawn,omitempty"`

	// deposited
	Deposited int64 `json:"deposited,omitempty"`

	// pending credit
	PendingCredit int64 `json:"pendingCredit,omitempty"`

	// pending debit
	PendingDebit int64 `json:"pendingDebit,omitempty"`

	// prev amount
	PrevAmount int64 `json:"prevAmount,omitempty"`

	// prev deposited
	PrevDeposited int64 `json:"prevDeposited,omitempty"`

	// prev timestamp
	// Format: date-time
	PrevTimestamp strfmt.DateTime `json:"prevTimestamp,omitempty"`

	// prev transfer in
	PrevTransferIn int64 `json:"prevTransferIn,omitempty"`

	// prev transfer out
	PrevTransferOut int64 `json:"prevTransferOut,omitempty"`

	// prev withdrawn
	PrevWithdrawn int64 `json:"prevWithdrawn,omitempty"`

	// script
	Script string `json:"script,omitempty"`

	// timestamp
	// Format: date-time
	Timestamp strfmt.DateTime `json:"timestamp,omitempty"`

	// transfer in
	TransferIn int64 `json:"transferIn,omitempty"`

	// transfer out
	TransferOut int64 `json:"transferOut,omitempty"`

	// withdrawal lock
	WithdrawalLock []string `json:"withdrawalLock"`

	// withdrawn
	Withdrawn int64 `json:"withdrawn,omitempty"`
}

// MarshalBinary interface implementation
func (m *Wallet) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Wallet) UnmarshalBinary(b []byte) error {
	var res Wallet
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
