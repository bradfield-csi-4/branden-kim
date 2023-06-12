.data

.equ SYS_EXIT, 93
.equ SYS_WRITE, 64
.equ STDOUT, 1
.equ SUCCESS, 0

MSG:
    .ascii "Hello world!\n"

.text

/* Our application's entry point. */
.globl _start
_start:
    /* syscall write(int fd, const void *buf, size_t count) */
    mov     x0, #1      /* fd := STDOUT_FILENO */
    adr     x1, =MSG    /* buf := msg */
    mov     x2, #13    /* count := len */
    mov     w8, #64     /* write is syscall #64 */
    svc     #0          /* invoke syscall */

    /* syscall exit(int status) */
    mov     x0, #0      /* status := 0 */
    mov     w8, #93     /* exit is syscall #93 */
    svc     #0          /* invoke syscall */
