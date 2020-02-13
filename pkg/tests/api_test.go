package tests

import (
	"database/sql"
	"errors"
	"github.com/JAbduvohidov/apm-ibank-core/pkg/core"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestAddClient(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)

	if !errors.Is(err, core.ErrLoginExist) {
		t.Error(err)
	}

	err = core.AddClient("Vasya", "vasya1", "1234", 999999999, db)

	if !errors.Is(err, core.ErrPhoneNumberExist) {
		t.Error(err)
	}

	err = core.AddClient("Piter", "piter", "1234", 888888888, db)

	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

}

func TestAddAccount(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddAccount(99999999, 200000, db)

	if err == nil {
		t.Errorf("expected error found: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)

	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(999999999, 200000, db)

	if err != nil {
		t.Errorf( "unexpected error at AddAccount: %v", err)
	}
}

func TestAddService(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddService("svet", db)

	if err != nil {
		t.Errorf("unexpected error at AddService: %v", err)
	}

	err = core.AddService("svet", db)

	if ok := errors.Is(err, core.ErrServiceExist); !ok {
		t.Errorf("expected error: %v, found: %v", core.ErrServiceExist, err)
	}
}

func TestAddATM(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddAtm("Oson", "location1", db)

	if err != nil {
		t.Errorf("unexpected error at AddAtm: %v", err)
	}

	err = core.AddAtm("Oson", "location1", db)

	if ok := errors.Is(err, core.ErrATMExist); !ok {
		t.Errorf("expected error: %v, found: %v", core.ErrATMExist, err)
	}
}

func TestLogin(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)

	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	phoneNumber, err := core.Login("", "", db)

	if phoneNumber != -1 {
		t.Errorf("expected: -1, found: %d", phoneNumber)
	}

	_, err = core.Login("vasya", "", db)

	if ok := errors.Is(err, core.ErrInvalidPass); !ok {
		t.Errorf("expected error: %v, found: %v", core.ErrInvalidPass, err)
	}

	phoneNumber, err = core.Login("vasya", "1234", db)

	if err != nil {
		t.Errorf("unexpected error at Login: %v", err)
	}

	if phoneNumber == -1 {
		t.Errorf("unexpected value: %d", phoneNumber)
	}
}

func TestGetListOfClientAccounts(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = core.GetListOfClientAccounts("", db)
	if err == nil {
		t.Errorf("expected error, found: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(999999999, 20000, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	accounts, err := core.GetListOfClientAccounts("vasya", db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfClientAccounts: %v", err)
	}

	if accounts == nil {
		t.Errorf("accounts list must not be nil, found: %v", accounts)
	}
}

func TestGetListOfATMs(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	listOfATMs, err := core.GetListOfATMs(db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfATMs: %v", err)
	}

	if listOfATMs != nil {
		t.Errorf("empty list must be equal to nil, found: %v", listOfATMs)
	}

	err = core.AddAtm("ATM1", "location1", db)
	if err != nil {
		t.Errorf("unexpected error at adding atm: %v", err)
	}

	listOfATMs, err = core.GetListOfATMs(db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfATMs: %v", err)
	}

	if listOfATMs == nil {
		t.Errorf("list of atms must not be empty: %v", listOfATMs)
	}
}

func TestSearchClientByName(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	clients, err := core.SearchClientByName("", db)
	if err != nil {
		t.Errorf("unexpected error at searching: %v", err)
	}

	if clients != nil {
		t.Errorf("empty list must be nil, found: %v", clients)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at adding client: %v", err)
	}

	clients, err = core.SearchClientByName("vasya", db)
	if err != nil {
		t.Errorf("unexpected error at searching: %v", err)
	}

	if clients == nil {
		t.Errorf("not empty list must not be equal to nil, found: %v", clients)
	}
}

func TestSearchClientByPhoneNumber(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	clients, err := core.SearchClientByPhoneNumber(0, db)
	if err != nil {
		t.Errorf("unexpected error at searching: %v", err)
	}

	if clients != nil {
		t.Errorf("empty list must be nil, found: %v", clients)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at adding client: %v", err)
	}

	clients, err = core.SearchClientByPhoneNumber(999999, db)
	if err != nil {
		t.Errorf("unexpected error at searching: %v", err)
	}

	if clients == nil {
		t.Errorf("not empty list must not be equal to nil, found: %v", clients)
	}
}

func TestGetListOfClients(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	clients, err := core.GetListOfClients(db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfClients: %v", err)
	}

	if clients != nil {
		t.Errorf("empty list must be nil, found: %v", clients)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at adding client: %v", err)
	}

	clients, err = core.GetListOfClients(db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfClients: %v", err)
	}

	if clients == nil {
		t.Errorf("not empty list must not be equal to nil, found: %v", clients)
	}
}

func TestGetListOfClientsFormatted(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	clients, err := core.GetListOfClientsFormatted(10, 0, db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfClientsFormatted: %v", err)
	}

	if clients != nil {
		t.Errorf("empty list must be nil, found: %v", clients)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	err = core.AddClient("Vasya", "vasya1", "1234", 999999998, db)
	if err != nil {
		t.Errorf("unexpected error at adding client: %v", err)
	}

	clients, err = core.GetListOfClientsFormatted(10, 0, db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfClientsFormatted: %v", err)
	}

	if len(clients) != 2 {
		t.Errorf("2 item slice must have size 2, found: %v", len(clients))
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	err = core.AddClient("Vasya", "vasya1", "1234", 999999998, db)
	err = core.AddClient("Vasya", "vasya2", "1234", 999999997, db)
	err = core.AddClient("Vasya", "vasya3", "1234", 999999996, db)
	err = core.AddClient("Vasya", "vasya4", "1234", 999999995, db)
	err = core.AddClient("Vasya", "vasya5", "1234", 999999994, db)
	err = core.AddClient("Vasya", "vasya6", "1234", 999999993, db)
	err = core.AddClient("Vasya", "vasya7", "1234", 999999992, db)
	err = core.AddClient("Vasya", "vasya8", "1234", 999999991, db)
	err = core.AddClient("Vasya", "vasya9", "1234", 999999990, db)
	err = core.AddClient("Vasya", "vasya10", "1234", 999999989, db)
	if err != nil {
		t.Errorf("unexpected error at adding client: %v", err)
	}

	clients, err = core.GetListOfClientsFormatted(10, 0, db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfClientsFormatted: %v", err)
	}

	if len(clients) != 10 {
		t.Errorf("11 item slice with limit 10 and offset 0 must have size 10, found: %v", len(clients))
	}
}

func TestGetListOfAccountsWithClients(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = core.GetListOfAccountsWithClients(db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfAccountsWithClients: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(999999999, 20000, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	accounts, err := core.GetListOfAccountsWithClients(db)
	if err != nil {
		t.Errorf("unexpected error at GetListOfAccountsWithClients: %v", err)
	}

	if accounts == nil {
		t.Errorf("accounts list must not be nil, found: %v", accounts)
	}
}

func TestPayForService(t *testing.T) {
	db, err := sql.Open("sqlite3", "db1.sqlite")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.PayForService("", 0, "", 0, db)
	if ok := errors.Is(err, core.ErrServiceNotExist); !ok {
		t.Errorf("expected error: %v, found: %v", core.ErrServiceNotExist, err)
	}

	err = core.AddService("Water", db)
	if err != nil {
		t.Errorf("unexpected error at AddService: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(999999999, 20000, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.PayForService("Water", 1, "vasya", 2000, db)
	if err != nil {
		t.Error(err)
	}
}

func TestChangeClientStatus(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya", "1234", 999999999, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.ChangeClientStatus(0, core.Active, db)
	if ok := errors.Is(err, core.ErrPhoneNumberNotExist); !ok {
		t.Errorf("expected error: %v, found: %v", core.ErrPhoneNumberNotExist, err)
	}

	err = core.ChangeClientStatus(999999999, core.Active, db)
	if err != nil {
		t.Errorf("unexpected error at ChangeClientStatus: %v", err)
	}
}

func TestTransferToByAccountId(t *testing.T) {
	db, err := sql.Open("sqlite3", "db2.sqlite")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya1", "1234", 1234, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddClient("Vasya", "vasya2", "1234", 5678, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(1234, 50000, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.AddAccount(5678, 0, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.TransferToByAccountId(2, "vasya1", 1, 50000, db)
	if err != nil {
		t.Errorf("expected empty error, found: %v", err)
	}
}

func TestTransferToByPhoneNumber(t *testing.T) {
	db, err := sql.Open("sqlite3", "db3.sqlite")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya1", "1234", 1234, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddClient("Vasya", "vasya2", "1234", 5678, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(1234, 50000, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.AddAccount(5678, 0, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.TransferToByPhoneNumber(5678, "vasya1", 1, 50000, db)
	if err != nil {
		t.Errorf("expected empty error, found: %v", err)
	}
}

func TestImportListOfClients(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var clients []core.Client

	err = core.ImportListOfClients(clients, db)
	if err != nil {
		t.Errorf("unexpected error at ImportListOfClients: %v", err)
	}
}

func TestImportListOfAccounts(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var accounts []core.AccountWithClientId

	err = core.ImportListOfAccounts(accounts, db)
	if err != nil {
		t.Errorf("unexpected error at ImportListOfAccounts: %v", err)
	}
}

func TestImportListOfATMs(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var atms []core.ATM

	err = core.ImportListOfATMs(atms, db)
	if err != nil {
		t.Errorf("unexpected error at ImportListOfATMs: %v", err)
	}
}

func TestGetJournalListFormatted(t *testing.T) {
	db, err := sql.Open("sqlite3", "db4.sqlite")
	if err != nil {
		t.Errorf("can't open db: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("can't close db: %v", err)
		}
	}()

	err = core.Init(db)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = core.AddClient("Vasya", "vasya1", "1234", 1234, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddClient("Vasya", "vasya2", "1234", 5678, db)
	if err != nil {
		t.Errorf("unexpected error at AddClient: %v", err)
	}

	err = core.AddAccount(1234, 50000, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.AddAccount(5678, 0, db)
	if err != nil {
		t.Errorf("unexpected error at AddAccount: %v", err)
	}

	err = core.TransferToByAccountId(2, "vasya1", 1, 50000, db)
	if err != nil {
		t.Errorf("expected empty error, found: %v", err)
	}

	journals, err := core.GetJournalListFormatted("vasya1", 10, 0, db)
	if err != nil {
		t.Errorf("unexpected error at GetJournalListFormatted: %v", err)
	}
	if journals == nil {
		t.Errorf("not empty list nust not be equal to nil")
	}
}