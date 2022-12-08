# NetCat

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.



## Usage/Examples

Clone the repository and start the server

```bash
  go run chatserver.go
```

then open multyple new terminals and to connect the server write

```bash
  nc localhost 9090
```

## Author

- [@aomashev](https://www.github.com)