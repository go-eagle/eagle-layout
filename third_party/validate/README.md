# protoc-gen-validate (PGV)

## Function 

add constraint rules for proto and gen `*.pb.validate.go`

## Examples

```protobuf
syntax = "proto3";

package examplepb;

import "validate/validate.proto";

message Person {
  uint64 id = 1 [(validate.rules).uint64.gt = 999];

  string email = 2 [(validate.rules).string.email = true];

  string name = 3 [(validate.rules).string = {
    pattern:   "^[^[0-9]A-Za-z]+( [^[0-9]A-Za-z]+)*$",
    max_bytes: 256,
  }];

  Location home = 4 [(validate.rules).message.required = true];

  message Location {
    double lat = 1 [(validate.rules).double = {gte: -90,  lte: 90}];
    double lng = 2 [(validate.rules).double = {gte: -180, lte: 180}];
  }
}
```

## Reference

* Usage docs: https://github.com/bufbuild/protoc-gen-validate?tab=readme-ov-file#usage
* Official docs:https://github.com/envoyproxy/protoc-gen-validate