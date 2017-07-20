# Next steps


## API limitations workaround
- If the interval is more than a week and the resolution is hourly the program should break that into smaller requests.
Result: `{"error":{"statuscode":400,"statusdesc":"Bad Request","errormessage":"Interval is too big for this resolution"}}`

## Features
- Better string replacemen

### Extend output possibilities
- Elasticsearch

### Interfaces improvements
- error handling and output
  - Properly handle mysql errors
### Extend functionality
- Support different resolutions (now only hourly, API supports 'hour, day, week')
