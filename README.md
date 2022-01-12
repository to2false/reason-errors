# Reason Errors

Extended from [go-kratos](https://github.com/go-kratos/kratos) errors

Difference: record the enum value as reason no
```proto
message Error {
  int32 code = 1;
  int32 reason_no = 2;
  string reason = 3;
  string message = 4;
  map<string, string> metadata = 5;
};
```
## Installation
```
go get github.com/to2false/reason-errors@v0.1.0
```

Relative Tool: https://github.com/to2false/protoc-gen-go-reason-errors
