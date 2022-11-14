# Sample Twilio-Go App (Twilio Go SDK)

An example application made by the Developer Experience - Developer Interfaces team to showcase the usage of the [Twilio Go SDK](https://github.com/twilio/twilio-go), [Programmable Messaging](https://www.twilio.com/docs/sms), and [Programmable Voice](https://www.twilio.com/docs/voice).

Made for the purpose of providing an example application to train individuals on how Golang and the Twilio Go SDK may be used to build a web application.

## Prerequisites

Before installing and running this application, please ensure that you complete the following prerequisites:

- Have a MacOS and Linux environment
- Install [Docker Desktop](https://www.docker.com/products/docker-desktop/) and [Docker Compose](https://docs.docker.com/compose/install/)
- Install [Go](https://go.dev/dl/)
- Install a Code Editor
    - We recommend [Visual Studio Code](https://code.visualstudio.com/) with the [Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) installed
- Obtain a [Twilio account](https://www.twilio.com/login) and phone number with Messaging and Voice enabled
- Set the following environment variables using your Twilio account SID, Auth Token, Twilio phone number, desired log level, and localhost tunnel URL (i.e. ngrok tunnel URL):
    - TWILIO_ACCOUNT_SID
    - TWILIO_AUTH_TOKEN
    - TWILIO_PHONE_NUMBER
    - LOG_LEVEL
    - BASE_URL
- Have [twilio-cli](https://www.twilio.com/docs/twilio-cli/quickstart) installed and logged into your account to easily update your twilio phone number configuration with a single CLI command
- Obtain access to the [Twilio Console](https://console.twilio.com/)
- Install/obtain a localhost tunnel to access your localhost server on the internet. We'll use [ngrok](https://ngrok.com/) in our examples.

## Build the application

To build the application, run `make build-app` to build the application binary in `out/bin/sample-twilio-go`

To build the application Docker image, run `make docker-build`. You will need to run this build if you are running the application with Docker Compose.

## Run the application with Docker Compose

To run the application, first start ngrok to forward HTTP traffic to port 8080 in a separate terminal:

```
ngrok http 8080
```

Copy the `Forwarding` URL from ngrok's output. The URL should look something like this:

```
https://1404-2601-282-1200-118-e1d0-ba12-d474-d43a.ngrok.io
```

Next, update your phone number's SMS webhook. You can do this using the [twilio-cli](https://www.twilio.com/docs/twilio-cli/quickstart) command-line tool, or through the Twilio Console on your web browser. The following sections will walkthrough each method.

To update your phone number's SMS webhook using the twilio-cli, just run the following CLI command in your terminal. Substitute `<phone_number>` with your Twilio phone number in [E.164](https://www.twilio.com/docs/glossary/what-e164) format, and `<forwarding_url>` with the forwarding URL from ngrok.

```
twilio phone-numbers:update <phone_number> \
--sms-url=<forwarding_url>/sms \
--sms-method=POST
```

**Alternatively**, you can configure your SMS webhook on the Twilio Console in your web browser. To do this, Open the Twilio Console [Phone Numbers tab](https://console.twilio.com/us1/develop/phone-numbers/manage/incoming), select your active phone number to open its Configurations page.

Once in the Configure page, under the `Messaging Service` section, set the `A MESSAGE COMES IN` URL field to the following, substituting the `<forwarding_url>` portion for the forwarding URL copied from ngrok's output:

```
<forwarding_url>/sms
```

It should look something like this:

```
https://1404-2601-282-1200-118-e1d0-ba12-d474-d43a.ngrok.io/sms
```

Next, click the `Save` button below in the Console to apply your Twilio messaging webhook URL.

After configuring your SMS webhook, start the application and its service dependencies with the following command (make sure that you run `make docker-build` first!):

```
make docker-compose
```

**To stop the application**, run `make docker-stop`.

## Interacting with the application

This application is a mock review rewards program, where a customer can submit a request to participate in leaving a review over-the-phone in exchange for rewards.

To initiate the review process, first navigate to the customer registration page at http://localhost:8080/register, fill out the form with your contact information, and submit.

Then navigate to the SMS Campaigns Control Panel page at http://localhost:8080/campaigns-control-panel and click the "Start Review Campaign" button. This will query the database for stored customer contact information and initiate a review invite for each customer found.

After leaving a review, you may playback your review recording in the Twilio Console [Call recordings tab](https://console.twilio.com/us1/monitor/logs/call-recordings?frameUrl=%2Fconsole%2Fvoice%2Frecordings%2Frecording-logs%3Fx-target-region%3Dus1) by clicking the play icon for each call recording.

To view the total number of completed calls made by your application (and phone number), navigate to the http://localhost:8080/call-total page.

## Further Configuration

The application uses [Zap](https://github.com/uber-go/zap) for structured, leveled logging. To configure the application's log level, set the `-loglevel` CLI option flag when starting the app to one of the following values (i.e. `-loglevel=debug`)

|Value|Description|
|---|---|
|debug|Debug level logs (lowest level logs)|
|info|Info level logs (default)|
|warn|Warn level logs|
|error|Error level logs|
|panic|Panic level logs|
|fatal|Fatal level logs (highest level logs)|

## Getting help

For any questions or assistance regarding this application or the Twilio Go SDK, please feel free to reach out to us at our [#help-dev-interfaces](https://twilio.slack.com/archives/CGQPL0RPH) Slack channel.

For Twilio Go SDK documentation, examples, and code snippets, please refer to the documentation on [pkg.go.dev/github.com/twilio/twilio-go](https://pkg.go.dev/github.com/twilio/twilio-go).

## Appendix

### Code Areas of Interest

Areas of interest in the code base that serve as examples for using the Twilio SDK and troubleshooting errors.

- Initializing the Twilio SDK client. A pointer to this client instance is then saved in each application's service `client` field for later use (i.e. [SMSService.client](https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/sms/sms_service.go#L20))

    https://github.com/twilio-labs/sample-twilio-go/blob/main/cmd/app/main.go#L115

- Sending an SMS message with the SDK
    
    https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/sms/sms_service.go#L84

- Initiating a voice call with the SDK

    https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/voice/voice_service.go#L38

- Generating TwiML using the SDK

    https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/message/message.go#L44

- Initializing the SDK Request Validator. A pointer to this request validator instance is then saved in the application controller `reqValidator` field for later use (i.e. [ReviewController.reqValidator](https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/controller/review_controller.go#L30))

    https://github.com/twilio-labs/sample-twilio-go/blob/main/cmd/app/main.go#L139

- Using the Request Validator to validate requests are coming from Twilio

    https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/controller/review_controller.go#L128

- Debugging errors returned in an API response while using the SDK

    https://github.com/twilio-labs/sample-twilio-go/blob/main/pkg/sms/sms_service.go#L102
