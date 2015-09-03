# sgdocloader
_sgdocloader_ is a CLI tool that makes it easy to add JSON documents to a [Sync Gateway](https://github.com/couchbase/sync_gateway/) instance.

## Installation
You can grab [pre-built binaries for Windows, OS X, and Linux](https://github.com/PaulCapestany/sgdocloader/releases/tag/v0.0.1) and drop them anywhere in your _$PATH_ and you should be set.

If you prefer to build from source (and already have Go installed):

`go get github.com/paulcapestany/sgdocloader`

## Usage
You're required to specify the Sync Gateway URL (including port) and bucket name, as well as the file(s) and/or directories that contain the JSON documents that you'd like to load.

Example command:

```shell
sgdocloader -u http://127.0.0.1:4984 -b mybucket path/to/docs
```

The names of the files/directories don't matter. The `_id` field for each doc will be auto-generated if not already included within the document(s). The heavy lifting is all done by the excellent [go-couch](https://github.com/tleyden/go-couch) package.
