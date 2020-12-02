# MC RCON 
Golang RCON implementation for remote execution of MC commands. [More info](https://wiki.vg/RCON)

It's based off of [Source RCON protocol implementation](https://developer.valvesoftware.com/wiki/Source_RCON_Protocol).

```
+--------------+--------+--------+
|    Field     |  Type  | Value  |
+--------------+--------+--------+
| Size         | int32  | varies | <- automatically calculated
| ID           | int32  | varies | <- default 0
| Type         | int32  | varies | <- automatically picked
| Body         | string | varies | <- the command being run
| Empty String | string | 0x00   | <- packet padding
+--------------+--------+--------+
```

### Example
The authentication has been separated from the client creation to allow for more flexibility. That being said, the client has to be authenticated before trying to send a command otherwise it will error out.

```go
  client, err := rcon.CreateClient("ip_address", "port")
  // handle err
  
  err = client.Authenticate("password")
  // handle err
  
  response, err := client.SendCommand("help")
  // handle err
  // more work with the response
```
Full example at [example.go](https://github.com/hum/mc-rcon/blob/main/example.go)

### TODO
- [ ] Properly parse the response to be readable
- [ ] Write tests
