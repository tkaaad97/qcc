SRCS=$(wildcard *.go)

qcc: $(SRCS)
	go build

test: qcc test.sh
	./test.sh

testc: qcc test/test.c
	cc -E -P -o tmp/test.i test/test.c
	./qcc tmp/test.i > tmp/test.s
	cc -o tmp/test tmp/test.s
	./tmp/test

clean:
	rm -f qcc
	rm -rf tmp

.PHONY: test testc clean
