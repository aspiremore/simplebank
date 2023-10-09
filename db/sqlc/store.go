package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
var TxKey = struct {}{}
func (s *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf(`rback err : %s , fn err : %s`, rerr, err)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResults struct {
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	Transfer    Transfer `json:"transfer"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResults, error) {
	result := TransferTxResults{}
	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(TxKey)

		//fmt.Println(txName, "create transfer from account")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{FromAccountID: arg.FromAccountID, ToAccountID: arg.ToAccountID,
			Amount: arg.Amount})
		if err != nil {
			return err
		}
		//fmt.Println(txName, "create entry to account")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: arg.ToAccountID, Amount: arg.Amount})
		if err != nil {
			return err
		}
		//fmt.Println(txName, "create entry from account")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{AccountID: arg.FromAccountID, Amount: -arg.Amount})
		if err != nil {
			return err
		}

		//fmt.Println(txName, "get from-account for update ")
		if arg.FromAccountID > arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
			fmt.Println(txName, "get to-account for update ")
			/*	fmt.Println("to account amount ", to_account.Balance)
				fmt.Println("from account amount ", from_account.Balance)*/

			//fmt.Println(txName, "update to account ")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			/*		fmt.Println("to account  after transfer amount ", result.ToAccount.Balance)
					fmt.Println("from account amount ", result.FromAccount.Balance)*/
			if err != nil {
				return err
			}
		} else {

			/*	fmt.Println("to account amount ", to_account.Balance)
				fmt.Println("from account amount ", from_account.Balance)*/

			//fmt.Println(txName, "update to account ")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			/*		fmt.Println("to account  after transfer amount ", result.ToAccount.Balance)
					fmt.Println("from account amount ", result.FromAccount.Balance)*/
			if err != nil {
				return err
			}
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})
			if err != nil {
				return err
			}
			fmt.Println(txName, "get to-account for update ")
		}

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
