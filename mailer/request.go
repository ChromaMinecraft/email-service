package mailer

type Request interface{}

type request struct {
	Host    string
	Port    int
	From    string
	To      []string
	Subject string
	Body    string
}

func NewRequest(opts ...func(*request)) *request {
	r := &request{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithHost(host string) func(*request) {
	return func(r *request) {
		r.Host = host
	}
}

func WithPort(port int) func(*request) {
	return func(r *request) {
		r.Port = port
	}
}

func WithFrom(from string) func(*request) {
	return func(r *request) {
		r.From = from
	}
}

func WithTo(to []string) func(*request) {
	return func(r *request) {
		r.To = to
	}
}

func WithSubject(subject string) func(*request) {
	return func(r *request) {
		r.Subject = subject
	}
}

func WithBody(body string) func(*request) {
	return func(r *request) {
		r.Body = body
	}
}
