package main
import (
 "ethos/altEthos"
 "ethos/syscall"
 "log"
)

func init( ) {
	SetupAccountsRpcGetBalanceReply(GetBalanceReply)
	SetupAccountsRpcTransferReply(TransferReply)
}

func GetBalanceReply(balance int64,status syscall.Status)(AccountsRpcProcedure){
	log.Printf("[Client1] Received balance response from server")
	if status != syscall.StatusOk {
		log.Printf("[Client1]Unable to reterive balance from server")
		return nil
	}
	log.Printf("[Client1] Received balance reply with balance= %v", balance)
	return nil
}

func TransferReply(status syscall.Status)(AccountsRpcProcedure){
	if status != syscall.StatusOk {
		log.Printf("[Client1] Server was unable to transfer the amount")
		return nil
	}
	log.Printf("[Client1] Received transfer success reply from server")
	return nil
}

func main() {

	log.Printf("STARTING ACCOUNTS CLIENT1....")
	altEthos.LogToDirectory( "test/AccountsRpcClient" )
	
	transfercalls := getTransferCalls() // get a list of transfer calls to be made to the server
	balanceCalls := getBalanceCalls() // get a list of getBalance calls to be made to the server

	log.Printf("[Client1] AccountsRpcClient : beforecall \n" )

	for i:=0 ; i <len(transfercalls); i++ {
		// Make the transfer call
		fd, status := altEthos.IpcRepeat("accountsRpc", "" , nil)
		if status != syscall.StatusOk {
			log.Printf( "Ipc failed : %v \n" , status )
			altEthos.Exit(status)
		}
	
		call := transfercalls[i]
		status = altEthos.ClientCall(fd , &call )
		if status != syscall.StatusOk {
			log.Printf( "[Client1]  clientCall failed : %v \n" , status )
			altEthos.Exit( status )
		}

		altEthos.Close( fd )
		//Make the balance check
		fd, status = altEthos.IpcRepeat("accountsRpc", "" , nil)
		if status != syscall.StatusOk {
			log.Printf( "Ipc failed : %v \n" , status )
			altEthos.Exit(status)
		}

		call2 := balanceCalls[i]
		status = altEthos.ClientCall(fd , &call2 )
		if status != syscall.StatusOk {
			log.Printf( "[Client1]  clientCall failed : %v \n" , status )
			altEthos.Exit( status )
		}

	}

	log.Printf( "[Client1] AccountsRpcClient : done \n" )
 }

 /*
 Get a list of transfer calls made by this client
 */
 func getTransferCalls() []AccountsRpcTransfer {
	calls := []AccountsRpcTransfer{
		{100001,100002,500},{100001,100003,300},{100007,100004,5000},{100005,100003,25000},{100002,100003,6000},{100005,100004,2500},{100004,100002,300},{100007,100002,4400},{100003,100006,15000},{100003,100009,500},{100001,100005,3601},{100006,100001,478},{100004,100005,25},{100008,100001,40},{100002,100002,30}}
	return calls
 }

 /*
 Get a list of balance calls made by this client
 */
 func getBalanceCalls() []AccountsRpcGetBalance {
	calls := []AccountsRpcGetBalance{{100001},{100001},{100007},{100005},{100002},{100005},{100004},{100007},{100003},{100003},{100001},{100006},{100004},{100008},{100002}}
	return calls
}