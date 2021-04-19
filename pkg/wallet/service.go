package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Mekhrona/wallet/pkg/types"
	"github.com/google/uuid"
)


var ErrPhoneRegistered=errors.New("phone already registered")
var ErrAmountMustBePositive=errors.New("amount must be greater than zero")
var ErrAccountNotFound=errors.New("account not found")
var ErrNotENoughBalance=errors.New("balance is not enough")
var ErrPaymentNotFound=errors.New("payment not found")
var ErrFavoriteNotFound=errors.New("favorite payment not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}
 
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
      if account.Phone==phone{
		  return nil, ErrPhoneRegistered
	  }	
	}

	s.nextAccountID=int64(uuid.New().ID())
	account:= &types.Account{
		ID:  s.nextAccountID,
		Phone:  phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts,account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error{
    if amount<=0{
		return ErrAmountMustBePositive
	}
	
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID==accountID{
			account=acc
			break
		}
		
	}

	if account==nil{
		return ErrAccountNotFound
	}

	account.Balance+=amount
	return nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID==accountID{
			account=acc
			break
		}
		
	}

	if account==nil{
		return nil,ErrAccountNotFound
	}
	return account,nil

}


func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error){
  if amount<=0{
	  return nil, ErrAmountMustBePositive
  }
  
  var account *types.Account
  for _, acc := range s.accounts {
	  if acc.ID==accountID{
		  account=acc
		  break
	  }				  
  }
  
  if account==nil {
	  return nil, ErrAccountNotFound
  }

  if account.Balance<amount{
	  return nil, ErrNotENoughBalance
  }

  account.Balance-=amount
  paymentID:=uuid.New().String()
  payment:=&types.Payment{
	  ID: paymentID,
	  AccountID: accountID,
	  Amount: amount,
	  Category: category,
	  Status: types.PaymentStatusInProgress,
  }

  s.payments=append(s.payments, payment)
  return payment,nil

}

func (s *Service) Reject(paymentID string) error{ 

	var payment *types.Payment
    for _, pay := range s.payments {
		if pay.ID==paymentID{
			payment=pay
			break
		}	
	}
	if payment==nil{
		return ErrAccountNotFound
	}

	payment.Status=types.PaymentStatusFail

	account, err:=s.FindAccountByID(payment.AccountID)
	if err!=nil{
		return ErrAccountNotFound
	}
	account.Balance+=payment.Amount
	return nil
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {

	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID==paymentID{
			payment=pay
			break
		}	
	}

	if payment==nil{
		return nil,ErrPaymentNotFound
	}
	return payment,nil
}


//Repeat fuctions allows to repeat already made payment 
func (s *Service)  Repeat(paymentID string) (*types.Payment, error){
  targetPayment,err:=s.FindPaymentByID(paymentID)
  if err != nil {
	  return nil, err
  }

  var account *types.Account
  for _, acc := range s.accounts {
	  if acc.ID==targetPayment.AccountID{
		  account=acc
		  break
	  }				  
  }
  
  if account==nil {
	  return nil, ErrAccountNotFound
  }

  if account.Balance<targetPayment.Amount{
	  return nil, ErrNotENoughBalance
  }

  account.Balance-=targetPayment.Amount

  newPaymentID:=uuid.New().String()
  newPayment:=&types.Payment{
	  ID: newPaymentID,
	  AccountID: targetPayment.AccountID,
	  Amount: targetPayment.Amount,
	  Category: targetPayment.Category,
	  Status: types.PaymentStatusInProgress,
  }

  s.payments=append(s.payments, newPayment)
  return newPayment,nil
}


func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error){
  targetPayment,err:=s.FindPaymentByID(paymentID)
  if err != nil {
	  return nil, ErrPaymentNotFound
  }
  newPaymentID:=uuid.New().String()
  favPayment:=&types.Favorite{
	  ID: newPaymentID,
	  AccountID: targetPayment.AccountID,
	  Name: name,
	  Amount: targetPayment.Amount,
	  Category: targetPayment.Category,

  }
 s.favorites=append(s.favorites, favPayment)
  return favPayment, nil
}


func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID==favoriteID{
			favorite=fav
			break
		}
		}

		if favorite==nil{
			return nil,ErrFavoriteNotFound
		}

		payment,err:=s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)

		return payment,err
}


//ExportToFile экспортирует аккаунт в файл
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)

	if err != nil {
		log.Print(err)

		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	for _, account := range s.accounts {
		ID := strconv.FormatInt(int64(account.ID), 10) + ";"
		phone := string(account.Phone) + ";"
		balance := strconv.FormatInt(int64(account.Balance), 10)
		text:=ID + phone + balance + "|"

		_, err = file.Write([]byte(text))

		if err != nil {
			log.Print(err)

			return err
		}
	}

	return err
}


func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)

	if err != nil {
		log.Print(err)

		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			log.Print(err)
		}
	}()

	result := make([]byte, 0)
	bufferf := make([]byte, 4)

	for {
		read, err := file.Read(bufferf)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)

			return err
		}

		result = append(result, bufferf[:read]...)
	}

	data := string(result)

	for _, line := range strings.Split(data, "|") {
		if len(line) == 0 {
			return err
		}

		item := strings.Split(line, ";")
		ID, err := strconv.ParseInt(item[0], 10, 64)

		if err != nil {
			return err
		}

		balance, err := strconv.ParseInt(item[2], 10, 64)

		if err != nil {
			return err
		}

		s.accounts = append(s.accounts, &types.Account{
			ID:      ID,
			Phone:   types.Phone(item[1]),
			Balance: types.Money(balance),
		})
	}

	return nil
}