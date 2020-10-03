qcc: main.go
	go build

test: qcc
	./test.sh

clean:
	rm -f qcc *.o
	rm -rf tmp

.PHONY: test clean
