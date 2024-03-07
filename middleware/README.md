# Middleware

## Retry
Retries a request on `http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout` or configurable http error code
The http request object is wrapped to buffer and ensure the request data stream is still available on the retry. Error responses a matching
error code are discarded when there is a retry option left

## Peek model
Optimistic approach to read the model name from the request payload when it is the first attribute of a json object. The request data is buffered so
that no data is lost for any following step. The model is stored in the `X-Model` so that upstream processes can read it from there.
When the model attribute is not the first or not read within the buffered data, the result is empty. 
Example request payload:
```json
{
  "model": "my-model",
  "messages": []
}
```

```
goos: darwin
goarch: arm64
pkg: github.com/alpe/netutils/common
BenchmarkDecoding/medium/PeekModel-12         	 1624220	       726.0 ns/op
BenchmarkDecoding/medium/JsonDecoder-12       	    8212	    144236 ns/op
BenchmarkDecoding/big/PeekModel-12            	 1661097	       727.7 ns/op
BenchmarkDecoding/big/JsonDecoder-12          	     433	   2735448 ns/op
BenchmarkDecoding/minimal/PeekModel-12        	 1779588	       678.8 ns/op
BenchmarkDecoding/minimal/JsonDecoder-12      	 2108025	       565.9 ns/op
BenchmarkDecoding/small/PeekModel-12          	 1651448	       722.1 ns/op
BenchmarkDecoding/small/JsonDecoder-12        	 1299910	       921.3 ns/op
```