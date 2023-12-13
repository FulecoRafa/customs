# Customs

A ~blazingly fast~ program to log HTTP requests midway

Logs are in a structured format, compatible with log analysis tools

## Usage

Run the program in your terminal using a port that is not being used
redirecting to the URL you want to send the request.

### Example
```bash
customs -o http -r 8765:myserver.com
```

You can always use the help command to get information

```bash
customs -h
```

### Flags

#### `-o|--output [format]`
The output format of logs. One of: curl; http (default "http")

#### `-r|--redirect [redirect string]`
Lists of ports redirecting to URLs in format `port:url`

#### `--debug`
Print debug logs
