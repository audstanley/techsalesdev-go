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
```

Then run the binary for your archetecture within the same folder as the .env
