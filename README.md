# IP Logger

A command line tool that lets you create redirection shortening links that record IP address and User agent of whoever visits them. 
This uses unix sockets to communicate with a small golang web server.

## Build Instructions

This should work on all platforms.

```
go build
```

## Examples

Start the web server, I put mine behind an nginx reverse proxy (although you will need to do some fun stuff with X-forwarded headers to get IP then).  
```
./ip-logger -server 0.0.0.0:8080
```

In another terminal, you can now register new links!

```
$ ./ip-logger https://youtube.com

http://127.0.0.1:8080/a/09a957b4a8
```


Using `ls` alone will list all URLS.  

```
./ip-logger ls

                                                 URLs
-------------------------------------------------------------------------------------------------------------------
| Label | Short URL                          | Destination                                        | Number Visits |
-------------------------------------------------------------------------------------------------------------------
|       | http://127.0.0.1:8080/a/2ee8832346 | https://google.co.nz                               | 4             |
-------------------------------------------------------------------------------------------------------------------
|       | http://127.0.0.1:8080/a/6bb5a80fc6 | https://google.co.nz/thisisalongboi/test/test/test | 0             |
-------------------------------------------------------------------------------------------------------------------
|       | http://127.0.0.1:8080/a/2b79cfbe4c | https://google.co.nz/test                          | 4             |
-------------------------------------------------------------------------------------------------------------------
|       | http://127.0.0.1:8080/a/471e26039d | https://google.co.nz/test                          | 1             |
-------------------------------------------------------------------------------------------------------------------
|       | http://127.0.0.1:8080/a/09a957b4a8 | https://youtube.com                                | 0             |
-------------------------------------------------------------------------------------------------------------------

```

However if you use ls on a single identity (which is the bytes at the end of `/a/`), then it'll print the visits of that url.

```
$ ./ip-logger ls 2b79cfbe4c
                                     URLs
------------------------------------------------------------------------------------------
| Label | Short URL                          | Destination               | Number Visits |
------------------------------------------------------------------------------------------
|       | http://127.0.0.1:8080/a/2b79cfbe4c | https://google.co.nz/test | 4             |
------------------------------------------------------------------------------------------
                                          Visits
------------------------------------------------------------------------------------------------------------
| Time            | IP              | UA                                                                   |
------------------------------------------------------------------------------------------------------------
| 30 Apr 21 23:37 | 127.0.0.1:45636 | Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0 |
------------------------------------------------------------------------------------------------------------
| 30 Apr 21 23:37 | 127.0.0.1:45636 | Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0 |
------------------------------------------------------------------------------------------------------------
| 30 Apr 21 23:37 | 127.0.0.1:45636 | Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0 |
------------------------------------------------------------------------------------------------------------
| 30 Apr 21 23:39 | 127.0.0.1:45668 | Mozilla/5.0 (X11; Linux x86_64; rv:88.0) Gecko/20100101 Firefox/88.0 |
------------------------------------------------------------------------------------------------------------

```

You can also `rm` values. This will delete the assocaited visits as well.

```
$ ./ip-logger rm 2b79cfbe4c
Deleted 2b79cfbe4c
```

Well thats all the features! Feel free to bug me for more indepth stuff
