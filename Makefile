bin: *.go */*.go go.*
	go build -o $@ -ldflags "-s -w"
	file $@
