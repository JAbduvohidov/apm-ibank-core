package core

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/JAbduvohidov/apm-ibank-core/pkg/queries"
	"strconv"
	"time"
)

var (
	ErrInvalidPass         = errors.New("invalid password")
	ErrLoginExist          = errors.New("login exits")
	ErrPhoneNumberExist    = errors.New("phone number exits")
	ErrPhoneNumberNotExist = errors.New("phone number does not exits")
	ErrServiceExist        = errors.New("service exits")
	ErrATMExist            = errors.New("atm exits")
	ErrClientIsLocked      = errors.New("client is locked")
	ErrServiceNotExist     = errors.New("service not found")
)

type QueryError struct {
	Query string
	Err   error
}

type DbError struct {
	Err error
}

type Account struct {
	Id      int64
	Balance float64
}

type AccountWithClientId struct {
	Id       int64
	ClientId int64
	Balance  float64
}

type ATM struct {
	Id       int64
	Name     string
	Location string
}

type Client struct {
	Id          int64
	Name        string
	Login       string
	Password    string
	PhoneNumber int64
	Status      string
}

type Journal struct {
	Id            int64
	Date          string
	Type          string
	TransferredTo string
	Amount        float64
}

const (
	Transfer = "transfer"
	Service  = "service"
	Active   = "active"
	Locked   = "locked"
	Clients  = "clients"
	ATMs     = "atms"
	Accounts = "accounts"
)

func (receiver *QueryError) Unwrap() error {
	return receiver.Err
}

func (receiver *QueryError) Error() string {
	return fmt.Sprintf("can't execute query %s", receiver.Err.Error())
}

func queryError(query string, err error) *QueryError {
	return &QueryError{Query: query, Err: err}
}

func (receiver *DbError) Error() string {
	return fmt.Sprintf("can't handle db operation: %v", receiver.Err.Error())
}

func (receiver *DbError) Unwrap() error {
	return receiver.Err
}

func dbError(err error) *DbError {
	return &DbError{Err: err}
}

func Init(db *sql.DB) (err error) {
	ddls := []string{queries.ClientsDDL, queries.AccountsDDL, queries.JournalDDL, queries.ServicesDDL, queries.AtmsDDL}
	for _, ddl := range ddls {
		_, err = db.Exec(ddl)
		if err != nil {
			return dbError(err)
		}
	}
	return nil
}

func checkClientExist(login string, phoneNumber int64, db *sql.DB) (err error) {
	var dbLogin string
	var dbPhoneNumber int64

	if login != "" {
		err = db.QueryRow(
			queries.LoginExistSQL,
			login).Scan(&dbLogin)

		if !errors.Is(err, sql.ErrNoRows) {
			return ErrLoginExist
		}
	}

	err = db.QueryRow(
		queries.PhoneNumberExistSQL,
		phoneNumber).Scan(&dbPhoneNumber)

	if !errors.Is(err, sql.ErrNoRows) {
		return ErrPhoneNumberExist
	}

	return nil
}

func checkServiceExist(name string, db *sql.DB) (err error) {
	var dbName string

	err = db.QueryRow(
		queries.ServiceExistSQL,
		name).Scan(&dbName)

	if !errors.Is(err, sql.ErrNoRows) {
		return ErrServiceExist
	}

	return nil
}

func checkATMExist(location string, db *sql.DB) (err error) {
	var dbLocation string

	err = db.QueryRow(
		queries.AtmExistSQL,
		location).Scan(&dbLocation)

	if !errors.Is(err, sql.ErrNoRows) {
		return ErrATMExist
	}

	return nil
}

func AddClient(name, login, password string, phoneNumber int64, db *sql.DB) (err error) {
	err = checkClientExist(login, phoneNumber, db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		queries.AddClientSQL,
		sql.Named("name", name),
		sql.Named("login", login),
		sql.Named("password", password),
		sql.Named("phone_number", phoneNumber),
	)
	if err != nil {
		return err
	}

	return nil
}

func AddAccount(phoneNumber, balance int64, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var clientId int64

	err = tx.QueryRow(
		queries.GetClientIdByPhoneNumberSQL,
		phoneNumber,
	).Scan(&clientId)
	if err != nil {
		return err
	}
	balance = balance * 100
	_, err = tx.Exec(
		queries.AddAccountSQL,
		sql.Named("client_id", clientId),
		sql.Named("balance", balance),
	)
	if err != nil {
		return err
	}

	return nil
}

func AddService(name string, db *sql.DB) (err error) {
	err = checkServiceExist(name, db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			_ = rollbackErr
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		queries.AddServiceSQL,
		sql.Named("name", name),
	)
	if err != nil {
		return err
	}

	return nil
}

func AddAtm(name, location string, db *sql.DB) (err error) {
	err = checkATMExist(location, db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.Exec(
		queries.AddAtmSQL,
		sql.Named("name", name),
		sql.Named("location", location),
	)
	if err != nil {
		return err
	}

	return nil
}

func Login(login, password string, db *sql.DB) (phoneNumber int64, err error) {
	var dbLogin, dbPassword, dbStatus string

	err = db.QueryRow(
		queries.LoginSQL,
		login).Scan(&dbLogin, &dbPassword, &phoneNumber, &dbStatus)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, nil
		}

		return -1, queryError(queries.LoginSQL, err)
	}

	if dbPassword != password {
		return -1, ErrInvalidPass
	}

	if dbStatus == Locked {
		return -1, ErrClientIsLocked
	}

	return phoneNumber, nil
}

func GetListOfClientAccounts(login string, db *sql.DB) (accounts []Account, err error) {
	var clientId int64
	err = db.QueryRow(
		queries.GetClientIdByLoginSQL,
		login,
	).Scan(&clientId)
	if err != nil {
		return nil, queryError(queries.GetClientIdByLoginSQL, err)
	}

	rows, err := db.Query(queries.GetClientAccountsSQL, clientId)
	if err != nil {
		return nil, queryError(queries.GetClientAccountsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			accounts, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		account := Account{}
		err = rows.Scan(&account.Id, &account.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		account.Balance /= 100.0
		accounts = append(accounts, account)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return accounts, nil
}

func GetJournalListFormatted(login string, limit, offset int64, db *sql.DB) (journals []Journal, err error) {
	var clientId int64
	err = db.QueryRow(
		queries.GetClientIdByLoginSQL,
		login,
	).Scan(&clientId)
	if err != nil {
		return nil, queryError(queries.GetClientIdByLoginSQL, err)
	}

	rows, err := db.Query(queries.GetJournalListFormattedSQL, clientId, limit, offset)
	if err != nil {
		return nil, queryError(queries.GetJournalListFormattedSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			journals, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		journal := Journal{}
		err = rows.Scan(&journal.Id, &journal.Date, &journal.Type, &journal.TransferredTo, &journal.Amount)
		if err != nil {
			return nil, dbError(err)
		}
		journal.Amount /= 100.0
		journals = append(journals, journal)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return journals, nil
}

func GetListOfATMs(db *sql.DB) (atms []ATM, err error) {
	rows, err := db.Query(queries.GetAllATMsSQL)
	if err != nil {
		return nil, queryError(queries.GetAllATMsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			atms, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		atm := ATM{}
		err = rows.Scan(&atm.Id, &atm.Name, &atm.Location)
		if err != nil {
			return nil, dbError(err)
		}
		atms = append(atms, atm)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return atms, nil
}

func SearchClientByName(name string, db *sql.DB) (clients []Client, err error) {
	name = "%" + name + "%"
	rows, err := db.Query(queries.SearchClientByName, name)
	if err != nil {
		return nil, queryError(queries.GetAllATMsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			clients, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		client := Client{}
		err = rows.Scan(&client.Id, &client.Name, &client.Login, &client.Password, &client.PhoneNumber, &client.Status)
		if err != nil {
			return nil, dbError(err)
		}
		clients = append(clients, client)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return clients, nil
}

func SearchClientByPhoneNumber(phoneNumber int64, db *sql.DB) (clients []Client, err error) {
	phoneNum := "%" + strconv.Itoa(int(phoneNumber)) + "%"
	rows, err := db.Query(queries.SearchClientByPhoneNumber, phoneNum)
	if err != nil {
		return nil, queryError(queries.GetAllATMsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			clients, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		client := Client{}
		err = rows.Scan(&client.Id, &client.Name, &client.Login, &client.Password, &client.PhoneNumber, &client.Status)
		if err != nil {
			return nil, dbError(err)
		}
		clients = append(clients, client)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return clients, nil
}

func GetListOfClients(db *sql.DB) (clients []Client, err error) {
	rows, err := db.Query(queries.GetListOfClientsSQL)
	if err != nil {
		return nil, queryError(queries.GetListOfClientsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			clients, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		client := Client{}
		err = rows.Scan(&client.Id, &client.Name, &client.Login, &client.Password, &client.PhoneNumber, &client.Status)
		if err != nil {
			return nil, dbError(err)
		}
		clients = append(clients, client)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return clients, nil
}

func GetListOfClientsFormatted(limit, offset int64, db *sql.DB) (clients []Client, err error) {
	rows, err := db.Query(queries.GetListOfClientsFormattedSQL, limit, offset)
	if err != nil {
		return nil, queryError(queries.GetListOfClientsFormattedSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			clients, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		client := Client{}
		err = rows.Scan(&client.Id, &client.Name, &client.Login, &client.Password, &client.PhoneNumber, &client.Status)
		if err != nil {
			return nil, dbError(err)
		}
		clients = append(clients, client)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return clients, nil
}

func GetListOfAccountsWithClients(db *sql.DB) (accountsWithClientIds []AccountWithClientId, err error) {
	rows, err := db.Query(queries.GetListOfAccountsSQL)
	if err != nil {
		return nil, queryError(queries.GetListOfAccountsSQL, err)
	}
	defer func() {
		if innerErr := rows.Close(); innerErr != nil {
			accountsWithClientIds, err = nil, dbError(innerErr)
		}
	}()

	for rows.Next() {
		accountWithClientId := AccountWithClientId{}
		err = rows.Scan(&accountWithClientId.Id, &accountWithClientId.ClientId, &accountWithClientId.Balance)
		if err != nil {
			return nil, dbError(err)
		}
		accountWithClientId.Balance /= 100.0
		accountsWithClientIds = append(accountsWithClientIds, accountWithClientId)
	}
	if rows.Err() != nil {
		return nil, dbError(rows.Err())
	}

	return accountsWithClientIds, nil
}

func PayForService(nameOfService string, accountId int64, login string, amount float64, db *sql.DB) (err error) {
	err = checkServiceExist(nameOfService, db)
	if !(errors.Is(err, ErrServiceExist)) {
		return ErrServiceNotExist
	}
	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	var clientId int64
	err = db.QueryRow(
		queries.GetClientIdByLoginSQL,
		login,
	).Scan(&clientId)
	if err != nil {
		return queryError(queries.GetClientIdByLoginSQL, err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	amount *= 100.0
	_, err = tx.Exec(
		queries.UpdateClientBalanceSQL,
		sql.Named("id", accountId),
		sql.Named("amount", -1*amount),
	)
	if err != nil {
		return err
	}

	dateAndTime := time.Now()

	_, err = tx.Exec(
		queries.AddToJournalSQL,
		sql.Named("date", dateAndTime.Format("01-02-2006 15:04:05")),
		sql.Named("client_id", clientId),
		sql.Named("type", Service),
		sql.Named("transferred_to", nameOfService),
		sql.Named("amount", amount),
	)
	if err != nil {
		return err
	}

	return nil
}

func TransferToByAccountId(targetAccountId int64, login string, accountId int64, amount float64, db *sql.DB) (err error) {
	var targetClientId int64
	err = db.QueryRow(
		queries.GetClientIdByAccountSQL,
		targetAccountId,
	).Scan(&targetClientId)
	if err != nil {
		return queryError(queries.GetClientIdByAccountSQL, err)
	}

	var targetClientStatus string
	err = db.QueryRow(
		queries.GetClientStatusSQL,
		targetClientId,
	).Scan(&targetClientStatus)
	if err != nil {
		return queryError(queries.GetClientStatusSQL, err)
	}

	if targetClientStatus == Locked {
		return ErrClientIsLocked
	}

	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var clientId int64
	err = db.QueryRow(
		queries.GetClientIdByLoginSQL,
		login,
	).Scan(&clientId)
	if err != nil {
		return queryError(queries.GetClientIdByLoginSQL, err)
	}

	amount *= 100.0

	_, err = tx.Exec(
		queries.UpdateClientBalanceSQL,
		sql.Named("id", accountId),
		sql.Named("amount", -1*amount),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		queries.UpdateClientBalanceSQL,
		sql.Named("id", targetAccountId),
		sql.Named("amount", amount),
	)
	if err != nil {
		return err
	}

	dateAndTime := time.Now()

	_, err = tx.Exec(
		queries.AddToJournalSQL,
		sql.Named("date", dateAndTime.Format("01-02-2006 15:04:05")),
		sql.Named("client_id", clientId),
		sql.Named("type", Transfer),
		sql.Named("transferred_to", targetAccountId),
		sql.Named("amount", amount),
	)
	if err != nil {
		return err
	}

	return nil
}

func TransferToByPhoneNumber(phoneNumber int64, login string, accountId int64, amount float64, db *sql.DB) (err error) {
	var targetClientStatus string
	err = db.QueryRow(
		queries.GetClientStatusByPhoneNumberSQL,
		phoneNumber,
	).Scan(&targetClientStatus)
	if err != nil {
		return queryError(queries.GetClientStatusByPhoneNumberSQL, err)
	}

	if targetClientStatus == Locked {
		return ErrClientIsLocked
	}

	var targetClientId int64
	err = db.QueryRow(
		queries.GetClientIdByPhoneNumberSQL,
		phoneNumber,
	).Scan(&targetClientId)
	if err != nil {
		return queryError(queries.GetClientIdByPhoneNumberSQL, err)
	}

	var targetAccountId int64
	err = db.QueryRow(
		queries.GetClientAccountIdSQL,
		targetClientId,
	).Scan(&targetAccountId)
	if err != nil {
		return queryError(queries.GetClientAccountIdSQL, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	var clientId int64
	err = db.QueryRow(
		queries.GetClientIdByLoginSQL,
		login,
	).Scan(&clientId)
	if err != nil {
		return queryError(queries.GetClientIdByLoginSQL, err)
	}

	amount *= 100.0

	_, err = tx.Exec(
		queries.UpdateClientBalanceSQL,
		sql.Named("id", accountId),
		sql.Named("amount", -1*amount),
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		queries.UpdateClientBalanceSQL,
		sql.Named("id", targetAccountId),
		sql.Named("amount", amount),
	)
	if err != nil {
		return err
	}

	dateAndTime := time.Now()

	_, err = tx.Exec(
		queries.AddToJournalSQL,
		sql.Named("date", dateAndTime.Format("01-02-2006 15:04:05")),
		sql.Named("client_id", clientId),
		sql.Named("type", Transfer),
		sql.Named("transferred_to", phoneNumber),
		sql.Named("amount", amount),
	)
	if err != nil {
		return err
	}

	return nil
}

func ImportListOfClients(clients []Client, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, client := range clients {
		_, err = tx.Exec(
			queries.UpdateListOfClientsSQL,
			sql.Named("id", client.Id),
			sql.Named("name", client.Name),
			sql.Named("login", client.Login),
			sql.Named("password", client.Password),
			sql.Named("phone_number", client.PhoneNumber),
			sql.Named("status", client.Status),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func ImportListOfAccounts(accountWithClientIds []AccountWithClientId, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, accountWithClientId := range accountWithClientIds {
		_, err = tx.Exec(
			queries.UpdateListOfAccountsWithClientIdsSQL,
			sql.Named("id", accountWithClientId.Id),
			sql.Named("client_id", accountWithClientId.ClientId),
			sql.Named("balance", accountWithClientId.Balance),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func ImportListOfATMs(atms []ATM, db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return dbError(err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, atm := range atms {
		_, err = tx.Exec(
			queries.UpdateListOfATMsSQL,
			sql.Named("id", atm.Id),
			sql.Named("name", atm.Name),
			sql.Named("location", atm.Location),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func ChangeClientStatus(phoneNumber int64, status string, db *sql.DB) (err error) {
	err = checkClientExist("", phoneNumber, db)
	if !errors.Is(err, ErrPhoneNumberExist) {
		return ErrPhoneNumberNotExist
	}

	_, err = db.Exec(queries.ChangeClientStatusSQL,
		sql.Named("status", status),
		sql.Named("phone_number", phoneNumber),
	)
	if err != nil {
		return dbError(err)
	}
	return nil
}