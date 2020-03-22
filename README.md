# Gmail Send

So you want to send email programmatically. You have a Gmail account, yes? Then why pay for some other service when you can send through gmail for free?
Obviously you can't send too many.

## Usage

```go
package main

import (
    "log"
    gmail "github.com/rohanthewiz/gmail_send"
)

func main() {
    cfg := gmail.GSMTPConfig{
        AccountEmail: "youraccount@gmail.com",
        Word:         "secret", // Can use an app password here (Enable MFA then setup app password)
        FromName:     "Gmail Send Test",
        Subject:      "Test Mail",
        ToAddrs:      []string{"recipient@gmail.com"},
        BCCs:         []string{"bccRecipient@gmail.com"},
        Body:        `<body>This <span style="font-style:italic">is</span> an <b>example</b> body.<br>It contains two lines.</body>`,
    }
    
    err := gmail.GmailSend(cfg)
    if err != nil {
        log.Println(err)
    }
}
```