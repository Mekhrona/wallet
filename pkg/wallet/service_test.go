package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Mekhrona/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct{
	*Service
}

func newTestService() *testService{
	return &testService{Service: &Service{}}
}

type testAccount struct{
	phone   types.Phone
	balance  types.Money
	payments []struct{
		amount  types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount=testAccount{
	phone: "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 1000_00, category:"auto"},
	},
}


func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment,error){
	//Registering user
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil,nil, fmt.Errorf("can't register account, error=%v", err)
	}

	//adding money to user's account
	err=s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil,nil, fmt.Errorf("can't deposit account, error=%v",err)
	}
	

	payments:=make([]*types.Payment, len(data.payments))
	for i, payment:= range data.payments{
		payments[i],err=s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error=%v", err)
		}
	}
	return account, payments, nil
}



func TestFindPaymentByID_success(t *testing.T) {
	//creating service
	s:=newTestService()
	_, payments, err:= s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}

	payment:=payments[0]
	got,err:=s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPAymentByID(): error=%v", err)
		return 	
	}

	if !reflect.DeepEqual(payment,got){
		t.Errorf("FindPAymentByID():wrong payment returned =%v", err)
	}
}


func TestService_FindPAymentByID_fail(t *testing.T) {
	//creating service
	s:=newTestService()
	_, _, err:= s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}	
	
	_, err1:=s.FindPaymentByID(uuid.New().String())
	if err1 == nil {
		t.Errorf("FindPAymentByID(): must return error, returned nil")
		return 	
	}

	if err1!=ErrPaymentNotFound{
		t.Errorf("FindPAymentByID():must return ErrPAymentNotFound, returned =%v", err1)
	}
}

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

func TestRepeat_success(t *testing.T) {
	s:=newTestService()
	_, payments, err:= s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}
	payment:=payments[0]
    repeatedPayment, err :=s.Repeat(payment.ID)
	
	if err != nil {
		t.Error(err)
		return 
	}
  if repeatedPayment.AccountID!=payment.AccountID{
	  t.Errorf("Two payments accountID's must be same, but expected %v and result %v ", repeatedPayment.AccountID,payment.AccountID )
  }
	
  if repeatedPayment.ID==payment.ID{
	t.Errorf("Two payments ID's must be different, but expected %v and result %v ", repeatedPayment.ID,payment.ID )
}
}


func TestFavoritePayment_success(t *testing.T) {
	s:=newTestService()
	_, payments, err:= s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}
	payment:=payments[0]
    favPayment, err :=s.FavoritePayment(payment.ID, "mobile")
	
	if err != nil {
		t.Error(err)
		return 
	}
   if favPayment.AccountID!=payment.AccountID{
	   t.Errorf("Favorite acccount ID, %v, is not matched with expected,%v",favPayment.AccountID,payment.AccountID )
   }

	
}

func TestPayFromFavorite_success(t *testing.T) {
	s:=newTestService()
	_, payments, err:= s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return 
	}
	payment:=payments[0]
    favPayment, err :=s.FavoritePayment(payment.ID, "mobile")
	
	if err != nil {
		t.Error(err)
		return 
	}

	resultPayment,err:=s.PayFromFavorite(favPayment.ID)
	if err != nil {
		t.Error(err)
		return 
	}

	if resultPayment.AccountID!=payment.AccountID{
		t.Errorf("Favorite payment ID,%v, and resulting payment ID do not match,%v", resultPayment.ID,payment.ID)
	}
	
}

