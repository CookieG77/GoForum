package functions

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
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
	PfpID  int
}

// EmailType is a type used to determine the type of the email
type EmailType string

// Constants used to determine the type of the email
const (
	ResetPasswordEmail EmailType = "reset_password" // Email used to reset the password
	VerifyEmailEmail   EmailType = "verify_email"   // Email used to verify the email
)

var EmailTypes = []EmailType{
	ResetPasswordEmail,
	VerifyEmailEmail,
}

// ThreadGoForum is a struct used to represent a thread in the GoForum
type ThreadGoForum struct {
	ThreadID     int
	ThreadName   string
	OwnerID      int
	CreationDate time.Time
}

// ThreadGoForumConfigs is a struct used to represent the configs of a thread in the GoForum
type ThreadGoForumConfigs struct {
	ThreadID                  int
	ThreadDescription         string
	ThreadIconID              int
	ThreadBannerID            int
	IsOpenToNonMembers        bool
	IsOpenToNonConnectedUsers bool
	AllowImages               bool
	AllowLinks                bool
	AllowTextFormatting       bool
}

type MediaType string

// Constants used to determine the type of the email
const (
	UserProfilePicture MediaType = "pfp"           // User profile picture
	ThreadIcon         MediaType = "thread_icon"   // Thread icon
	ThreadBanner       MediaType = "thread_banner" // Thread banner
)

var MediaTypes = []MediaType{
	UserProfilePicture,
	ThreadIcon,
	ThreadBanner,
}

type MediaLinks struct {
	MediaID      int
	MediaType    MediaType
	MediaAddress string
}

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
	checkIfAlreadyInDB := "SELECT email FROM Users WHERE email = ?"
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

// GetUser returns the user
func GetUser(r *http.Request) User {
	email := GetUserEmail(r)
	getUser := "SELECT * FROM Users WHERE email = ?"
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
		err := rows.Scan(
			&user.UserID,
			&user.Email,
			&user.Username,
			&user.Firstname,
			&user.Lastname,
			&user.PasswordHash,
			&user.EmailVerified,
			&user.OAuthProvider,
			&user.OAuthID,
			&user.CreatedAt,
		)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetUser: %v\n", err)
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
	checkRights := "SELECT rights_level FROM Moderation WHERE user_id = (SELECT user_id FROM Users WHERE email = ?)"
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
			ErrorPrintf("Error scanning the rows in GetUserRank: %v\n", err)
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
	getUserConfig := "SELECT * FROM UserConfigs WHERE user_id = (SELECT user_id FROM Users WHERE email = ?)"
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
		err := rows.Scan(&userConfigs.UserID, &userConfigs.Lang, &userConfigs.Theme, &userConfigs.PfpID)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetUserConfig: %v\n", err)
			return UserConfigs{}
		}
		return userConfigs
	}
	return UserConfigs{}
}

// SaveUserConfig saves the user configs
// Returns an error if there is one
func SaveUserConfig(userConfigs UserConfigs) error {
	email := GetEmailFromID(userConfigs.UserID)
	saveUserConfig := "UPDATE UserConfigs SET lang = ?, theme = ?, pfp_id = ? WHERE user_id = (SELECT user_id FROM Users WHERE email = ?)"
	_, err := db.Exec(saveUserConfig, userConfigs.Lang, userConfigs.Theme, userConfigs.PfpID, email)
	if err != nil {
		ErrorPrintf("Error saving the user configs: %v\n", err)
		return err
	}
	return nil
}

// CheckIfEmailLinkedToOAuth checks if the email is already linked to an OAuth account
// Returns true and the OAuth provider as a string if the email is linked to an OAuth provider
// Returns false and an empty string otherwise
func CheckIfEmailLinkedToOAuth(email string) (bool, string) {
	checkIfLinkedToOAuth := "SELECT oauth_provider FROM Users WHERE email = ?"
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
			ErrorPrintf("Error scanning the rows in CheckIfEmailLinkedToOAuth: %v\n", err)
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
	checkIfAlreadyInDB := "SELECT username FROM Users WHERE username = ?"
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
	getUsername := "SELECT username FROM Users WHERE email = ?"
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
			ErrorPrintf("Error scanning the rows in GetUsernameFromEmail: %v\n", err)
			return ""
		}
		return username
	}
	return ""
}

// GetEmailFromUsername returns the email from the username
func GetEmailFromUsername(username string) string {
	getEmail := "SELECT email FROM Users WHERE username = ?"
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
			ErrorPrintf("Error scanning the rows in GetEmailFromUsername: %v\n", err)
			return ""
		}
		return email
	}
	return ""
}

// GetEmailFromID returns the email from the user id
func GetEmailFromID(userID int) string {
	getEmail := "SELECT email FROM Users WHERE user_id = ?"
	rows, err := db.Query(getEmail, userID)
	if err != nil {
		ErrorPrintf("Error getting the email from the user id: %v\n", err)
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
			ErrorPrintf("Error scanning the rows in GetEmailFromID: %v\n", err)
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
	checkEmailVerified := "SELECT email_verified FROM Users WHERE email = ?"
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
			ErrorPrintf("Error scanning the rows in IsUserVerified: %v\n", err)
			return false
		}
		return emailVerified
	}
	return false
}

// VerifyEmail verifies the email of the user.
// Returns an error if there is one.
func VerifyEmail(email string) error {
	verifyEmail := "UPDATE Users SET email_verified = TRUE WHERE email = ?"
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
	checkPassword := "SELECT password_hash FROM Users WHERE email = ?"
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
			ErrorPrintf("Error scanning the rows in CheckUserConnectingWMail: %v\n", err)
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
	checkPassword := "SELECT password_hash FROM Users WHERE username = ?"
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
			ErrorPrintf("Error scanning the rows in CheckUserConnectingWUsername: %v\n", err)
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
// As well as in the 'UserConfigs' table.
// Returns an error if there is one.
func AddUser(email, username, firstname, lastname, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		ErrorPrintf("Error hashing the password: %v\n", err)
		return err
	}
	insertUser := "INSERT INTO Users (email, username, firstname, lastname, password_hash) VALUES (?, ?, ?, ?, ?)"
	_, err = db.Exec(insertUser, email, username, firstname, lastname, hashedPassword)
	if err != nil {
		ErrorPrintf("Error inserting the user into the 'Users' database: %v\n", err)
		return err
	}
	insertUserConfigs := "INSERT INTO UserConfigs (user_id) VALUES ((SELECT user_id FROM Users WHERE email = ?))"
	_, err = db.Exec(insertUserConfigs, email)
	if err != nil {
		ErrorPrintf("Error inserting the user into the 'UserConfigs' table: %v\n", err)
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
	changePassword := "UPDATE Users SET password_hash = ? WHERE email = ?"
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
// Also gives the user his pfp and check if he is verified.
func GiveUserHisRights(PageInfo *map[string]interface{}, r *http.Request) {
	if IsAuthenticated(r) {
		(*PageInfo)["IsAuthenticated"] = true
		(*PageInfo)["IsAddressVerified"] = false
		(*PageInfo)["IsAdmin"] = false
		(*PageInfo)["IsModerator"] = false

		// Check if the user is an admin or a moderator
		email := GetUserEmail(r)
		checkRights := "SELECT rights_level FROM Moderation WHERE user_id = (SELECT user_id FROM Users WHERE email = ?)"
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
				ErrorPrintf("Error scanning the rows in GiveUserHisRights: %v\n", err)
				return
			}
			if rightsLevel == 1 {
				(*PageInfo)["IsModerator"] = true
			} else if rightsLevel == 2 {
				(*PageInfo)["IsAdmin"] = true
			}
		}

		// Check if the email is verified
		checkEmailVerified := "SELECT email_verified FROM Users WHERE email = ?"
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
				ErrorPrintf("Error scanning the rows in GiveUserHisRights: %v\n", err)
				return
			}
			if emailVerified {
				(*PageInfo)["IsAddressVerified"] = true
			}
		}

		// Give the user his Pfp
		(*PageInfo)["UserPfpPath"] = GetMediaLinkFromID(GetUserConfig(r).PfpID).MediaAddress
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
	checkIfAlreadyInDB := "SELECT user_id FROM Moderation WHERE user_id = (SELECT user_id FROM Users WHERE email = ?)"
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
	insertUser := "INSERT INTO Moderation (user_id, rights_level) VALUES ((SELECT user_id FROM Users WHERE email = ?), ?)"
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
	insertEmailIdentification := "INSERT INTO EmailIdentification (email_id, user_id, email_type) VALUES (?, (SELECT user_id FROM Users WHERE email = ?), ?)"
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
	removeEmailIdentification := "DELETE FROM EmailIdentification WHERE user_id = (SELECT user_id FROM Users WHERE email = ?) AND email_type = ?"
	_, err := db.Exec(removeEmailIdentification, email, string(emailType))
	if err != nil {
		ErrorPrintf("Error removing the email identification from the database: %v\n", err)
		return err
	}
	return nil
}

// RemoveEmailIdentificationsOlderThan removes the email identification older than the given time (in minutes) from the database.
// Returns an error if there is one.
func RemoveEmailIdentificationsOlderThan(lifetime int) error {
	removeEmailIdentification := "DELETE FROM EmailIdentification WHERE creation_date < datetime('now', '-%d minutes')"
	_, err := db.Exec(fmt.Sprintf(removeEmailIdentification, lifetime))
	if err != nil {
		ErrorPrintf("Error removing the email identification from the database: %v\n", err)
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
	getEmail := "SELECT email FROM Users WHERE user_id = (SELECT user_id FROM EmailIdentification WHERE email_id = ?)"
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
			ErrorPrintf("Error scanning the rows in GetEmailFromEmailIdentification: %v\n", err)
			return ""
		}
		return email
	}
	return ""
}

// AutoDeleteOldEmailIdentification remove the old email identifications from the database every 1 minutes by default
// To disable it, set the environment variable 'AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS' to 'false'
// To change the interval, set the environment variable 'AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS_INTERVAL' to the desired interval in minutes
// To change the max age of the email identifications, set the environment variable 'EMAIL_IDENTIFICATIONS_MAX_AGE' to the desired max age in minutes
func AutoDeleteOldEmailIdentification() {
	if os.Getenv("AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS") == "false" {
		InfoPrintln("Auto delete old email identifications was disabled")
		return
	}
	interval := 1
	if os.Getenv("AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS_INTERVAL") != "" {
		_, err := fmt.Sscanf(os.Getenv("AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS_INTERVAL"), "%d", &interval)
		if err != nil {
			ErrorPrintf("Error parsing the interval AUTO_DELETE_OLD_EMAIL_IDENTIFICATIONS_INTERVAL : %v\n", err)
			interval = 1
		}
	}
	emailMaxAge := 10
	if os.Getenv("EMAIL_IDENTIFICATIONS_MAX_AGE") != "" {
		_, err := fmt.Sscanf(os.Getenv("EMAIL_IDENTIFICATIONS_MAX_AGE"), "%d", &emailMaxAge)
		if err != nil {
			ErrorPrintf("Error parsing the interval EMAIL_IDENTIFICATIONS_MAX_AGE : %v\n", err)
			emailMaxAge = 1
		}
	}
	InfoPrintf("Auto delete old email identifications interval is set %d minute(s) with a lifetime of %d minute(s)\n", interval, emailMaxAge)
	for {
		err := RemoveEmailIdentificationsOlderThan(emailMaxAge)
		if err != nil {
			ErrorPrintf("Error removing old email identifications: %v\n", err)
			return
		}
		DebugPrintln("Old email identifications removed")
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

// AddThread adds a thread to the database.
// Returns an error if there is one.
func AddThread(threadName string, owner User, description string) error {
	insertThread := "INSERT INTO ThreadGoForum (thread_name, owner_id, creation_date) VALUES (?, ?, ?)"
	_, err := db.Exec(insertThread, threadName, owner.UserID, time.Now())
	if err != nil {
		ErrorPrintf("Error inserting the thread into the database: %v\n", err)
		return err
	}
	insertThreadConfig := "INSERT INTO ThreadGoForumConfigs (thread_id, thread_description) VALUES ((SELECT thread_id FROM ThreadGoForum WHERE thread_name = ?), ?)"
	_, err = db.Exec(insertThreadConfig, threadName, description)
	if err != nil {
		ErrorPrintf("Error inserting the thread config into the database: %v\n", err)
		return err
	}
	return nil
}

// GetThreadFromName returns the ThreadGoForum from the thread name
func GetThreadFromName(threadName string) ThreadGoForum {
	getThread := "SELECT * FROM ThreadGoForum WHERE thread_name = ?"
	rows, err := db.Query(getThread, threadName)
	if err != nil {
		ErrorPrintf("Error getting the thread from the name: %v\n", err)
		return ThreadGoForum{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var thread ThreadGoForum
		err := rows.Scan(&thread.ThreadID, &thread.ThreadName, &thread.OwnerID, &thread.CreationDate)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetThreadFromName: %v\n", err)
			return ThreadGoForum{}
		}
		return thread
	}
	return ThreadGoForum{}
}

// GetAllThreads returns a slice of all the threads in the database
func GetAllThreads() []ThreadGoForum {
	getAllThreads := "SELECT * FROM ThreadGoForum"
	rows, err := db.Query(getAllThreads)
	if err != nil {
		ErrorPrintf("Error getting all the threads: %v\n", err)
		return []ThreadGoForum{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	var threads []ThreadGoForum
	for rows.Next() {
		var thread ThreadGoForum
		err := rows.Scan(&thread.ThreadID, &thread.ThreadName, &thread.OwnerID, &thread.CreationDate)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetAllThreads: %v\n", err)
			return []ThreadGoForum{}
		}
		threads = append(threads, thread)
	}
	return threads
}

// GetThreadConfigsFromID returns the ThreadGoForumConfigs from the thread id
func GetThreadConfigsFromID(threadID int) ThreadGoForumConfigs {
	getThreadConfig := "SELECT * FROM ThreadGoForumConfigs WHERE thread_id = ?"
	rows, err := db.Query(getThreadConfig, threadID)
	if err != nil {
		ErrorPrintf("Error getting the thread config from the id: %v\n", err)
		return ThreadGoForumConfigs{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var threadConfig ThreadGoForumConfigs
		err := rows.Scan(
			&threadConfig.ThreadID,
			&threadConfig.ThreadDescription,
			&threadConfig.ThreadIconID,
			&threadConfig.ThreadBannerID,
			&threadConfig.IsOpenToNonMembers,
			&threadConfig.IsOpenToNonConnectedUsers,
			&threadConfig.AllowImages,
			&threadConfig.AllowLinks,
			&threadConfig.AllowTextFormatting,
		)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetThreadConfigsFromID: %v\n", err)
			return ThreadGoForumConfigs{}
		}
		return threadConfig
	}
	return ThreadGoForumConfigs{}
}

// GetThreadConfigFromThread returns the ThreadGoForumConfigs from the thread
func GetThreadConfigFromThread(thread ThreadGoForum) ThreadGoForumConfigs {
	return GetThreadConfigsFromID(thread.ThreadID)
}

// IsThreadNameValid checks if the thread name is valid
// Thread name must be at least 5 characters long
// Thread name must be at most 50 characters long
// Thread name must only contain letters, numbers, underscores and hyphens
func IsThreadNameValid(threadName string) bool {
	threadNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]{5,50}$`)
	return threadNameRegex.MatchString(threadName)
}

// IsThreadDescriptionValid checks if the thread description is valid
// Thread description must be at least 20 characters long
// Thread description must be at most 500 characters long
// Thread description must only contain letters, numbers, underscores, hyphens, spaces, punctuation and most special characters
func IsThreadDescriptionValid(threadDescription string) bool {
	threadDescriptionRegex := regexp.MustCompile(`^[a-zA-Z0-9 _\-.,;:!?(){}\[\]<>@#$%^&*+=~|\\"'/]{20,500}$`)
	return threadDescriptionRegex.MatchString(threadDescription)
}

// UpdateThreadConfigs updates the thread configs
// Returns an error if there is one
func UpdateThreadConfigs(threadConfigs ThreadGoForumConfigs) error {
	updateThreadConfig := `
		UPDATE ThreadGoForumConfigs SET
			thread_description = ?,
			thread_icon_id = ?,
			thread_banner_id = ?,
			is_open_to_non_members = ?,
			is_open_to_non_connected_users = ?,
			allow_images = ?,
			allow_links = ?,
			allow_text_formatting = ?
		WHERE thread_id = ?
		`
	_, err := db.Exec(updateThreadConfig,
		threadConfigs.ThreadDescription,
		threadConfigs.ThreadIconID,
		threadConfigs.ThreadBannerID,
		threadConfigs.IsOpenToNonMembers,
		threadConfigs.IsOpenToNonConnectedUsers,
		threadConfigs.AllowImages,
		threadConfigs.AllowLinks,
		threadConfigs.AllowTextFormatting,
		threadConfigs.ThreadID)
	if err != nil {
		ErrorPrintf("Error updating the thread configs: %v\n", err)
		return err
	}
	return nil
}

// CheckIfThreadNameExists checks if the thread name is already in the database
func CheckIfThreadNameExists(threadName string) bool {
	checkIfAlreadyInDB := "SELECT thread_name FROM ThreadGoForum WHERE thread_name = ?"
	rows, err := db.Query(checkIfAlreadyInDB, threadName)
	if err != nil {
		ErrorPrintf("Error checking if the thread name is already in the database: %v\n", err)
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

// IsThreadOwner checks if the user is the owner of the thread
func IsThreadOwner(thread ThreadGoForum, r *http.Request) bool {
	email := GetUserEmail(r)
	checkIfOwner := "SELECT thread_name FROM ThreadGoForum WHERE thread_name = ? AND owner_id = (SELECT user_id FROM Users WHERE email = ?)"
	rows, err := db.Query(checkIfOwner, thread.ThreadName, email)
	if err != nil {
		ErrorPrintf("Error checking if the user is the owner of the thread: %v\n", err)
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

// IsThreadMember checks if the user is a member of the thread
func IsThreadMember(thread ThreadGoForum, r *http.Request) bool {
	email := GetUserEmail(r)
	checkIfMember := "SELECT thread_id FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = (SELECT user_id FROM Users WHERE email = ?)"
	rows, err := db.Query(checkIfMember, thread.ThreadID, email)
	if err != nil {
		ErrorPrintf("Error checking if the user is a member of the thread: %v\n", err)
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

// GetThreadMemberRightsLevel returns the rights level of the user in the given thread
func GetThreadMemberRightsLevel(thread ThreadGoForum, r *http.Request) int {
	email := GetUserEmail(r)
	checkIfModerator := "SELECT rights_level FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = (SELECT user_id FROM Users WHERE email = ?)"
	rows, err := db.Query(checkIfModerator, thread.ThreadID, email)
	if err != nil {
		ErrorPrintf("Error getting the user rights level: %v\n", err)
		return 0
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var rightLevel int
		err := rows.Scan(rightLevel)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetThreadMemberRightsLevel: %v\n", err)
			return 0
		}
		return rightLevel
	}
	return 0
}

// IsThreadModerator checks if the user is a moderator of the given thread
func IsThreadModerator(thread ThreadGoForum, r *http.Request) bool {
	rightLevel := GetThreadMemberRightsLevel(thread, r)
	if rightLevel == 1 {
		return true
	}
	return false
}

// IsThreadAdmin checks if the user is an admin of the given thread
func IsThreadAdmin(thread ThreadGoForum, r *http.Request) bool {
	rightLevel := GetThreadMemberRightsLevel(thread, r)
	if rightLevel == 2 {
		return true
	}
	return false
}

// GetMediaLinkFromID returns the media link from the media id
// Returns an empty MediaLinks struct if there is an error
func GetMediaLinkFromID(mediaID int) MediaLinks {
	getMediaLink := "SELECT * FROM MediaLinks WHERE media_id = ?"
	rows, err := db.Query(getMediaLink, mediaID)
	if err != nil {
		ErrorPrintf("Error getting the media link from the id: %v\n", err)
		return MediaLinks{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var mediaLink MediaLinks
		err := rows.Scan(&mediaLink.MediaID, &mediaLink.MediaType, &mediaLink.MediaAddress)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetMediaLinkFromID: %v\n", err)
			return MediaLinks{}
		}
		return mediaLink
	}
	return MediaLinks{}
}

// AddMediaLink adds a media link to the database
// Returns the media id and an error if there is one
func AddMediaLink(mediaType MediaType, mediaAddress string) (int, error) {
	insertMediaLink := "INSERT INTO MediaLinks (media_type, media_address) VALUES (?, ?)"
	res, err := db.Exec(insertMediaLink, string(mediaType), mediaAddress)
	if err != nil {
		ErrorPrintf("Error inserting the media link into the database: %v\n", err)
		return 0, err
	}
	mediaID, err := res.LastInsertId()
	if err != nil {
		ErrorPrintf("Error getting the last insert id: %v\n", err)
		return 0, err
	}
	return int(mediaID), nil
}

// UpdateMediaLink updates the media link in the database
// Returns an error if there is one
func UpdateMediaLink(media MediaLinks) error {
	updateMediaLink := "UPDATE MediaLinks SET media_type = ?, media_address = ? WHERE media_id = ?"
	_, err := db.Exec(updateMediaLink, string(media.MediaType), media.MediaAddress, media.MediaID)
	if err != nil {
		ErrorPrintf("Error updating the media link in the database: %v\n", err)
		return err
	}
	return nil
}

// InitDatabase initialises the database.
// It creates the tables if they do not exist.
func InitDatabase() {
	UserTableTableSQL := `
		CREATE TABLE IF NOT EXISTS Users (
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
	_, err := db.Exec(UserTableTableSQL)
	if err != nil {
		ErrorPrintf("Error creating Users table: %v\n", err)
		return
	}

	UserConfigsTableSQL := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS UserConfigs (
			user_id INTEGER PRIMARY KEY,
			lang TEXT DEFAULT '%s' NOT NULL,
			theme TEXT DEFAULT '%s' NOT NULL,
			pfp_id INTEGER DEFAULT 1 NOT NULL,
			FOREIGN KEY (user_id) REFERENCES Users(id)
		);
		`, string(DefaultLang),
		string(DefaultTheme))
	_, err = db.Exec(UserConfigsTableSQL)
	if err != nil {
		ErrorPrintf("Error creating UserConfigs table: %v\n", err)
		return
	}

	// the 'Moderation' table only contains the id of the user who has admin/moderator rights
	// the 'rights_level' column is used to determine the rights level of the user
	// 0 = user; 1 = moderator; 2 = admin
	ModerationTableSQL := `
		CREATE TABLE IF NOT EXISTS Moderation (
			user_id INTEGER PRIMARY KEY,
			rights_level INTEGER DEFAULT 0 NOT NULL,
			FOREIGN KEY (user_id) REFERENCES Users(user_id)
		);
		`
	_, err = db.Exec(ModerationTableSQL)
	if err != nil {
		ErrorPrintf("Error creating Moderation table: %v\n", err)
		return
	}

	// the 'EmailIdentificationTable' table only contains the id of a user and the id of a link from an email
	// the 'email_id' column is used to determine the email id (it's a unique identifier, it's a 64 characters long hexadecimal string)
	// the 'email_type' column is used to determine the type of the email (it can be 'reset_password', 'verify_email' or 'other')
	// the 'user_id' column is used to determine the user id of the user who has this email (it's a foreign key to the 'Users' table)
	EmailIdentificationTableSQL := `
		CREATE TABLE EmailIdentification (
			email_id TEXT PRIMARY KEY UNIQUE,
			user_id INTEGER NOT NULL,
			email_type TEXT NOT NULL,
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES Users(user_id)
		);
		`
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

	ThreadGoForumTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadGoForum (
		    thread_id INTEGER PRIMARY KEY AUTOINCREMENT,
		    thread_name TEXT NOT NULL UNIQUE,
		    owner_id INTEGER NOT NULL,
		    creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (owner_id) REFERENCES Users(user_id)
		);
		`
	_, err = db.Exec(ThreadGoForumTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadGoForum table: %v\n", err)
		return
	}

	ThreadGoForumConfigsTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadGoForumConfigs (
		    thread_id INTEGER PRIMARY KEY UNIQUE,
		    thread_description TEXT NOT NULL,
		    thread_icon_id INTEGER DEFAULT 2 NOT NULL,
		    thread_banner_id INTEGER DEFAULT 3 NOT NULL,
		    is_open_to_non_members BOOLEAN DEFAULT TRUE NOT NULL,
		    is_open_to_non_connected_Users BOOLEAN DEFAULT TRUE NOT NULL,
		    allow_images BOOLEAN DEFAULT TRUE NOT NULL,
		    allow_links BOOLEAN DEFAULT TRUE NOT NULL,
		    allow_text_formatting BOOLEAN DEFAULT TRUE NOT NULL		    
		);
		`
	_, err = db.Exec(ThreadGoForumConfigsTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadGoForumConfigs table: %v\n", err)
		return
	}

	// The 'ThreadGoForumMembers' represents the members of a thread
	ThreadGoForumMembersTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadGoForumMembers (
			user_id INTEGER NOT NULL,
			thread_id INTEGER NOT NULL,
			rights_level INTEGER DEFAULT 0 NOT NULL,
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, thread_id),
		    FOREIGN KEY (user_id) REFERENCES Users(user_id)
		);
		`
	_, err = db.Exec(ThreadGoForumMembersTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadGoForumMembers table: %v\n", err)
		return
	}

	// The 'MediaLinks' table represents the media links (images, videos, etc.) that are shared in the threads
	// For now, we only will do images as stated in the project instructions
	MediaLinksTableSQL := `
		CREATE TABLE IF NOT EXISTS MediaLinks (
			media_id INTEGER PRIMARY KEY AUTOINCREMENT,
			media_type TEXT NOT NULL,
			media_address TEXT NOT NULL UNIQUE
		);
	`
	_, err = db.Exec(MediaLinksTableSQL)
	if err != nil {
		ErrorPrintf("Error creating MediaLinks table: %v\n", err)
		return
	}

	// Repairing the database just in case
	// If a user doesn't have a row in UserConfigs, we add it
	insertMissingUserConfigs := `
		INSERT INTO UserConfigs (user_id)
		SELECT user_id FROM Users
		WHERE user_id NOT IN (SELECT user_id FROM UserConfigs)
		`
	_, err = db.Exec(insertMissingUserConfigs)
	if err != nil {
		ErrorPrintf("Error inserting missing user configs: %v\n", err)
		return
	}

	// Add the default media links (example: default pfp, default thread icon, etc.)
	query := "SELECT COUNT(*) FROM MediaLinks"
	var count int
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error checking if MediaLinks table is empty: %v\n", err)
	}
	// If the table is empty, we insert the default media links
	if count == 0 {
		insertDefaultMediaLinks := fmt.Sprintf(`
		INSERT INTO MediaLinks (media_type, media_address) VALUES
			('%s', '/img/default_user_icon.png'),
			('%s', '/img/default_thread_icon.png'),
			('%s', '/img/default_thread_banner.gif');
		`,
			string(UserProfilePicture),
			string(ThreadIcon),
			string(ThreadBanner),
		)

		_, err = db.Exec(insertDefaultMediaLinks)
		if err != nil {
			ErrorPrintf("Error inserting default media links: %v\n", err)
			return
		}
	}

	// Starting the auto delete of the old email identifications
	go AutoDeleteOldEmailIdentification()

	InfoPrintln("Database initialised")
}

// FillDatabase fills the database with test data.
func FillDatabase() {
	// TODO : fill the database with test data for development testing and demonstration purposes

	// A test thread
	err := AddThread("TestThread", User{UserID: 1}, "This is a test thread ! :P  (o_o)")
	if err != nil {
		ErrorPrintf("Error adding thread TestThread: %v\n", err)
		return
	}

	// A test thread with must be connected
	err = AddThread("TestThread2", User{UserID: 1}, "This is an other test thread where you must be connected ! (►__◄)")
	if err != nil {
		ErrorPrintf("Error adding thread TestThread2: %v\n", err)
		return
	}
	TestThread2Configs := GetThreadConfigFromThread(GetThreadFromName("TestThread2"))
	TestThread2Configs.IsOpenToNonConnectedUsers = false
	err = UpdateThreadConfigs(TestThread2Configs)
	if err != nil {
		ErrorPrintf("Error updating thread TestThread2: %v\n", err)
		return
	}

	// A test thread with must be a member
	err = AddThread("TestThread3", User{UserID: 1}, "This is also an other test thread where you must be a member ! (◕‿-)")
	if err != nil {
		ErrorPrintf("Error adding thread TestThread3: %v\n", err)
		return
	}
	TestThread3Configs := GetThreadConfigFromThread(GetThreadFromName("TestThread3"))
	TestThread3Configs.IsOpenToNonMembers = false
	err = UpdateThreadConfigs(TestThread3Configs)
	if err != nil {
		ErrorPrintf("Error updating thread TestThread3 configs: %v\n", err)
		return
	}
}
