package runtime

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// DashboardChecker contains logic allowing to check if Kyma instance is accessible
type DashboardChecker struct {
	client http.Client
	log    logrus.FieldLogger
}

func NewDashboardChecker(clientHttp http.Client, log logrus.FieldLogger) *DashboardChecker {
	clientHttp.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &DashboardChecker{
		client: clientHttp,
		log:    log,
	}
}

// AssertRedirectedToBusola sends request to the dashboardUrl and expects to be redirected to the Busola.
func (c *DashboardChecker) AssertRedirectedToBusola(dashboardURL string, busolaurl string) error {
	targetURL := c.buildTargetURL(dashboardURL)

	c.log.Infof("Calling the dashboard URL: %s", targetURL)
	resp, err := c.client.Get(targetURL)
	if err != nil {
		return errors.Wrapf(err, "while calling dashboard '%s'", dashboardURL)
	}

	if err = checkStatusCode(resp, http.StatusFound); err != nil {
		return err
	}

	if location, err := resp.Location(); err != nil {
		return errors.Wrap(err, "while getting response location")
	} else if location.Hostname() != busolaurl {
		return errors.Errorf("request was wrongly redirected: %s", location.String())
	}

	c.warnOnError(resp.Body.Close())
	c.log.Info("Successful response from the dashboard URL")

	return nil
}

// Kyma console URL won't redirect us to the UUA logging page, to achieve that we must call dex with a set of parameters
// state and nonce params are faked
func (c *DashboardChecker) buildTargetURL(dashboardURL string) string {
	consoleHost := strings.Split(strings.Split(dashboardURL, "//")[1], "/")[0]

	u := url.URL{
		Scheme: "https",
		Host:   consoleHost,
		Path:   "console-redirect",
	}

	return u.String()
}

func checkStatusCode(resp *http.Response, expectedStatusCode int) error {
	if resp.StatusCode != expectedStatusCode {
		// limited buff to ready only ~4kb, so big response will not blowup our component
		body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 4096))
		if err != nil {
			body = []byte(fmt.Sprintf("cannot read body, got error: %s", err))
		}
		return errors.Errorf("got unexpected status code, want %d, got %d, url: %s, body: %s",
			expectedStatusCode, resp.StatusCode, resp.Request.URL.String(), body)
	}
	return nil
}

func (c *DashboardChecker) warnOnError(err error) {
	if err != nil {
		c.log.Warn(err.Error())
	}
}
