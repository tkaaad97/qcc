#define WITH_STR(expr) expr, #expr

int assert_int(int expected, int actual, char *expr) {
    if (actual == expected) {
        printf("OK %s => %d\n", expr, actual);
    } else {
        printf("\x1b[31mNG\x1b[0m %s => %d expected, but got %d\n", expr, expected, actual);
        exit(1);
        return 1;
    }
    return 0;
}

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
    return 0;
}
