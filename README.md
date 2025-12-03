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

5. Transaction
```bash
{
  "src": "n3",
  "dest": "n4",
  "body": {
    "type": "txn",
    "txn": [
      ["w", 2, 5],
      ["w", 4, 19],
      ["r", 2, null]
    ]
  }
}
```

```bash
{
  "src": "n3",
  "dest": "n4",
  "body": {
    "type": "txn",
    "txn": [
      ["w", 1, 10],
      ["w", 2, 5],
      ["r", 1, null],
      ["w", 2, 7],
      ["r", 2, null],
      ["w", 10, 999],
      ["r", 10, null],
      ["r", 3, null],
      ["w", 3, 42],
      ["r", 3, null],
      ["w", 4, 19],
      ["w", 4, 20],
      ["r", 4, null],
      ["w", 100, 500],
      ["r", 100, null],
      ["r", 999, null]
    ]
  }
}
```
