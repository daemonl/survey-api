Mini Survey Response API
------------------------

Allows public users to submit and view results to a pet survey.


## TODO:

### JSON Schema / OAS validation

- Submitting `age: 0` and omitting `age` is equivalent, and valid.
- Attempting to set age to a string will throw 500.
- There is no documentation for the expected API.

### Integration testing for the Mongo and S3 stores

There is very limited testable logic in the implementations themselves, however
there is complexity in the interaction. 

### Server Shutdown

- Server will close all connections immediately on exit signal
- Docker's default exit signal may not actually shut down the server
