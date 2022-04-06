package main
import (
	"ethos/syscall"
	"ethos/altEthos"
	"strconv"
	"log"
)

var path = "/user/"+altEthos.GetUser()+"/accounts/"

 func init() {
	SetupAccountsRpcGetBalance(GetBalance)
	SetupAccountsRpcTransfer(Transfer)
	createAccounts(10)
 }

 func GetBalance(accNum int64)(AccountsRpcProcedure) {
	log.Printf("[SERVER] Accounts server recieved balance request for account number: %v",accNum)
	account, status := GetAccountBalanceUtil(accNum)
	log.Printf("[SERVER] Balance in account %v is %v",accNum, account.balance)

	if status != syscall.StatusOk{
		log.Printf("[SERVER] Account with number %v does not exist",accNum)
		return &AccountsRpcGetBalanceReply{0, status}
	}
	log.Printf("Calling the client reply with %v and %v",account.balance, status)
	return &AccountsRpcGetBalanceReply{account.balance, status} 
 }

 func Transfer(from int64, to int64, amount int64)(AccountsRpcProcedure) {
	log.Printf("[SERVER] Accounts server recieved transfer request from account %v to account %v amount %v",from,to,amount)

	var fromAccount BankAccount
	var toAccount BankAccount
	var status syscall.Status

	//Get the 'from account'
	fromAccount, status = GetAccountBalanceUtil(from)
	if status != syscall.StatusOk {
		log.Printf("[SERVER] Account with number %v does not exist",from)
		return &AccountsRpcTransferReply{status}
	}

	//Get the 'to account'
	toAccount, status = GetAccountBalanceUtil(to)
	if status != syscall.StatusOk {
		log.Printf("[SERVER] Account with number %v does not exist",to)
		return &AccountsRpcTransferReply{status}
	}

	//check if the account has sufficient balance
	if fromAccount.balance < amount {
		log.Printf("[SERVER] Insufficient balance in account %v, balance: %v",from, fromAccount.balance)
		status = syscall.StatusFail
		return &AccountsRpcTransferReply{status}
	}

	// reduce the balance of 'from account'
	fromAccount.balance = fromAccount.balance - amount
	//sace the fromAccount
	status = saveAccount(fromAccount)
	if status != syscall.StatusOk {
		log.Printf("[SERVER] Error updating from account %v while transferring",from)
		return &AccountsRpcTransferReply{status}
	}
	// increase and save the balance of 'to account'
	toAccount.balance = toAccount.balance + amount
	status = saveAccount(toAccount)

	//If increasing the balance of `to account` fails, reset the original balance in the 'from account'
	if status != syscall.StatusOk {
		log.Printf("[SERVER] Error updating to account %v while transferring",to)
		// reset the from account to the original amount. Try 5 times
		fromAccount.balance = fromAccount.balance + amount
		for i :=0 ; i<5; i++ {
			resetstatus := saveAccount(fromAccount)
			if resetstatus == syscall.StatusOk {
				log.Printf("[SERVER] From account %v amount has been reset to original balance %v",from, fromAccount.balance)
				break
			}		
		}
		return &AccountsRpcTransferReply{status}
	}
	return &AccountsRpcTransferReply{status} 
 }

 /*
* Util function to fetch an Account file from the accounts directory
*/
func GetAccountBalanceUtil(accNum int64)(BankAccount, syscall.Status){
	var account BankAccount
	status := altEthos.Read( path+ strconv.Itoa(int(accNum)) , &account)
	return account, status
}

/*
* Util function to save an Account into the accounts directory
*/
func saveAccount(account BankAccount)(syscall.Status){
	status := altEthos.Write( path+ strconv.Itoa(int(account.number)), &account)
	return status
}

/*
* Create the number of accounts specified as the count argument
*/
func createAccounts(count int){

	// create an initial set of accounts with accout numbers and current balance
	//[Note] account owners are not included in the account information currently and clients can access all accounts
	accounts := [6]BankAccount{{100001,10000}, {100002,2345}, {100003,500}, {100004,55000}, {100005,123589}, {100006,66000}}

	if count < 0 {
		log.Fatalf("[SERVER] Please enter a valid value (>0) for number of boxes")
	}
	// check if the accounts storage location is already present.
	status := altEthos.DirectoryCreate(path, &accounts[0], "all")

	if status != syscall.StatusOk {
		log.Fatalf("[SERVER] Error creating root directry for accounts at %v",path)
	}

	for i :=0 ; i<len(accounts); i++ {

		account := accounts[i]

		var accNum int64
		accNum = account.number
		// check if the accNum already exists.
		accPath := path+ strconv.Itoa(int(accNum))
		_, status = altEthos.GetFileInformation(accPath)
		//stop the generation if the account with number does not exist
		if status == syscall.StatusOk{
			log.Fatalf("[SERVER] account with account number %v at already exists",accNum)
			continue
		}

		status = altEthos.Write(accPath, &account)
		
		if status != syscall.StatusOk{
			log.Fatalf("[SERVER] Error writing to %v %v\n",path, status)
		}else {
			log.Printf("[SERVER] Successfully created account %v at %v with balance %v",accNum, path, account.balance)
		}
	}

}

 func main(){
	
	log.Printf("STARTING ACCOUNTS SERVER....")
	altEthos.LogToDirectory( "test/AccountsRpcService" )
	//Advertise the Accounts server
	listeningFd, status := altEthos.Advertise("accountsRpc")
	if status != syscall.StatusOk {
		log.Printf( "Advertising service failed : %s \n" , status )
		altEthos.Exit(status)
	}
   //Listen to new connections being made to the server
	for {
		_ , fd , status := altEthos.Import(listeningFd)
		if status != syscall.StatusOk {
			log.Printf( "Error calling Import : %v\n" , status)
			altEthos.Exit(status)
		}
	   
		log.Printf( "AccountsRpcService : new connection accepted \n" )
	   
		t := AccountsRpc{}
		altEthos.Handle(fd, &t)// hands the request received
	}
 }