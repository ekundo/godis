# What is Godis?

Godis is simple implementation of Redis-like in-memory cache written in Go.

## Features


* A key-value storage that supports strings, lists and hashes as stored values.
* Automatic expiration of a key. TTL can be optionally assigned to a key.
* The following commands are supported:

  |common|string|list|hash| 
  |:---|:---|:---|:---|
  |[command](https://redis.io/commands/command)|[get](https://redis.io/commands/get)|[llen](https://redis.io/commands/llen)|[hkeys](https://redis.io/commands/hkeys)|
  |[type](https://redis.io/commands/type)|[set](https://redis.io/commands/set)|[lindex](https://redis.io/commands/lindex)|[hexists](https://redis.io/commands/hexists)|
  |[keys](https://redis.io/commands/keys)| |[lpop](https://redis.io/commands/lpop)|[hget](https://redis.io/commands/hget)|
  |[exists](https://redis.io/commands/exists)| |[lpush](https://redis.io/commands/lpush)|[hgetall](https://redis.io/commands/hgetall)
  |[del](https://redis.io/commands/del)| |[rpop](https://redis.io/commands/rpop)|[hset](https://redis.io/commands/hset)|
  |[expire](https://redis.io/commands/expire)| |[rpush](https://redis.io/commands/rpush)|[hdel](https://redis.io/commands/hdel)|
  |[expireat](https://redis.io/commands/expireat)| |[lset](https://redis.io/commands/lset)| |
  |[pexpire](https://redis.io/commands/pexpire)| | |
  |[pexpireat](https://redis.io/commands/pexpireat)| | |
  |[ttl](https://redis.io/commands/ttl)| | |
  |[pttl](https://redis.io/commands/pttl)| | |
  |[persist](https://redis.io/commands/persist)| | |

  The commands have the same signature and work exactly the same as corresponding commands in Redis bellow version 2.4 
  (Godis does not support multikey and multivalue commands).
* Both [REdis Serialization Protocol (RESP)](https://redis.io/topics/protocol) and [Inline Commands Protocol](https://redis.io/topics/protocol#inline-commands) are supported. 
Thus simple `telnet` and `redis-cli` clients can be used to play with Godis. 
Moreover, performance can be tested with `redis-benchmark`. 
* Golang API client
* Write-ahead logging based persistence. 
Log compaction is not yet implemented and format of the log is not optimized due to reasons of simplicity.       
 

## Installation

To build Godis as a Docker image install and run [Docker](https://docs.docker.com/engine/installation/) and use the following command:  
```docker build github.com/ekundo/godis```

## Playing with Godis

As a first step start Godis server with:  
```docker run -v $(pwd):/work -p 2121:2121 -it --rm godis```

Then send commands to server using any tcp client like `telnet` 
```
$ telnet 127.0.0.1 2121
set foo bar
+OK
get foo
$3
bar
```
Although `redis-cli` is the best choice:
```
$ redis-cli -p 2121 -h 127.0.0.1
redis> set foo bar
OK
redis> get foo
"bar"
```
Use `command` to get the list of available commands. 
Details of each command from the list can be found at [redis command guide](http://redis.io/commands).  