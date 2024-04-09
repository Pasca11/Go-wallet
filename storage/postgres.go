package storage

import (
	"fmt"
	"github.com/Pasca11/types"
	"github.com/jmoiron/sqlx"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	psqlConn := "user=postgres password=postgres dbname=gowallet sslmode=disable"
	db, err := sqlx.Open("postgres", psqlConn)
	if err != nil {
		return nil, fmt.Errorf("%w: Cannot open connection to %s", err, "gowallet")
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStorage) createAccountTable() error {
	stmt, err := s.db.Preparex(`
		CREATE TABLE IF NOT EXISTS account(
		    id int PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
		    first_name varchar(50),
		    last_name varchar(50),
		    patronymic varchar(50),
		    balance int,
		    password varchar not null  ,
		    created_at timestamp
		)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil
	}
	return nil
}

func (s *PostgresStorage) CreateAccount(account *types.Account) (*types.Account, error) {
	stmt := "INSERT INTO account(first_name, last_name, balance, patronymic, created_at, password) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *"
	row := s.db.QueryRowx(stmt, account.FirstName, account.LastName, account.Balance, account.Patronymic, account.CreatedAt, account.Password)
	acc := &types.Account{}
	err := row.StructScan(acc)
	//_, err := s.db.Exec(stmt, account.FirstName, account.LastName, account.Balance, account.Patronymic, account.CreatedAt)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	stmt := "DELETE FROM account WHERE id=$1"
	_, err := s.db.Exec(stmt, id)
	return err
}

func (s *PostgresStorage) UpdateAccount(account *types.Account) error {
	return nil
}

func (s *PostgresStorage) GetAccountByID(id int) (*types.Account, error) {
	stmt := "SELECT * FROM account WHERE id=$1"
	var acc types.Account
	err := s.db.Get(&acc, stmt, id)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (s *PostgresStorage) GetAccountByWallet(wallet int) (*types.Account, error) {
	stmt := "SELECT * FROM account WHERE wallet=$1"
	var acc types.Account
	err := s.db.Get(&acc, stmt, wallet)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

func (s *PostgresStorage) GetAllAccounts() ([]*types.Account, error) {
	qr := "SELECT * FROM account"
	var res []*types.Account
	err := s.db.Select(&res, qr)
	return res, err
}

func (s *PostgresStorage) Transfer(req types.TransferRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := "UPDATE account SET balance=balance-$1 WhERE wallet=$2"
	_, err = tx.Exec(stmt, req.Amount, req.From)
	if err != nil {
		return err
	}
	stmt = "UPDATE account SET balance=balance+$1 WHERE wallet=$2"
	_, err = tx.Exec(stmt, req.Amount, req.To)
	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}
