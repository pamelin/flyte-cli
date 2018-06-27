package cmd

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"github.com/HotelsDotCom/flyte/flytepath"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"github.com/HotelsDotCom/flyte/httputil"
	"fmt"
)

func TestUploadDs_ShouldUploadDsFromFile(t *testing.T) {
	//given
	rec := struct {
		reqURL          string
		reqMethod       string
		reqContentType  string
		fileBody        []byte
		fileContentType string
	}{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.reqURL = r.URL.String()
		rec.reqMethod = r.Method
		rec.reqContentType = r.Header.Get(httputil.HeaderContentType)

		f, h, err := r.FormFile("value")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rec.fileBody, err = ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}

		rec.fileContentType = h.Header.Get(httputil.HeaderContentType)

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	dsFile := "./testdata/env.json"

	//when
	output, err := executeCommand("upload", "ds", "-f", dsFile, "--url", ts.URL)
	require.NoError(t, err)

	//then
	assert.Equal(t, flytepath.DatastorePath+"/env", rec.reqURL)
	assert.Equal(t, http.MethodPut, rec.reqMethod)
	assert.Contains(t, rec.reqContentType, "multipart/form-data; boundary=")

	wantContent, err := ioutil.ReadFile(dsFile)
	require.NoError(t, err)

	assert.Equal(t, wantContent, rec.fileBody)
	assert.Equal(t, httputil.MediaTypeJson, rec.fileContentType)

	l := fmt.Sprintf("Location: %s%s/%s", ts.URL, flytepath.DatastorePath, "env")
	assert.Contains(t, output, l)
}

func TestUploadDs_ShouldUploadDsFromFileWithDefaultsOverriddenByFlags(t *testing.T) {
	//given
	rec := struct {
		reqURL      string
		description string
		contentType string
	}{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.reqURL = r.URL.String()

		f, h, err := r.FormFile("value")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		rec.description = r.Form.Get("description")
		rec.contentType = h.Header.Get(httputil.HeaderContentType)

		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	file := "./testdata/env.json"
	name := "my-data"
	description := "This is my data"
	contentType := "text/plain; charset=us-ascii"

	//when
	output, err := executeCommand("upload", "ds",
		"-f", file,
		"--name", name,
		"--content-type", contentType,
		"--description", description,
		"--url", ts.URL)
	require.NoError(t, err)

	//then
	assert.Equal(t, flytepath.DatastorePath+"/"+name, rec.reqURL)
	assert.Equal(t, description, rec.description)
	assert.Equal(t, contentType, rec.contentType)

	l := fmt.Sprintf("Location: %s%s/%s", ts.URL, flytepath.DatastorePath, name)
	assert.Contains(t, output, l)
}

func TestUploadDs_ShouldCreateResource(t *testing.T) {
	//given
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	//when
	output, err := executeCommand("upload", "ds", "-f", "./testdata/env.json", "--url", ts.URL)
	require.NoError(t, err)

	//then
	l := fmt.Sprintf("Location: %s%s/%s", ts.URL, flytepath.DatastorePath, "env")
	assert.Contains(t, output, l)
}

func TestUploadDs_ShouldUpdateResource(t *testing.T) {
	//given
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	//when
	output, err := executeCommand("upload", "ds", "-f", "./testdata/env.json", "--url", ts.URL)
	require.NoError(t, err)

	//then
	l := fmt.Sprintf("Location: %s%s/%s", ts.URL, flytepath.DatastorePath, "env")
	assert.Contains(t, output, l)
}

func TestUploadDs_ShouldErrorForNon201Or204Response(t *testing.T) {
	//given
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	//when
	_, err := executeCommand("upload", "ds", "-f", "./testdata/env.json", "--url", ts.URL)

	//then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404 Not Found")
}
