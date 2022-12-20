# error introduction
With a customized error interface can make the error message more fine, which to handle error action according to different error types.  


There are several tips for use error:

**tips1: `return nil` to `error`**
```
func testError() error {
	var fErr *fileErr
	var flag bool = true
	if flag {
		return fErr
	}

	return fErr
}
```

The `testError` will always return the object non nil, cause for a interface the nil is mean the type and value is nil. For `testError` the `type=&fileErr, value=nil`, the caller will get non nil from `testError`.

More detail analysis as [Why is my nil error value not equal to nil?](https://go.dev/doc/faq#nil_error).

**tips2: handle error with customize handler**
Like `appHandler` here to handle error with a common place.
```
func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ah(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
```
