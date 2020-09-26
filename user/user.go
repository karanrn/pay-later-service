package user

import (
	"database/sql"
	"errors"

	"github.com/karanrn/pay-later-service/db"
)

// User type for user information
type User struct {
	UserID      string
	EmailID     string
	CreditLimit float64
	CreditSpent float64
	Limit       bool
}

// Search searches user basis userID
func Search(userID string) (user User, err error) {
	db := db.Connect()
	defer db.Close()

	searchStmt, err := db.Query("SELECT user_id, email_id, credit_limit, credit_spent, limit_reached FROM user WHERE user_id=?", userID)
	if err != nil {
		return User{}, errors.New(err.Error())
	}

	user = User{}
	for searchStmt.Next() {
		err = searchStmt.Scan(&user.UserID, &user.EmailID, &user.CreditLimit, &user.CreditSpent, &user.Limit)
		if err != nil {
			return User{}, errors.New(err.Error())
		}
	}
	if (user == User{}) {
		return User{}, errors.New("User not found")
	}

	return user, nil
}

// Add adds user to the database
func Add(user User) (err error) {
	db := db.Connect()
	defer db.Close()

	// Set defaults for Credit_Spent (0) and Limit (False)
	insertStmt, err := db.Prepare("INSERT INTO user (user_id, email_id, credit_limit) VALUES (?, ?, ?)")
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = insertStmt.Exec(user.UserID, user.EmailID, user.CreditLimit)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// CreditUpdate updates creditSpent for the user
func CreditUpdate(user User, amt float64, payback bool) (err error) {
	db := db.Connect()
	defer db.Close()

	// Accept and update user's payback or transaxtion amount
	var creditUsed float64
	if payback {
		// Payback
		creditUsed = user.CreditSpent - amt
	} else {
		// Transaction, credits used
		creditUsed = user.CreditSpent + amt
	}
	_, err = db.Query("UPDATE user SET credit_spent = ? WHERE user_id = ?", creditUsed, user.UserID)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// ListUsers lists all the users from the system
// atCreditLimit = true will get all users that have thier credit limit reached
func ListUsers(atCreditLimit bool) (users []User, err error) {
	db := db.Connect()
	defer db.Close()

	var selectStmt *sql.Rows
	if atCreditLimit {
		selectStmt, err = db.Query("SELECT user_id FROM user WHERE credit_limit = credit_spent ORDER BY user_id;")
		if err != nil {
			return []User{}, errors.New(err.Error())
		}
	} else {
		// Get all users
		selectStmt, err = db.Query("SELECT user_id, email_id, credit_limit, credit_spent, limit_reached FROM user ORDER BY user_id;")
		if err != nil {
			return []User{}, errors.New(err.Error())
		}
	}

	users = []User{}
	for selectStmt.Next() {
		user := User{}
		err = selectStmt.Scan(&user.UserID, &user.EmailID, &user.CreditLimit, &user.CreditSpent, &user.Limit)
		if err != nil {
			return []User{}, errors.New(err.Error())
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		return []User{}, errors.New("No users in the system")
	}

	return users, nil
}
