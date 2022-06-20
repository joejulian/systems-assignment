# Keyvalue API

## Server

### Building the Server

If you have go installed, you can build the server with:

```bash
make server
```

If you don't have go installed, but you do have docker installed, you can build the server with:

```bash
make server.docker
```

### Server Usage

The keyvalue api server is a simple http server that can be used to store and retrieve key-value pairs where the key is a uuid.

```text
Usage:
  server [flags]

Flags:
      --config string   config file (default is $HOME/.server.yaml)
      --data string     data file (default "example.data")
  -h, --help            help for server
      --host string     host to listen on (default "localhost")
      --port int        port to listen on (default 8080)
```

Example:

```bash
bin/server --host 0.0.0.0 --port 80 --data example.data
```

### Server Configuration

You may save the config file as `.server.yaml` in your home directory.

Example `.server.yaml`:

```yaml
data: /mnt/shared/example.data
host: 0.0.0.0
port: 8080
```

## Client

### Building the Client

If you have go installed, you can build the server with:

```bash
make client
```

If you don't have go installed, but you do have docker installed, you can build the server with:

```bash
make client.docker
```

### Client Usage

The keyvalue api client is a simple http client that can be used to  retrieve values from the api server.

```text
Usage:
  client get [flags]

Flags:
  -h, --help              help for get
      --host string       The host of the api server (default "localhost")
      --json              Return the value as json
  -k, --key stringArray   The key(uuid) to get
      --port int          The port of the api server (default 8080)

Global Flags:
      --config string   config file (default is $HOME/.client.yaml)
```

Examples:

```bash
bin/client get --host localhost --port 8080 \
    --key 72cf48f7-d604-4423-b012-4e2329f5117b
```

```bash
bin/client get --host localhost --port 8080 \
    --key 72cf48f7-d604-4423-b012-4e2329f5117b --json
```

```bash
bin/client get --host localhost --port 8080 \
    --key 72cf48f7-d604-4423-b012-4e2329f5117b \
    --key f8cc65e1-2474-44c7-8502-4fd1f9ed7ea3
```

### Client Configuration

You may save the config file as `.client.yaml` in your home directory.

Example `.client.yaml`:

```yaml
get:
    host: 0.0.0.0
    port: 8080
```

This would allow you to use the client without specifying the host and port.

```bash
bin/client get --key 72cf48f7-d604-4423-b012-4e2329f5117b
```
