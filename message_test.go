package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMessage(t *testing.T) {
	//m, err := ParseMessage("{\"log\":\"time=\\\"2017-04-07T02:22:09Z\\\" level=info msg=\\\"response completed\\\" go.version=go1.7.3 http.request.host=\\\"localhost:5000\\\" http.request.id=a3e73015-2d9a-479b-9bb0-d995eabc1730 http.request.method=GET http.request.remoteaddr=\\\"172.17.0.1:39604\\\" http.request.uri=\\\"/v2/busybox/manifests/latest\\\" http.request.useragent=\\\"docker/1.12.6 go/go1.6.3 git-commit/d5236f0 kernel/4.9.16-coreos-r1 os/linux arch/amd64 UpstreamClient(Docker-Client/1.12.6 \\\\\\\\(linux\\\\\\\\))\\\" http.response.contenttype=\\\"application/vnd.docker.distribution.manifest.v2+json\\\" http.response.duration=2.131865ms http.response.status=200 http.response.written=527 instance.id=2643fd78-8819-4205-9c63-bd6aec78b7d0 version=v2.6.0 \\n\",\"stream\":\"stdout\",\"time\":\"2017-04-07T02:22:09.302292883Z\"}")
	m, err := ParseMessage("{\"log\":\"time=\\\"2017-04-09T09:14:03.819998833Z\\\" level=info msg=\\\"response completed\\\" go.version=go1.7.3 http.request.host=docker.paddlepaddle.org http.request.id=8db3d408-8d99-4b0a-9430-34b61e6e8f63 http.request.method=GET http.request.remoteaddr=47.90.80.226 http.request.uri=\\\"/v2/paddle/manifests/0.10.0rc2\\\" http.request.useragent=\\\"docker/1.12.6 go/go1.6.3 git-commit/d5236f0 kernel/4.10.4-coreos-r1 os/linux arch/amd64 UpstreamClient(Docker-Client/1.12.6 \\\\\\\\(linux\\\\\\\\))\\\" http.response.contenttype=\\\"application/vnd.docker.distribution.manifest.v2+json\\\" http.response.duration=5.87701ms http.response.status=200 http.response.written=1588 instance.id=85cc17d4-f6d2-436c-bd2a-264a2296f438 service=registry version=v2.6.0 \\n\",\"stream\":\"stdout\",\"time\":\"2017-04-09T09:14:03.820237214Z\"}")
	assert.Nil(t, err)
	assert.True(t, strings.EqualFold(m.ImageName, "paddle"))

}
