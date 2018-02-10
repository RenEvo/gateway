Gateway
=======

A dedicated http content server with backend reverse proxy support.

## But why did you create this?

I have used nginx for years to host websites and use it as an API gateway. My front end "SPA" was completely static, then I spent a ton of time copy/pasting all of my backend routes and upstreams into the conf files. There is a better way. Additionally, NGINX has crap for monitoring unless you pay for support (personal websites aren't really cost effective), so I wanted a way to find out how my servers were performing, without adding a lot of effort (just give me a dashboard).

## Primary Project Goals

Host static files in memory.
Read backends from consul or docker and automatically map the routes for them.
Automatically retry backend requests when the server is unavailable (can't connect).
Deploy with my application (FROM gateway; COPY site /var/www/public)
Provide a useful modern dashboard for viewing my website.
Integrate with external monitoring, alerting, and analysis software (ELK, TICK, etc..)
Scale to "thousands" of requests per second. I don't have high hopes, as I don't need more than that, but feel free to help me improve it for better perf.
Remove my dependency on nginx.
Run on any platform (well, whatever GO supports).

## sample.yaml

This configuration file represents the overall wanted specification for the server and will be provided as a "full example". The server will have a way to dump a "default" yaml file 
using a command line argument.

## Running

```bash
go run main.go
```

### Command Line Options

These are in the projects main.go, but provided here for reference.

| Command | Description | Default |
|---------|-------------|---------|
|`-config`| Path to a configuration yaml file | NA - uses gateway default configurations |


### Environmental Variables

These are setting that allow you to tweak some of the underlying system behavior.

| Variable                  | Description               | Default           |
|---------------------------|---------------------------|-------------------|
|`GATEWAY_DEBUG`| When set to true, will output debug logging| false |
|`GATEWAY_SITE_MEMORY_FILE_DISABLE`| The web hosting will not read files into memory and only serve from disk| false |
|`GATEWAY_SITE_MEMORY_FILE_MAX_SIZE`| The maximum file size to put in memory for the web hosting, if you have more memory, use it! | 2mb |


## Credits

* `favicon.ico` in default ./public/www taken from https://www.shareicon.net/balancing-elastic-copy-networking-compute-load-92244

