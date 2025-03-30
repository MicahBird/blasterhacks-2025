#!/run/current-system/sw/bin/bash

curl -X POST http://localhost:11434/api/chat -H "Content-Type: application/json" -d '{
  "model": "llama3.2:3b",
  "messages": [
    {"role": "system", "content": "You are a machine that will translate a linux command given by a user to catagory in a JSON output. The only possible catagories are: pythondev, javadev, jsdev, cppdev, rustdev, csharpdev, phpdev, golangdev, swiftdev, kotlindev, sysdev, webdev, and securitydev. You must only provide one of the provided categories and nothing else. Following this format is critical"},
    {"role": "user", "content": "cargo build"}
  ],
  "stream": false,
  "format": {
    "type": "object",
    "properties": {
      "catagory": {
        "type": "string"
      }
    },
    "required": [
      "catagory"
    ]
  }
}'
