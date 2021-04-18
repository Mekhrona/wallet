package main

import (
	"fmt"
	"log"
	"os"

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


	accountNot,errNot:=svc.FindAccountByID(account2.ID)

	if errNot!=nil{
		fmt.Println(errNot)
	}
	
	fmt.Println(&accountNot.ID)


	file, err:=os.Create("cmd/readme.txt")
	if err!=nil{
		return 
	}
	defer func(){
		err:=file.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	

errFile:=svc.ExportToFile("data/accounts.txt")

fmt.Print(errFile)


		
}
