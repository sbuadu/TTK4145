
from threading import Thread
import Queue



def someThreadFunction1(q):
	for x in range(1000000):
		while  not q.empty :
		 	i = q.get() 
		 	q.put(i+1)

def someThreadFunction2(q):
	for x in range(1000000):
		while  not q.empty :
		 	i = q.get() 
		 	q.put(i-1)

# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")



def main():	

	q = Queue.LifoQueue(1)
	q.put(0); 

	someThread1 = Thread(target = someThreadFunction1, args = (q,))
	someThread2 = Thread(target = someThreadFunction2, args = (q,))
	someThread1.start()
	someThread2.start()
    
	someThread1.join()
	someThread2.join()
	print q.get()


main()
