// ~~~line comment~~~
/**
 * multi
 * line
 * comment
 */
#define WITH_STR(expr) expr, #expr
#define GEN_CASE(name, expected, body) int test_case_function_ ## name () { body }\
    int exec_test_case_ ## name () { assert_int(expected, test_case_function_ ## name(), #body); }
#define EXEC_CASE(name) exec_test_case_ ## name ()

/**
 * some functions
 */
int foo() { printf("OK\n"); }

int add(int a, int b) { return a + b; }

int alloc4(int **p, int a, int b, int c, int d) {
    *p = malloc(32);
    (*p)[0] = a;
    (*p)[1] = b;
    (*p)[2] = c;
    (*p)[3] = d;
}

/**
 * assert_int
 */
int assert_int(int expected, int actual, char *expr) {
    if (actual == expected) {
        printf("OK %s => %d\n", expr, actual);
    } else {
        printf("NG %s => %d expected, but got %d\n", expr, expected, actual);
        exit(1);
        return 1;
    }
    return 0;
}

GEN_CASE(t1, 2, int a; int b; int c; a=1;b=1;c=a+b;c;)
GEN_CASE(t2, 5, int a; int b; int c; int d; a=1;b=2+3;c=-1;a+b;d=a+b+c;)
GEN_CASE(t3, 2, int foo;int bar;int cdr; foo=1;bar=1;cdr=foo+bar;cdr;)
GEN_CASE(t4, 123, int a; int b; int c; 1+1;a=1;return 123;b=3;c=a+b;c;)
GEN_CASE(t5, 5, int a1; int a2; int a3; a1=1;a2=2;a3=a1+a2;a3*2-1;)
GEN_CASE(t6, 0, if(0)2;)
GEN_CASE(t7, 2, if(1)2;)
GEN_CASE(t8, 4, int a; int b; a = 1; b = 2; if(a < b) 4; else 5;)
GEN_CASE(t9, 5, int a; int b; a = 10; b = 2; if(a < b) 4; else 5;)
GEN_CASE(t10, 3, int a; int b; a = 1; if(1) a = 2; b = 9; if(0 == 0) b = 1; a + b;)
GEN_CASE(t11, 10, int b; b = 0; while(b<6) b = b + 5; return b;)
GEN_CASE(t12, 20, int a; int b; b = 0; for(a = 0; a < 10; a = a + 1) b = b + 2; return b;)
GEN_CASE(t13, 2, int a; a = 0; { a = a + 1; a = a + 1; return a; })
GEN_CASE(t14, 42, if(1){} return 42;)
GEN_CASE(t15, 50, int a; int b; a = b = 0; while(a < 10){ a = a + 1; b = b + 5; } return b;)
GEN_CASE(t16, 2, return add(1,1);)
GEN_CASE(t17, 7, int a; int b; a = 1; b = 2; add(a, b) + 4;)
GEN_CASE(t18, 4, int a; int* b; a = 0; b = &a; *b = 3; a + 1;)
GEN_CASE(t19, 4, int *p; alloc4(&p, 1, 2, 4, 8); *(p + 2);)
GEN_CASE(t20, 8, int *p; alloc4(&p, 1, 2, 4, 8); int *q; q = p + 3; return *q;)
GEN_CASE(t21, 4, sizeof(1);)
GEN_CASE(t22, 4, int a; a = 2; sizeof(a);)
GEN_CASE(t23, 8, int a; a = 3; sizeof(&a);)
int fun(int a, int b) { return a + b; }
GEN_CASE(t24, 4, return sizeof(fun(1, 2));)
GEN_CASE(t25, 42, int a[3]; int *b; int c; c = 41; b = &c; *b = (*b) + 1; return c;)
GEN_CASE(t26, 2, int a[3]; *a = 2; return *a;)
GEN_CASE(t27, 2, int a[3]; *(a + 1) = 2; return *(a + 1);)
GEN_CASE(t28, 3, int a[2]; a[0] = 2; a[1] = 1; return a[0] + a[1];)
int x; int y;
GEN_CASE(t29, 11, x = 2; y = 9; return x + y;)
GEN_CASE(t30, 179, char x[3]; x[0] = -1; x[1] = 2; int y; y = 180; return y + x[0];)
char gx[3];
GEN_CASE(t31, 9, gx[0] = -2; gx[1] = 1; gx[2] = 10; return gx[0] + gx[1] + gx[2];)
GEN_CASE(t32, 104, char *a; a = "hello"; a[0];)

int main() {
    printf("start test\n");
    assert_int(0, WITH_STR(0));
    assert_int(42, WITH_STR(42));
    assert_int(21, WITH_STR(5+20-4));
    assert_int(47, WITH_STR(5+6*7));
    assert_int(15, WITH_STR(5*(9-6)));
    assert_int(4, WITH_STR((3+5)/2));
    assert_int(10, WITH_STR(-10+20));
    assert_int(1, WITH_STR((-1-2)/-3));
    assert_int(1, WITH_STR(99==99));
    assert_int(0, WITH_STR(1==0));
    assert_int(0, WITH_STR(99!=99));
    assert_int(1, WITH_STR(1!=0));
    assert_int(1, WITH_STR((5*(9-6)==15)==1));
    assert_int(1, WITH_STR(-1<989));
    assert_int(0, WITH_STR(1000-1<=989));
    assert_int(1, WITH_STR(1000-11<=989));
    assert_int(1, WITH_STR(1+1+1+1>1));
    assert_int(0, WITH_STR(1>1+1+1+1));
    assert_int(1, WITH_STR(1+1+1+1>=1));
    assert_int(0, WITH_STR(1>=1+1+1+1));
    assert_int(1, WITH_STR(1+1+1+1>=1+1+1+1));

    EXEC_CASE(t1);
    EXEC_CASE(t2);
    EXEC_CASE(t3);
    EXEC_CASE(t4);
    EXEC_CASE(t5);
    EXEC_CASE(t6);
    EXEC_CASE(t7);
    EXEC_CASE(t8);
    EXEC_CASE(t9);
    EXEC_CASE(t10);
    EXEC_CASE(t11);
    EXEC_CASE(t12);
    EXEC_CASE(t13);
    EXEC_CASE(t14);
    EXEC_CASE(t15);
    EXEC_CASE(t16);
    EXEC_CASE(t17);
    EXEC_CASE(t18);
    EXEC_CASE(t19);
    EXEC_CASE(t20);
    EXEC_CASE(t21);
    EXEC_CASE(t22);
    EXEC_CASE(t23);
    EXEC_CASE(t24);
    EXEC_CASE(t25);
    EXEC_CASE(t26);
    EXEC_CASE(t27);
    EXEC_CASE(t28);
    EXEC_CASE(t29);
//    EXEC_CASE(t30);
//    EXEC_CASE(t31);
//    EXEC_CASE(t32);
    return 0;
}
