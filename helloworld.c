
#include <pthread.h>
#include <stdio.h>

int i = 0; 

// Note the return type: void*
void* addMillion(i){
    i+=1; 
    return NULL 
}

void* subtracktMillion(i){
    i-=1; 
    return NULL 
}


int main(){
    pthread_t thread_1;
    pthread_t thread_2; 

     pthread_create(&thread_1, NULL, addMillion(), i);
    // Arguments to a thread would be passed here ---------^
     pthread_create(&thread_2, NULL, subtrackMillion(),i);
    // Arguments to a thread would be passed here ---------^
      
    
    pthread_join(thread_1, NULL);
        pthread_join(thread_2, NULL);
    printf(i);
    return 0;
    
}