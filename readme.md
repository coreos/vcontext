vcontext: context for validation

Go's `encoding/json` library does not support retrieving information about where in the source
structures occur. This information is useful when validating user supplied configuration and
wanting to give line and column information about where the errors are.

This project aims to allow decoding to a structure that contains metadata about where each object,
array, etc occurs in a json blob. It also aims to keep that structure agnostic so other languages
with similar primitive structures (e.g. yaml) can be added later.

This project is not stable.
