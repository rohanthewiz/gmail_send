package gmail_send

import "testing"

func TestGmailSend(t *testing.T) {
	type args struct {
		cfg GSMTPConfig
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Send Gmail",
			args: struct{cfg GSMTPConfig}{
				cfg: GSMTPConfig{
				AccountEmail: "youraccount@gmail.com",
				Word:         "secret", // Can use an app password here (Enable MFA then setup app password)
				FromName:     "Gmail Send Test",
				Subject:      "Test Mail",
				ToAddrs:      []string{"recipient@gmail.com"},
				BCCs:         []string{"bccRecipient@gmail.com"},
				Body:        `<body>This <span style="font-style:italic">is</span> an <b>example</b> body.<br>It contains two lines.</body>`,
			}},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := GmailSend(tc.args.cfg); (err != nil) != tc.wantErr {
				t.Errorf("GmailSend() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}