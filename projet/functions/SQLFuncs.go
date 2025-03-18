package functions

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"regexp"
)

var databaseInitialised = false
var db *sql.DB

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
		InfoPrintf("Database initialised\n")
		databaseInitialised = true

		// Debug func call to fill the database with test data
		InitDatabase()
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

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// IsDatabaseInitialised returns whether the database connection is initialised
func IsDatabaseInitialised() bool {
	return databaseInitialised
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

// GetUserRank returns the rank of the user
// 0 = user; 1 = moderator; 2 = admin
func GetUserRank(r *http.Request) int {
	email := GetUserEmail(r)
	checkRights := "SELECT rights_level FROM Moderation WHERE id = (SELECT id FROM users WHERE email = ?)"
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

// CheckIfEmailLinkedToOAuth checks if the email is already linked to an OAuth account
// Returns a boolean and the OAuth provider if the email is linked
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
		var provider string
		err := rows.Scan(&provider)
		if err != nil {
			ErrorPrintf("Error scanning the rows: %v\n", err)
			return false, ""
		}
		return true, provider
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

// CheckPasswordStrength checks if the password is strong enough
func CheckPasswordStrength(password string) bool {
	// Check if the password is at least 8 characters long and at most 64 characters long
	if len(password) < 8 {
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
		ErrorPrintf("Error inserting the user into the database: %v\n", err)
		return err
	}
	return nil
}

// SetSessionCookie sets the session cookie for the user.
// The session cookie is used to store the id of the user.
// The cookie is stored in the session store.
// Returns an error if there is one.
func SetSessionCookie(w http.ResponseWriter, r *http.Request, email string) error {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return err
	}
	session.Values["email"] = email
	err = session.Save(r, w)
	if err != nil {
		ErrorPrintf("Error saving the session: %v\n", err)
		return err
	}
	DebugPrintf("Setting the session cookie for email: %v\n", email)
	return nil
}

// GetSessionCookie returns the session cookie for the user.
// Returns the session cookie and an error if there is one.
func GetSessionCookie(r *http.Request) (string, error) {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return "", err
	}
	email := session.Values["email"].(string)
	return email, nil
}

// EmptySessionCookie empties the session cookie for the user.
// Returns an error if there is one.
func EmptySessionCookie(w http.ResponseWriter, r *http.Request) error {
	session, err := GetSession(r)
	if err != nil {
		ErrorPrintf("Error getting the session: %v\n", err)
		return err
	}
	session.Values["email"] = ""
	err = session.Save(r, w)
	if err != nil {
		ErrorPrintf("Error saving the session: %v\n", err)
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
		(*PageInfo)["IsAdmin"] = false
		(*PageInfo)["IsModerator"] = false
		email := GetUserEmail(r)
		checkRights := "SELECT rights_level FROM Moderation WHERE id = (SELECT id FROM users WHERE email = ?)"
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
	}
	(*PageInfo)["IsAuthenticated"] = false
}

// AddUserToModeration adds the user to the Moderation table if he is not already in it.
// The user is added with the rights level given as parameter.
// Returns an error if there is one.
func AddUserToModeration(email string, rightsLevel int) error {
	checkIfAlreadyInDB := "SELECT id FROM Moderation WHERE id = (SELECT id FROM users WHERE email = ?)"
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
	insertUser := "INSERT INTO Moderation (id, rights_level) VALUES ((SELECT id FROM users WHERE email = ?), ?)"
	_, err = db.Exec(insertUser, email, rightsLevel)
	if err != nil {
		ErrorPrintf("Error inserting the user into the Moderation table: %v\n", err)
		return err
	}
	return nil
}

// InitDatabase initialises the database.
// It creates the tables if they do not exist.
func InitDatabase() {
	UserTableSQL := `CREATE TABLE IF NOT EXISTS users (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	email TEXT NOT NULL UNIQUE,
    	username TEXT NOT NULL UNIQUE,
    	firstname TEXT NOT NULL,
    	lastname TEXT NOT NULL,
    	password_hash TEXT,
    	email_verified BOOLEAN DEFAULT FALSE,
    	oauth_provider TEXT,
    	oauth_id TEXT,
    	creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`

	// the 'Moderation' table only contains the id of the user who has admin/moderator rights
	// the 'rights_level' column is used to determine the rights level of the user
	// 0 = user; 1 = moderator; 2 = admin
	ModerationTableSQL := `CREATE TABLE IF NOT EXISTS Moderation (
    	id INTEGER PRIMARY KEY,
    	rights_level INTEGER DEFAULT 0 NOT NULL,
    	FOREIGN KEY (id) REFERENCES users(id)
		);`

	_, err := db.Exec(UserTableSQL)
	if err != nil {
		ErrorPrintf("Error creating users table: %v\n", err)
		return
	}
	_, err = db.Exec(ModerationTableSQL)
	if err != nil {
		ErrorPrintf("Error creating Moderation table: %v\n", err)
		return
	}
}

// FillDatabase fills the database with test data.
func FillDatabase() {
	// TODO : fill the database with test data for development testing and demonstration purposes
}
