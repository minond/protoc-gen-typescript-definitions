`protoc-gen-typescript-definitions` is a Protocol Buffers compiler plugin that
generates TypeScript type definitions for messages.

Given the following protobuf message you will get a file with the TypeScript
definition below it:

**user.proto**

```proto
syntax = "proto3";

message User {
  string guid             = 1;
  string name             = 2;
  int age                 = 3;
  int64 dob               = 4;
  map<string, int> scores = 5;
}
```

**user.d.ts**

```typescript
type User = {
  guid: string
  name: string
  age: number
  dob: number
  scores: Dictionary<string, number>
}
```
