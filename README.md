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

#### `-l|--logs [format]`
Log format. One of: json; kv (default "kv")

##### `kv`
```
time=2023-12-13T07:55:33.780-03:00 level=DEBUG msg="Starting application" ports="[7070 -> localhost:6969]" outputFormat=http
```

##### `json`
```
{"time":"2023-12-13T07:55:40.611606-03:00","level":"DEBUG","msg":"Starting application","ports":[{"Port":7070,"Destination":"localhost:6969"}],"outputFormat":"http"}
```

#### `-r|--redirect [redirect string]`
Lists of ports redirecting to URLs in format `port:url`

#### `--debug`
Print debug logs
