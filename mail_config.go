package gmail_send


type GSMTPConfig struct {
	AccountEmail string
	Word string
	FromName string
	Subject string
	ToAddrs []string
	BCCs []string
	Body string
}
