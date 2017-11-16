# MemDB

This is an attempt to create a very simple in-memory database using inverted indices.

It is primarily an experiment to learn how to code in Go, so should neither be relied upon in production nor treated too harshly when looking at the code!

## Setting Up

Start the server by running:

```bash
go run main.go
```

## Storing Documents

To store a document, make a HTTP `PUT` request with the JSON document as the request body to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to store.

## Retrieving Documents

To retrieve a document, make a HTTP `GET` request to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to retrieve.

## Viewing Index Statistics

To view index statistics, make a HTTP `GET` request to `http://localhost:9999/_stats`.