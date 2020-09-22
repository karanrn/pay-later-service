package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type user struct {
	userID      string
	emailID     string
	creditLimit int64
	creditSpent int64
	limit       bool
}

type merchant struct {
	merchantID string
	emailID    string
	discount   int64
}

func main() {

	// Users and Merchants
	var users []user
	var merchants []merchant

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
	new := flag.String("new", "", newUsage)
	/*
		update := flag.String("update", "", updateUsage)
		report := flag.String("report", "", reportUsage)
	*/
	flag.Parse()

	if *new != "" {
		switch *new {
		case "user":
			credit, _ := strconv.ParseInt(os.Args[5], 10, 64)
			users = append(users, user{userID: os.Args[3], emailID: os.Args[4], creditLimit: credit})
			fmt.Printf("%s(%d)\n", os.Args[3], credit)
		case "merchant":
			discount, _ := strconv.ParseInt(strings.TrimRight(os.Args[5], "%"), 10, 64)
			merchants = append(merchants, merchant{merchantID: os.Args[3], emailID: os.Args[4], discount: discount})
			fmt.Printf("%s(%d)\n", os.Args[3], discount)
		case "txn":
			userID, err := searchUser(users, os.Args[3])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			merchantID, err := searchMerchant(merchants, os.Args[4])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			txnAmt, _ := strconv.ParseInt(os.Args[5], 10, 64)
			if txnAmt > (users[userID].creditLimit - users[userID].creditSpent) {
				users[userID].limit = true
				fmt.Println("Rejected! (credit limit)")
			} else {
				users[userID].creditSpent += txnAmt * (merchants[merchantID].discount / 100.0)
			}
		default:
			fmt.Println("Pass the valid subcommand: [user, merchant, txn]")
		}
	}

	fmt.Println(users)
	fmt.Println(merchants)
}

func searchUser(users []user, id string) (index int, err error) {
	for i, user := range users {
		if id == user.userID {
			return i, nil
		}
	}
	return -1, errors.New("User does not exist")
}

func searchMerchant(merchants []merchant, id string) (index int, err error) {
	for i, merchant := range merchants {
		if id == merchant.merchantID {
			return i, nil
		}
	}
	return -1, errors.New("Merchant does not exist")
}
