# Nothing to see here

This directory contains a version of the main project but sans-Temporal. It is entirely illustrative and meant to _very, very roughly_ approximate patterns that might be used to mitigate failure without using Temporal.

Those patterns include (but are not limited to):
* Exponential backoff retries
* Manual checking error statuses to determine retryability
* Use of a distributed keystore to maintain "where are we" state
* Skipping or resuming steps based on that state
* Parallel channels to be able to catch timeouts
