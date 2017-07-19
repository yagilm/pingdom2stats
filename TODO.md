# Next steps


## API limitations workaround
- If the interval is more than a week and the resolution is hourly the program should break that into smaller requests.
Result: `{"error":{"statuscode":400,"statusdesc":"Bad Request","errormessage":"Interval is too big for this resolution"}}`

## Features

### Extend output possibilities
- Elasticsearch
- Raw MySQL statements

### Interfaces improvements
- error handling and output

### Extend functionality
- Support different resolutions (now only hourly, API supports 'hour, day, week')
