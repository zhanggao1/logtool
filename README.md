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

The purpose for this tool is to calculate the 90%, 95% and 99% percentile response, so we can just figure out all the top 10% of the most longest response time,
For the numbers of the log files may be very big, we write code with golang to embrace a parallel file read capability easily,
Min-heap is used as the data structure to host the top 10% response time.

## Steps

This tool works as follow:

 1. Read all the files in /var/log/httpd/*.log and figure out all the valid response time for READ API, calculated the total size for all the valid response, save all the times to a temporary file
 2. Calculated the heap size based on the total size.
 `max-heap-size = total-size * 10%`
 3. Build up the min-heap with max size of max-heap-size, and push all the response time value to the heap, pop node when the heap size is great than max-heap-size
 4. Pop the top value as the result after push all the data
 5. Shrink the heap size to the corresponding max size for other percentile and pop the top value as result for 
`new-heap-size = total-size * Y%`

## Performance
The main performance cost for each step is:

 1. O(N) for read all files N represents for the size of all the valid response time
 2. O(1)
 3. time complexity is O(NlogM) N represents for the size of all the valid response time, M for the heap size which is 10% * N space complexity is O(N)
 4. O(1)
 5. O(1)
 6. O(N)

As a conclusion, the total complexity is O(NlogN) for time complexity and O(N) for space complexity

## Install
 1. Clone this repo to path $GOPATH/src/github.com
 2. Run make

## Bugs

## TODO

