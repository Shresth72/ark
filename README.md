# Ark
Distributed Key Value Store

There are multiple libraries

## Jepsen
Testing framework for distributed nodes

## Types of Message Events

### 1. Init
```json
{"src":"n1","dest":"n2","body":{"type":"init","node_id":"n1","node_ids":["n1","n2"]}}
```

### 2. Echo
```json
{"src":"n1","dest":"n2","body":{"type":"echo","msg_id":1,"echo":"hello"}}
```

### 3. UniqueIds
```json
{"src":"n1","dest":"n2","body":{"type":"generate"}}
```

### 4. Grow Counter

**Add**
```json
{"src":"n1","dest":"n2","body":{"type":"add","delta":5}}
```

**Read**
```json
{"src":"n1","dest":"n2","body":{"type":"read"}}
```

### 5. Transaction

**Txn Example 1**
```json
{"src":"n3","dest":"n4","body":{"type":"txn","txn":[["w",2,5],["w",4,19],["r",2,null]]}}
```

**Txn Example 2**
```json
{"src":"n3","dest":"n4","body":{"type":"txn","txn":[["w",1,10],["w",2",5],["r",1,null],["w",2,7],["r",2,null],["w",10,999],["r",10,null],["r",3,null],["w",3,42],["r",3,null],["w",4,19],["w",4,20],["r",4,null],["w",100,500],["r",100,null],["r",999,null]]}}
```

### 6. Consistent Log

**Send Log 1**
```json
{"src":"c1","dest":"n1","body":{"type":"send","key":"log1","msg":47}}
```

```json
{"src":"c1","dest":"n1","body":{"type":"send","key":"log2","msg":99}}
```

```json
{"src":"c1","dest":"n1","body":{"type":"send","key":"log1","msg":69}}
```

**Poll Logs**
```json
{"src":"n1","dest":"c1","body":{"type":"poll","offsets":{"log1":1000, "log2": 2000, "log1": 1001}}}
```

**Commit Offsets**
```json
{"src":"n1","dest":"c1","body":{"type":"commit_offsets","offsets":{"log1":1001,"log2":2000}}}
```

**List Committed Offsets**
```json
{"src":"c1","dest":"n1","body":{"type":"list_committed_offsets","keys":["log1","log2"]}}
```
