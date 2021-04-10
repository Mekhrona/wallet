package wallet

import (
	"testing"
	"github.com/Mekhrona/wallet/pkg/types"
)


func TestExistingID(t *testing.T) {

	accountss := [] *types.Account {
		{
			ID: 56,
		},
		{
			ID: 79,
		},
	}
	svc:=  Service{
		accounts: accountss ,
		payments: nil,
		}
	
	resultAccount, err:=svc.FindAccountByID(56)
	
	if resultAccount.ID!=56{
		t.Errorf( "invalid result")
	}	
	
	if err!=nil{	
		t.Errorf( "invalid result")
	}
}

func TestPaymentFound(t *testing.T) {

	accounts := [] *types.Account {
		{
			ID: 56,
			Balance: 1000,
		},
		{
			ID: 79,
			Balance: 20000,
		},
		{
			ID: 85,
			Balance: 9800,
		},
	}

	svc:=  Service{
		accounts: accounts ,
		payments: nil,
		}
	payment,err:=svc.Pay(85, 1000,"clothes")
	
	if err!=nil && payment==nil{	
		t.Errorf( "invalid result")
	}

    err1:=svc.Reject(payment.ID)

	if err1!=nil{	
		t.Errorf( "invalid result")
	}	
}

func TestPaymentNotFound(t *testing.T) {

	accounts := [] *types.Account {
		{
			ID: 56,
			Balance: 1000,
		},
		{
			ID: 79,
			Balance: 20000,
		},
		{
			ID: 85,
			Balance: 9800,
		},
	}

	svc:=  Service{
		accounts: accounts ,
		payments: nil,
		}
	payment,err:=svc.Pay(85, 1000,"clothes")
	
	if err!=nil && payment==nil{	
		t.Errorf( "invalid result")
	}

    err1:=svc.Reject("58")

	if err1==nil{	
		t.Errorf( "invalid result")
	}	
}