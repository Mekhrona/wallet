package main

import (
	"fmt"

	"github.com/Mekhrona/wallet/pkg/types"
	"github.com/Mekhrona/wallet/pkg/wallet"
)



func main() {
     svc:=&wallet.Service{}
	account1,err1 :=svc.RegisterAccount("+992000000001")
	account2, err2:=svc.RegisterAccount("+992000000008")
	account3, err3:=svc.RegisterAccount("+992000000074")
	accounts:=[] *types.Account{account1,account2, account3}

	for _, acc := range accounts {
		fmt.Println(acc.ID)
	}
	
	fmt.Println(err1,err2,err3)

	account,err:=svc.FindAccountByID(2)

	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(account.ID)

	accountNot,errNot:=svc.FindAccountByID(5)

	if errNot!=nil{
		fmt.Println(err)
	}
	
	fmt.Println(accountNot.ID)

}
