package merchant

import (
	"errors"

	"github.com/karanrn/pay-later-service/db"
)

// Merchant holds information for merchant
type Merchant struct {
	MerchantID      string
	EmailID         string
	Discount        float64
	DiscountOffered float64
}

// Search searches merchant basis merchantID
func Search(merchantID string) (merchant Merchant, err error) {
	db := db.Connect()
	defer db.Close()

	searchStmt, err := db.Query("SELECT merchant_id, email_id, discount, discount_offered FROM merchant WHERE merchant_id=?", merchantID)
	if err != nil {
		return Merchant{}, errors.New(err.Error())
	}

	merchant = Merchant{}
	for searchStmt.Next() {
		err = searchStmt.Scan(&merchant.MerchantID, &merchant.EmailID, &merchant.Discount, &merchant.DiscountOffered)
		if err != nil {
			return Merchant{}, errors.New(err.Error())
		}
	}
	if (merchant == Merchant{}) {
		return Merchant{}, errors.New("merchant not found")
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
func UpdateDiscount(merchant Merchant, newDiscount float64) (err error) {
	db := db.Connect()
	defer db.Close()

	// Update discount percentage
	_, err = db.Query("UPDATE merchant SET discount = ? WHERE merchant_id = ?", newDiscount, merchant.MerchantID)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// TxnUpdate updates discount offered to the user from the merchant
func TxnUpdate(merchant Merchant) (err error) {
	db := db.Connect()
	defer db.Close()

	// Update discountOffered to the user
	_, err = db.Query("UPDATE merchant SET discount_offered = ? WHERE merchant_id = ?",
		(merchant.DiscountOffered + merchant.Discount), merchant.MerchantID)
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
