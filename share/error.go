package share

import "errors"

// Now, all the errors are in one place, and you can easily see them all at once.
var ErrUnauthorized = errors.New("unauthorized")
var ErrRecordNotFound = errors.New("record not found")
