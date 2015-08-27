# sgdocloader
_sgdocloader_ is a CLI tool that makes it easy to add JSON documents to a [Sync Gateway](https://github.com/couchbase/sync_gateway/) instance.

## Usage
You're required to specify the Sync Gateway URL (including port) and bucket name, as well as the file(s) and/or directories that contain the JSON documents that you'd like to load.

Example:

```shell
sgdocloader -u http://127.0.0.1:4984 -b mybucket path/to/docs
```

The `_id` field for each doc will be auto-generated if not already included within the document(s).
