package main

import (
	"flag"
	"fmt"
	"os"
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

	// New subcommands or flags
	userCmd := flag.NewFlagSet("user", flag.ContinueOnError)
	userID := userCmd.String("id", "", "User ID: alphanumeric value. ex: -user user101")
	email := userCmd.String("email", "", "Email ID of user: xyz@domain.com. ex: -email abc@xyz.com")
	creditLimit := userCmd.Int64("credit-limit", 0, "Credit limit for the user, default is 0. ex: -credit-limit 1000")

	merchantCmd := flag.NewFlagSet("merchant", flag.ContinueOnError)
	merchantID := merchantCmd.String("id", "", "Merchant ID: alphanumeric value. ex: -merchant m101")
	mEmail := merchantCmd.String("email", "", "Email ID of merchant: xyz@domain.com. ex: -email abc@xyz.com")
	discount := merchantCmd.String("discount", "0%", "Discount offered by merchant. ex: -discount 5%")

	txn := flag.NewFlagSet("txn", flag.ContinueOnError)
	tUser := txn.String("user-id", "", "User involved in transaction. ex: -user-id user101")
	tMerchant := txn.String("merchant-id", "", "Merchant involved in transaction. ex: -merchant-id m101")
	tAmt := txn.Int64("amt", 0, "Transaction amount. ex: -amt 100")

	if len(os.Args[1:]) < 1 {
		fmt.Println("You must pass sub command - [user, merchant, txn]")
		return
	}

	switch os.Args[1] {
	case "user":
		if err := userCmd.Parse(os.Args[2:]); err == nil {
			if *userID != "" {
				emailCheck := isEmailValid(*email)
				if !emailCheck {
					fmt.Printf("%s - Not a valid email address.", *email)
					return
				}
				if *creditLimit == 0 {
					fmt.Printf("Credit Limit cannot be zero (Default: zero).")
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
		}
	case "merchant":
		if err := merchantCmd.Parse(os.Args[2:]); err == nil {
			if *merchantID != "" {
				newMerchant.MerchantID = *merchantID

				emailCheck := isEmailValid(*mEmail)
				if !emailCheck {
					fmt.Printf("%s - Not a valid email address.", *email)
					return
				}
				newMerchant.EmailID = *mEmail

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
	case "txn":
		if err := txn.Parse(os.Args[2:]); err == nil {
			if *tUser != "" && *tMerchant != "" {
				// TO DO - implement transaction logic
				fmt.Printf("Transaction amount: %d", *tAmt)
			}
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
