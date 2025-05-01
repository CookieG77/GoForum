package functions

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	mr "math/rand"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
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

type SimplifiedUser struct {
	Username    string
	PfpAddress  string
	RightsLevel int
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
	UserProfilePicture   MediaType = "pfp"             // User profile picture
	ThreadIcon           MediaType = "thread_icon"     // Thread icon
	ThreadBanner         MediaType = "thread_banner"   // Thread banner
	ThreadMessagePicture MediaType = "message_picture" // Thread message picture
)

var MediaTypes = []MediaType{
	UserProfilePicture,
	ThreadIcon,
	ThreadBanner,
	ThreadMessagePicture,
}

type MediaLink struct {
	MediaID      int
	MediaType    MediaType
	MediaAddress string
	CreationDate time.Time
}

type FormattedThreadMessage struct {
	MessageID      int       `json:"message_id"`
	MessageTitle   string    `json:"message_title"`
	MessageContent string    `json:"message_content"`
	WasEdited      bool      `json:"was_edited"`
	CreationDate   time.Time `json:"creation_date"`
	UserName       string    `json:"user_name"`
	UserPfpAddress string    `json:"user_pfp_address"`
	Upvotes        int       `json:"up_votes"`
	Downvotes      int       `json:"down_votes"`
	MediaLinks     []string  `json:"media_links"`
	MessageTags    []string  `json:"message_tags"`
	VoteState      int       `json:"vote_state"`
}

var OrderingList = []string{"asc", "desc", "popular", "unpopular"}

type ThreadTag struct {
	TagID    int
	ThreadID int
	TagName  string
	TagColor string
}

// PossibleMessageOrderingList is a list of possible message ordering
var PossibleMessageOrderingList = []string{"asc", "desc", "popular", "unpopular"}

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
// Must be authenticated to get the user
func GetUserEmail(r *http.Request) string {
	session, err := GetSessionCookie(r)
	if err != nil {
		ErrorPrintf("Error getting the user email: %v\n", err)
		return ""
	}
	email := session.Values["email"]
	if email == nil {
		return ""
	}
	return email.(string)
}

// GetUser returns the user
// Must be authenticated to get the user
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

// UpdateUserConfig saves the user configs
// Returns an error if there is one
func UpdateUserConfig(userConfigs UserConfigs) error {
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

// GetUserFromUsername returns the email from the username
func GetUserFromUsername(username string) (User, error) {
	getEmail := "SELECT * FROM Users WHERE username = ?"
	rows, err := db.Query(getEmail, username)
	if err != nil {
		ErrorPrintf("Error getting the email from the username: %v\n", err)
		return User{}, err
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
			ErrorPrintf("Error scanning the rows in GetUserFromUsername: %v\n", err)
			return User{}, err
		}
		return user, nil
	}
	return User{}, fmt.Errorf("user not found")
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
	(*PageInfo)["IsAuthenticated"] = false
	(*PageInfo)["IsAddressVerified"] = false
	if IsAuthenticated(r) {
		(*PageInfo)["IsAuthenticated"] = true

		// Check if the user is an admin or a moderator
		user := GetUser(r)

		// Check if the email is verified
		checkEmailVerified := "SELECT email_verified FROM Users WHERE user_id = ?"
		rows, err := db.Query(checkEmailVerified, user.UserID)
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
		userConfig := GetUserConfig(r)
		(*PageInfo)["UserPfpPath"] = GetMediaLinkFromID(userConfig.PfpID).MediaAddress
		return
	}
}

// CreateEmailIdentificationLink creates a link between an email and an email id.
// Returns the email id and an error if there is one.
func CreateEmailIdentificationLink(email string, emailType EmailType) (string, error) {
	emailID := uuid.New().String() // Generate a new UUID for the email id for it to be unique
	insertEmailIdentification := "INSERT INTO EmailIdentification (email_id, user_id, email_type) VALUES (?, (SELECT user_id FROM Users WHERE email = ?), ?)"
	_, err := db.Exec(insertEmailIdentification, emailID, email, string(emailType))
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
func AddThread(owner User, threadName string, description string) error {
	insertThread := "INSERT INTO ThreadGoForum (thread_name, owner_id, creation_date) VALUES (?, ?, ?)"
	_, err := db.Exec(insertThread, threadName, owner.UserID, time.Now())
	if err != nil {
		ErrorPrintf("Error inserting the thread into the database: %v\n", err)
		return err
	}
	// Insert the thread config into the ThreadGoForumConfigs table
	insertThreadConfig := "INSERT INTO ThreadGoForumConfigs (thread_id, thread_description) VALUES ((SELECT thread_id FROM ThreadGoForum WHERE thread_name = ?), ?)"
	_, err = db.Exec(insertThreadConfig, threadName, description)
	if err != nil {
		ErrorPrintf("Error inserting the thread config into the database: %v\n", err)
		return err
	}
	// Insert the thread owner into the ThreadGoForumMembers table
	insertThreadOwner := "INSERT INTO ThreadGoForumMembers (thread_id, user_id, rights_level) VALUES ((SELECT thread_id FROM ThreadGoForum WHERE thread_name = ?), ?, 3)"
	_, err = db.Exec(insertThreadOwner, threadName, owner.UserID)
	if err != nil {
		ErrorPrintf("Error inserting the thread owner into the database: %v\n", err)
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

// AlsoIsThreadOwner checks if the user is the owner of the thread (variant checking of ownership)
func AlsoIsThreadOwner(thread ThreadGoForum, user User) bool {
	checkIfOwner := "SELECT thread_name FROM ThreadGoForum WHERE thread_name = ? AND owner_id = ?"
	rows, err := db.Query(checkIfOwner, thread.ThreadName, user.UserID)
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
	user := GetUser(r)
	return IsUserInThread(thread, user)
}

// IsUserInThread checks if the user is in the thread
func IsUserInThread(thread ThreadGoForum, user User) bool {
	checkIfMember := "SELECT thread_id FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = (SELECT user_id FROM Users WHERE email = ?)"
	rows, err := db.Query(checkIfMember, thread.ThreadID, user.Email)
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

// IsUserBannedFromThread checks if the user is banned from the thread
func IsUserBannedFromThread(thread ThreadGoForum, user User) bool {
	checkifBanned := `SELECT rights_level FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = ?`
	rows, err := db.Query(checkifBanned, thread.ThreadID, user.UserID)
	if err != nil {
		ErrorPrintf("Error checking if the user is banned from the thread: %v\n", err)
		return true // Assume the user is banned if there is an error
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
			ErrorPrintf("Error scanning the rows in IsUserBannedFromThread: %v\n", err)
			return true // Assume the user is banned if there is an error
		}
		if rightsLevel >= 0 {
			return false
		}
	}
	return true
}

// GetThreadMemberRightsLevel returns the rights level of the user in the given thread
func GetThreadMemberRightsLevel(thread ThreadGoForum, user User) int {
	checkIfModerator := "SELECT rights_level FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = ?"
	rows, err := db.Query(checkIfModerator, thread.ThreadID, user.UserID)
	if err != nil {
		ErrorPrintf("Error checking if the user is a member of the thread: %v\n", err)
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
			ErrorPrintf("Error scanning the rows in GetThreadMemberRightsLevel: %v\n", err)
			return 0
		}
		return rightsLevel
	}
	return 0
}

// IsThreadModerator checks if the user is a moderator of the given thread
func IsThreadModerator(thread ThreadGoForum, user User) bool {
	rightLevel := GetThreadMemberRightsLevel(thread, user)
	if rightLevel == 1 {
		return true
	}
	return false
}

// IsThreadAdmin checks if the user is an admin of the given thread
func IsThreadAdmin(thread ThreadGoForum, user User) bool {
	rightLevel := GetThreadMemberRightsLevel(thread, user)
	if rightLevel == 2 {
		return true
	}
	return false
}

// IsThreadOwner checks if the user is the owner of the given thread
func IsThreadOwner(thread ThreadGoForum, user User) bool {
	rightLevel := GetThreadMemberRightsLevel(thread, user)
	if rightLevel == 3 {
		return true
	}
	return false
}

// JoinThread adds the user to the thread
// Returns an error if there is one
func JoinThread(thread ThreadGoForum, user User) error {
	insertThreadMember := "INSERT INTO ThreadGoForumMembers (thread_id, user_id) VALUES (?, ?)"
	_, err := db.Exec(insertThreadMember, thread.ThreadID, user.UserID)
	if err != nil {
		ErrorPrintf("Error inserting the user into the thread: %v\n", err)
		return err
	}
	InfoPrintf("User %s joined the thread %s\n", user.Email, thread.ThreadName)
	return nil
}

// LeaveThread removes the user from the thread
// Returns an error if there is one
func LeaveThread(thread ThreadGoForum, user User) error {
	removeThreadMember := "DELETE FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = ?"
	_, err := db.Exec(removeThreadMember, thread.ThreadID, user.UserID)
	if err != nil {
		ErrorPrintf("Error removing the user from the thread: %v\n", err)
		return err
	}
	InfoPrintf("User %s left the thread %s\n", user.Email, thread.ThreadName)
	return nil
}

// AddUserToThread adds the user to the thread
// Returns an error if there is one
func AddUserToThread(thread ThreadGoForum, user User, rightLevel int) error {
	if rightLevel == -1 {
		InfoPrintf("Adding user %s to the thread %s as a banned member\n", user.Email, thread.ThreadName)
	} else if rightLevel == 0 {
		InfoPrintf("Adding user %s to the thread %s as a member\n", user.Email, thread.ThreadName)
	} else if rightLevel == 1 {
		InfoPrintf("Adding user %s to the thread %s as a moderator\n", user.Email, thread.ThreadName)
	} else if rightLevel == 2 {
		InfoPrintf("Adding user %s to the thread %s as an admin\n", user.Email, thread.ThreadName)
	} else if rightLevel == 3 {
		InfoPrintf("Adding user %s to the thread %s as a owner\n", user.Email, thread.ThreadName)
	} else {
		ErrorPrintf("Error: unknown rights level %d\n", rightLevel)
		return fmt.Errorf("unknown rights level %d", rightLevel)
	}
	insertThreadOwner := "INSERT INTO ThreadGoForumMembers (thread_id, user_id, rights_level) VALUES ((SELECT thread_id FROM ThreadGoForum WHERE thread_name = ?), ?, ?)"
	_, err := db.Exec(insertThreadOwner, thread.ThreadName, user.UserID, rightLevel)
	if err != nil {
		ErrorPrintf("Error inserting the thread owner into the database: %v\n", err)
		return err
	}
	return nil
}

// IsAMediaType checks if the string is a media type
func IsAMediaType(mediaType string) bool {
	for _, v := range MediaTypes {
		if string(v) == mediaType {
			return true
		}
	}
	return false
}

// GetMediaTypeFromString returns the media type from the string
// Returns an error if the media type is not valid
func GetMediaTypeFromString(mediaType string) (MediaType, error) {
	for _, v := range MediaTypes {
		if string(v) == mediaType {
			return v, nil
		}
	}
	return "", fmt.Errorf("invalid media type: %s", mediaType)
}

// GetMediaLinkFromID returns the media link from the media id
// Returns an empty MediaLink struct if there is an error
func GetMediaLinkFromID(mediaID int) MediaLink {
	getMediaLink := "SELECT * FROM MediaLink WHERE media_id = ?"
	rows, err := db.Query(getMediaLink, mediaID)
	if err != nil {
		ErrorPrintf("Error getting the media link from the id: %v\n", err)
		return MediaLink{}
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var mediaLink MediaLink
		err := rows.Scan(&mediaLink.MediaID, &mediaLink.MediaType, &mediaLink.MediaAddress, &mediaLink.CreationDate)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetMediaLinkFromID: %v\n", err)
			return MediaLink{}
		}
		return mediaLink
	}
	return MediaLink{}
}

// AddMediaLink adds a media link to the database
// Returns the media id and an error if there is one
func AddMediaLink(mediaType MediaType, mediaAddress string) (int, error) {
	insertMediaLink := "INSERT INTO MediaLink (media_type, media_address) VALUES (?, ?)"
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

// NewMediaWithIdIsValid checks if the media with the given id is valid
// It verifies if the id exists in the MediaLink table
// And if it's a ThreadMessagePicture and not yet in ThreadMessageMediaLinks table
func NewMediaWithIdIsValid(id int) bool {
	// Check if the media id exists in the MediaLink table
	getMediaLink := "SELECT media_id FROM MediaLink WHERE media_id = ? AND media_type = ?"
	rows, err := db.Query(getMediaLink, id, string(ThreadMessagePicture))
	if err != nil {
		ErrorPrintf("Error getting the media link from the id: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if !rows.Next() {
		return false
	}
	// Check if not already in ThreadMessageMediaLinks
	checkIfAlreadyInDB := "SELECT media_id FROM ThreadMessageMediaLinks WHERE media_id = ?"
	rows, err = db.Query(checkIfAlreadyInDB, id)
	if err != nil {
		ErrorPrintf("Error checking if the media link is already in the database: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		return false
	}
	return true
}

// UpdateMediaLink updates the media link in the database
// Returns an error if there is one
func UpdateMediaLink(media MediaLink) error {
	updateMediaLink := "UPDATE MediaLink SET media_type = ?, media_address = ? WHERE media_id = ?"
	_, err := db.Exec(updateMediaLink, string(media.MediaType), media.MediaAddress, media.MediaID)
	if err != nil {
		ErrorPrintf("Error updating the media link in the database: %v\n", err)
		return err
	}
	return nil
}

// GetMediaLinkFullPath returns the full path of the media link
func GetMediaLinkFullPath(mediaLink MediaLink) string {
	return path.Join(uploadFolder, imgUploadSubFolder, mediaLink.MediaAddress)
}

// DeleteUselessMediaLinks removes the media links that are not used in any message and are of the type ThreadMessagePicture and are older than 1 hour
// Returns an error if there is one
func DeleteUselessMediaLinks() error {
	// Delete the media links files before deleting their link from the database
	getUselessMediaLinks := "SELECT media_address FROM MediaLink WHERE media_id NOT IN (SELECT media_id FROM ThreadMessageMediaLinks) AND media_type = ? AND creation_date < datetime('now', '-1 hour')"
	rows, err := db.Query(getUselessMediaLinks, string(ThreadMessagePicture))
	if err != nil {
		ErrorPrintf("Error getting the useless media links: %v\n", err)
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	var mediaAddress string
	for rows.Next() {
		err := rows.Scan(&mediaAddress)
		if err != nil {
			ErrorPrintf("Error scanning the rows in DeleteUselessMediaLinks: %v\n", err)
			continue
		}
		fullPath := GetMediaLinkFullPath(MediaLink{MediaAddress: mediaAddress})
		res := RemoveImg(fullPath)
		if res == false {
			ErrorPrintf("Error removing the media link file: %v\n", err)
			continue
		}
		DebugPrintf("Removed the media link file: %s\n", fullPath)
	}

	// Delete the media links
	deleteUselessMediaLinks := "DELETE FROM MediaLink WHERE media_id NOT IN (SELECT media_id FROM ThreadMessageMediaLinks) AND media_type = ? AND creation_date < datetime('now', '-1 hour')"
	_, err = db.Exec(deleteUselessMediaLinks, string(ThreadMessagePicture))
	if err != nil {
		ErrorPrintf("Error deleting the useless media links: %v\n", err)
		return err
	}
	return nil
}

// AutoDeleteUselessMediaLinks removes the media links that are not used in any message and that are older than 1 minute
func AutoDeleteUselessMediaLinks() {
	if os.Getenv("AUTO_DELETE_USELESS_MEDIA_LINKS") == "false" {
		InfoPrintln("Auto delete useless media links was disabled")
		return
	}
	interval := 1
	var err error // We have to define it here so we can use it in the 'if' statement
	if os.Getenv("AUTO_DELETE_USELESS_MEDIA_LINKS_INTERVAL") != "" {
		interval, err = strconv.Atoi(os.Getenv("AUTO_DELETE_USELESS_MEDIA_LINKS_INTERVAL"))
		if err != nil {
			ErrorPrintf("Error parsing the interval AUTO_DELETE_USELESS_MEDIA_LINKS_INTERVAL : %v\n", err)
			interval = 1
		}
	}
	InfoPrintf("Auto delete useless media links interval is set %d minute(s)\n", interval)
	for {
		err := DeleteUselessMediaLinks()
		if err != nil {
			ErrorPrintf("Error removing useless media links: %v\n", err)
			return
		}
		DebugPrintln("Useless media links removed")
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

// IsUserAllowedToSendMessageInThread checks if the user is allowed to send a message in the thread
// Returns true if the user is allowed to send a message and false otherwise
func IsUserAllowedToSendMessageInThread(thread ThreadGoForum, user User) bool {
	if thread.ThreadID <= 0 {
		return false
	}
	if thread.OwnerID == user.UserID {
		return true
	}
	if IsUserInThread(thread, user) {
		if IsUserBannedFromThread(thread, user) {
			return false
		}
		return true
	}
	return false
}

// IsUserAllowedToEditMessageInThread checks if the user is allowed to edit a message in the thread
// Returns true if the user is allowed to edit a message and false otherwise
func IsUserAllowedToEditMessageInThread(thread ThreadGoForum, user User, messageID int) bool {
	if thread.ThreadID <= 0 {
		return false
	}
	if IsUserInThread(thread, user) {
		if IsUserBannedFromThread(thread, user) {
			return false
		}
		checkIfOwner := "SELECT message_id FROM ThreadMessages WHERE thread_id = ? AND message_id = ? AND user_id = ?"
		rows, err := db.Query(checkIfOwner, thread.ThreadID, messageID, user.UserID)
		if err != nil {
			ErrorPrintf("Error checking if the user is the owner of the message: %v\n", err)
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
	}
	return false
}

// AddMessageInThread adds a message to the thread
// Returns an error if there is one
func AddMessageInThread(thread ThreadGoForum, title string, content string, user User, mediaLinksID ...int) (int, error) {
	insertMessage := "INSERT INTO ThreadMessages (thread_id, user_id, message_title, message_content) VALUES (?, ?, ?, ?)"
	res, err := db.Exec(insertMessage, thread.ThreadID, user.UserID, title, content)
	if err != nil {
		ErrorPrintf("Error inserting the message into the database: %v\n", err)
		return -1, err
	}
	messageID, err := res.LastInsertId()
	if err != nil {
		ErrorPrintf("Error getting the last insert id: %v\n", err)
		return -1, err
	}
	// Add the media links to the message
	for _, mediaLinkID := range mediaLinksID {
		insertMediaLink := "INSERT INTO ThreadMessageMediaLinks (message_id, media_id) VALUES (?, ?)"
		_, err = db.Exec(insertMediaLink, messageID, mediaLinkID)
		if err != nil {
			ErrorPrintf("Error inserting the media link into the message: %v\n", err)
			return -1, err
		}
		DebugPrintf("Media link %d added to message %d\n", mediaLinkID, messageID)
	}
	return int(messageID), nil
}

// IsMessageTitleValid checks if the message title is valid
// Message title must be at least 5 characters long
// Message title must be at most 50 characters long
// Message title must only contain letters, numbers, underscores, hyphens, spaces, punctuation, most special characters, accents and emojis
func IsMessageTitleValid(messageTitle string) bool {
	messageTitleRegex := regexp.MustCompile(`^[a-zA-Z0-9 _\-.,;:!?(){}\[\]<>@#$%^&*+=~|\\"'/éèêëôçàâäïîùûü]{5,50}$`)
	return messageTitleRegex.MatchString(messageTitle)
}

// IsMessageContentValid checks if the message content is valid
// Message content must be at least 5 characters long
// Message content must be at most 500 characters long
// Message content must only contain letters, numbers, underscores, hyphens, spaces, punctuation, most special characters, accents and emojis
func IsMessageContentValid(messageContent string) bool {
	messageContentRegex := regexp.MustCompile(`^[a-zA-Z0-9 _\-.,;:!?(){}\[\]<>@#$%^&*+=~|\\"'/éèêëôçàâäïîùûü]{5,500}$`)
	return messageContentRegex.MatchString(messageContent)
}

// RemoveMessageFromThread removes the message from the thread
// Returns an error if there is one
func RemoveMessageFromThread(thread ThreadGoForum, messageID int) error {
	removeMessage := "DELETE FROM ThreadMessages WHERE thread_id = ? AND message_id = ?"
	_, err := db.Exec(removeMessage, thread.ThreadID, messageID)
	if err != nil {
		ErrorPrintf("Error removing the message from the database: %v\n", err)
		return err
	}
	return nil
}

// EditMessageFromThread edits the message in the thread
// Returns an error if there is one
func EditMessageFromThread(thread ThreadGoForum, messageID int, newTitle string, newContent string) error {
	editMessage := "UPDATE ThreadMessages SET message_title = ? , message_content = ? , was_edited = true WHERE thread_id = ? AND message_id = ?"
	_, err := db.Exec(editMessage, newTitle, newContent, thread.ThreadID, messageID)
	if err != nil {
		ErrorPrintf("Error editing the message in the database: %v\n", err)
		return err
	}
	return nil
}

// RemoveMediaLinkFromMessage removes a media link from a message
// Returns an error if there is one
func RemoveMediaLinkFromMessage(messageID int, mediaID int) error {
	removeMediaLink := "DELETE FROM MessagesMediaLinks WHERE message_id = ? AND media_id = ?"
	_, err := db.Exec(removeMediaLink, messageID, mediaID)
	if err != nil {
		ErrorPrintf("Error removing the media link from the message: %v\n", err)
		return err
	}
	return nil
}

// MediaExistsInMessage checks if the media link exists in the message
// Returns true if the media link exists and false otherwise
func MediaExistsInMessage(messageID int, mediaID int) bool {
	checkIfMediaExists := "SELECT media_id FROM ThreadMessageMediaLinks WHERE message_id = ? AND media_id = ?"
	rows, err := db.Query(checkIfMediaExists, messageID, mediaID)
	if err != nil {
		ErrorPrintf("Error checking if the media link exists in the message: %v\n", err)
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

// GetNumberOfMessagesInThread returns the number of messages in the thread
// Returns the number of messages and an error if there is one
func GetNumberOfMessagesInThread(thread ThreadGoForum) (int, error) {
	getNumberOfMessages := "SELECT COUNT(*) FROM ThreadMessages WHERE thread_id = ?"
	rows, err := db.Query(getNumberOfMessages, thread.ThreadID)
	if err != nil {
		ErrorPrintf("Error getting the number of messages in the thread: %v\n", err)
		return 0, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var numberOfMessages int
		err := rows.Scan(&numberOfMessages)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetNumberOfMessagesInThread: %v\n", err)
			return 0, err
		}
		return numberOfMessages, nil
	}
	return 0, nil
}

// MessageExistsInThread checks if the message exists in the thread
// Returns true if the message exists and false otherwise
func MessageExistsInThread(thread ThreadGoForum, messageID int) bool {
	checkIfMessageExists := "SELECT message_id FROM ThreadMessages WHERE thread_id = ? AND message_id = ?"
	rows, err := db.Query(checkIfMessageExists, thread.ThreadID, messageID)
	if err != nil {
		ErrorPrintf("Error checking if the message exists in the thread: %v\n", err)
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

// IsUserAllowedToDeleteMessage checks if the user is allowed to delete the message
// Returns true if the user is allowed to delete the message and false otherwise
// A user is allowed to delete a message if he is the owner of the message or if he is the owner or an admin of the thread
func IsUserAllowedToDeleteMessage(thread ThreadGoForum, user User, messageID int) bool {
	if thread.OwnerID == user.UserID {
		return true
	}
	if IsThreadOwner(thread, user) {
		return true
	}
	if IsThreadAdmin(thread, user) {
		return true
	}
	// Check if the user is the owner of the message
	checkIfMessageOwner := "SELECT user_id FROM ThreadMessages WHERE thread_id = ? AND message_id = ?"
	rows, err := db.Query(checkIfMessageOwner, thread.ThreadID, messageID)
	if err != nil {
		ErrorPrintf("Error checking if the user is the owner of the message: %v\n", err)
		return false
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var messageOwnerID int
		err := rows.Scan(&messageOwnerID)
		if err != nil {
			ErrorPrintf("Error scanning the rows in IsUserAllowedToDeleteMessage: %v\n", err)
			return false
		}
		if messageOwnerID == user.UserID {
			return true
		}
	}
	return false
}

// HasUserAlreadyVoted checks if the user has already voted for the message
// Returns 0 if the user has not voted, 1 if the user has upvoted and -1 if the user has downvoted
func HasUserAlreadyVoted(user User, messageID int) int {
	checkIfUserVoted := "SELECT is_upvote FROM ThreadMessageVotes WHERE user_id = ? AND message_id = ?"
	rows, err := db.Query(checkIfUserVoted, user.UserID, messageID)
	if err != nil {
		ErrorPrintf("Error checking if the user has already voted: %v\n", err)
		return 0
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows in HasUserAlreadyVoted: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var isUpvote bool
		err := rows.Scan(&isUpvote)
		if err != nil {
			ErrorPrintf("Error scanning the rows in HasUserAlreadyVoted: %v\n", err)
			return 0
		}
		if isUpvote {
			return 1
		} else {
			return -1
		}
	}
	return 0
}

// GetMessagesFromThread returns the messages from the thread
// Returns a slice of messages and an error if there is one
// The messages are ordered by the given order (from the OrderingList)
// The offset is used to paginate the messages
// By default the function returns a maximum of 10 messages or is equal to the environment variable 'MAX_MESSAGES_PER_PAGE_LOAD'
func GetMessagesFromThread(thread ThreadGoForum, offset int, order string) ([]FormattedThreadMessage, error) {
	return GetMessagesFromThreadWithPOV(thread, offset, order, User{})
}

// GetMessagesFromThreadWithPOV returns the messages from the thread viewed from the point of view of the user
// Returns a slice of messages and an error if there is one
// The messages are ordered by the given order (from the OrderingList)
// The offset is used to paginate the messages
// By default the function returns a maximum of 10 messages or is equal to the environment variable 'MAX_MESSAGES_PER_PAGE_LOAD'
func GetMessagesFromThreadWithPOV(thread ThreadGoForum, offset int, order string, user User) ([]FormattedThreadMessage, error) {
	// Check if there is still incompleteMessages to load
	numberOfMessages, err := GetNumberOfMessagesInThread(thread)
	if err != nil {
		ErrorPrintf("Error getting the number of incompleteMessages in the thread: %v\n", err)
		return nil, err
	}
	if offset >= numberOfMessages {
		return nil, nil
	}

	// Get the max incompleteMessages per page load from the environment variable
	maxMessagesPerPageLoad := 10
	if os.Getenv("MAX_MESSAGES_PER_PAGE_LOAD") != "" {
		var err error
		maxMessagesPerPageLoad, err = strconv.Atoi(os.Getenv("MAX_MESSAGES_PER_PAGE_LOAD"))
		if err != nil {
			ErrorPrintf("Error parsing the max incompleteMessages per page load: %v\n", err)
			maxMessagesPerPageLoad = 10
		}
	}
	var getMessages string
	switch order {
	case "desc": // descending order
		getMessages = "SELECT * FROM ViewThreadMessagesWithVotes WHERE thread_name = ? ORDER BY creation_date DESC LIMIT ? OFFSET ?"
		break
	case "popular": // popular order
		getMessages = "SELECT * FROM ViewThreadMessagesWithVotes WHERE thread_name = ? ORDER BY (upvotes - downvotes) DESC LIMIT ? OFFSET ?"
		break
	case "unpopular": // unpopular order
		getMessages = "SELECT * FROM ViewThreadMessagesWithVotes WHERE thread_name = ? ORDER BY (upvotes - downvotes) ASC LIMIT ? OFFSET ?"
		break
	default: // ascending order
		getMessages = "SELECT * FROM ViewThreadMessagesWithVotes WHERE thread_name = ? ORDER BY creation_date ASC LIMIT ? OFFSET ?"
		break
	}
	rows, err := db.Query(getMessages, thread.ThreadName, maxMessagesPerPageLoad, offset)
	if err != nil {
		ErrorPrintf("Error getting all the incompleteMessages from the thread: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	var incompleteMessages []FormattedThreadMessage
	for rows.Next() {
		var message FormattedThreadMessage
		err := rows.Scan(
			&message.MessageID,
			&message.MessageTitle,
			&message.MessageContent,
			&message.WasEdited,
			&message.CreationDate,
			&message.UserName,
			&message.UserPfpAddress,
			&message.Upvotes,
			&message.Downvotes)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetMessagesFromThread: %v\n", err)
			return nil, err
		}
		incompleteMessages = append(incompleteMessages, message)
	}
	// Get the media links for each message
	var Messages []FormattedThreadMessage
	for _, message := range incompleteMessages {
		getMessageMediaLinks := "SELECT ml.media_id, ml.media_type, ml.media_address, ml.creation_date FROM ThreadMessageMediaLinks tmml JOIN MediaLink ml ON tmml.media_id = ml.media_id WHERE tmml.message_id = ?;"
		rows, err := db.Query(getMessageMediaLinks, message.MessageID)
		if err != nil {
			ErrorPrintf("Error getting the incompleteMessages from the thread: %v\n", err)
			return nil, err
		}
		for rows.Next() {
			var mediaLink MediaLink
			err := rows.Scan(&mediaLink.MediaID, &mediaLink.MediaType, &mediaLink.MediaAddress, &mediaLink.CreationDate)
			if err != nil {
				ErrorPrintf("Error scanning the rows in GetMessagesFromThread: %v\n", err)
				return nil, err
			}
			message.MediaLinks = append(message.MediaLinks, mediaLink.MediaAddress)
		}
		// Sets FormattedThreadMessage.VoteState to -1 if the user disliked the message
		// Sets FormattedThreadMessage.VoteState to 1 if the user liked the message
		// Sets FormattedThreadMessage.VoteState to 0 if the user has not voted
		if user.UserID != 0 {
			message.VoteState = HasUserAlreadyVoted(user, message.MessageID)
		} else {
			message.VoteState = 0
		}
		Messages = append(Messages, message)
		err = rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}
	return Messages, nil
}

// GetUserRankInThread returns the rights_level of the user in the thread
// ( 0 = member, 1 = moderator, 2 = admin, 3 = owner, -1 = banned )
// Returns the rights_level of the user in the thread
func GetUserRankInThread(thread ThreadGoForum, user User) int {
	getUserRank := "SELECT rights_level FROM ThreadGoForumMembers WHERE thread_id = ? AND user_id = ?"
	rows, err := db.Query(getUserRank, thread.ThreadID, user.UserID)
	if err != nil {
		ErrorPrintf("Error getting the rank of the user in the thread: %v\n", err)
		return 0
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		var rank int
		err := rows.Scan(&rank)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetUserRankInThread: %v\n", err)
			return 0
		}
		return rank
	}
	return 0
}

// GetThreadModerationTeam returns the moderation team of the thread
// Returns a slice of SimplifiedUser
func GetThreadModerationTeam(forum ThreadGoForum) []SimplifiedUser {
	getThreadModerationTeam := `
		SELECT u.username, ml.media_address, tgm.rights_level
		FROM ThreadGoForumMembers tgm
		    JOIN Users u ON tgm.user_id = u.user_id
		    JOIN UserConfigs uc ON u.user_id = uc.user_id
		    JOIN MediaLink ml ON uc.pfp_id = ml.media_id
		WHERE tgm.thread_id = ? AND tgm.rights_level > 0
		ORDER BY tgm.rights_level DESC;`
	rows, err := db.Query(getThreadModerationTeam, forum.ThreadID)
	if err != nil {
		ErrorPrintf("Error getting the thread moderation team: %v\n", err)
		return nil
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	var moderationTeam []SimplifiedUser
	for rows.Next() {
		var user SimplifiedUser
		err := rows.Scan(&user.Username, &user.PfpAddress, &user.RightsLevel)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetThreadModerationTeam: %v\n", err)
			return nil
		}
		moderationTeam = append(moderationTeam, user)
	}
	return moderationTeam
}

// GetUserThreads returns the threads where the user is a member
// Returns a slice of threads or nil if there is an error
func GetUserThreads(user User) []ThreadGoForum {
	getUserThreads := "SELECT tg.thread_id, tg.thread_name, tg.owner_id, tg.creation_date FROM ThreadGoForumMembers tgm JOIN ThreadGoForum tg ON tgm.thread_id = tg.thread_id WHERE tgm.user_id = ?"
	rows, err := db.Query(getUserThreads, user.UserID)
	if err != nil {
		ErrorPrintf("Error getting the threads from the user: %v\n", err)
		return nil
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
			ErrorPrintf("Error scanning the rows in GetUserThreads: %v\n", err)
			return nil
		}
		threads = append(threads, thread)
	}
	return threads
}

// ThreadMessageAddVote adds a vote to the message
// Returns an error if there is one
func ThreadMessageAddVote(messageID int, userID int, voteType bool) error {
	addVote := "INSERT INTO ThreadMessageVotes (message_id, user_id, is_upvote) VALUES (?, ?, ?)"
	_, err := db.Exec(addVote, messageID, userID, voteType)
	if err != nil {
		ErrorPrintf("Error adding the vote to the message: %v\n", err)
		return err
	}
	return nil
}

// ThreadMessageUpVote adds a vote to the message
// Returns an error if there is one
func ThreadMessageUpVote(messageID int, userID int) error {
	return ThreadMessageAddVote(messageID, userID, true)
}

// ThreadMessageDownVote adds a downvote to the message
// Returns an error if there is one
func ThreadMessageDownVote(messageID int, userID int) error {
	return ThreadMessageAddVote(messageID, userID, false)
}

// ThreadMessageRemoveVote removes the vote from the message
// Returns an error if there is one
func ThreadMessageRemoveVote(messageID int, userID int) error {
	removeVote := "DELETE FROM ThreadMessageVotes WHERE message_id = ? AND user_id = ?"
	_, err := db.Exec(removeVote, messageID, userID)
	if err != nil {
		ErrorPrintf("Error removing the vote from the message: %v\n", err)
		return err
	}
	return nil
}

// ThreadMessageUpdateVote updates the vote of the message
// Returns an error if there is one
// It updates the vote of the message to the new vote
// if voteType is true, it means the user upvoted the message
// if voteType is false, it means the user downvoted the message
func ThreadMessageUpdateVote(messageID int, userID int, voteType bool) error {
	updateVote := "UPDATE ThreadMessageVotes SET is_upvote = ? WHERE message_id = ? AND user_id = ?"
	_, err := db.Exec(updateVote, voteType, messageID, userID)
	if err != nil {
		ErrorPrintf("Error updating the vote of the message: %v\n", err)
		return err
	}
	return nil
}

func isValidHexColor(color string) bool {
	// Check if the color is a valid hexadecimal color code
	// The color must start with # and be followed by 6 or 3 hexadecimal digits
	if len(color) != 7 && len(color) != 4 {
		return false
	}
	if color[0] != '#' {
		return false
	}
	for i := 1; i < len(color); i++ {
		if !((color[i] >= '0' && color[i] <= '9') || (color[i] >= 'a' && color[i] <= 'f') || (color[i] >= 'A' && color[i] <= 'F')) {
			return false
		}
	}
	return true
}

// IsTagAlreadyExists checks if the tag already exists in the thread
// Returns true if the tag already exists and false otherwise
func IsTagAlreadyExists(thread ThreadGoForum, tagName string) (bool, error) {
	getTag := "SELECT tag_name FROM ThreadGoForumTags WHERE thread_id = ? AND tag_name = ?"
	rows, err := db.Query(getTag, thread.ThreadID, tagName)
	if err != nil {
		ErrorPrintf("Error checking if the tag already exists: %v\n", err)
		return false, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

// IsTagIDAssociatedWithThread checks if the tag id is associated with the thread
// Returns true if the tag id is associated with the thread and false otherwise
func IsTagIDAssociatedWithThread(thread ThreadGoForum, tagID int) (bool, error) {
	getTag := "SELECT tag_id FROM ThreadGoForumTags WHERE thread_id = ? AND tag_id = ?"
	rows, err := db.Query(getTag, thread.ThreadID, tagID)
	if err != nil {
		ErrorPrintf("Error checking if the tag id is associated with the thread: %v\n", err)
		return false, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

// AddThreadTag adds a tag to the thread
// Returns an error if there is one
// The tag name must be unique in the thread
// The tag color must be a valid hexadecimal color code
// The tag name must be at least 3 characters long
// The tag name must be at most 30 characters long
func AddThreadTag(thread ThreadGoForum, tagName string, tagColor string) error {
	if len(tagName) < 3 || len(tagName) > 30 {
		return fmt.Errorf("tag name must be between 3 and 30 characters long")
	}
	if !isValidHexColor(tagColor) {
		return fmt.Errorf("tag color must be a valid hexadecimal color code")
	}
	exists, err := IsTagAlreadyExists(thread, tagName)
	if err != nil {
		ErrorPrintf("Error checking if the tag already exists: %v\n", err)
		return err
	}
	if exists {
		return fmt.Errorf("tag name already exists in the thread")
	}
	insertTag := "INSERT INTO ThreadGoForumTags (thread_id, tag_name, tag_color) VALUES (?, ?, ?)"
	_, err = db.Exec(insertTag, thread.ThreadID, tagName, tagColor)
	if err != nil {
		ErrorPrintf("Error adding the tag to the thread: %v\n", err)
		return err
	}
	return nil
}

// RemoveThreadTag removes the tag from the thread
// Returns an error if there is one
func RemoveThreadTag(thread ThreadGoForum, tagName string) error {
	removeTag := "DELETE FROM ThreadGoForumTags WHERE thread_id = ? AND tag_name = ?"
	_, err := db.Exec(removeTag, thread.ThreadID, tagName)
	if err != nil {
		ErrorPrintf("Error removing the tag from the thread: %v\n", err)
		return err
	}
	return nil
}

// GetThreadTags returns the tags of the thread
// Returns a slice of tags and an error if there is one
func GetThreadTags(thread ThreadGoForum) ([]ThreadTag, error) {
	getTags := "SELECT * FROM ThreadGoForumTags WHERE thread_id = ?"
	rows, err := db.Query(getTags, thread.ThreadID)
	if err != nil {
		ErrorPrintf("Error getting the tags from the thread: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	var tags []ThreadTag
	for rows.Next() {
		var tag ThreadTag
		err := rows.Scan(&tag.TagID, &tag.ThreadID, &tag.TagName, &tag.TagColor)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetThreadTags: %v\n", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

// AddTagToMessage adds the tag to the message
// Returns an error if there is one
func AddTagToMessage(messageID int, tags int) error {
	addTag := "INSERT INTO ThreadMessageTags (message_id, tag_id) VALUES (?, ?)"
	_, err := db.Exec(addTag, messageID, tags)
	if err != nil {
		ErrorPrintf("Error adding the tag to the message: %v\n", err)
		return err
	}
	return nil
}

// GetMessageTags returns the tags of the message
// Returns a slice of tags and an error if there is one
func GetMessageTags(messageID int) ([]ThreadTag, error) {
	getTags := "SELECT tt.tag_id, tt.thread_id, tt.tag_name, tt.tag_color FROM ThreadGoForumTags tt JOIN ThreadMessageTags tmt ON tt.tag_id = tmt.tag_id WHERE tmt.message_id = ?"
	rows, err := db.Query(getTags, messageID)
	if err != nil {
		ErrorPrintf("Error getting the tags from the message: %v\n", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			ErrorPrintf("Error closing the rows: %v\n", err)
		}
	}(rows)
	var tags []ThreadTag
	for rows.Next() {
		var tag ThreadTag
		err := rows.Scan(&tag.TagID, &tag.ThreadID, &tag.TagName, &tag.TagColor)
		if err != nil {
			ErrorPrintf("Error scanning the rows in GetMessageTags: %v\n", err)
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
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

	// The 'ThreadGoForumTags' table represents the tags of a thread
	// the tag_color column is used to determine the color of the tag (it's a hexadecimal color code, e.g. #FF0000)
	ThreadGoForumTagsTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadGoForumTags (
		    tag_id INTEGER PRIMARY KEY AUTOINCREMENT,
		    thread_id INTEGER NOT NULL,
		    tag_name TEXT NOT NULL,
		    tag_color TEXT NOT NULL,
		    FOREIGN KEY (thread_id) REFERENCES ThreadGoForum(thread_id)
		);
		`
	_, err = db.Exec(ThreadGoForumTagsTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadGoForumTags table: %v\n", err)
		return
	}

	// The 'ThreadGoForumMembers' represents the members of a thread
	// if the 'right_level' is -1, the user is banned from the thread
	// if the 'rights_level' is 0, the user is a normal member
	// if the 'rights_level' is 1, the user is a moderator
	// if the 'rights_level' is 2, the user is an admin
	// if the 'rights_level' is 3, the user is the owner of the thread
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

	// The 'MediaLink' table represents the media links (images, videos, etc.) that are shared in the threads
	// For now, we only will do images as stated in the project instructions
	MediaLinksTableSQL := `
		CREATE TABLE IF NOT EXISTS MediaLink (
			media_id INTEGER PRIMARY KEY AUTOINCREMENT,
			media_type TEXT NOT NULL,
			media_address TEXT NOT NULL UNIQUE,
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		`
	_, err = db.Exec(MediaLinksTableSQL)
	if err != nil {
		ErrorPrintf("Error creating MediaLink table: %v\n", err)
		return
	}

	// The 'ThreadMessages' table represents the messages that are sent in the threads
	ThreadMessagesTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadMessages (
			message_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			thread_id INTEGER NOT NULL,
			message_title TEXT NOT NULL,
			message_content TEXT NOT NULL,
			was_edited BOOLEAN DEFAULT FALSE NOT NULL,
			creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		    FOREIGN KEY (user_id) REFERENCES Users(user_id),
		    FOREIGN KEY (thread_id) REFERENCES ThreadGoForum(thread_id)
		);
		`
	_, err = db.Exec(ThreadMessagesTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadMessages table: %v\n", err)
		return
	}

	// The 'ThreadMessageMediaLinks' table represents the media links that are shared in the messages
	ThreadMessageMediaLinksTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadMessageMediaLinks (
		    message_id INTEGER NOT NULL,
		    media_id INTEGER NOT NULL,
		    FOREIGN KEY (message_id) REFERENCES ThreadMessages(message_id),
		    FOREIGN KEY (media_id) REFERENCES MediaLink(media_id),
		    PRIMARY KEY (message_id, media_id)
		);
		`
	_, err = db.Exec(ThreadMessageMediaLinksTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadMessageMediaLinks table: %v\n", err)
		return
	}
	ThreadMessageVotesTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadMessageVotes (
		    message_id INTEGER NOT NULL,
		    user_id INTEGER NOT NULL,
		    is_upvote BOOLEAN NOT NULL
		);
		`
	_, err = db.Exec(ThreadMessageVotesTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadMessageVotes table: %v\n", err)
		return
	}

	ThreadMessageTagsTableSQL := `
		CREATE TABLE IF NOT EXISTS ThreadMessageTags (
		    message_id INTEGER NOT NULL,
		    tag_id INTEGER NOT NULL,
		    FOREIGN KEY (message_id) REFERENCES ThreadMessages(message_id),
		    FOREIGN KEY (tag_id) REFERENCES ThreadGoForumTags(tag_id),
		    PRIMARY KEY (message_id, tag_id)
		);
		`
	_, err = db.Exec(ThreadMessageTagsTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ThreadMessageTags table: %v\n", err)
		return
	}

	ViewThreadMessageWithLikesTableSQL := `
		CREATE VIEW IF NOT EXISTS ViewThreadMessagesWithVotes AS
		SELECT 
			tm.message_id,
			tm.message_title,
			tm.message_content,
			tm.was_edited,
			tm.creation_date,
			u.username,
			ml.media_address AS pfp_media_address,
			COALESCE(v.upvotes, 0) AS upvotes,
			COALESCE(v.downvotes, 0) AS downvotes
		FROM ThreadMessages tm
		JOIN ThreadGoForum tg ON tm.thread_id = tg.thread_id
		JOIN Users u ON tm.user_id = u.user_id
		LEFT JOIN UserConfigs uc ON u.user_id = uc.user_id
		LEFT JOIN MediaLink ml ON uc.pfp_id = ml.media_id
		LEFT JOIN (
			SELECT 
				message_id,
				SUM(CASE WHEN is_upvote = 1 THEN 1 ELSE 0 END) AS upvotes,
				SUM(CASE WHEN is_upvote = 0 THEN 1 ELSE 0 END) AS downvotes
			FROM ThreadMessageVotes
			GROUP BY message_id
		) v ON tm.message_id = v.message_id;
		`
	_, err = db.Exec(ViewThreadMessageWithLikesTableSQL)
	if err != nil {
		ErrorPrintf("Error creating ViewThreadMessagesWithVotes view: %v\n", err)
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
	query := "SELECT COUNT(*) FROM MediaLink"
	var count int
	err = db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error checking if MediaLink table is empty: %v\n", err)
	}
	// If the table is empty, we insert the default media links
	if count == 0 {
		// clone the default media files from the assets folder to the media folder
		defaultMediaFiles := []string{
			"default_user_icon.png",
			"default_thread_icon.png",
			"default_thread_banner.gif",
		}

		for _, file := range defaultMediaFiles {
			origin := fmt.Sprintf("statics/img/%s", file)
			destination := fmt.Sprintf("%s/%s", GetImgUploadFolder(), file)
			err := copyFile(origin, destination)
			if err != nil {
				ErrorPrintf("Error copying default media file %s: %v\n", file, err)
				return
			}
		}

		// Insert default media links
		insertDefaultMediaLinks := fmt.Sprintf(`
		INSERT INTO MediaLink (media_type, media_address) VALUES
			('%s', 'default_user_icon.png'),
			('%s', 'default_thread_icon.png'),
			('%s', 'default_thread_banner.gif');
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

	// Starting the auto delete of the useless media links
	go AutoDeleteUselessMediaLinks()

	InfoPrintln("Database initialised")
}

// copyFile copies a file from src to dst.
// It creates a new file in the destination path and copies the contents of the source file to it.
// If the destination file already exists, it will be overwritten.
// It returns an error if there is one.
func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Flush the destination file to ensure all data is written
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}

// FillDatabase fills the database with test data.
func FillDatabase() {
	// TODO : fill the database with test data for development testing and demonstration purposes

	// A test thread
	err := AddThread(User{UserID: 1}, "TestThread", "This is a test thread ! :P  (o_o)")
	if err != nil {
		ErrorPrintf("Error adding thread TestThread: %v\n", err)
		return
	}

	// A test thread with must be connected
	err = AddThread(User{UserID: 1}, "TestThread2", "This is an other test thread where you must be connected ! (►__◄)")
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
	err = AddThread(User{UserID: 1}, "TestThread3", "This is also an other test thread where you must be a member ! (◕‿-)")
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

	// Adding fake users
	for i := 0; i < 15; i++ {
		err := AddUser(fmt.Sprintf("fakeuser%d@fakemailservice.com", i), fmt.Sprintf("fakeuser%d", i), "Fake", "User", "password")
		if err != nil {
			ErrorPrintf("Error adding fake user %d: %v\n", i, err)
			return
		}
		err = VerifyEmail(fmt.Sprintf("fakeuser%d@fakemailservice.com", i))
		if err != nil {
			ErrorPrintf("Error verifying fake user %d: %v\n", i, err)
			return
		}
	}

	// Adding fake messages
	for i := 0; i < 50; i++ {
		// Randomly add a media link to the message
		var mediaIDs []int
		if mr.Intn(2) == 0 {
			mediaID, err := AddMediaLink(ThreadMessagePicture, fmt.Sprintf("https://fakeimage%d.com", i))
			if err != nil {
				ErrorPrintf("Error adding media link %d: %v\n", i, err)
				return
			}
			mediaIDs = append(mediaIDs, mediaID)
			if mr.Intn(2) == 0 {
				mediaID, err := AddMediaLink(ThreadMessagePicture, fmt.Sprintf("https://secondfakeimage%d.com", i))
				if err != nil {
					ErrorPrintf("Error adding media link %d: %v\n", i, err)
					return
				}
				mediaIDs = append(mediaIDs, mediaID)
			}
		}
		_, err := AddMessageInThread(
			GetThreadFromName("TestThread"),
			fmt.Sprintf("Test message %d title", i),
			fmt.Sprintf("This is a test %d message ", i),
			User{UserID: (i % 15) + 1},
			mediaIDs...)
		if err != nil {
			ErrorPrintf("Error adding fake message %d: %v\n", i, err)
			return
		}
		time.Sleep(250 * time.Millisecond)
	}

	// Adding fake upvotes / downvotes to the messages
	for i := 0; i < 15; i++ { // loop for each user
		for j := 0; j < 50; j++ { // loop for each message
			// Randomly upvote or downvote the message
			if mr.Intn(2) == 0 {
				err := ThreadMessageUpVote(j+1, i+1)
				if err != nil {
					ErrorPrintf("Error adding upvote %d for user %d: %v\n", j, i, err)
					return
				}
			} else {
				err := ThreadMessageDownVote(j+1, i+1)
				if err != nil {
					ErrorPrintf("Error adding downvote %d for user %d: %v\n", j, i, err)
					return
				}
			}
		}
	}

	// Adding fake Moderators / Admin to the threads
	for i := 1; i < 15; i++ {
		err := AddUserToThread(GetThreadFromName("TestThread"), User{UserID: i + 1, Username: fmt.Sprintf("fakeuser%d", i)}, mr.Intn(3))
		if err != nil {
			ErrorPrintf("Error adding user/moderator/admin %d to thread TestThread: %v\n", i, err)
			return
		}
	}
}
