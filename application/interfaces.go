package application

type UseCases interface {
	// Shorten receives a URL and returns a shortened version of it
	// e.g. http://www.google.com -> abc123
	Shorten(url string) (string, error)

	// Expand receives a shortened URL and returns the original URL
	// e.g. abc123 -> http://www.google.com
	Expand(url string) (string, error)
}

type Dal interface {
	// Save receives a URL and returns a shortened version of it
	Save(shortened, url string) error

	// Find receives a shortened URL and returns the original URL
	Find(short string) (string, error)
}
