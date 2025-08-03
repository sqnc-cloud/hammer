# hammer

a cli tool to backup your mongodb on aws s3

## Building

```shell
  go build -o hammer cmd/main.go
```

## Running

```shell
hammer backup --upload --connection mongodb://localhost:27017 --database hammer
```

## Example of use

```shell
 AWS_ACCESS_KEY_ID=example AWS_SECRET_ACCESS_KEY=example AWS_REGION=us-east-1 AWS_BUCKET=example hammer backup --upload --connection mongodb://localhost:27017 --database example
```
