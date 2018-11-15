# cloudgo

本cloudgo应用与老师博客给出的样例基本相同，仅新添加了一个http.HandlerFunc。由于应用极其简单，因此未使用任何框架。

程序功能说明：
搭建一个服务端，访问端口在启动程序时使用-p(flag)设置。之后可以使用两种类型的路径访问服务端
```
/{act}/{id}/{time}:根据输入的time(次数)重复输出 time 次 act id

/find/{id}:输出request的url和"Cannot find user " + id
```

测试结果：

curl测试`/{act}/{id}/{time}`
```
C:\Users\liyike622\Downloads\Apache24\bin>curl -v http://localhost:9090/beat/liyike/3
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 9090 (#0)
> GET /beat/liyike/3 HTTP/1.1
> Host: localhost:9090
> User-Agent: curl/7.55.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=UTF-8
< Date: Thu, 15 Nov 2018 21:16:57 GMT
< Content-Length: 105
<
{
  "Repeate": "3"
}
{
  "Test": "beat liyike"
}
{
  "Test": "beat liyike"
}
{
  "Test": "beat liyike"
}
* Connection #0 to host localhost left intact
```

curl测试`/find/{id}`
```
C:\Users\liyike622\Downloads\Apache24\bin>curl -v http://localhost:9090/find/liyike
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 9090 (#0)
> GET /find/liyike HTTP/1.1
> Host: localhost:9090
> User-Agent: curl/7.55.1
> Accept: */*
>
< HTTP/1.1 200 OK
< Date: Thu, 15 Nov 2018 21:18:41 GMT
< Content-Length: 75
< Content-Type: text/plain; charset=utf-8
<
Find request to localhost:9090/find/liyike
Result: Cannot find user liyike
* Connection #0 to host localhost left intact
```

curl测试`不可处理的路由`
```
C:\Users\liyike622\Downloads\Apache24\bin>curl -v http://localhost:9090/liyike
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 9090 (#0)
> GET /liyike HTTP/1.1
> Host: localhost:9090
> User-Agent: curl/7.55.1
> Accept: */*
>
< HTTP/1.1 404 Not Found
< Content-Type: text/plain; charset=utf-8
< X-Content-Type-Options: nosniff
< Date: Thu, 15 Nov 2018 21:19:23 GMT
< Content-Length: 19
<
404 page not found
* Connection #0 to host localhost left intact
```

ab测试
`(参考链接:https://www.cnblogs.com/gumuzi/p/5617232.html)`
```
    测试使用参数
    -n requests     Number of requests to perform
    -c concurrency  Number of multiple requests to make at a time
```
```
测试结果及结果分析

C:\Users\liyike622\Downloads\Apache24\bin>ab -n 1000 -c 100 http://localhost:9090/hello/liyike/3
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/
//以上为apache的版本信息

Benchmarking localhost (be patient)
Completed 100 requests
Completed 200 requests
Completed 300 requests
Completed 400 requests
Completed 500 requests
Completed 600 requests
Completed 700 requests
Completed 800 requests
Completed 900 requests
Completed 1000 requests
Finished 1000 requests
//以上内容显示测试完成度，在请求数量很多时会分行显示当前完成数量。


Server Software:
//被测试的服务器所用的软件信息，此处无
Server Hostname:        localhost
//被测主机名
Server Port:            9090
//被测主机的服务端口号

Document Path:          /hello/liyike/3
//请求的具体文件
Document Length:        108 bytes
//请求的文件大小

Concurrency Level:      100
//并发级别，也就是并发数，请求中－c参数指定的数量
Time taken for tests:   1.294 seconds
//本次测试总共花费的时间
Complete requests:      1000
//本次测试总共发起的请求数量，请求中－n参数指定的数量
Failed requests:        0
//失败的请求数量。因网络原因或服务器性能原因，发起的请求并不一定全部成功。
//通过该数值和Complete requests相除可以计算请求的失败率，作为测试结果的重要参考。
Total transferred:      232000 bytes
//总共传输的数据量，指的是ab从被测服务器接收到的总数据量，包括输出到客户端的文本内容和请求头信息。
HTML transferred:       108000 bytes
//从服务器接收到的index.html文件的总大小，等于Document Length＊Complete requests 
//即108 bytes ＊ 100 ＝ 108000 bytes
Requests per second:    772.56 [#/sec] (mean)
//平均(mean)每秒完成的请求数
Time per request:       129.440 [ms] (mean)
//从用户角度看，完成一个请求所需要的时间（因用户数量不止一个，服务器完成100个请求，平均每个用户才接收到一个完整的返回，所以该值是下一项数值的100倍。）
Time per request:       1.294 [ms] (mean, across all concurrent requests)
//服务器完成一个请求的时间
Transfer rate:          175.03 [Kbytes/sec] received
//网络传输速度。对于大文件的请求测试，这个值很容易成为系统瓶颈所在
//要确定该值是不是瓶颈，需要了解客户端和被测服务器之间的网络情况
//包括网络带宽和网卡速度等信息

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.4      0       1
Processing:    18  127  27.1    129     172
Waiting:        1  127  27.1    128     172
Total:         18  127  27.1    129     172
//这几行组成的表格主要是针对响应时间也就是第一个Time per request进行细分和统计。
//一个请求的响应时间可以分成网络链接（Connect），系统处理（Processing）和等待（Waiting）三个部分。
//表中min表示最小值； mean表示平均值；median表示中位数； max表示最大值了。
//[+/-sd]表示标准差（Standard Deviation） ，也称均方差（mean square error），表示数据的离散程度
//数值越大表示数据越分散，系统响应时间越不稳定。 

//需要注意的是表中的Total并不等于前三行数据相加，因为前三行的数据并不是在同一个请求中采集到的
//可能某个请求的网络延迟最短，但是系统处理时间又是最长。所以Total是从整个请求所需要的时间的角度来统计的。

Percentage of the requests served within a certain time (ms)
  50%    129
  66%    139
  75%    150
  80%    154
  90%    161
  95%    165
  98%    167
  99%    169
 100%    172 (longest request)
//这个表第一行表示有50%的请求都是在129ms内完成的，可以看到这个值是比较接近平均系统响应时间
//（第一个Time per request:       129.440 [ms] (mean) ）

//以此类推，所有请求（100%）的时间都是小于等于172ms的，也就是表中最后一行的数据是时间最长的那个请求（longest request）。
```