package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lensesio/tableprinter"
	"github.com/spf13/cobra"
)

// getFegCmd represents the getFeg command
var getFegCmd = &cobra.Command{
	Use:   "getFeg",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getFeg()
	},
}

func init() {
	rootCmd.AddCommand(getFegCmd)
}

type Feg struct {
	NetworksID []string `header:"networkid"`
}

func getFeg() {
	url := "https://api.galaxy.nitin.com/magma/v1/feg_lte"
	responseBytes := getFegData(url)

	var feg []Feg

	if err := json.Unmarshal(responseBytes, &feg); err != nil {
		log.Printf("Could not unmarshal reponseBytes. %v", err)
	}
	// log.Printf("%+v", feg)
	tableprinter.Print(os.Stdout, feg)
}

func fegCheckError(err error, hdr string) {
	if err != nil {
		fmt.Printf("[%s] Fatal error: %v\n", hdr, err.Error())
		os.Exit(1)
	}
}

func getFegData(baseAPI string) []byte {

	cert, err := tls.LoadX509KeyPair("/home/nitin/Project/magma/cmd/admin_operator.pem", "/home/nitin/Project/magma/cmd/admin_operator.key.pem")
	fegCheckError(err, "loadcert")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
			ClientAuth:         tls.RequireAnyClientCert,
		},
	}
	client := &http.Client{Transport: tr}

	request, err := http.NewRequest(
		http.MethodGet, //method
		baseAPI,        //url
		nil,            //body
	)

	if err != nil {
		log.Printf("Could not request a magma api. %v", err)
	}

	request.Header.Add("Accept", "application/json")

	response, err := client.Do(request)
	if err != nil {
		log.Printf("Could not make a request. %v", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body. %v", err)
	}

	return responseBytes
}
