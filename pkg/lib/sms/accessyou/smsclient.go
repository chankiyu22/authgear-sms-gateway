package accessyou

import (
	"log/slog"
	"net/http"
	"regexp"

	"github.com/authgear/authgear-sms-gateway/pkg/lib/sms/smsclient"
)

type AccessYouClient struct {
	BaseUrl   string
	Client    *http.Client
	AccountNo string
	User      string
	Pwd       string
	Sender    string
	Logger    *slog.Logger
}

func NewAccessYouClient(
	httpClient *http.Client,
	baseUrl string,
	accountNo string,
	user string,
	pwd string,
	sender string,
	logger *slog.Logger,
) *AccessYouClient {
	if baseUrl == "" {
		baseUrl = "http://sms.accessyou-anyip.com"
	}
	return &AccessYouClient{
		Client:    httpClient,
		BaseUrl:   baseUrl,
		AccountNo: accountNo,
		User:      user,
		Pwd:       pwd,
		Sender:    sender,
		Logger:    logger,
	}
}

var plusHyphensRegexp = regexp.MustCompile(`[\+\-]+`)

func fixPhoneNumber(phoneNumber string) string {
	// Access you phone number should have no + and -
	return plusHyphensRegexp.ReplaceAllString(phoneNumber, "")
}

func (n *AccessYouClient) Send(options *smsclient.SendOptions) (*smsclient.SendResult, error) {
	to := fixPhoneNumber(string(options.To))

	respData, sendSMSResponse, err := SendSMS(
		n.Client,
		n.BaseUrl,
		n.AccountNo,
		n.User,
		n.Pwd,
		n.Sender,
		to,
		options.Body,
	)
	if err != nil {
		return nil, err
	}

	return &smsclient.SendResult{
		ClientResponse: respData,
		Success:        sendSMSResponse.Status == "100",
	}, nil
}

var _ smsclient.RawClient = &AccessYouClient{}
