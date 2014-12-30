package elasticthought

import (
	"fmt"
	"sync"
	"testing"

	"github.com/couchbaselabs/go.assert"
	"github.com/tleyden/fakehttp"
)

// starting port to use for fakehttp servers
var port = 4444

// since I'm not sure if NextPort() can be called by multiple threads,
// protect it with a mutex
var portMutex = &sync.Mutex{}

// the fakehttp library doesn't provide an easy way to shutdown an http
// server, so as a workaround, run each fake http server on a unique port
// for each test.
func NextPort() int {
	portMutex.Lock()
	defer portMutex.Unlock()
	port2use := port
	port += 1
	return port2use
}

func TestInsertClassifier(t *testing.T) {

	testServer := fakehttp.NewHTTPServerWithPort(NextPort())
	testServer.Start()

	// response when go-couch tries to see that the server is up
	testServer.Response(200, jsonHeaders(), `{"version": "fake"}`)

	// response when go-couch check is db exists
	testServer.Response(200, jsonHeaders(), `{"db_name": "db"}`)

	// insert succeeds
	testServer.Response(200, jsonHeaders(), `{"id": "classifier", "rev": "bar", "ok": true}`)

	configuration := NewDefaultConfiguration()
	configuration.DbUrl = fmt.Sprintf("%v/db", testServer.URL)

	classifier := NewClassifier()
	classifier.SpecificationUrl = "http://s3.com/proto.txt"
	classifier.TrainingJobID = "123"
	classifier.Configuration = *configuration

	err := classifier.Insert()

	assert.True(t, err == nil)
	assert.Equals(t, "classifier", classifier.Id)
	assert.Equals(t, "bar", classifier.Revision)

}
