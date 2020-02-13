package queries

const ClientsDDL = `CREATE TABLE IF NOT EXISTS clients
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    name         TEXT    NOT NULL,
    login        TEXT    NOT NULL UNIQUE,
    password     TEXT    NOT NULL,
    phone_number INTEGER NOT NULL UNIQUE,
    status       TEXT NOT NULL
);`

const JournalDDL = `CREATE TABLE IF NOT EXISTS journal
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    date           TEXT    NOT NULL,
    client_id      INTEGER NOT NULL REFERENCES clients,
    type           TEXT    NOT NULL,
    transferred_to TEXT    NOT NULL,
    amount         INTEGER NOT NULL check ( amount > 0 )
);`

const AccountsDDL = `CREATE TABLE IF NOT EXISTS accounts
(
    id        INTEGER PRIMARY KEY AUTOINCREMENT,
    client_id INTEGER NOT NULL REFERENCES clients,
    balance   INTEGER NOT NULL check ( balance >= 0 )
);`

const ServicesDDL = `CREATE TABLE IF NOT EXISTS services
(
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT    NOT NULL UNIQUE
);`

const AtmsDDL = `CREATE TABLE IF NOT EXISTS atms
(
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL,
    location TEXT NOT NULL UNIQUE
);`

const AddClientSQL = `INSERT INTO clients(name, login, password, phone_number, status)
VALUES (:name, :login, :password, :phone_number, 'active');`

const AddAccountSQL = `INSERT INTO accounts(client_id, balance)
VALUES (:client_id, :balance);`

const AddServiceSQL = `INSERT INTO services(name)
VALUES (:name);`

const AddAtmSQL = `INSERT INTO atms(name, location)
VALUES (:name, :location);`

const AddToJournalSQL = `INSERT INTO journal(date, client_id, type, transferred_to, amount)
VALUES (:date, :client_id, :type, :transferred_to, :amount);`

const LoginSQL = `SELECT login, password, phone_number, status
FROM clients
WHERE login = ?;`

const LoginExistSQL = `SELECT login
FROM clients
WHERE login = ?;`

const PhoneNumberExistSQL = `SELECT phone_number
FROM clients
WHERE phone_number = ?;`

const ServiceExistSQL = `SELECT name
FROM services
WHERE name = ?;`

const AtmExistSQL = `SELECT location
FROM atms
WHERE location = ?;`

const GetClientIdByPhoneNumberSQL = `SELECT id
FROM clients
WHERE phone_number = ?;`

const GetClientIdByLoginSQL = `SELECT id
FROM clients
WHERE login = ?;`

const GetClientIdByAccountSQL = `SELECT client_id
FROM accounts
WHERE id = ?;`

const GetClientStatusSQL = `SELECT status
FROM clients
WHERE id = ?;`

const GetClientStatusByPhoneNumberSQL = `SELECT status
FROM clients
WHERE phone_number = ?;`

const GetClientAccountsSQL = `SELECT id, balance
FROM accounts
WHERE client_id = ?;`

const GetListOfAccountsSQL = `SELECT id, client_id, balance
FROM accounts;`

const GetListOfClientsSQL = `SELECT id, name, login, password, phone_number, status
FROM clients;`

const GetListOfClientsFormattedSQL = `SELECT id, name, login, password, phone_number, status
FROM clients ORDER BY name DESC LIMIT ? OFFSET ?;`

const GetJournalListFormattedSQL = `SELECT id, date, type, transferred_to, amount
FROM journal
WHERE client_id = ? ORDER BY date
LIMIT ? OFFSET ?;`

const GetClientAccountIdSQL = `SELECT id
FROM accounts
WHERE client_id = ? LIMIT 1;`

const GetAllATMsSQL = `SELECT id, name, location
FROM atms;`

const UpdateClientBalanceSQL = `UPDATE accounts
SET balance = balance + :amount
WHERE id = :id;`

const UpdateListOfClientsSQL = `INSERT INTO clients (id, name, login, password, phone_number, status)
VALUES (:id, :name, :login, :password, :phone_number, :status)
ON CONFLICT (login)
    DO UPDATE SET name=excluded.name,
                  login=excluded.login,
                  password=excluded.password,
                  phone_number=excluded.phone_number,
                  status=excluded.status;`

const UpdateListOfAccountsWithClientIdsSQL = `INSERT INTO accounts (id, client_id, balance)
VALUES (:id, :client_id, :balance)
ON CONFLICT (id)
    DO UPDATE SET id=excluded.id,
                  client_id=excluded.client_id,
                  balance=excluded.balance;`

const UpdateListOfATMsSQL = `INSERT OR
REPLACE INTO atms (id, name, location)
VALUES (:id, :name, :location);`

const ChangeClientStatusSQL = `UPDATE clients
SET status = :status
WHERE phone_number = :phone_number;`

const SearchClientByName = `SELECT id, name, login, password, phone_number, status
FROM clients
WHERE name LIKE ?
ORDER BY id;`

const SearchClientByPhoneNumber = `SELECT id, name, login, password, phone_number, status
FROM clients
WHERE phone_number LIKE ?
ORDER BY id;`