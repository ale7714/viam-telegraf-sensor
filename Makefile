bin: *.go */*.go go.*
	go build -o $@ -ldflags "-s -w" -tags osusergo,netgo
	file $@

module.tar.gz: bin
	tar -czf module.tar.gz run.sh bin viam-telegraf.conf

clean:
	rm -rf bin/* 
	rm -f module*.tar.gz  