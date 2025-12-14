package repository

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"wisetech-lms-api/internal/database"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB initializes an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}

	// Use the schema from internal/database/database.go
	_, err = db.Exec(database.SqliteSchema)
	if err != nil {
		t.Fatalf("Failed to create tables using SqliteSchema: %v", err)
	}

	return db
}

// teardownTestDB closes the database connection.
func teardownTestDB(db *sql.DB) {
	db.Close()
}

func TestCreateLenderAndAccount(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewAuthRepository(db)

	// Test case 1: Successful creation
	businessName := "Lender Business"
	email := "lender@example.com"
	phone := "123-456-7890"
	username := "lenderuser"
	passwordHash := "hashedpassword"
	interestRate := 5.0 // 5%

	accountID, err := repo.CreateLenderAndAccount(businessName, email, phone, username, passwordHash, interestRate)
	if err != nil {
		t.Fatalf("CreateLenderAndAccount failed: %v", err)
	}
	if accountID == 0 {
		t.Error("Expected a non-zero account ID, got 0")
	}

	// Verify account and lender exist
	var retrievedUsername string
	var retrievedLenderID int
	err = db.QueryRow("SELECT Username, Lender_ID FROM Accounts WHERE Account_ID = ?", accountID).Scan(&retrievedUsername, &retrievedLenderID)
	if err != nil {
		t.Fatalf("Failed to retrieve account: %v", err)
	}
	if retrievedUsername != username {
		t.Errorf("Expected username '%s', got '%s'", username, retrievedUsername)
	}
	if retrievedLenderID == 0 {
		t.Error("Expected a non-zero lender ID in account, got 0")
	}

	var retrievedBusinessName string
	var retrievedEmail string
	err = db.QueryRow("SELECT Business_Name, Email FROM Lenders WHERE Lender_ID = ?", retrievedLenderID).Scan(&retrievedBusinessName, &retrievedEmail)
	if err != nil {
		t.Fatalf("Failed to retrieve lender: %v", err)
	}
	if retrievedBusinessName != businessName {
		t.Errorf("Expected business name '%s', got '%s'", businessName, retrievedBusinessName)
	}
	if retrievedEmail != email {
		t.Errorf("Expected email '%s', got '%s'", email, retrievedEmail)
	}

	// Test case 2: Transaction rollback on duplicate username
	// Attempt to create with existing username, which should fail on Account insertion
	_, err = repo.CreateLenderAndAccount("Another Business", "another@example.com", "987-654-3210", username, "anotherhash", 6.0)
	if err == nil {
		t.Fatal("Expected error for duplicate username, got nil")
	}

	// Verify no new lender or account was created beyond the first successful one
	var accountCount int
	err = db.QueryRow("SELECT COUNT(*) FROM Accounts WHERE Username = 'anotheruser'").Scan(&accountCount)
	if err != nil {
		t.Fatalf("Failed to count accounts: %v", err)
	}
	if accountCount != 0 {
		t.Error("Expected no new account 'anotheruser' to be created after failed transaction, but found one")
	}
	var lenderCount int
	err = db.QueryRow("SELECT COUNT(*) FROM Lenders WHERE Email = 'another@example.com'").Scan(&lenderCount)
	if err != nil {
		t.Fatalf("Failed to count lenders: %v", err)
	}
	if lenderCount != 0 {
		t.Error("Expected no new lender to be created after failed transaction, but found one")
	}
}

func TestGetAccountByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewAuthRepository(db)

	// Seed data for a lender and account
	businessName := "Test Lender"
	email := "test@example.com"
	phone := "111-222-3333"
	username := "testuser"
	passwordHash := "hashedpass"
	interestRate := 7.5

	_, err := repo.CreateLenderAndAccount(businessName, email, phone, username, passwordHash, interestRate)
	if err != nil {
		t.Fatalf("Failed to seed lender and account: %v", err)
	}

	// Test case 1: Account found
	account, err := repo.GetAccountByUsername(username)
	if err != nil {
		t.Fatalf("GetAccountByUsername failed: %v", err)
	}
	if account == nil {
		t.Fatal("Expected account, got nil")
	}
	if account.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, account.Username)
	}
	if account.LenderID == 0 {
		t.Error("Expected a non-zero LenderID for the account, got 0")
	}

	// Test case 2: Account not found
	account, err = repo.GetAccountByUsername("nonexistent")
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Expected ErrAccountNotFound, got %v", err)
	}
	if account != nil {
		t.Error("Expected nil account for nonexistent user, got non-nil")
	}
}

func TestGetAccountByID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewAuthRepository(db)

	// Seed data for a lender and account
	businessName := "Test Lender 2"
	email := "test2@example.com"
	phone := "444-555-6666"
	username := "testuser2"
	passwordHash := "hashedpass2"
	interestRate := 8.0

	seededAccountID, err := repo.CreateLenderAndAccount(businessName, email, phone, username, passwordHash, interestRate)
	if err != nil {
		t.Fatalf("Failed to seed lender and account: %v", err)
	}

	// Test case 1: Account found
	account, err := repo.GetAccountByID(seededAccountID)
	if err != nil {
		t.Fatalf("GetAccountByID failed: %v", err)
	}
	if account == nil {
		t.Fatal("Expected account, got nil")
	}
	if account.AccountID != seededAccountID {
		t.Errorf("Expected account ID %d, got %d", seededAccountID, account.AccountID)
	}
	if account.Username != username {
		t.Errorf("Expected username '%s', got '%s'", username, account.Username)
	}

	// Test case 2: Account not found
	account, err = repo.GetAccountByID(99999) // Non-existent ID
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Expected ErrAccountNotFound, got %v", err)
	}
	if account != nil {
		t.Error("Expected nil account for nonexistent ID, got non-nil")
	}
}

func TestGetLenderByAccountID(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewAuthRepository(db)

	// Seed data for a lender and account
	businessName := "Lender Inc."
	email := "lenderinc@example.com"
	phone := "777-888-9999"
	username := "lenderuserinc"
	passwordHash := "hashedpassinc"
	interestRate := 6.5

	seededAccountID, err := repo.CreateLenderAndAccount(businessName, email, phone, username, passwordHash, interestRate)
	if err != nil {
		t.Fatalf("Failed to seed lender and account: %v", err)
	}

	// Test case 1: Lender found for an account
	lender, err := repo.GetLenderByAccountID(seededAccountID)
	if err != nil {
		t.Fatalf("GetLenderByAccountID failed: %v", err)
	}
	if lender == nil {
		t.Fatal("Expected lender, got nil")
	}
	// Verify some fields
	if lender.BusinessName != businessName {
		t.Errorf("Expected business name '%s', got '%s'", businessName, lender.BusinessName)
	}
	if lender.Email != email {
		t.Errorf("Expected email '%s', got '%s'", email, lender.Email)
	}
	if lender.InterestRatePercent != interestRate {
		t.Errorf("Expected interest rate %.1f, got %.1f", interestRate, lender.InterestRatePercent)
	}

	// Test case 2: Account not found for the given account ID
	lender, err = repo.GetLenderByAccountID(99999) // Non-existent account ID
	if !errors.Is(err, ErrAccountNotFound) {
		t.Errorf("Expected ErrAccountNotFound, got %v", err)
	}
	if lender != nil {
		t.Error("Expected nil lender for nonexistent account ID, got non-nil")
	}

	// Test case 3: Account exists but no corresponding lender (this scenario shouldn't happen with FK constraints,
	// but testing robustness)
	// For this test, we need to manually insert an account without a corresponding lender, which violates FK
	// So, a more realistic test would be if a non-lender account exists, it won't return a lender.
	// Since all accounts are linked to lenders via Lender_ID, an account without a lender_id would not exist.
	// We'll rely on the ErrAccountNotFound for non-existent Lender_ID from the join implicitly.
}

func TestUpdateLastLogin(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(db)

	repo := NewAuthRepository(db)

	// Seed data for a lender and account
	businessName := "Updater Lender"
	email := "updater@example.com"
	phone := "000-111-2222"
	username := "updateruser"
	passwordHash := "hashedpassupdate"
	interestRate := 4.0

	seededAccountID, err := repo.CreateLenderAndAccount(businessName, email, phone, username, passwordHash, interestRate)
	if err != nil {
		t.Fatalf("Failed to seed lender and account: %v", err)
	}

	// Initial check: Last_Login should be null
	var lastLogin sql.NullTime
	err = db.QueryRow("SELECT Last_Login FROM Accounts WHERE Account_ID = ?", seededAccountID).Scan(&lastLogin)
	if err != nil {
		t.Fatalf("Failed to query Last_Login: %v", err)
	}
	if lastLogin.Valid {
		t.Error("Expected Last_Login to be NULL initially, but it's valid")
	}

	// Update last login
	err = repo.UpdateLastLogin(seededAccountID)
	if err != nil {
		t.Fatalf("UpdateLastLogin failed: %v", err)
	}

	// Verify update
	err = db.QueryRow("SELECT Last_Login FROM Accounts WHERE Account_ID = ?", seededAccountID).Scan(&lastLogin)
	if err != nil {
		t.Fatalf("Failed to query Last_Login after update: %v", err)
	}
	if !lastLogin.Valid {
		t.Error("Expected Last_Login to be valid after update, but it's NULL")
	}
	// Check if the updated time is recent (within a reasonable margin)
	if time.Since(lastLogin.Time) > 5*time.Second { // Allowing for some processing time
		t.Errorf("Last_Login was not updated to a recent time. Expected within 5s, got %v ago", time.Since(lastLogin.Time))
	}

	// Test updating a non-existent account (should not return an error from the function, but not update anything)
	err = repo.UpdateLastLogin(99999)
	if err != nil {
		t.Errorf("UpdateLastLogin for non-existent account returned unexpected error: %v", err)
	}
	// Verify no error for non-existent account means no record was touched.
	// This is implicit as the function simply returns nil if no rows are affected by the update.
}
