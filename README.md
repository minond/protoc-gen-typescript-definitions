`protoc-gen-typescript-definitions` is a Protocol Buffers compiler plugin that
generates TypeScript type definitions for messages. No encoding/decoding code
is generated, this will only generate the type declarations for messages. For
example, given the protobuf file on the left you will get the type definitions
on the right:

<table>
<tr>
<td><b>definitions/log.proto</b></td>
<td><b>typings/log.d.ts</b></td>
</tr>
<tr>
<td valign="top">

```proto
syntax = "proto3";

package log;

message Log {
  string guid              = 1;
  string text              = 2;
  map<string, string> data = 3;
  int64 createdOn          = 4;
  string createdBy         = 5;
  int64 updatedOn          = 6;
  string updatedBy         = 7;
  int64 deletedOn          = 8;
  string deletedBy         = 9;
}

message LogCreateRequest {
  string id       = 1;
  string text     = 2;
  int64 createdOn = 3;
  int64 updatedOn = 4;
}

message LogCreateResponse {
  string id = 1;
  Log log   = 2;
}
```

</td>
<td valign="top">

```typescript
export type Log = {
  guid: string
  text: string
  data: Map<string, string>
  createdOn: number
  createdBy: string
  updatedOn: number
  updatedBy: string
  deletedOn: number
  deletedBy: string
}

export type LogCreateRequest = {
  id: string
  text: string
  createdOn: number
  updatedOn: number
}

export type LogCreateResponse = {
  id: string
  log: {
    guid: string
    text: string
    data: Map<string, string>
    createdOn: number
    createdBy: string
    updatedOn: number
    updatedBy: string
    deletedOn: number
    deletedBy: string
  }
}
```

</td>
</tr>
</table>
