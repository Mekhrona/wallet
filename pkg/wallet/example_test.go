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
	svc:= Service {
		accounts: accountss ,
		}
	
 
	
	resultAccount, err:=svc.FindAccountByID(56)
	
	if resultAccount.ID!=56{
		t.Errorf( "invalid result")
	}
		
	
	if err!=nil{	
		t.Errorf( "invalid result")
	}


	
}