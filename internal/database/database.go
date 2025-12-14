package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"wisetech-lms-api/internal/config"
)

const SqliteSchema = `
-- Lenders Table
CREATE TABLE IF NOT EXISTS Lenders (
    Lender_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Business_Name TEXT NOT NULL,
    Phone_Number TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    Interest_Rate_Percent REAL NOT NULL CHECK (Interest_Rate_Percent >= 0 AND Interest_Rate_Percent <= 100),
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Is_Active INTEGER DEFAULT 1
);

-- Borrowers Table
CREATE TABLE IF NOT EXISTS Borrowers (
    Borrower_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Fullnames TEXT NOT NULL,
    Email TEXT NOT NULL UNIQUE,
    Phone_Number TEXT NOT NULL,
    Residence TEXT,
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Is_Active INTEGER DEFAULT 1
);

-- Accounts Table
CREATE TABLE IF NOT EXISTS Accounts (
    Account_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Lender_ID INTEGER NOT NULL REFERENCES Lenders(Lender_ID) ON DELETE CASCADE,
    Username TEXT NOT NULL UNIQUE,
    Password_Hash TEXT NOT NULL,
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Last_Login DATETIME,
    Is_Locked INTEGER DEFAULT 0
);

-- Plans Table
CREATE TABLE IF NOT EXISTS Plans (
    Plan_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Plan TEXT NOT NULL,
    Price REAL NOT NULL CHECK (Price >= 0),
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Is_Active INTEGER DEFAULT 1
);

-- Lender_Ledger Table
CREATE TABLE IF NOT EXISTS Lender_Ledger (
    Ledger_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Lender_ID INTEGER NOT NULL REFERENCES Lenders(Lender_ID) ON DELETE CASCADE,
    Plan_ID INTEGER NOT NULL REFERENCES Plans(Plan_ID) ON DELETE RESTRICT,
    Status TEXT NOT NULL CHECK (Status IN ('active', 'inactive', 'suspended', 'expired')),
    Start_Date DATETIME DEFAULT CURRENT_TIMESTAMP,
    End_Date DATETIME,
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Loans Table
CREATE TABLE IF NOT EXISTS Loans (
    Loan_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Borrower_ID INTEGER NOT NULL REFERENCES Borrowers(Borrower_ID) ON DELETE RESTRICT,
    Lender_ID INTEGER NOT NULL REFERENCES Lenders(Lender_ID) ON DELETE RESTRICT,
    Months_To_Pay INTEGER NOT NULL CHECK (Months_To_Pay > 0),
    Payment_Status TEXT NOT NULL CHECK (Payment_Status IN ('pending', 'active', 'paid', 'defaulted', 'cancelled')),
    Amount REAL NOT NULL CHECK (Amount > 0),
    Interest_Rate REAL NOT NULL CHECK (Interest_Rate >= 0),
    Monthly_Payment REAL,
    Start_Date DATE NOT NULL,
    End_Date DATE,
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Recipets Table
CREATE TABLE IF NOT EXISTS Recipets (
    Recipet_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Loan_ID INTEGER NOT NULL REFERENCES Loans(Loan_ID) ON DELETE CASCADE,
    Timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    Status TEXT NOT NULL CHECK (Status IN ('paid', 'pending', 'failed', 'refunded')),
    Amount REAL NOT NULL CHECK (Amount > 0),
    Payment_Method TEXT,
    Transaction_Reference TEXT UNIQUE,
    Notes TEXT
);

-- File Table
CREATE TABLE IF NOT EXISTS File (
    File_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Lender_ID INTEGER NOT NULL REFERENCES Lenders(Lender_ID) ON DELETE CASCADE,
    Value TEXT NOT NULL,
    File_Type TEXT,
    File_Size INTEGER,
    Original_Filename TEXT,
    Uploaded_At DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Text Table
CREATE TABLE IF NOT EXISTS Text (
    Text_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Lender_ID INTEGER NOT NULL REFERENCES Lenders(Lender_ID) ON DELETE CASCADE,
    Value TEXT NOT NULL,
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Number Table
CREATE TABLE IF NOT EXISTS Number (
    Number_ID INTEGER PRIMARY KEY AUTOINCREMENT,
    Lender_ID INTEGER NOT NULL REFERENCES Lenders(Lender_ID) ON DELETE CASCADE,
    Value REAL NOT NULL,
    Created_At DATETIME DEFAULT CURRENT_TIMESTAMP,
    Updated_At DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_accounts_lender_id ON Accounts(Lender_ID);
CREATE INDEX IF NOT EXISTS idx_lender_ledger_lender_id ON Lender_Ledger(Lender_ID);
CREATE INDEX IF NOT EXISTS idx_loans_borrower_id ON Loans(Borrower_ID);

-- Triggers to update the Updated_At timestamp
CREATE TRIGGER IF NOT EXISTS update_lenders_updated_at AFTER UPDATE ON Lenders
FOR EACH ROW
BEGIN
    UPDATE Lenders SET Updated_At = CURRENT_TIMESTAMP WHERE Lender_ID = OLD.Lender_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_borrowers_updated_at AFTER UPDATE ON Borrowers
FOR EACH ROW
BEGIN
    UPDATE Borrowers SET Updated_At = CURRENT_TIMESTAMP WHERE Borrower_ID = OLD.Borrower_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_accounts_updated_at AFTER UPDATE ON Accounts
FOR EACH ROW
BEGIN
    UPDATE Accounts SET Updated_At = CURRENT_TIMESTAMP WHERE Account_ID = OLD.Account_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_plans_updated_at AFTER UPDATE ON Plans
FOR EACH ROW
BEGIN
    UPDATE Plans SET Updated_At = CURRENT_TIMESTAMP WHERE Plan_ID = OLD.Plan_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_lender_ledger_updated_at AFTER UPDATE ON Lender_Ledger
FOR EACH ROW
BEGIN
    UPDATE Lender_Ledger SET Updated_At = CURRENT_TIMESTAMP WHERE Ledger_ID = OLD.Ledger_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_loans_updated_at AFTER UPDATE ON Loans
FOR EACH ROW
BEGIN
    UPDATE Loans SET Updated_At = CURRENT_TIMESTAMP WHERE Loan_ID = OLD.Loan_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_text_updated_at AFTER UPDATE ON Text
FOR EACH ROW
BEGIN
    UPDATE Text SET Updated_At = CURRENT_TIMESTAMP WHERE Text_ID = OLD.Text_ID;
END;

CREATE TRIGGER IF NOT EXISTS update_number_updated_at AFTER UPDATE ON Number
FOR EACH ROW
BEGIN
    UPDATE Number SET Updated_At = CURRENT_TIMESTAMP WHERE Number_ID = OLD.Number_ID;
END;
`

// NewConnection creates a new database connection
func NewConnection(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open database: %w", err)
	}

	// Ping the database to verify the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

// InitializeSchema creates the database schema if it doesn't exist
func InitializeSchema(db *sql.DB) error {
	_, err := db.Exec(SqliteSchema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	log.Println("Database schema initialized successfully")
	return nil
}
