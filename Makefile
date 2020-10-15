SRCS=$(wildcard *.go)

qcc: $(SRCS)
	go build

test: qcc test.sh
	./test.sh

clean:
	rm -f qcc
	rm -rf tmp

.PHONY: test clean
