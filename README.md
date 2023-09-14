## Introduction

`gnoigo` is a go gNOI client library that provides convenience functions for
accessing gNOI operations.

To build and execute the `gnoigo` unit tests, run the following:

```
go build ./...
go test ./...
```

## Usage

Following is an example of how to use `gnoigo` API for performing a Ping
Operation which is present in the `System` module.

*   Create the `gnoigo` clients object.

    ```
    conn, err := grpc.DialContext(ctx, "host")
    if err != nil {
        return err
    }
    clients := gnoigo.NewClients(conn)
    ```

*   Create the PingOperation object with inputs like source and destination.

    ```
    pingOp := system.NewPingOperation().Destination("1.2.3.4").Source("5.6.7.8")
    ```

*   Call the Execute operation to perform the Ping operation.

    ```
    response, err := gnoigo.Execute(ctx, clients, pingOp)
    ```

    In this example `response` will be of type `PingResponse`.
