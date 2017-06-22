# gossip

`gossip` watches folders for new files and moves / copies them to an AWS S3 bucket. It is meant to be used as a service
/ daemon.

## Configuration

Most of the configuration is done through a JSON file. A sample is provided below:

```json
{
  "paths": [
    {
      "root": "/home/alex/photos",
      "delete": true,
      "delay": 2
    },
    {
      "root": "/home/alex/music",
      "recursive": true,
      "bucket": "my-music-backup"
    }
  ],
  "global": {
    "bucket": "default-bucket",
    "delay": 5
  },
  "credentials": {
    "access": "AWS_ACCESS_KEY",
    "secret": "AWS_SECRET_KEY"
  }
}
```

`paths` is an array of paths to watch. Each path has the following options:
* `root` (required): The path on the file system to watch
* `delete` (optional): Whether or not to delete the file after uploading it to S3 (default: false)
* `delay` (optional): How long to wait after a write event before initializing the upload (see `global`)
* `recursive` (optional): Whether or not the path should be watched recursively. (default: false)
* `bucket` (optional): The name of the AWS S3 bucket to upload the files to (see `global`)

`global` defines default values for both `bucket` and `delay` which will be used for each path that does not explicitly
define a value for these settings. `delay` must be at least 1 (second) in order to avoid the upload starting while the
file is still being written / saved on the local file system. `bucket` is required if one or more paths do not have an
explicit value. `delay` will default to 1 if missing.

`credentials` is the AWS `access` and `secret` key used. Alternatively, if you would prefer not to write those values
to a file, you can set them as environment variables: `AWS_ACCESS` and `AWS_SECRET`.

## Usage

`gossip [-config file]`

If `-config` is missing or `file` does not exist, `gossip` will look at the `GOSSIP_CONFIG` environment variables for
a file path and use it if it exists.