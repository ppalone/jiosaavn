package jiosaavn

import "fmt"

// Search Option
type SearchOption func(opts *searchOptions)

// Search Options
type searchOptions struct {
	page  int
	limit int
	query string
}

func (o *searchOptions) validate() error {
	if len(o.query) == 0 {
		return fmt.Errorf("search query cannot be empty")
	}

	if o.limit < 10 || o.limit > 40 {
		return fmt.Errorf("limit must be between 10 and 40")
	}

	return nil
}

// defaultSearchOpts returns the default search options
func defaultSearchOpts() *searchOptions {
	return &searchOptions{
		page:  1,
		limit: 10,
	}
}

// WithPage sets the page search option
func WithPage(page int) SearchOption {
	return func(opts *searchOptions) {
		opts.page = page
	}
}

// WithLimit sets the limit search option
func WithLimit(limit int) SearchOption {
	return func(opts *searchOptions) {
		opts.limit = limit
	}
}
