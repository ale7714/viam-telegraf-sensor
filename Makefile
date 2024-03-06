bin: *.go */*.go go.*
	go build -o $@ -ldflags "-s -w" -tags osusergo,netgo
	file $@

module.tar.gz: clean bin
	tar -czf module.tar.gz run.sh bin viam-telegraf.conf

clean:
	rm -f bin
	rm -f module*.tar.gz  