# Templify

With Templify you can easily send emails and SMS notifications with or without the use of templates. This gives you an easy way to save, fill and send those templates without having to worry about anything else.

## Features

- Great developer experience due to extensive OpenAPI Spec with examples
- Send SMS with custom text
- Send emails with custom text (optional: with attachments)
- Save templates with placeholders
- Fill placeholders in templates
- Send filled templates via email or SMS

### Upcoming Features

- PDF templating
- Send emails with generated PDF files attached
- Previewing filled templates

## Installation

To send SMS and emails with Templify you will need a Twilio account. Create one here:
- [Twilio](https://login.twilio.com/u/signup?state=hKFo2SBVcU9sbDN5UzM4VEJnUjhRQTh2M3l3SkJ3aXR6djlJT6Fur3VuaXZlcnNhbC1sb2dpbqN0aWTZIExhQkVZbVhZZWpsRDB3eHVWRVBFMjBsSS0ycEJScDU3o2NpZNkgTW05M1lTTDVSclpmNzdobUlKZFI3QktZYjZPOXV1cks)

In your Twilio account overview you will find your account SID and Auth Token.

For sending Emails Templify is using Sendgrid. With your Twilio account you can log into Sendgrid as well and create an API Key in your account overview.
- [Sendgrid](https://login.sendgrid.com/login/)

Go ahead and create a copy of the `.env.example` file. Rename it to `.env`. Now you can fill in your environment variables (including your SID, Auth Token and API Key).

This project is used with docker. To get started take a look at the `taskfile.yml`.
