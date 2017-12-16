# MemDB

![](https://travis-ci.org/D-L-M/mem-db.svg?branch=master)

MemDB is a simple in-memory database management system which allows storing and searching of unstructured JSON documents using inverted indices.

---

Please note that MemDB is **not** production ready — it is maintained exclusively as a personal Golang learning project. It is probably good enough to support a small web service, but I wouldn't recommend it.

If you like the sound of this project and don't have prior experience with such databases, I encourage you to explore [Elasticsearch](https://www.elastic.co/products/elasticsearch) for use in production applications.

---

## Setting Up

Start the server by running:

```bash
go run ./src/main.go --port=XXX
```

If the `port` argument is omitted, MemDB will fall back to port 9999.

## Authentication

All requests must be made with Basic authentication. The default username and password are `root` and `password`, respectively, which form the following header:

```
Authorization: Basic cm9vdDpwYXNzd29yZA==
```

## Storing Documents

To store a document, make a HTTP `PUT` request with the JSON document as the request body to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to store.

Alternatively you can omit the ID to have one randomly generated for the document.

## Retrieving Documents

To retrieve a document, make a HTTP `GET` request to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to retrieve.

## Searching

To search for documents, make a HTTP `GET` or `POST` request to `http://localhost:9999/_search` with a JSON body describing the search criteria, for example:

```javascript
{
  "and":
    [
      {"equals": {"fieldName": 123}}
    ]
}
```

The top-most node of the JSON request must always be represented by an `and` or `or` key that contains an array of criteria that must either all be satisfied (`and`) or at least one of which must be satisfied (`or`).

The top-most node of each criterion object can be one of the following: `equals`, `not_equals`, `contains`, `not_contains` — the 'contains' options allow searching of individual words within string fields.

Field names should be given in dot-notation, with numeric array indices removed. For example:

```javascript
{
  "field1": "value", // field1
  "field2":
    {
      "field3": "value" // field2.field3
    },
  "field4":
    [
      {
        "field5": "value" // field4.field5
      },
      {
        "field6": "value" // field4.field6
      }
    ]
}
```

To buid up complex criteria, a further `and` or `or` criteria set can be nested. For example:

```javascript
{
  "and":
    [
      {"equals": {"age": 30}},
      {
        "or":
          [
            {"equals": {"address.country": "Wales"}},
            {"contains": {"full_name": "Smith"}}
          ]
      }
    ]
}
```

In this example, documents would match where the field `age` was equal to 30 and either the `address.country` field was equal to 'Wales' or the `full_name` field contained the word 'Smith' (note that string searches are case-insensitive).

If required, criteria can be nested many levels deep.

## Removing Documents

To remove an individual document, make a HTTP `DELETE` request to `http://localhost:9999/{id}`, where `{id}` is the unique identifier of the document to remove.

To remove multiple documents, make a HTTP `GET` or `POST` request to `http://localhost:9999/_delete` with a JSON body describing the search criteria, as per the 'Searching' section.

To remove all documents, make a HTTP `DELETE` request to `http://localhost:9999/_all`.

## Viewing Index Statistics

To view index statistics, make a HTTP `GET` request to `http://localhost:9999/_stats`.

## Testing

To run the project's unit tests, simply run:

```bash
npm test
```