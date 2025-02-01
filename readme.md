This is a really simple SMTP server built from the socket layer in an effort to practice Go.

This initial version is a dumped down version from the [RFC 821](https://www.rfc-editor.org/rfc/rfc821)
There are a few things that can be improved here to make it work even better.
Things like error handling, mailboxes, verifying received content from connecting
clients are things that can really improve this.

Despite its simplicity, it successfully receives emails—I was even able to send one to it from my Gmail account! ;)

The plan is to extend this to be complaint with [RFC 5321](https://www.rfc-editor.org/rfc/rfc5321.html) and finally,
making this available over API (as a temp mail service) or to extend it to be a mail server
for local development. Picking up emails and rendering them in a nice UI.


### Some kinda roadmap
- Improve error handling, data validation, those bits.
- Provision mailboxes into which emails are written (Simple sqlite db)
- Parse email. Get important stuff like subjects, CC, and others from the email headers
- API Layer:
  - Create mailboxes, just records of accounts in the DB
  - Viewing content in user mailbox
- Validating email recipients.