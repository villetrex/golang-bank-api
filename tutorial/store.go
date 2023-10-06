package tutorial

import "database/sql"

type Store interface {
	Querier
	 TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*query
	db *sql.DB
}

func NewStore(db *sql.DB) SQLStore {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}


// below function is not an exported function as it begins with lower case 'e'
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nill);
	if err != nill {
		return err
	}
	q := New(tx)
	err = fn(q);
	if err != nil {
		if rbErr := tx.Rollback() {
			return fmt.Error("tx err %v rb err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountid int64 `json:"from_accountid"`
	ToAccountid int64 `json:"to_accountid"`
	Amount int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	FromAccount Account `json:"from_account"`
	ToAccount Account `json:"to_account"`
	FromEntry Entry `json:"from_entry"`
	ToEntry Entry `json:"to_entry"`
}

func addMoney(
    ctx context.Context,
    q *Queries,
    accountID1 int64,
    amount1 int64,
    accountID2 int64,
    amount2 int64,
) (account1 Account, account2 Account, err error) {
    account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        ID:     accountID1,
        Amount: amount1,
    })
    if err != nil {
        return // because we're using named returned values, simply writing 'return' is equivalent to return account1, account2, err
    }

    account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
        ID:     accountID2,
        Amount: amount2,
    })
    return
}

var txKey = struct{}{}

func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult = &TransferTxResult{};
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.vallue(txKey)
		fmt.Println(txName, "create Transfer")
		result.Transfer , err = q.CreateTransfer(ctx,createTransferParams{
			fromAccountId: arg.FromAccountid,
			toAccountId: arg.ToAccountid,
		})
		if err != nill {
			return err
		}

		account1, err := q.GetAccountForUpdate(ctx,arg.FromAccpountId);
		if err != nil {
			return err
		}
		account2, err := q.GetAccountForUpdate(ctx,arg.ToAccountid);
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
            result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
        } else {
            result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
        }
		return err
	})
	return result, err
}


