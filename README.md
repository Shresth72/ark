# Ark
Distributed Key Value Store

There are multiple libraries

## Jepsen
Testing framework for distributed nodes

Types of Message Events

1. Init
```bash
{"src":"n1","dest":"n2","body":{"type":"init","node_id":"n1","node_ids":["n1","n2"]}}
```

2. Echo
```bash
{
  "src": "n1",
  "dest": "n2",
  "body": {
    "type": "echo",
    "msg_id": 1,
    "echo": "hello"
  }
}
```

3. UniqueIds
```bash
{
  "src": "n1",
  "dest": "n2",
  "body": {
    "type": "generate"
  }
}
```

4. Grow Counter
- Add Counter
```bash
{
  "src": "n1",
  "dest": "n2",
  "body": {
    "type": "add",
    "delta": 5
  }
}
```

- Read Counter
```bash
{
  "src": "n1",
  "dest": "n2",
  "body": {
    "type": "read"
  }
}
```
