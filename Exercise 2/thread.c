// gcc -std=gnu99 -Wall -g -o thread thread.c -lpthread
#include <pthread.h>
#include <stdio.h>
#include <unistd.h>

int i = 0; 
pthread_mutex_t lock;
// Note the return type: void*
void* addMillion(){
      pthread_mutex_lock(&lock);
	for (int j = 0; j < 1000000; j++){
    		i+=1;
	} 
       pthread_mutex_unlock(&lock);
    	return NULL; 
}

void* subtracktMillion(){
       pthread_mutex_lock(&lock);
    	for (int x = 0; x < 1000000; x++){
		i-=1;
	} 
       pthread_mutex_unlock(&lock);
	return NULL; 
}


int main(){
    	pthread_t thread_1;
    	pthread_t thread_2; 
        
     	pthread_create(&thread_1, NULL, addMillion, NULL);
    	// Arguments to a thread would be passed here ---------^
     	pthread_create(&thread_2, NULL, subtracktMillion, NULL);
    	// Arguments to a thread would be passed here ---------^
    	  
    
    	pthread_join(thread_1, NULL);
        pthread_join(thread_2, NULL);
    	printf("%d\n", i);
          pthread_mutex_destroy(&lock);

    	return 0;
    
}
