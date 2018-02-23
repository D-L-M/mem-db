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
go run ./src/main.go --port=XXXX
```

If the `port` argument is omitted, MemDB will fall back to port 9999.

## Running Multiple Nodes

It is possible to configure MemDB to operate on multiple nodes that share the same home directory (e.g. an EFS filesystem mounted as the home directory of multiple EC2 instances).

To achieve this, simply provide the hostname of the instance and the hostnames of all other instances as flags when starting the application:

```bash
go run ./src/main.go --port=9999 \
--hostname=http://192.168.1.1:9999 \
--peers=http://192.168.1.2:9999,http://192.168.1.3:9999
```

The host and peer names need to be accessible to each other, but do not need to be accessible from the Internet; they can be provided as domain names, public IP addresses or local IP addresses.

If you omit the `hostname` flag, `http://127.0.0.1:XXXX` will be assumed, where `XXXX` is the port of the node being started (falling back to port 9999 if not provided).

It is also possible to use a custom directory for shared storage by providing the directory as a flag:

```bash
go run ./src/main.go --base-directory=/path/to/storage
```

## Authentication

All requests must be made with Basic authentication. The default username and password are `root` and `password`, respectively, which form the following header:

```
Authorization: Basic cm9vdDpwYXNzd29yZA==
```

The `root` user can also create a new user (or update an existing user's password) by making a HTTP `POST` or `PUT` request to `http://localhost:9999/_user` with a JSON body describing the credentials:

```javascript
{
  "username": "foo",
  "password": "bar",
  "action": "create" // Or 'update'
}
```

The same body can also be used to delete a user, by setting `delete` as the action and omitting the password.

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

By default, 25 records will be returned, although this can be altered by providing query string parameters such as `http://localhost:9999/_search?size=20&from=60`.

### Statistics

You can also request a list of significant terms from a field in the filtered results by appending the following query string parameters to a search URL: `http://localhost:9999/_search?&significant_terms_field=description&significant_terms_threshold=300&significant_terms_minimum=25`.

In the above example, `significant_terms_field` is the dot-notation name of the field to get significant terms from, `significant_terms_threshold` is the percentage by which the terms should be significant (in this example, 300% or 3x more common than the background data) and `significant_terms_minimum` is the minimum percentage of matching documents a term must occur in to be included.

If omitted, the threshold will default to 200% and the minimum will default to the threshold divided by 100 (so at least 2 occurrences for the default threshold). Sensible custom figures for the minimum value range between 2 and 25.

It is also possible to provide a negative threshold number to identify insignificant terms.

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