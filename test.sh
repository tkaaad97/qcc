#!/bin/bash

assert() {
  expected="$1"
  input="$2"

  ./qcc "$input" > tmp/test.s
  cc -o tmp/test tmp/test.s
  ./tmp/test
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

mkdir -p tmp

assert 0 0
assert 42 42

echo OK
