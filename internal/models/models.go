package models

import (
	"database/sql"
	"time"
)

// Lender represents the Lenders table
type Lender struct {
	LenderID            int       `json:"lender_id"`
	BusinessName        string    `json:"business_name"`
	PhoneNumber         string    `json:"phone_number"`
	Email               string    `json:"email"`
	InterestRatePercent float64   `json:"interest_rate_percent"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	IsActive            bool      `json:"is_active"` // SQLite stores BOOL as INTEGER, 0 for false, 1 for true
}

// Borrower represents the Borrowers table
type Borrower struct {
	BorrowerID  int            `json:"borrower_id"`
	Fullnames   string         `json:"fullnames"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Residence   sql.NullString `json:"residence"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	IsActive    bool           `json:"is_active"`
}

// Account represents the Accounts table
type Account struct {
	AccountID    int          `json:"account_id"`
	LenderID     int          `json:"lender_id"` // Foreign key to Lenders table
	Username     string       `json:"username"`
	PasswordHash string       `json:"-"` // Do not expose password hash
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	LastLogin    sql.NullTime `json:"last_login"`
	IsLocked     bool         `json:"is_locked"` // SQLite stores BOOL as INTEGER, 0 for false, 1 for true
}

// Plan represents the Plans table
type Plan struct {
	PlanID    int       `json:"plan_id"`
	Plan      string    `json:"plan"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

// LenderLedger represents the Lender_Ledger table
type LenderLedger struct {
	LedgerID  int          `json:"ledger_id"`
	LenderID  int          `json:"lender_id"`
	PlanID    int          `json:"plan_id"`
	Status    string       `json:"status"`
	StartDate time.Time    `json:"start_date"`
	EndDate   sql.NullTime `json:"end_date"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// Loan represents the Loans table
type Loan struct {
	LoanID         int             `json:"loan_id"`
	BorrowerID     int             `json:"borrower_id"`
	LenderID       int             `json:"lender_id"`
	MonthsToPay    int             `json:"months_to_pay"`
	PaymentStatus  string          `json:"payment_status"`
	Amount         float64         `json:"amount"`
	InterestRate   float64         `json:"interest_rate"` // Note: This is an interest rate for the loan, distinct from Lender's base interest rate
	MonthlyPayment sql.NullFloat64 `json:"monthly_payment"`
	StartDate      time.Time       `json:"start_date"`
	EndDate        sql.NullTime    `json:"end_date"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// Receipt represents the Recipets table
type Receipt struct {
	ReceiptID            int            `json:"receipt_id"`
	LoanID               int            `json:"loan_id"`
	Timestamp            time.Time      `json:"timestamp"`
	Status               string         `json:"status"`
	Amount               float64        `json:"amount"`
	PaymentMethod        sql.NullString `json:"payment_method"`
	TransactionReference sql.NullString `json:"transaction_reference"`
	Notes                sql.NullString `json:"notes"`
}

// File represents the File table
type File struct {
	FileID           int            `json:"file_id"`
	LenderID         int            `json:"lender_id"`
	Value            string         `json:"value"`
	FileType         sql.NullString `json:"file_type"`
	FileSize         sql.NullInt64  `json:"file_size"`
	OriginalFilename sql.NullString `json:"original_filename"`
	UploadedAt       time.Time      `json:"uploaded_at"`
}

// Text represents the Text table
type Text struct {
	TextID    int       `json:"text_id"`
	LenderID  int       `json:"lender_id"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Number represents the Number table
type Number struct {
	NumberID  int       `json:"number_id"`
	LenderID  int       `json:"lender_id"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}