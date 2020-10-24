#!/bin/bash

assertExpr() {
  expected="$1"
  input="$2"

  ./qcc "main(){${input}}" > tmp/test.s
  cc -o tmp/test tmp/test.s tmp/external.o
  ./tmp/test
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assertStdout() {
  expected="$1"
  input="$2"

  ./qcc "$input" > tmp/test.s
  cc -o tmp/test tmp/test.s tmp/external.o
  actual="$(./tmp/test)"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

mkdir -p tmp

echo '#include <stdio.h>
int foo() { printf("OK\n"); }
int add(int a, int b) { return a + b; }
' > tmp/external.c
cc -c -o tmp/external.o tmp/external.c

assertExpr 0 '0;'
assertExpr 42 '42;'
assertExpr 21 '5+20-4;'
assertExpr 47 '5+6*7;'
assertExpr 15 '5*(9-6);'
assertExpr 4 '(3+5)/2;'
assertExpr 10 '-10+20;'
assertExpr 1 '(-1-2)/-3;'
assertExpr 1 '99==99;'
assertExpr 0 '1==0;'
assertExpr 0 '99!=99;'
assertExpr 1 '1!=0;'
assertExpr 1 '(5*(9-6)==15)==1;'
assertExpr 1 '-1<989;'
assertExpr 0 '1000-1<989;'
assertExpr 1 '-1<=989;'
assertExpr 0 '1000-1<=989;'
assertExpr 1 '1000-11<=989;'
assertExpr 1 '1+1+1+1>1;'
assertExpr 0 '1>1+1+1+1;'
assertExpr 1 '1+1+1+1>=1;'
assertExpr 0 '1>=1+1+1+1;'
assertExpr 1 '1+1+1+1>=1+1+1+1;'
assertExpr 2 'a=1;b=1;c=a+b;c;'
assertExpr 5 'a=1;b=2+3;c=-1;a+b;d=a+b+c;'
assertExpr 2 'foo=1;bar=1;cdr=foo+bar;cdr;'
assertExpr 123 '1+1;a=1;return 123;b=3;c=a+b;c;'
assertExpr 5 'a1=1;a2=2;a3=a1+a2;a3*2-1;'
assertExpr 0 'if(0)2;'
assertExpr 2 'if(1)2;'
assertExpr 4 'a = 1; b = 2; if(a < b) 4; else 5;'
assertExpr 5 'a = 10; b = 2; if(a < b) 4; else 5;'
assertExpr 3 'a = 1; if(1) a = 2; b = 9; if(0 == 0) b = 1; a + b;'
assertExpr 10 'b = 0; while(b<6) b = b + 5; return b;'
assertExpr 20 'b = 0; for(a = 0; a < 10; a = a + 1) b = b + 2; return b;'
assertExpr 2 'a = 0; { a = a + 1; a = a + 1; return a; }'
assertExpr 42 'if(1){} return 42;'
assertExpr 50 'a = b = 0; while(a < 10){ a = a + 1; b = b + 5; } return b;'
assertStdout "OK" 'main(){foo();}'
assertExpr 2 'return add(1,1);'
assertExpr 7 'a = 1; b = 2; add(a, b) + 4;'

echo OK
