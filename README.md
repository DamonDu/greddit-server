This project show you how to host an API server on AWS Lambda Function.
**DO NOT USE IT ON PRODUCTION!!!**
### Tech Stack
* Golang / Gofiber for backend
* SQLite for storage
* AWS / Serverless Framework for deployment
* Zig for cgo cross-compile
### Local Setup
1. Create `.env.local` file like `.env.example`
2. To run this project locally, run this command in root directory:
```shell
go run ./cmd/server/main.go
```
### Deploy
1. [Install and setup Serverless Framework and AWS.](https://www.serverless.com/framework/docs/providers/aws/guide/installation)
2. [Install Zig.](https://ziglang.org/learn/getting-started/)
3. Run this command to deploy lambda function:
```shell
make deploy
```