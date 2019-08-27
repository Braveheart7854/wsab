# wsab
websocket 压测工具

# 参数
- -n   : 总请求数
- -c   : 并发请求数
- -url : ws或者wss地址

# 使用样例
- [root@xxxx ~]# wsab -n 20 -c 20 -url ws://127.0.0.1:7777/ws

# 返回结果
- total request :  20               （总请求数量）
- spend time :  0.005833 s          （总花费时间）
- qps : 3428.77 [#/sec]             （qps）
- per request min time :  2.303 ms   (单个请求所用最小时间)
- per request max time :  3.602 ms   (单个请求所用最大时间)
- connect success:  20              （websocket连接成功数）
- connect fail:  0                  （websocket连接失败数）
- send message success:  20         （发送数据包成功数）
- send message fail:  0             （发送数据包失败数）
- accept message success:  20       （接收数据包成功数）
- accept message fail:  0           （接收数据包失败数）