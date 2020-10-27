#!/bin/bash

assertProgram() {
  expected="$1"
  input="$2"

  ./qcc "${input}" > tmp/test.s
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
assertExpr() {
  expected="$1"
  input="$2"

  ./qcc "int main(){${input}}" > tmp/test.s
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
#include<stdlib.h>
int foo() { printf("OK\n"); }
int add(int a, int b) { return a + b; }
void alloc4(int **p, int a, int b, int c, int d) {
    *p = (int*)malloc(32);
    (*p)[0] = a;
    (*p)[1] = b;
    (*p)[2] = c;
    (*p)[3] = d;
}
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
assertExpr 2 'int a; int b; int c; a=1;b=1;c=a+b;c;'
assertExpr 5 'int a; int b; int c; int d; a=1;b=2+3;c=-1;a+b;d=a+b+c;'
assertExpr 2 'int foo;int bar;int cdr; foo=1;bar=1;cdr=foo+bar;cdr;'
assertExpr 123 'int a; int b; int c; 1+1;a=1;return 123;b=3;c=a+b;c;'
assertExpr 5 'int a1; int a2; int a3; a1=1;a2=2;a3=a1+a2;a3*2-1;'
assertExpr 0 'if(0)2;'
assertExpr 2 'if(1)2;'
assertExpr 4 'int a; int b; a = 1; b = 2; if(a < b) 4; else 5;'
assertExpr 5 'int a; int b; a = 10; b = 2; if(a < b) 4; else 5;'
assertExpr 3 'int a; int b; a = 1; if(1) a = 2; b = 9; if(0 == 0) b = 1; a + b;'
assertExpr 10 'int b; b = 0; while(b<6) b = b + 5; return b;'
assertExpr 20 'int a; int b; b = 0; for(a = 0; a < 10; a = a + 1) b = b + 2; return b;'
assertExpr 2 'int a; a = 0; { a = a + 1; a = a + 1; return a; }'
assertExpr 42 'if(1){} return 42;'
assertExpr 50 'int a; int b; a = b = 0; while(a < 10){ a = a + 1; b = b + 5; } return b;'
assertStdout "OK" 'int main(){foo();}'
assertExpr 2 'return add(1,1);'
assertExpr 7 'int a; int b; a = 1; b = 2; add(a, b) + 4;'
assertExpr 4 'int a; int* b; a = 0; b = &a; *b = 3; a + 1;'
assertExpr 4 'int *p; alloc4(&p, 1, 2, 4, 8); *(p + 2);'
assertExpr 8 'int *p; alloc4(&p, 1, 2, 4, 8); int *q; q = p + 3; return *q;'

echo OK
