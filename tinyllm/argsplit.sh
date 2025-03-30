#!/run/current-system/sw/bin/bash

curl -X POST http://localhost:11435/api/chat -H "Content-Type: application/json" -d '{
  "model": "llama3.2:3b",
  "messages": [
    {"role": "system", "content": "you are a conversion tool that simply takes an standard linux command as input and properly splits the arguments, respecting flags, quoted strings, positional arguments, and everything else. The outputs must exactly match the inputs, any alterations to the given commands is detrimental to the system. The input is simply a string that contains only the command needed and nothing else, the output is only the split command and nothing else."},
    {"role": "user", "content": "iptables -A POSTROUTING -t nat -p tcp -d 192.168.1.200 --dport 8080 -j MASQUERADE"}
  ],
  "options": { "num_thread": 24 },
  "stream": false,
  "format": {
    "type": "object",
    "properties": {
      "base": {
        "type": "string"
      },
      "args": {
	"type": "array",
	"items": {
  	  "type": "string"
	}
      }
    },
    "required": [
      "base",
      "args"
    ]
  }
}'
