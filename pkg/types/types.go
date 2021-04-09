package types

// Money представляет собой денежную сумму в минимальных единицах (центы, копейки, дирамы и т.д.)
type Money int64

// Currency представляет собой категорию, в которой был совершён платёж
type Category string

// Status представляет собой статус платежа.
type Status string

//Предопределённые статусы платежей.
const (
	StatusOk Status = "OK"
	StatusFail Status = "FAIL"
	StatusInProgress Status = "INPROGRESS"
)

//Payment представляет информацию о платеже.
type Payment struct{
	ID int
	Amount Money
	Category Category
	Status Status
}

type Phone string

//Account представляет информацию о счёте пользователя
type Account struct {
ID      int64
Phone   Phone
Balance Money

}