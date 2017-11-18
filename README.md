# MemDB

![](https://travis-ci.org/D-L-M/mem-db.svg?branch=master)

MemDB is an attempt to create a very simple in-memory database using inverted indices.

It is primarily an experiment to learn how to code in Go, so should neither be relied upon in production nor treated too harshly when looking at the code!

## Setting Up

Start the server by running:

```bash
go run ./src/main.go
```

MemDB will listen for TCP connections on port 9999.

## Storing Documents

To store a document, make a HTTP `PUT` request with the JSON document as the request body to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to store.

Alternatively you can omit the ID to have one randomly generated for the document.

## Retrieving Documents

To retrieve a document, make a HTTP `GET` request to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to retrieve.

## Viewing Index Statistics

To view index statistics, make a HTTP `GET` request to `http://localhost:9999/_stats`.

## Testing

To run the project's unit tests, simply run:

```bash
npm test
```