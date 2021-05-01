# techsalesdev-go
This is for techsales.dev backend

Check the releases page to run locally.

You need a .env file which consits of the following variables:

```bash
SMTP_ACCOUNT=someEmailAccountForSendingEmails@gmail.com
SMTP_PASS=theSMTPPassword
REDIS_PASS=thePasswordForRedis
REDIS_ENDPOINT=127.0.0.1:6379
API_DEV_PORT=8084
ACCESS_TOKEN_SECRET=SomeAccessTokenThatCanBeUsedAcrossMicroServices
ETHERIUM_NETWORK=https://someEtheriumNetwork.com
```

Then run the binary for your archetecture within the same folder as the .env

## Requirements

* The backend is written specifically for Redis, and you can spin up a Redis Container with Docker Compose

```bash
docker-compose up -d;
```

* We've hard coded the images that are automactically pushed into the Redis database as base64 for quick access, therefore it is required to have a .products/photos/ folder
in order to run the application.  Details of what images are needed can be found in handlers/createProducts.go filenames **must** match the pointers.

* Once Redis is running, and the .env file is filled out, and photos are populated matching the hard coded handlers/createProducts.go file, everything should work as expected.

