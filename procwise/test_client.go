package procwise

import (
	"io"
	"net/http"
	"strings"
)

type TestClient struct {
	Requests []*http.Request
}

func NewTestClient() *TestClient {
	return &TestClient{Requests: make([]*http.Request, 0)}
}

func (c *TestClient) Do(r *http.Request) (*http.Response, error) {
	c.Requests = append(c.Requests, r)
	var body string

	switch r.URL.Path {
	case "/users/sign_in.json":
		body = `{"report_frequency":"day","email":"super@proctorexam.com","role":"superuser","id":6,"institute_id":1,"student_number":null,"name":"ProctorExam SuperUser","created_at":"2022-07-07T11:00:16.677+02:00","updated_at":"2022-07-07T11:00:16.677+02:00","api_token":"lrOs9Pd7amlSSIDMYkqd3g","secret_key":"yR0-V7_T_O7jJ5o2GSSpVJaMkaiDCri8AkzBdy19ylg","receives_reports":true,"first_visit_guide":false,"global_proctor":true,"global_reviewer":true,"max_proctoring_sessions":null,"disabled":false,"proctoring_sessions_count":0,"deleted":false,"logo_image":"/logo_images/original/missing.png","institute_name":"ProctorExam"}`
	default:
		body = "{}"
	}
	return &http.Response{StatusCode: 200, Status: "200 OK", Body: io.NopCloser(strings.NewReader(body))}, nil
}
