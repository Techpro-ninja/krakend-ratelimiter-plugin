# Custom KrakenD Ratelimiter Plugin

This is a custom Krakend Ratelimiter Plugin for demonstration.

Read the full blog here [https://techpro.ninja/krakend-plugin-golang-tutorial/](https://techpro.ninja/krakend-plugin-golang-tutorial/)

## Building the plugin Steps

1. Using docker **krakend/builder:2.2.1** to generate the `.so` file. _The version of builder should be same as the Krakend deployment or else the plugin will not work_.

```sh
docker run -it -v "some_path:/go/src/krakend-ratelimiter-plugin" -w /go/src/krakend-ratelimiter-plugin krakend/builder:2.2.1 go build -buildmode=plugin -o krakend-ratelimiter-plugin.so .

```

Once built, copy the `krakend-ratelimiter-plugin.so` file to **tests** folder from the Docker container

2. Build Docker image

```sh
cd tests
docker build . -t test-v1
```

3. Running the Docker image

```sh
docker run -it -p 8001:8001 test-v1
```

At this point, you can check

1. [http://localhost:8001/v1/users/](http://localhost:8001/v1/users/)
2. [http://localhost:8001/v1/posts/](http://localhost:8001/v1/posts/)

To check bucket details
[http://localhost:8001/\_\_bucket-tracker](http://localhost:8001/__bucket-tracker)

# Krakend.json configuration

This is a **service level** plugin meaning it resides in the root of **krakend.json** file

```json
"plugin": {
    "pattern": ".so",
    "folder": "/etc/krakend/plugins/"
},
    ...

"extra_config":{
    "plugin/http-server":{
        "name":["krakend-ratelimiter-plugin"],
        "krakend-ratelimiter-plugin":{
            "trackerPath": "/__bucket-tracker"
        }
    }
}
```

Configurable values are `trackerPath`. Do not change anything else.

## trackerPath

Path to track the bucket capacity and token stock details inside the plugin. For example, `https://<SERVER_ADDRESS>/__bucket-tracker`.

### Note: If you change the plugin name from `krakend-ratelimiter-plugin`, do remeber to change the following line in `main.go` file as well

```go
var pluginName = "krakend-ratelimiter-plugin"
```

Also, remember to change the values in `Dockefile` in **tests** folder
