# Response time analyse tool 

## Description:

This tool is dedicated to analyse a API log file of a PlayerItem micro-service.
It read all log files stored at /var/log/httpd/*.log
the format for each line is:


    P_ADDRESS [timestamp] "HTTP_VERB URI" HTTP_ERROR_CODE RESPONSE_TIME_IN_MILLISECONDS

Example:

    10.2.3.4 [2018/13/10:14:02:39] "GET /api/playeritems?playerId=3" 200 1230

this tool calculates the 90%, 95% and 99% percentile response time for READ
API requests, based on ALL log files stored in /var/log/httpd/*.log
the output for the tool as follows:

    90% of requests return a response in X ms
    95% of requests return a response in Y ms
    99% of requests return a response in Z ms

## Solution

The purpose for this tool is to calculate the 90%, 95% and 99% percentile response, so we can just figure out all the top 10% of the most longest response time. 
for the numbers of the log files may be very big, we write code with golang to embrace a parallel file read capability easily
min-heap is used as the data structure to host the top 10% time for its $O(LogN)$ time-complexity for push and pop

## Steps

This tool works as follow:

 1. Read all the files in /var/log/httpd/*.log and figure out all the valid response time for READ API, calculated the total size for all the valid response, save all the times to a temporary file
 2. Calculated the heap size based on the total size.
 `max-heap-size = total-size * 10%`
 3. Build up the min-heap and push all the response time value which is great than the top value of the heap
 4. Pop value if heap size exceed the max-heap-size
 5. Pop the top value as the result after push all the data
 6. Shrink the heap size to the corresponding size for other 
percentile and pop the top value as other results:
`new-heap-size = total-size * Y%`

## Performance
The main performance cost for each step is:

 1. $O(N)$ for read all files N represents for the size of all the valid response time
 2. $O(1)$
 3. time complexity is $O(NlogM)$ N represents for the size of all the valid response time, M for the heap size which is 10% * N space complexity is $O(N)$
 4. $O(1)$
 5. $O(1)$
 6. $O(N)$

As a conclusion, the total complexity is $O(NlogN)$ for time complexity and O(N) for space complexity

## File Struct
files are organized follows golang convention:
 - **heap/heap.go** 
	 - file for package heap, provide heap container and related algorithm
 - **main/main.go**
	 - file for package main provide log file
   scan and the main process for the tool

## Bugs

## TODO

