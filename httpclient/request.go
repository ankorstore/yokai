package httpclient

import "net/http"

const (
	HeaderXRequestId  = "x-request-id"
	HeaderTraceParent = "traceparent"
)

// CopyRequestHeaders performs a copy of a specified list of headers between two [http.Request].
func CopyRequestHeaders(source *http.Request, dest *http.Request, headers ...string) {
	for _, header := range headers {
		canonicalHeader := http.CanonicalHeaderKey(header)

		for _, value := range source.Header.Values(canonicalHeader) {
			dest.Header.Add(canonicalHeader, value)
		}
	}
}

// CopyObservabilityRequestHeaders performs a copy of x-request-id and traceparent headers between two [http.Request].
func CopyObservabilityRequestHeaders(source *http.Request, dest *http.Request) {
	CopyRequestHeaders(source, dest, HeaderXRequestId, HeaderTraceParent)
}
