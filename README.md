# tamaki

This is a slack bot which works on GAE. 

# Deploy

0. Download the App Engine SDK for Go

https://cloud.google.com/appengine/docs/standard/go/download

1. Deploy it to your application.

```goapp deploy -application [YOUR_PROJECT_ID] -version [YOUR_VERSION_ID] .```

2. Set environment variables on your application.

```appcfg.py update . -E SLACK_API_TOKEN:"XXXXX...." -E REDIS_URL:"XXXXX..."  -A [YOUR_PROJECT_ID] -V [YOUR_VERSION_ID]```