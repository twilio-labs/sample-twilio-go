package message

import "fmt"

const (
	// Static messages
	GREETING = "Hello from Twilio Resorts and Spas. Thank you for" +
		" reaching out to our review rewards program, where you can leave a" +
		" review over the phone and receive $50 Twilbucks off your next stay!"

	PARTICIPATION_INVITE = "Would you like to leave a review today?" +
		" (yes/no)"
	PARTICIPATION_ACCEPT_RESPONSE = "Thank you for choosing to participate." +
		" Your Twilbucks will be available in your account after leaving your" +
		" review."
	PARTICIPATION_DECLINE_RESPONSE = "We're sorry to hear that. Please reach" +
		" out to us in the future if you wish to participate."
	PARTICIPATION_INVITE_FALLBACK = "Sorry, I did not understand that. Please" +
		` say "yes" or "no".`

	ASK_FOR_NAME = `Please enter your name.`

	CALL_NOTIFICATION = `You will receive a call shortly to leave your review.`

	GOODBYE                = `Thank you for reaching out to us. Goodbye.`
	PARTICIPATION_THANKYOU = "We have received your review. Thank you for" +
		" participating!"

	// Templates
	GREETING_WITH_NAME_TEMPLATE = `Hello, %s!`

	// TwiML templates
	GREETING_SAY_TWIML_TEMPLATE = `
	<Say>Ahoy! Greetings from Twilio Resorts and Spas.</Say>
`
	INSTRUCTIONS_SAY_TWIML = `
	<Say>Thank you for participating in our review rewards program, where you
will receive 50 Twilbucks by leaving a review of your recent stay! You may leave
a review up to 1 minute long. Please leave your review after the beep.</Say>
`
	RECORD_VERB_TWIML = `
	<Record timeout="10" maxLength="60" playBeep="true"></Record>	
`
)

func GetHelloMessage(name string) string {
	return fmt.Sprintf(GREETING_WITH_NAME_TEMPLATE, name)
}

func GetReviewGreetingAndInstructionsTwiML() string {
	var twiml string = "<Response>"
	twiml += GREETING_SAY_TWIML_TEMPLATE
	twiml += INSTRUCTIONS_SAY_TWIML
	twiml += RECORD_VERB_TWIML
	twiml += "</Response>"
	return twiml
}
