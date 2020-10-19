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

assert 0 '0;'
assert 42 '42;'
assert 21 '5+20-4;'
assert 47 '5+6*7;'
assert 15 '5*(9-6);'
assert 4 '(3+5)/2;'
assert 10 '-10+20;'
assert 1 '(-1-2)/-3;'
assert 1 '99==99;'
assert 0 '1==0;'
assert 0 '99!=99;'
assert 1 '1!=0;'
assert 1 '(5*(9-6)==15)==1;'
assert 1 '-1<989;'
assert 0 '1000-1<989;'
assert 1 '-1<=989;'
assert 0 '1000-1<=989;'
assert 1 '1000-11<=989;'
assert 1 '1+1+1+1>1;'
assert 0 '1>1+1+1+1;'
assert 1 '1+1+1+1>=1;'
assert 0 '1>=1+1+1+1;'
assert 1 '1+1+1+1>=1+1+1+1;'
assert 2 'a=1;b=1;c=a+b;c;'
assert 5 'a=1;b=2+3;c=-1;a+b;d=a+b+c;'
assert 2 'foo=1;bar=1;cdr=foo+bar;cdr;'
assert 123 '1+1;a=1;return 123;b=3;c=a+b;c;'
assert 5 'a1=1;a2=2;a3=a1+a2;a3*2-1;'
assert 0 'if(0)2;'
assert 2 'if(1)2;'
assert 4 'a = 1; b = 2; if(a < b) 4; else 5;'
assert 5 'a = 10; b = 2; if(a < b) 4; else 5;'
assert 3 'a = 1; if(1) a = 2; b = 9; if(0 == 0) b = 1; a + b;'
assert 10 'b = 0; while(b<6) b = b + 5; return b;'

echo OK
