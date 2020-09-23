package merchant

import (
	"errors"

	"github.com/karanrn/pay-later-service/db"
)

// Merchant holds information for merchant
type Merchant struct {
	MerchantID string
	EmailID    string
	Discount   int64
}

// Search searches merchant basis merchantID
func Search(merchantID string) (merchant Merchant, err error) {
	db := db.Connect()
	defer db.Close()

	searchStmt, err := db.Query("SELECT merchant_id, email_id, discount FROM merchant WHERE merchant_id=?", merchantID)
	if err != nil {
		return Merchant{}, errors.New(err.Error())
	}

	merchant = Merchant{}
	for searchStmt.Next() {
		err = searchStmt.Scan(&merchant.MerchantID, &merchant.EmailID, &merchant.Discount)
		if err != nil {
			return Merchant{}, errors.New(err.Error())
		}
	}
	return merchant, nil
}

// Add adds merchant to the database
func Add(merchant Merchant) (err error) {
	db := db.Connect()
	defer db.Close()

	insertStmt, err := db.Prepare("INSERT INTO merchant (merchant_id, email_id, discount) VALUES (?, ?, ?)")
	if err != nil {
		return errors.New(err.Error())
	}
	_, err = insertStmt.Exec(merchant.MerchantID, merchant.EmailID, merchant.Discount)
	if err != nil {
		return errors.New(err.Error())
	}
	return nil
}

// UpdateDiscount updates discount offer of the merchant
func UpdateDiscount(merchantID string, newDiscount int64) (err error) {
	db := db.Connect()
	defer db.Close()

	selectStmt, err := db.Query("SELECT merchant_id, email_id, discount FROM merchant WHERE merchant_id=?", merchantID)
	if err != nil {
		return errors.New(err.Error())
	}

	merchant := Merchant{}
	for selectStmt.Next() {
		err = selectStmt.Scan(&merchant.MerchantID, &merchant.EmailID, &merchant.Discount)
		if err != nil {
			return errors.New(err.Error())
		}
	}

	// Update discount percentage
	_, err = db.Query("UPDATE merchant SET discount = ? WHERE merchant_id = ?", newDiscount, merchant.MerchantID)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
