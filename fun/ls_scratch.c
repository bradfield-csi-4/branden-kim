#include <stdio.h>
#include <sys/syscall.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h> 
#include <dirent.h>     /* Defines DT_* constants */
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h> 
extern int errno;

struct linux_dirent {
    unsigned long  d_ino;     /* Inode number */
    unsigned long  d_off;     /* Offset to next linux_dirent */
    unsigned short d_reclen;  /* Length of this linux_dirent */
    char           d_name[];  /* Filename (null-terminated) */
                        /* length is actually (d_reclen - 2 -
                           offsetof(struct linux_dirent, d_name) */
    /*
    char           pad;       // Zero padding byte
    char           d_type;    // File type (only since Linux 2.6.4;
                              // offset is (d_reclen - 1))
    */
};

#define BUFFER_SIZE 1024
#define MAX_ENTRIES 500

int main(int argc, char **argv) {
    // Variable Declarations 
    int fd, nread, nlines;
    char buf[BUFFER_SIZE];
    char *lineptr[MAX_ENTRIES];
    struct linux_dirent *d;
    int bpos;

    // For parsing out the path from the inputted command
    // Will increase this later so that options / flags can be added
    if (argc > 2) {
        printf("Too many arguments provided please supply just the path");
        return 0;
    }
    
    // This is the path that is supposed to be where the directory is supposed to be `ls`ed
    char* path = argv[1];
    
    fd = openat(AT_FDCWD, path, O_RDONLY|O_NONBLOCK|O_CLOEXEC|O_DIRECTORY);
      
    if (fd ==-1) 
    { 
        // print which type of error have in a code 
        printf("Error Number % d\n", errno); 
          
        // print program detail "Success or failure" 
        perror("Program");                 
    } 
    
    nlines = 0;
    for ( ; ; ) {
        nread = syscall(SYS_getdents64, fd, buf, BUFFER_SIZE);

        if (nread == -1) {
            // print which type of error have in a code 
            printf("Error Number % d\n", errno); 
              
            // print program detail "Success or failure" 
            perror("Program");                 
        }

        if (nread == 0) {
            break;
        }

        for (bpos = 0; bpos < nread;) {
            char *p;
            d = (struct linux_dirent *) (buf + bpos); 
            if (d -> d_name[0] == '\n') {
                p = &(d -> d_name[1]);
            }
            else {
                p = d -> d_name;
            }
            lineptr[nlines++] = p;
            bpos += d -> d_reclen;
        } 
    }
    
    close(fd);

    int i;
    for (i = 0; i < nlines; i++) {
        printf("%s  ", lineptr[i]);
    }
    printf("\n");
    
    return 0;
}












