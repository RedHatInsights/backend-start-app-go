# Testing

_Here we cover how to test our Go code, we are only covering unit testing of the code._

_Integration and system testing of ConsoleDot application should be done by a [IQE](https://insights-qe.pages.redhat.com/iqe-core-docs/tutorial/index.html) plugin_

We have an application that has documented API and on request writes data to a database.

## Testing strategy

Testing strategy chosen for this template is:

* Tests are as isolated as possible
  * isolate between the tests
  * isolate the testing layer from other application layers
* Test are testing the public API of a package
  * avoid using internal functions
  * public interface means interface of given layer, not a user interface (API)

We have two suites of tests.
One is classical Go unit tests, second is for database tests.

## Stubbing

Now we can dive deeper on the idea of multiple implementation for a code layer.
We have an interface for our DAO, which gives us an option to add another implementation.
This other implementation will be storing the data in memory.
Stubbed layer makes tests significantly faster and simplifies the DAO code to a bare minimum.

We also keep the stub in a context.
Thanks to context isolation, the tests are runnable in parallel without leaking data.

The minimal stub reservation looks like this.
We are skipping implementation of the context setter and getter here.

```go
func init() {
    dao.GetHelloDao = getHelloDao
}

type helloDaoStub struct {
    store []*models.Hello
}

func getHelloDao(ctx context.Context) dao.HelloDao {
	return getHelloDaoStub(ctx)
}

func (x *helloDaoStub) List(ctx context.Context, limit, offset int64) ([]*models.Hello, error) {
	return x.store, nil
}

func (x *helloDaoStub) Record(ctx context.Context, hello *models.Hello) error {
	hello.ID = int64(len(x.store))
	x.store = append(x.store, hello)
	return nil
}
```

Now we are ready to write a handler test isolated from the underlying database.

## Handler tests

Golang has a many features that make testing easier.
We will cover the most basic, but there is much more!
Ask in the community, not all the features are as well documented as go production code.

Let see the simplest test we can have to get into testing.
The following example shows how to set up a context with:
* the prepared DAO stub,
* stubbed request using directly http package testing helper,
* response mock also using http package helper

We are testing whether the empty response is rendered correctly.


```go
package services_test

func TestListHellos(t *testing.T) {
	t.Run("handles empty database well", func(t *testing.T) {
		ctx := stub.WithHelloDao(context.Background())

		req, err := http.NewRequestWithContext(ctx, "GET", "/api/template/hellos", nil)
		require.NoError(t, err, "failed to create request")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(services.ListHellos)
		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code, "Wrong status code")
		assert.Equal(t, "[]\n", rr.Body.String())
	})
}
```

We have already made few design decisions:
* the test lives in a separate package
  * go recognizes a suffix `_test` as special package and allows two packages in the same folder
  * this makes it impossible to test internal unexported functions
* We use a `TestXX` function to test function `XX`
* We are nesting test cases using `t.Run`
* We are using [stretchr/testify](github.com/stretchr/testify) to add a bit of syntactical sugger.

## Unit test suite

Now we can add a make target `make test` that runs our tests.

## Database test suite

We are stubbing database to speed up and isolate our unit tests.
This is great, but how do we test our DAO tests.
There are ways to test with in-memory databases,
tho here we have decided to use real database.
Database is a major integration, and it is not so slow we can't run these tests with ease,
when limited to testing DAO methods.

We have a testing main in `internal/dao/tests/main.go`
and environment setup and teardown code in `internal/dao/tests/environments.go`.

All the other files are testing files and best practice is to have a file per DAO.

We aim at full coverage, to make sure our SQL queries are correct.
As every test, if you are doing anything special in your DAO method,
be sure to test for it.

Database test suite can be run `make test-database` and we can set it up in a separate suite in CI.
