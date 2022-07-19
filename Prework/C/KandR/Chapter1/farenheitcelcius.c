#include <stdio.h>

int main() {
    float fahr, celsius;
    float lower, upper, step;

    lower = 0.0;
    upper = 300.0;
    step = 20.0;

    fahr = lower;
    printf("######## Farenheit to Celsius Table ########\n");
    while (fahr <= upper) {
        celsius = (5.0/9.0) * (fahr - 32.0);
        printf("%3.0f %6.2f\n", fahr, celsius);
        fahr = fahr + step;
    }

    printf("######## Celsius to Farenheit Table ########\n");
    celsius = lower;
    
    while (celsius <= upper) {
        fahr = celsius * (9.0 / 5.0) + 32.0;
        printf("%6.2f %3.0f\n", celsius, fahr);
        celsius = celsius + step;
    }
}