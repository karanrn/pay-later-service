package main

import (
	"flag"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/karanrn/pay-later-service/merchant"
	"github.com/karanrn/pay-later-service/user"
)

func main() {
	// Users and Merchants
	var newUser user.User
	var newMerchant merchant.Merchant

	newUsage := `For adding a new user/merchant: 
	new user user_id email_id credit-limit
	ex: new user user1 u1@users.com 500
	For adding a new merchant: 
	new merchant merchant_id discount-percentage
	ex: new merchant m1 email_id 10%
	For new transaction:
	new txn user_id merchant_id txn-amount0
	ex: new txn user1 m1 150`
	/*
		updateUsage := `For updating merchant information - discount-percentage:
		update merchant merchant_id new_discount_percentage
		ex: update merchant m1 5%`
		reportUsage := `For getting report of the system:
		report [discount dues total-dues users-at-credit-limit] [user_id merchant_id]
		ex: report discount m1
			report dues user1`
	*/
	// New subcommands or flags
	userID := flag.String("user", "", newUsage)
	email := flag.String("email", "", newUsage)
	creditLimit := flag.Int64("credit-limit", 0, newUsage)

	merchantID := flag.String("merchant", "", newUsage)
	discount := flag.String("discount", "0%", newUsage)
	flag.Parse()

	if *userID != "" {
		emailCheck := isEmailValid(*email)
		if !emailCheck {
			fmt.Printf("%s - Not a valid email address.", *email)
			return
		}

		newUser.UserID = *userID
		newUser.EmailID = *email
		newUser.CreditLimit = *creditLimit

		err := user.Add(newUser)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	if *merchantID != "" {
		newMerchant.MerchantID = *merchantID

		emailCheck := isEmailValid(*email)
		if !emailCheck {
			fmt.Printf("%s - Not a valid email address.", *email)
			return
		}
		newMerchant.EmailID = *email

		disAmt, err := strconv.ParseInt(strings.TrimRight(*discount, "%"), 10, 64)
		if err != nil {
			fmt.Printf("%s - Invalid discount type. Should of format 10%%", *discount)
		}
		newMerchant.Discount = disAmt

		err = merchant.Add(newMerchant)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(email) < 3 && len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}
