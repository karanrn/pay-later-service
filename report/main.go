package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/karanrn/pay-later-service/merchant"
	"github.com/karanrn/pay-later-service/user"
)

func main() {
	discount := flag.NewFlagSet("discount", flag.ContinueOnError)
	merchID := discount.String("id", "", "Merchant Id for whom discount report will be generated. ex: -id m101")

	dues := flag.NewFlagSet("dues", flag.ContinueOnError)
	userID := dues.String("id", "", "User ID for whom dues report will be generated. ex: -id u2")

	if len(os.Args[1:]) < 1 {
		fmt.Println("You must pass sub command - [discount, dues, total-dues, users-at-credit-limit]")
		return
	}

	switch os.Args[1] {
	case "discount":
		if err := discount.Parse(os.Args[2:]); err == nil {
			if *merchID != "" {
				reportMerchant, err := merchant.Search(*merchID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Printf("%.2f\n", reportMerchant.DiscountOffered)
			} else {
				fmt.Println("merchant ID needs to be supplied.")
			}
		}
	case "dues":
		if err := dues.Parse(os.Args[2:]); err == nil {
			if *userID != "" {
				reportUser, err := user.Search(*userID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				fmt.Printf("%.2f\n", reportUser.CreditSpent)
			} else {
				fmt.Println("user ID needs to be supplied.")
			}
		}
	}
}
