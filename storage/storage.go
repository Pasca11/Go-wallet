package storage

import (
	"github.com/Pasca11/types"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(account *types.Account) (*types.Account, error)
	DeleteAccount(id int) error
	GetAllAccounts() ([]*types.Account, error)
	GetAccountByID(id int) (*types.Account, error)
	UpdateAccount(account *types.Account) error
	GetAccountByWallet(wallet int) (*types.Account, error)
	Transfer(req types.TransferRequest) error
}
