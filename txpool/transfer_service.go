package txpool

//todo
type TxPeriodicTransferService struct {
	txRepository      TransactionRepository
	leaderRepository  LeaderRepository
	messageDispatcher MessageDispatcher
}

//todo 이 함수가 call되었을 때 조건에 맞는 tx를 leader에게 전송하는 로직 추가
//todo infra의 timeout_service에 이 함수를 등록, timeout_service가 시간단위로 이 함수를 실행
func (t TxPeriodicTransferService) TransferTxToLeader() {

}