package test

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Rw(t *testing.T){
	Stkbtc()
}

func TestCreateTransaction(t *testing.T) {
	transaction, err := CreateTransaction("5HusYj2b2x4nroApgfvaSfKYZhRbKFH41bVyPooymbC6KfgSXdD", "1KKKK6N21XKo48zWKuQKXdvSsCf95ibHFa", 91234, "6be2ae9d55572430f633365cad9af2c72a2a0a9d2e13f38339eff333ce3f527f")
	if err != nil {
		fmt.Println(err)
		return
	}
	data, _ := json.Marshal(transaction)
	fmt.Println(string(data))
	fmt.Println(transaction.SignedTx)
}