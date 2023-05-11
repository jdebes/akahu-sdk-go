package akahu

import (
	"net/url"
	"time"
)

func paramsWithDateRange(startTime, endTime time.Time) url.Values {
	queryParams := url.Values{}
	queryParams.Add("start", startTime.Format(time.RFC3339))
	queryParams.Add("end", endTime.Format(time.RFC3339))

	return queryParams
}

func pathWithParams(path string, values url.Values) string {
	encodedPath, _ := url.Parse(path)
	encodedPath.RawQuery = values.Encode()

	return encodedPath.String()
}
