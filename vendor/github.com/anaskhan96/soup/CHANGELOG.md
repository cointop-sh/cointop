## v1.1

### Added

- Cookies can be added to the HTTP request, either via the `Cookies` map or the `Cookie()` function
- Function `GetWithClient()` provides the ability to send the request with a custom HTTP client
- Function `FindStrict()` finds the first instance of the mentioned tag with the exact matching values of the provided attribute (previously `Find()`)
- Function `FindAllStrict()` finds all the instances of the mentioned tag with the exact matching values of the attributes (previously `FindAll()`)

## Changed

- Function `Find()` now finds the first instance of the mentioned tag with any matching values of the provided attribute.
- Function `FindAll()` now finds all the instances of the mentioned tag with any matching values of the provided attribute.

---