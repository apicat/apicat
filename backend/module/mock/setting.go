package mock

type Option func(m *MockServer)

func WithApiUrl(url string) Option {
	return func(m *MockServer) {
		m.ApiUrl = url
	}
}

func WithApiPath(path string) Option {
	return func(m *MockServer) {
		m.ApiPath = path
	}
}
