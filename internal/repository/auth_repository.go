package repository

import (
	"database/sql"
	"errors"
	"time"

	"wisetech-lms-api/internal/models"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrLenderNotFound  = errors.New("lender not found")
)

// AuthRepository defines the interface for authentication-related database operations.
type AuthRepository interface {
	CreateLenderAndAccount(businessName, email, phone, username, passwordHash string, interestRate float64) (int, error)
	GetAccountByUsername(username string) (*models.Account, error)
	GetAccountByID(accountID int) (*models.Account, error)
	GetLenderByAccountID(accountID int) (*models.Lender, error)
	UpdateLastLogin(accountID int) error
}

// authRepository implements AuthRepository using a SQLite database connection.
type authRepository struct {
	db *sql.DB
}

// NewAuthRepository creates a new AuthRepository instance.
func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

// CreateLenderAndAccount creates a new lender and an associated account within a transaction.
func (r *authRepository) CreateLenderAndAccount(businessName, email, phone, username, passwordHash string, interestRate float64) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Rollback on error or if Commit fails

	now := time.Now()

	// Insert into Lenders table first
	stmtLender, err := tx.Prepare("INSERT INTO Lenders (Business_Name, Phone_Number, Email, Interest_Rate_Percent, Created_At, Updated_At) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmtLender.Close()

	resLender, err := stmtLender.Exec(businessName, phone, email, interestRate, now, now)
	if err != nil {
		return 0, err
	}

	lenderID, err := resLender.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Insert into Accounts table
	stmtAccount, err := tx.Prepare("INSERT INTO Accounts (Lender_ID, Username, Password_Hash, Created_At, Updated_At) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmtAccount.Close()

	resAccount, err := stmtAccount.Exec(lenderID, username, passwordHash, now, now)
	if err != nil {
		return 0, err
	}

	accountID, err := resAccount.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(accountID), tx.Commit()
}

// GetAccountByUsername retrieves an account by its username.
func (r *authRepository) GetAccountByUsername(username string) (*models.Account, error) {
	var account models.Account
	query := `SELECT Account_ID, Lender_ID, Username, Password_Hash, Created_At, Updated_At, Last_Login, Is_Locked FROM Accounts WHERE Username = ?`
	err := r.db.QueryRow(query, username).Scan(
		&account.AccountID,
		&account.LenderID,
		&account.Username,
		&account.PasswordHash,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLogin,
		&account.IsLocked,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// GetAccountByID retrieves an account by its ID.
func (r *authRepository) GetAccountByID(accountID int) (*models.Account, error) {
	var account models.Account
	query := `SELECT Account_ID, Lender_ID, Username, Password_Hash, Created_At, Updated_At, Last_Login, Is_Locked FROM Accounts WHERE Account_ID = ?`
	err := r.db.QueryRow(query, accountID).Scan(
		&account.AccountID,
		&account.LenderID,
		&account.Username,
		&account.PasswordHash,
		&account.CreatedAt,
		&account.UpdatedAt,
		&account.LastLogin,
		&account.IsLocked,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

// GetLenderByAccountID retrieves a lender by its account ID.
func (r *authRepository) GetLenderByAccountID(accountID int) (*models.Lender, error) {
	var lender models.Lender
	var lenderID int

	// First, get the Lender_ID from the Accounts table using the Account_ID
	err := r.db.QueryRow("SELECT Lender_ID FROM Accounts WHERE Account_ID = ?", accountID).Scan(&lenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrAccountNotFound // Account not found for the given accountID
		}
		return nil, err
	}

	// Then, retrieve the lender details using the Lender_ID
	query := `SELECT Lender_ID, Business_Name, Phone_Number, Email, Interest_Rate_Percent, Created_At, Updated_At, Is_Active FROM Lenders WHERE Lender_ID = ?`
	err = r.db.QueryRow(query, lenderID).Scan(
		&lender.LenderID,
		&lender.BusinessName,
		&lender.PhoneNumber,
		&lender.Email,
		&lender.InterestRatePercent,
		&lender.CreatedAt,
		&lender.UpdatedAt,
		&lender.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// This case should ideally not happen if the foreign key constraint is properly enforced
			// and an account always has a corresponding lender.
			return nil, ErrLenderNotFound
		}
		return nil, err
	}
	return &lender, nil
}

// UpdateLastLogin updates the Last_Login timestamp for a given account.
func (r *authRepository) UpdateLastLogin(accountID int) error {
	stmt, err := r.db.Prepare("UPDATE Accounts SET Last_Login = ? WHERE Account_ID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), accountID)
	return err
}