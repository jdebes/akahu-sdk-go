package akahu

import (
	"net/url"
	"time"
)

func buildDateRangePath(path string, startTime, endTime time.Time) string {
	queryParams := url.Values{}
	queryParams.Add("start", startTime.Format(time.RFC3339))
	queryParams.Add("end", endTime.Format(time.RFC3339))

	urlPath, _ := url.Parse(path)
	urlPath.RawQuery = queryParams.Encode()

	return urlPath.String()
}
