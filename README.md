# elc-bin-notifier
Retrieves next bin collections for East Lothian Council and sends day-before notifications via Pushover.

## Prerequisites
1. The full URL of the ELC bin collection details for your area.

This can be retrieved by inspecting the URL (using your browsers developer tools) that is hit after [entering your postcode and address here](http://collectiondates.eastlothian.gov.uk/your-calendar).
It should look something like this:

`http://collectiondates.eastlothian.gov.uk/ajax/your-calendar/load-recycling-2022.asp?id=ELC-XXXXXX`

1. A Pushover licence (or free trial).

## Environment Variables
These should be configured before running the application.
- `API_URL`: The ELC bin collection URL retrieved using the above instructions.
- `PUSHOVER_TOKEN`: The API token of your Pushover application.
- `PUSHOVER_TARGET`: The Pushover target user/group for the day-before notification.

## Usage
- Clone the repository
- Build using the Go compiler `go build -o elc-bin-notifier`
- Run the application manually `./elc-bin-notifier`
- Or (more usefully) run on a schedule using `cron`
