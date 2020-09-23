package user

import (
	"errors"

	"github.com/karanrn/pay-later-service/db"
)

// User type for user information
type User struct {
	UserID      string
	EmailID     string
	CreditLimit int64
	CreditSpent int64
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

// Payback updates creditSpent for the user
func Payback(userID string, paybackAmt int64) (err error) {
	db := db.Connect()
	defer db.Close()

	selectStmt, err := db.Query("SELECT user_id, email_id, credit_limit, credit_spent, limit_reached FROM user WHERE user_id=?", userID)
	if err != nil {
		return errors.New(err.Error())
	}

	user := User{}
	for selectStmt.Next() {
		err = selectStmt.Scan(&user.UserID, &user.EmailID, &user.CreditLimit, &user.CreditSpent, &user.Limit)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	// Accept and update user's payback
	payback := user.CreditSpent - paybackAmt
	_, err = db.Query("UPDATE user SET credit_spent = ? WHERE user_id = ?", payback, user.UserID)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
