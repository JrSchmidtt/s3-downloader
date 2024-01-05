# S3 Bucket Downloader

This script is designed to list objects in a specified AWS S3 bucket and download each object into a folder with the same name as the bucket. It utilizes the AWS SDK for Go (`github.com/aws/aws-sdk-go`) to interact with the AWS S3 service.

## Prerequisites

- Go installed on your machine.
- AWS credentials (Access Key ID and Secret Access Key) with permission to access the specified S3 bucket.
- Set the following environment variables:

  ```bash
  export AWS_ACCESS_KEY_ID=your-access-key-id
  export AWS_SECRET_ACCESS_KEY=your-secret-access-key
  export AWS_REGION=your-region

## installation

```bash
go mod download
```

## Usage

```bash
go run main.go -b your-bucket-name
```

## Options

| Option | Description | Required |
| ------ | ----------- | -------- |
| -b     | The name of the S3 bucket to download objects from. | Yes |

## Contributing

1. [Fork the repository](https://github.com/JrSchmidtt/s3-downloader/fork)!
2. Clone your fork.
3. Create your feature branch: `git checkout -b my-new-feature`
4. Commit your changes: `git commit -am 'Add some feature'`
5. Push to the branch: `git push origin my-new-feature`
6. Submit a pull request :D