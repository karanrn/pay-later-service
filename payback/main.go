package main

import (
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/karanrn/pay-later-service/user"
)

func main() {

	pUser := flag.String("user-id", "", "User paying back. ex: -user-id user101")
	pAmt := flag.Float64("amt", 0.0, "Payback amount. ex: -amt 500")
	flag.Parse()

	if *pUser != "" && *pAmt != 0.0 {
		pbUser, err := user.Search(*pUser)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		// Accept the payback
		err = user.CreditUpdate(pbUser, *pAmt, true)
		if err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println("Both User and payback Amount should be supplied.\nRun payback -h")
	}
}
