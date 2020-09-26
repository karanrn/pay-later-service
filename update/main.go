package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/karanrn/pay-later-service/merchant"
)

func main() {
	merch := flag.NewFlagSet("merchant", flag.ContinueOnError)
	merchID := merch.String("id", "", "Merchant id of merchant updating. ex: -id m101")
	discount := merch.String("discount", "0%", "Discount value to be updated. ex: -discount 10%")

	if len(os.Args[1:]) < 1 {
		fmt.Println("You must pass sub command - [merchant]")
		return
	}

	switch os.Args[1] {
	case "merchant":
		if err := merch.Parse(os.Args[2:]); err == nil {
			if *merchID != "" && *discount != "0%" {
				updateMerchant, err := merchant.Search(*merchID)
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				disAmt, err := strconv.ParseFloat(strings.TrimRight(*discount, "%"), 64)
				if err != nil {
					fmt.Println("Discount amount is not in a valid format.")
					return
				}

				err = merchant.UpdateDiscount(updateMerchant, disAmt)
				if err != nil {
					fmt.Println(err.Error())
				}
			} else {
				fmt.Println("Pass merchID and Discount for update")
			}
		}
	}
}
