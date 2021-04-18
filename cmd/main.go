package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
	

	for _, account := range accounts {

		_, err=file.Write([]byte (strconv.FormatInt(int64(account.ID),10)))
		if err!=nil{
			return 
		}

		_, err=file.Write([]byte (";"))
		if err!=nil{
			return 
		}

		_, err=file.Write([]byte (account.Phone))
		if err!=nil{
			return 
		}


		_, err=file.Write([]byte ("|"))
		if err!=nil{
			return 
		}
	}


		
}
