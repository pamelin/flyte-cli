package cmd

import (
	"testing"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"net/http"
	"strings"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
	"github.com/HotelsDotCom/flyte/httputil"
	"github.com/HotelsDotCom/flyte/flytepath"
)

type requestRec struct {
	request http.Request
	body    []byte
}

func TestUploadCommand_ShouldUploadFlowFromJsonFile(t *testing.T) {
	//given
	rec := requestRec{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.request = *r
		rec.body, _ = ioutil.ReadAll(r.Body)
		w.Header().Set("Location", flytepath.FlowsPath+"/my-flow")
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	flowFile := "./testdata/my-flow.json"
	host := strings.Replace(ts.URL, "http://", "", -1)

	//when
	output, err := executeCommand("upload", flowFile, "--host="+host)
	require.NoError(t, err)

	//then
	assert.Equal(t, flytepath.FlowsPath, rec.request.URL.String())
	assert.Equal(t, http.MethodPost, rec.request.Method)
	assert.Equal(t, httputil.MediaTypeJson, rec.request.Header.Get(httputil.HeaderContentType))

	wantBody, err := ioutil.ReadFile(flowFile)
	require.NoError(t, err)
	assert.Equal(t, wantBody, rec.body)

	assert.Contains(t, output, "Location: "+flytepath.FlowsPath+"/my-flow")
}

func TestUploadCommand_ShouldFailWhenFlyteHostReturnsNon201(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	host := strings.Replace(ts.URL, "http://", "", -1)

	_, err := executeCommand("upload", "./testdata/my-flow.json", "--host="+host)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "cannot upload flow\nHTTP/1.1 400 Bad Request")
}

func TestUploadCommand_ShouldFailForNonJsonOrYamlFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	host := strings.Replace(ts.URL, "http://", "", -1)

	_, err := executeCommand("upload", "./testdata/my-flow.haha", "--host="+host)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "cannot upload flow: unsupported file format .haha")
}

func TestUploadCommand_ShouldUploadFlowFromYamlFile(t *testing.T) {
	//given
	rec := requestRec{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.request = *r
		rec.body, _ = ioutil.ReadAll(r.Body)
		w.Header().Set("Location", flytepath.FlowsPath+"/my-flow")
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	flowFile := "./testdata/my-flow.yaml"
	host := strings.Replace(ts.URL, "http://", "", -1)

	//when
	output, err := executeCommand("upload", flowFile, "--host="+host)
	require.NoError(t, err)

	//then
	assert.Equal(t, flytepath.FlowsPath, rec.request.URL.String())
	assert.Equal(t, http.MethodPost, rec.request.Method)
	assert.Equal(t, httputil.MediaTypeYaml, rec.request.Header.Get(httputil.HeaderContentType))

	wantBody, err := ioutil.ReadFile(flowFile)
	require.NoError(t, err)
	assert.Equal(t, wantBody, rec.body)

	assert.Contains(t, output, "Location: "+flytepath.FlowsPath+"/my-flow")
}

func TestUploadCommand_ShouldUploadFlowFromYmlFile(t *testing.T) {
	//given
	rec := requestRec{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.request = *r
		rec.body, _ = ioutil.ReadAll(r.Body)
		w.Header().Set("Location", flytepath.FlowsPath+"/my-flow")
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	flowFile := "./testdata/my-flow.yml"
	host := strings.Replace(ts.URL, "http://", "", -1)

	//when
	output, err := executeCommand("upload", flowFile, "--host="+host)
	require.NoError(t, err)

	//then
	assert.Equal(t, flytepath.FlowsPath, rec.request.URL.String())
	assert.Equal(t, http.MethodPost, rec.request.Method)
	assert.Equal(t, httputil.MediaTypeYaml, rec.request.Header.Get(httputil.HeaderContentType))

	wantBody, err := ioutil.ReadFile(flowFile)
	require.NoError(t, err)
	assert.Equal(t, wantBody, rec.body)

	assert.Contains(t, output, "Location: "+flytepath.FlowsPath+"/my-flow")
}
