VERSION=0.0.9

bin: bin/telegraf-sensor-linux bin/telegraf-sensor-darwin

bin/telegraf-sensor-linux: *.go */*.go go.*
	GOOS="linux" go build -o $@ -ldflags "-s -w"
	file $@

bin/telegraf-sensor-darwin: *.go */*.go go.*
	GOOS="darwin" go build -o $@ -ldflags "-s -w"
	file $@

module-upload: clean bin
	tar -czf module-linux.tar.gz run.sh bin/telegraf-sensor-linux viam-telegraf.conf
	tar -czf module-darwin.tar.gz run.sh bin/telegraf-sensor-darwin viam-telegraf.conf
	viam module upload --version $(VERSION) --platform darwin/any module-darwin.tar.gz
	viam module upload --version $(VERSION) --platform linux/any module-linux.tar.gz

clean:
	rm -rf bin/* 
	rm -f module*.tar.gz  