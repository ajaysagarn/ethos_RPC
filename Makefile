export UDIR= .
export GOC = x86_64-xen-ethos-6g
export GOL = x86_64-xen-ethos-6l
export ETN2GO = etn2go
export ET2G   = et2g
export EG2GO  = eg2go

export GOARCH = amd64
export TARGET_ARCH = x86_64
export GOETHOSINCLUDE=/usr/lib64/go/pkg/ethos_$(GOARCH)
export GOLINUXINCLUDE=/usr/lib64/go/pkg/linux_$(GOARCH)

export ETHOSROOT=server/rootfs
export MINIMALTDROOT=server/minimaltdfs

.PHONY: all install
all: AccountService AccountClient1

accountsRpc.go: accountsRpc.t
	$(ETN2GO) . accountsRpc main $^

AccountService: AccountService.go accountsRpc.go
	ethosGo $^ 

AccountClient1: AccountClient1.go accountsRpc.go
	ethosGo $^ 

AccountClient2: AccountClient2.go accountsRpc.go
	ethosGo $^ 

# install types, service
install: AccountService AccountClient1 AccountClient2
	sudo rm -rf server/
	(ethosParams server && cd server/ && ethosMinimaltdBuilder)
	echo 7 > server/param/sleepTime
	ethosTypeInstall accountsRpc
	ethosDirCreate $(ETHOSROOT)/services/accountsRpc $(ETHOSROOT)/types/spec/accountsRpc/AccountsRpc all
	install -D AccountService AccountClient1 AccountClient2 $(ETHOSROOT)/programs
	ethosStringEncode /programs/AccountService > $(ETHOSROOT)/etc/init/services/AccountService
	ethosStringEncode /programs/AccountClient1 > $(ETHOSROOT)/etc/init/services/AccountClient1
	ethosStringEncode /programs/AccountClient2 > $(ETHOSROOT)/etc/init/services/AccountClient2


# remove build artifacts
clean:
	rm -rf accountsRpc/ accountsRpcIndex/
	rm -f accountsRpc.go
	rm -f AccountClient1
	rm -f AccountClient2
	rm -f AccountService
	rm -f AccountClient1.goo.ethos
	rm -f AccountClient2.goo.ethos
	rm -f AccountService.goo.ethos
	sudo rm -rf server/
