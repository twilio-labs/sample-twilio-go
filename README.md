# Review Rewards Example App (Twilio Go SDK)

An example application made by the Developer Experience - Developer Interfaces team to showcase the usage of the [Twilio Go SDK](https://github.com/twilio/twilio-go), [Programmable Messaging](https://www.twilio.com/docs/sms), and [Programmable Voice](https://www.twilio.com/docs/voice).

Made for the purpose of providing an example application to train individuals on how Golang and the Twilio Go SDK may be used to build a web application.

## Prerequisites

Before installing and running this application, please ensure that you complete the following prerequisites:

- Have a MacOS and Linux environment
- Install [Go](https://go.dev/dl/)
- Install a Code Editor
    - We recommend [Visual Studio Code](https://code.visualstudio.com/) with the [Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) installed
- Obtain a [Twilio account](https://www.twilio.com/login) and phone number with Messaging and Voice enabled
- Save your Twilio account SID and Auth Token in the following environment variables:
    - TWILIO_ACCOUNT_SID
    - TWILIO_AUTH_TOKEN
- Obtain access to the [Twilio Console](https://console.twilio.com/)
- Install/obtain a localhost tunnel to access your localhost server on the internet. We'll use [ngrok](https://ngrok.com/) in our examples.
- (**Optional**) Have [twilio-cli](https://www.twilio.com/docs/twilio-cli/quickstart) installed and logged into your account to easily update your twilio phone number configuration with a single CLI command

## Build the application

To build the application, run `go build`

## Run the application

To run the application, first start ngrok to forward HTTP traffic to port 8080 in a separate terminal:

```
ngrok http 8080
```

Copy the `Forwarding` URL from ngrok's output. The URL should look something like this:

```
https://1404-2601-282-1200-118-e1d0-ba12-d474-d43a.ngrok.io
```

Next, update your phone number's SMS webhook. You can either do this on your web browser in the Twilio Console, or using [twilio-cli](https://www.twilio.com/docs/twilio-cli/quickstart). The following sections will walkthrough each method.

To update your phone number's SMS webhook in the Twilio Console, open the Twilio Console [Phone Numbers tab](https://console.twilio.com/us1/develop/phone-numbers/manage/incoming), select your active phone number to open its Configurations page.

Once in the Configure page, under the `Messaging Service` section, set the `A MESSAGE COMES IN` URL field to the following, substituting the `<forwarding_url>` portion for the forwarding URL copied from ngrok's output:

```
<forwarding_url>/sms
```

It should look something like this:

```
https://1404-2601-282-1200-118-e1d0-ba12-d474-d43a.ngrok.io/sms
```

Next, click the `Save` button below in the Console to apply your Twilio messaging webhook URL so that SMS messages are forwarded through ngrok to your application.

**Alternatively**, if you have [twilio-cli](https://www.twilio.com/docs/twilio-cli/quickstart) installed and setup for your account, you can run the following to configure your SMS webhook in a single CLI command. Just substitute `<phone_number>` with your Twilio phone number in [E.164](https://www.twilio.com/docs/glossary/what-e164) format, and `<forwarding_url>` with the forwarding URL from ngrok. 

```
twilio phone-numbers:update <phone_number> \
--sms-url=<forwarding_url>/sms \
--sms-method=POST
```

After configuring your SMS webhook, start the application with the following command, replacing `<phone_number>` with your Twilio phone number in E.164 format and `<forwarding_url>` with your ngrok forwarding URL:

```
./review-rewards-example-app -from=<phone_number> -url=<forwarding_url>
```

Your command should look something like this:

```
./review-rewards-example-app -from=+15555555555 -url=https://1404-2601-282-1200-118-e1d0-ba12-d474-d43a.ngrok.io
```

## Interacting with the application

To initiate the review process, send an SMS message to your phone number.

After leaving a review, you may playback your review recording in the Twilio Console [Call recordings tab](https://console.twilio.com/us1/monitor/logs/call-recordings?frameUrl=%2Fconsole%2Fvoice%2Frecordings%2Frecording-logs%3Fx-target-region%3Dus1) by clicking the play icon for each call recording.

To view the total number of completed calls made by your application (and phone number), navigate to the http://localhost:8080/call-total page.

## Getting help

For any questions or assistance regarding this application or the Twilio Go SDK, please feel free to reach out to us at our [#help-dev-interfaces](https://twilio.slack.com/archives/CGQPL0RPH) Slack channel.

For Twilio Go SDK documentation, examples, and code snippets, please refer to the documentation on [pkg.go.dev/github.com/twilio/twilio-go](https://pkg.go.dev/github.com/twilio/twilio-go).
