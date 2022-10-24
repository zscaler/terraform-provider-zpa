package main

import (
	"flag"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/zscaler/zscaler-sdk-go/zpa"
	"github.com/zscaler/zscaler-sdk-go/zpa/services/appconnectorcontroller"
)

func main() {
	var log = logrus.New()

	//rootCmd.PersistentFlags().StringVar(&ids, "ids", "", "Which application connectors you wish to delete")
	ids := flag.String("ids", "", "-ids \"122,2233,3322\"")
	flag.Parse()
	if ids == nil {
		log.Error("no application connector is specified")
		return
	}
	appIDs := strings.Split(*ids, ",")
	if *ids == "" || len(appIDs) == 0 {
		log.Error("no application connector is specified")
		return
	}
	zpa_client_id := os.Getenv("ZPA_CLIENT_ID")
	zpa_client_secret := os.Getenv("ZPA_CLIENT_SECRET")
	zpa_customer_id := os.Getenv("ZPA_CUSTOMER_ID")
	zpa_cloud := os.Getenv("ZPA_CLOUD")
	config, err := zpa.NewConfig(zpa_client_id, zpa_client_secret, zpa_customer_id, zpa_cloud, "app-connector-bulk-cli")
	if err != nil {
		log.Errorf("creating config failed: %v\n", err)
		return
	}
	zpaClient := zpa.NewClient(config)
	appConnService := appconnectorcontroller.New(zpaClient)
	_, err = appConnService.BulkDelete(appIDs)
	if err != nil {
		log.Errorf("deleting app connectors failed: %v\n", err)
		return
	}
	log.Infof("deleted app connectors: %v\n", appIDs)
}
