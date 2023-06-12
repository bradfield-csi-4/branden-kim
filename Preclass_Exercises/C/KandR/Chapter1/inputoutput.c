#include <stdio.h>

// int main() {
//     int c;

//     while ((c = getchar()) != EOF) {
//         putchar(c);
//     }

//     return 0;
// }

int main() {

    int c;

    while (c = (getchar() != EOF)) {
        if ((c == 0) || (c == 1)) {
            putchar('e');
        }
    }

    return 0;
}

// int main() {

//     int c;

//     while(1){
//         c = getchar();
        
//         if (c == EOF) {
//             putchar(c);
//             break;
//         }
//     }

//     return 0;
// }