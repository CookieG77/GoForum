package functions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
	"time"
)

var databaseInitialised = false
var db *sql.DB

// User is a struct used to represent a user
type User struct {
	UserID        int
	Email         string
	Username      string
	Firstname     string
	Lastname      string
	PasswordHash  sql.NullString
	EmailVerified bool
	OAuthProvider sql.NullString
	OAuthID       sql.NullString
	CreatedAt     time.Time
}

// UserConfigs is a struct used to represent the user configs
type UserConfigs struct {
	UserID int
	Lang   string
	Theme  string
}

// EmailType is a type used to determine the type of the email
type EmailType string

// Constants used to determine the type of the email
const (
	ResetPasswordEmail EmailType = "reset_password" // Email used to reset the password
	VerifyEmailEmail   EmailType = "verify_email"   // Email used to verify the email
)

// InitDatabaseConnection initialises the database connection
func InitDatabaseConnection() {
	if !databaseInitialised {
		if os.Getenv("DB_URL") == "" {
			ErrorPrintf("DB_URL environment variable not set\n")
			return
		}
		testDB, err := sql.Open(os.Getenv("DB_URL"), os.Getenv("DB_NAME"))
		if err != nil {
			ErrorPrintf("Error opening database: %v\n", err)
			return
		}
		err = testDB.Ping()
		if err != nil {
			ErrorPrintf("Error pinging database: %v\n", err)
			return
		}
		db = testDB
		InfoPrintf("Database connection initialised\n")
		databaseInitialised = true

		// Initialise the database (create the tables if they do not exist and repair the database if needed)
		InitDatabase()

		// Debug func call to fill the database with test data
		FillDatabase()
	}
}

// CloseDatabase closes the database connection
func CloseDatabase() {
	if databaseInitialised {
		InfoPrintf("Database closed\n")
		err := db.Close()
		if err != nil {
			ErrorPrintf("Error closing database: %v\n", err)
			return
		}
		databaseInitialised = false
	}
}

// IsEmailValid checks if the email is valid
func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// CheckIfEmailExists checks if the email is already in the database
func CheckIfEmailExists(email string) bool {
	checkIfAlreadyInDB := "SELECT email FROM users WHERE email = ?"
	rows, err := db.Query(checkIfAlreadyInDB, email)
	if err != nil {
		ErrorPrintf("Error checking if the email is already in the database: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		return true
	}
	return false
}

// GetUserEmail returns the email of the user
func GetUserEmail(r *http.Request) string {
	email, err := GetSessionCookie(r)
	if err != nil {
		ErrorPrintf("Error getting the user email: %v\n", err)
		return ""
	}
	return email
}

// GetUser returns the user configs
func GetUser(r *http.Request) User {
	email := GetUserEmail(r)
	getUser := "SELECT * FROM users WHERE email = ?"
	rows, err := db.Query(getUser, email)
	if err != nil {
		ErrorPrintf("Error getting the user: %v\n", err)
		return User{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var user User
		err := rows.Scan(&user.UserID, &user.Email, &user.Username, &user.Firstname, &user.Lastname, &user.PasswordHash, &user.EmailVerified, &user.OAuthProvider, &user.OAuthID, &user.CreatedAt)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return User{}
		}
		return user
	}
	return User{}
}

// GetUserRank returns the rank of the user
// 0 = user; 1 = moderator; 2 = admin
func GetUserRank(r *http.Request) int {
	email := GetUserEmail(r)
	checkRights := "SELECT rights_level FROM Moderation WHERE user_id = (SELECT user_id FROM users WHERE email = ?)"
	rows, err := db.Query(checkRights, email)
	if err != nil {
		ErrorPrintf("Error checking the user rights: %v\n", err)
		return 0
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var rightsLevel int
		err := rows.Scan(&rightsLevel)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return 0
		}
		return rightsLevel
	}
	return 0
}

// GetUserRankString returns the rank of the user as a string
func GetUserRankString(r *http.Request) string {
	switch GetUserRank(r) {
	case 1:
		return "moderator"
	case 2:
		return "admin"
	default:
		return "user"
	}
}

// GetUserConfig returns the user configs
func GetUserConfig(r *http.Request) UserConfigs {
	email := GetUserEmail(r)
	getUserConfig := "SELECT * FROM user_configs WHERE user_id = (SELECT user_id FROM users WHERE email = ?)"
	rows, err := db.Query(getUserConfig, email)
	if err != nil {
		ErrorPrintf("Error getting the user configs: %v\n", err)
		return UserConfigs{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var userConfigs UserConfigs
		err := rows.Scan(&userConfigs.UserID, &userConfigs.Lang, &userConfigs.Theme)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return UserConfigs{}
		}
		return userConfigs
	}
	return UserConfigs{}
}

// CheckIfEmailLinkedToOAuth checks if the email is already linked to an OAuth account
// Returns true and the OAuth provider as a string if the email is linked to an OAuth provider
// Returns false and an empty string otherwise
func CheckIfEmailLinkedToOAuth(email string) (bool, string) {
	checkIfLinkedToOAuth := "SELECT oauth_provider FROM users WHERE email = ?"
	rows, err := db.Query(checkIfLinkedToOAuth, email)
	if err != nil {
		ErrorPrintf("Error checking if the email is linked to an OAuth account: %v\n", err)
		return false, ""
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var provider sql.NullString
		err := rows.Scan(&provider)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return false, ""
		}
		if provider.Valid {
			return true, provider.String
		}
		return false, ""
	}
	return false, ""
}

// CheckIfUsernameExists checks if the username is already in the database
func CheckIfUsernameExists(username string) bool {
	checkIfAlreadyInDB := "SELECT username FROM users WHERE username = ?"
	rows, err := db.Query(checkIfAlreadyInDB, username)
	if err != nil {
		ErrorPrintf("Error checking if the username is already in the database: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		return true
	}
	return false
}

// IsUsernameValid checks if the username is valid
// Username must be at least 3 characters long
// Username must be at most 20 characters long
// Username must only contain letters, numbers, underscores and hyphens
func IsUsernameValid(username string) bool {
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{3,20}$`)
	return usernameRegex.MatchString(username)
}

// GetUsernameFromEmail returns the username from the email
func GetUsernameFromEmail(email string) string {
	getUsername := "SELECT username FROM users WHERE email = ?"
	rows, err := db.Query(getUsername, email)
	if err != nil {
		ErrorPrintf("Error getting the username from the email: %v\n", err)
		return ""
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return ""
		}
		return username
	}
	return ""
}

// GetEmailFromUsername returns the email from the username
func GetEmailFromUsername(username string) string {
	getEmail := "SELECT email FROM users WHERE username = ?"
	rows, err := db.Query(getEmail, username)
	if err != nil {
		ErrorPrintf("Error getting the email from the username: %v\n", err)
		return ""
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return ""
		}
		return email
	}
	return ""
}

// IsUserVerified checks if the user is verified, i.e. if the email is verified.
// Returns true if the user is verified and false otherwise.
func IsUserVerified(r *http.Request) bool {
	email := GetUserEmail(r)
	checkEmailVerified := "SELECT email_verified FROM users WHERE email = ?"
	rows, err := db.Query(checkEmailVerified, email)
	if err != nil {
		ErrorPrintf("Error checking if the email is verified: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var emailVerified bool
		err := rows.Scan(&emailVerified)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return false
		}
		return emailVerified
	}
	return false
}

// VerifyEmail verifies the email of the user.
// Returns an error if there is one.
func VerifyEmail(email string) error {
	verifyEmail := "UPDATE users SET email_verified = TRUE WHERE email = ?"
	_, err := db.Exec(verifyEmail, email)
	if err != nil {
		ErrorPrintf("Error verifying the email: %v\n", err)
		return err
	}
	return nil
}

// CheckPasswordStrength checks if the password is strong enough
func CheckPasswordStrength(password string) bool {
	// Check if the password is at least 8 characters long and at most 64 characters long
	if len(password) < 8 {
		return false
	}
	if len(password) > 64 {
		return false
	}
	// Check if the password contains at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		return false
	}
	// Check if the password contains at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		return false
	}
	// Check if the password contains at least one digit
	digitRegex := regexp.MustCompile(`[0-9]`)
	if !digitRegex.MatchString(password) {
		return false
	}
	return true
}

// HashPassword hashes the password using bcrypt.
// Returns the hashed password and an error if there is one.
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedBytes), err
}

// CheckPasswordHash checks if the password matches the hashed password.
// Returns true if the password matches the hashed password and false otherwise.
func checkPasswordHash(password, hashedPassword string) bool {
	/*hash, err := HashPassword(password)
	if err != nil {
		ErrorPrintf("Error hashing the password: %v\n", err)
		return false
	}
	DebugPrintf("password = %s\n\thashedPassword ----- = %s\n\tstoredhashedPassword = %s", password, hash, hashedPassword)*/
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CheckUserConnectingWMail checks if the values are correct when connecting using email.
// Returns true if the values are correct and false otherwise.
// Returns an error if there is one.
func CheckUserConnectingWMail(email, password string) (bool, error) {
	if !CheckIfEmailExists(email) {
		return false, nil
	}
	checkPassword := "SELECT password_hash FROM users WHERE email = ?"
	rows, err := db.Query(checkPassword, email)
	if err != nil {
		ErrorPrintf("Error checking the password: %v\n", err)
		return false, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var hashedPassword string
		err := rows.Scan(&hashedPassword)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return false, err
		}
		return checkPasswordHash(password, hashedPassword), nil
	}
	return false, nil
}

// CheckUserConnectingWUsername checks if the values are correct when connecting using the username.
// Returns true if the values are correct and false otherwise.
// Returns an error if there is one.
func CheckUserConnectingWUsername(username, password string) (bool, error) {
	if !CheckIfUsernameExists(username) {
		return false, nil
	}
	checkPassword := "SELECT password_hash FROM users WHERE username = ?"
	rows, err := db.Query(checkPassword, username)
	if err != nil {
		ErrorPrintf("Error checking the password: %v\n", err)
		return false, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var hashedPassword string
		err := rows.Scan(&hashedPassword)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return false, err
		}
		return checkPasswordHash(password, hashedPassword), nil
	}
	return false, nil
}

// GetConnectionMethod returns the connection method used by the user.
// Returns "email" if the user connected with email, "oauth" if the user connected with OAuth and "username" if the user connected with username.
// Also returns the provider if the user connected with OAuth (empty string otherwise).
// Returns an empty string if the connection method is not valid.
func GetConnectionMethod(emailOrUsername string) (string, string) {
	if b, provider := CheckIfEmailLinkedToOAuth(emailOrUsername); b {
		return "oauth", provider
	}
	if CheckIfEmailExists(emailOrUsername) {
		return "email", ""
	}
	if CheckIfUsernameExists(emailOrUsername) {
		return "username", ""
	}
	return "", ""
}

// AddUser adds a user to the database.
// As well as in the 'user_configs' table.
// Returns an error if there is one.
func AddUser(email, username, firstname, lastname, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		ErrorPrintf("Error hashing the password: %v\n", err)
		return err
	}
	insertUser := "INSERT INTO users (email, username, firstname, lastname, password_hash) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(insertUser, email, username, firstname, lastname, hashedPassword)
	if err != nil {
		ErrorPrintf("Error inserting the user into the 'users' database: %v\n", err)
		return err
	}
	insertUserConfigs := "INSERT INTO user_configs (user_id) VALUES ((SELECT user_id FROM users WHERE email = ?))"
	_, err = db.Exec(insertUserConfigs, email)
	if err != nil {
		ErrorPrintf("Error inserting the user into the 'user_configs' table: %v\n", err)
		return err
	}
	return nil
}

// ChangeUserPassword changes the password of the user with the given mail with the given password.
// Returns an error if there is one.
func ChangeUserPassword(userMail, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		ErrorPrintf("Error hashing the password: %v\n", err)
		return err
	}
	changePassword := "UPDATE users SET password_hash = ? WHERE email = ?"
	_, err = db.Exec(changePassword, hashedPassword, userMail)
	if err != nil {
		ErrorPrintf("Error changing the user password: %v\n", err)
		return err
	}
	return nil
}

// IsAuthenticated checks if the user is authenticated.
// Returns true if the user is authenticated and false otherwise.
func IsAuthenticated(r *http.Request) bool {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return false
	}
	if session.Values["email"] == nil {
		return false
	}
	return CheckIfEmailExists(session.Values["email"].(string))
}

// GiveUserHisRights gives the user his admin/moderator rights.
func GiveUserHisRights(PageInfo *map[string]interface{}, r *http.Request) {
	if IsAuthenticated(r) {
		(*PageInfo)["IsAuthenticated"] = true
		(*PageInfo)["IsAddressVerified"] = false
		(*PageInfo)["IsAdmin"] = false
		(*PageInfo)["IsModerator"] = false

		// Check if the user is an admin or a moderator
		email := GetUserEmail(r)
		checkRights := "SELECT rights_level FROM Moderation WHERE user_id = (SELECT user_id FROM users WHERE email = ?)"
		rows, err := db.Query(checkRights, email)
		if err != nil {
			ErrorPrintf("Error checking the user rights: %v\n", err)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				ErrorPrintf("Error closing the rows: %v\n", err)
			}
		}(rows)
		if rows.Next() {
			var rightsLevel int
			err := rows.Scan(&rightsLevel)
			if err != nil {
				ErrorPrintf("Error scanning the rows: %v\n", err)
				return
			}
			if rightsLevel == 1 {
				(*PageInfo)["IsModerator"] = true
			} else if rightsLevel == 2 {
				(*PageInfo)["IsAdmin"] = true
			}
		}

		// Check if the email is verified
		checkEmailVerified := "SELECT email_verified FROM users WHERE email = ?"
		rows, err = db.Query(checkEmailVerified, email)
		if err != nil {
			ErrorPrintf("Error checking if the email is verified: %v\n", err)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				ErrorPrintf("Error closing the rows: %v\n", err)
			}
		}(rows)
		if rows.Next() {
			var emailVerified bool
			err := rows.Scan(&emailVerified)
			if err != nil {
				ErrorPrintf("Error scanning the rows: %v\n", err)
				return
			}
			if emailVerified {
				(*PageInfo)["IsAddressVerified"] = true
				DebugPrintln("Email verified is true")
			}
		}
		return
	}
	(*PageInfo)["IsAuthenticated"] = false
	(*PageInfo)["IsAddressVerified"] = false
}

// AddUserToModeration adds the user to the Moderation table if he is not already in it.
// The user is added with the rights level given as parameter.
// Returns an error if there is one.
func AddUserToModeration(email string, rightsLevel int) error {
	InfoPrintf("Adding user %s with email %s to moderation: %d\n", GetUsernameFromEmail(email), email, rightsLevel)
	checkIfAlreadyInDB := "SELECT user_id FROM Moderation WHERE user_id = (SELECT user_id FROM users WHERE email = ?)"
	rows, err := db.Query(checkIfAlreadyInDB, email)
	if err != nil {
		ErrorPrintf("Error checking if the user is already in the Moderation table: %v\n", err)
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		DebugPrintf("User %s is already in the Moderation table\n", email)
		return nil
	}
	insertUser := "INSERT INTO Moderation (user_id, rights_level) VALUES ((SELECT user_id FROM users WHERE email = ?), ?)"
	_, err = db.Exec(insertUser, email, rightsLevel)
	if err != nil {
		ErrorPrintf("Error inserting the user into the Moderation table: %v\n", err)
		return err
	}
	return nil
}

// RandomHexString generates a random hexadecimal string of length n.
// Returns the random hexadecimal string and an error if there is one.
func RandomHexString(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("la longueur doit être supérieure à zéro")
	}

	bytes := make([]byte, (n+1)/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	hexStr := hex.EncodeToString(bytes)
	return hexStr[:n], nil // trim any excess
}

// CreateEmailIdentificationLink creates a link between an email and an email id.
// Returns the email id and an error if there is one.
func CreateEmailIdentificationLink(email string, emailType EmailType) (string, error) {
	emailID, err := RandomHexString(64)
	if err != nil {
		ErrorPrintf("Error generating the email id: %v\n", err)
		return "", err
	}
	insertEmailIdentification := "INSERT INTO EmailIdentification (email_id, user_id, email_type) VALUES (?, (SELECT user_id FROM users WHERE email = ?), ?)"
	_, err = db.Exec(insertEmailIdentification, emailID, email, string(emailType))
	if err != nil {
		ErrorPrintf("Error inserting the email identification into the database: %v\n", err)
		return "", err
	}
	return emailID, nil
}

// RemoveEmailIdentificationWithID removes the email identification with the given user_id from the database.
// Returns an error if there is one.
func RemoveEmailIdentificationWithID(emailID string) error {
	removeEmailIdentification := "DELETE FROM EmailIdentification WHERE email_id = ?"
	_, err := db.Exec(removeEmailIdentification, emailID)
	if err != nil {
		ErrorPrintf("Error removing the email identification from the database: %v\n", err)
		return err
	}
	return nil
}

// RemoveEmailIdentificationForUser removes the email identification for the user with the given email and type from the database.
// Returns an error if there is one.
func RemoveEmailIdentificationForUser(email string, emailType EmailType) error {
	removeEmailIdentification := "DELETE FROM EmailIdentification WHERE user_id = (SELECT user_id FROM users WHERE email = ?) AND email_type = ?"
	_, err := db.Exec(removeEmailIdentification, email, string(emailType))
	if err != nil {
		ErrorPrintf("Error removing the email identification from the database: %v\n", err)
		return err
	}
	return nil
}

// RemoveOldEmailIdentifications removes the old email identifications from the database.
// Returns an error if there is one.
func RemoveOldEmailIdentifications() error {
	removeOldEmailIdentifications := "DELETE FROM EmailIdentification WHERE creation_date < datetime('now', '-1 day')"
	_, err := db.Exec(removeOldEmailIdentifications)
	if err != nil {
		ErrorPrintf("Error removing the old email identifications from the database: %v\n", err)
		return err
	}
	return nil
}

// CheckEmailIdentification checks if the email id with the given id and type is in the database.
// Returns true if the email id is in the database and false otherwise.
func CheckEmailIdentification(emailID string, emailType EmailType) bool {
	checkEmailIdentification := "SELECT email_id FROM EmailIdentification WHERE email_id = ? AND email_type = ?"
	rows, err := db.Query(checkEmailIdentification, emailID, string(emailType))
	if err != nil {
		ErrorPrintf("Error checking the email identification: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		return true
	}
	return false
}

// GetEmailFromEmailIdentification returns the email from the email id.
func GetEmailFromEmailIdentification(emailID string) string {
	getEmail := "SELECT email FROM users WHERE user_id = (SELECT user_id FROM EmailIdentification WHERE email_id = ?)"
	rows, err := db.Query(getEmail, emailID)
	if err != nil {
		ErrorPrintf("Error getting the email from the email identification: %v\n", err)
		return ""
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return ""
		}
		return email
	}
	return ""
}

// InitDatabase initialises the database.
// It creates the tables if they do not exist.
func InitDatabase() {
	UserTableSQL := `
		CREATE TABLE IF NOT EXISTS users (
    	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    	email TEXT NOT NULL UNIQUE,
    	username TEXT NOT NULL UNIQUE,
    	firstname TEXT NOT NULL,
    	lastname TEXT NOT NULL,
    	password_hash TEXT,
    	email_verified BOOLEAN DEFAULT FALSE,
    	oauth_provider TEXT,
    	oauth_id TEXT,
    	creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
		`

	UserConfigsSQL := `
		CREATE TABLE IF NOT EXISTS user_configs (
    	user_id INTEGER PRIMARY KEY,
    	lang TEXT DEFAULT 'en' NOT NULL,
    	style TEXT DEFAULT 'light' NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES users(id)
		);
		`

	// the 'Moderation' table only contains the id of the user who has admin/moderator rights
	// the 'rights_level' column is used to determine the rights level of the user
	// 0 = user; 1 = moderator; 2 = admin
	ModerationTableSQL := `
		CREATE TABLE IF NOT EXISTS Moderation (
    	user_id INTEGER PRIMARY KEY,
    	rights_level INTEGER DEFAULT 0 NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
		`

	// the 'EmailIdentificationTable' table only contains the id of a user and the id of a link from an email
	// the 'email_id' column is used to determine the email id (it's a unique identifier, it's a 64 characters long hexadecimal string)
	// the 'email_type' column is used to determine the type of the email (it can be 'reset_password', 'verify_email' or 'other')
	// the 'user_id' column is used to determine the user id of the user who has this email (it's a foreign key to the 'users' table)
	EmailIdentificationTableSQL := `
		CREATE TABLE EmailIdentification (
    	email_id TEXT PRIMARY KEY UNIQUE,
    	user_id INTEGER NOT NULL,
    	email_type TEXT NOT NULL,
    	creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    	FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
		`

	_, err := db.Exec(UserTableSQL)
	if err != nil {
		ErrorPrintf("Error creating users table: %v\n", err)
		return
	}
	_, err = db.Exec(UserConfigsSQL)
	if err != nil {
		ErrorPrintf("Error creating user_configs table: %v\n", err)
		return
	}
	_, err = db.Exec(ModerationTableSQL)
	if err != nil {
		ErrorPrintf("Error creating Moderation table: %v\n", err)
		return
	}

	// Before creating the EmailIdentification table, we want to be sure that if it already exists, it is empty
	_, err = db.Exec("DROP TABLE IF EXISTS EmailIdentification")
	if err != nil {
		ErrorPrintf("Error dropping EmailIdentification table: %v\n", err)
		return
	}
	_, err = db.Exec(EmailIdentificationTableSQL)
	if err != nil {
		ErrorPrintf("Error creating EmailIdentification table: %v\n", err)
		return
	}

	// Repairing the database just in case
	// If a user doesn't have a row in user_configs, we add it
	insertMissingUserConfigs := `
		INSERT INTO user_configs (user_id)
		SELECT user_id FROM users
		WHERE user_id NOT IN (SELECT user_id FROM user_configs)
		`
	_, err = db.Exec(insertMissingUserConfigs)
	if err != nil {
		ErrorPrintf("Error inserting missing user configs: %v\n", err)
		return
	}

	// Remove the old email identifications from the database
	err = RemoveOldEmailIdentifications()
	if err != nil {
		ErrorPrintf("Error removing old email identifications: %v\n", err)
		return
	}

	InfoPrintln("Database initialised")
}

// FillDatabase fills the database with test data.
func FillDatabase() {
	// TODO : fill the database with test data for development testing and demonstration purposes
}
