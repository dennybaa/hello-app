[![](https://img.shields.io/circleci/build/github/dennybaa/hello-app.svg)](https://circleci.com/gh/dennybaa/hello-app/tree/master) [![](https://img.shields.io/github/release/dennybaa/hello-app.svg)](https://github.com/dennybaa/hello-app/releases)

# helloapp a simple Golang REST api exampe

hellapp creates a user in the database (MongoDB is used) and stores his date of birth. Notably it greats a user and provides days left to his birthday.


## Endpoints

### `/hello/:username` PUT

Creates or updates a user.

```bash
# creates foobar user
curl -XPUT -d'{"dateofbirth": "1980-10-05"}' http://localhost:8080/hello/foobar
# on success 204
```

### `/hello/:username` GET

Outputs days left to user's birthday.

```bash
curl http://localhost:8080/hello/foobar
# on success 200
# {"message": "Hello, foobar! Your birthday is in xxx day(s)"}
```


## Environment

- **MONGODB_URI** - a valid mongodb uri, ex: `mongodb://localhost:27017`
- **MOGNODB_DATABASE** - database name (default: `hello`)
- **MONGODB_CONNTIMEOUT** - time to wait for connection to establish (default: `15s`)
- **PORT** - port to listen on (default: `8080`)

### Usage

Use `docker-compose.yml` to bring and play up with the application.
