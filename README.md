# MemDB

This is an attempt to create a very simple in-memory database using inverted indices.

It is primarily an experiment to learn how to code in Go, so should neither be relied upon in production nor treated too harshly when looking at the code!

## Usage

* Start the server by running `go run main.go`;
* `PUT` a JSON document to `http://localhost:9999/ID`, where `ID` is the unique identifier of the document.