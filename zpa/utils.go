package zpa

import (
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/segmentgroup"
)

func ValidateLatitude(val interface{}, _ string) (warns []string, errs []error) {
	v, _ := strconv.ParseFloat(val.(string), 64)
	if v < -90 || v > 90 {
		errs = append(errs, fmt.Errorf("latitude must be between -90 and 90"))
	}
	return
}

func ValidateLongitude(val interface{}, _ string) (warns []string, errs []error) {
	v, _ := strconv.ParseFloat(val.(string), 64)
	if v < -180 || v > 180 {
		errs = append(errs, fmt.Errorf("longitude must be between -180 and 180"))
	}
	return
}

func DiffSuppressFuncCoordinate(_, old, new string, _ *schema.ResourceData) bool {
	o, _ := strconv.ParseFloat(old, 64)
	n, _ := strconv.ParseFloat(new, 64)
	return math.Round(o*1000000)/1000000 == math.Round(n*1000000)/1000000
}

func SetToStringSlice(d *schema.Set) []string {
	list := d.List()
	return ListToStringSlice(list)
}

func SetToStringList(d *schema.ResourceData, key string) []string {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []string{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []string{}
	}
	return SetToStringSlice(set)
}

func ListToStringSlice(v []interface{}) []string {
	if len(v) == 0 {
		return []string{}
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case nil:
			ans[i] = ""
		case string:
			ans[i] = x
		}
	}

	return ans
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func MergeSchema(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	final := make(map[string]*schema.Schema)
	for _, s := range schemas {
		for k, v := range s {
			final[k] = v
		}
	}
	return final
}

/*
func convertPortsToListString(portRangeLst []common.NetworkPorts) []string {
	portRanges := make([]string, len(portRangeLst)*2)
	for i := range portRangeLst {
		portRanges[2*i] = portRangeLst[i].From
		portRanges[2*i+1] = portRangeLst[i].To
	}
	return portRanges
}
*/

/*
// Function to duplicate each port range value in the list
func duplicatePortRanges(portRangeLst []interface{}) []string {
	portRanges := make([]string, 0, len(portRangeLst)*2)
	for _, v := range portRangeLst {
		port := v.(string)
		portRanges = append(portRanges, port, port)
	}
	return portRanges
}

func convertToPortRange(portRangeLst []interface{}) []string {
	portRanges := make([]string, len(portRangeLst))
	for i := range portRanges {
		portRanges[i] = portRangeLst[i].(string)
	}
	return portRanges
}
*/

// func expandAppSegmentNetworkPorts(d *schema.ResourceData, key string) []common.NetworkPorts {
// 	var ports []common.NetworkPorts
// 	if portInterface, ok := d.GetOk(key); ok {
// 		portSet := portInterface.(*schema.Set)
// 		ports = make([]common.NetworkPorts, len(portSet.List()))
// 		for i, val := range portSet.List() {
// 			portItem := val.(map[string]interface{})
// 			ports[i] = common.NetworkPorts{
// 				From: portItem["from"].(string),
// 				To:   portItem["to"].(string),
// 			}
// 		}
// 	}
// 	return ports
// }

func expandAppSegmentNetworkPorts(d *schema.ResourceData, key string) []common.NetworkPorts {
	var ports []common.NetworkPorts
	if portInterface, ok := d.GetOk(key); ok {
		portList := portInterface.([]interface{}) // TypeList ensures the order is preserved
		ports = make([]common.NetworkPorts, len(portList))
		for i, val := range portList {
			portItem := val.(map[string]interface{})
			ports[i] = common.NetworkPorts{
				From: portItem["from"].(string),
				To:   portItem["to"].(string),
			}
		}
	}
	return ports
}

func sliceHasCommon(s1, s2 []string) (bool, string) {
	for _, i1 := range s1 {
		for _, i2 := range s2 {
			if i1 == i2 {
				return true, i1
			}
		}
	}
	return false, ""
}

/*
func isSameSlice(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
*/

func expandStringInSlice(d *schema.ResourceData, key string) []string {
	applicationSegments := d.Get(key).([]interface{})
	applicationSegmentList := make([]string, len(applicationSegments))
	for i, applicationSegment := range applicationSegments {
		applicationSegmentList[i] = applicationSegment.(string)
	}

	return applicationSegmentList
}

func validateAppPorts(selectConnectorCloseToApp bool, udpAppPortRange []common.NetworkPorts, udpPortRanges []string) error {
	if selectConnectorCloseToApp {
		if len(udpAppPortRange) > 0 || len(udpPortRanges) > 0 {
			return errors.New("the protocol configuration for the application is invalid. App Connector Closest to App supports only TCP applications")
		}
	}
	return nil
}

// createValidResourceName converts the given name to a valid Terraform resource name
func createValidResourceName(name string) string {
	return strings.ReplaceAll(name, " ", "_")
}

func GetString(v interface{}) string {
	if v == nil {
		return ""
	}
	str, ok := v.(string)
	if ok {
		return str
	}
	return fmt.Sprintf("%v", v)
}

// Helper to safely extract bool values from map
func GetBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

// Converts an epoch time (in seconds, represented as a string) to a human-readable format.
func epochToRFC1123(epochStr string, useRFC1123Z bool) (string, error) {
	epoch, err := strconv.ParseInt(epochStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("failed to parse epoch time: %s", err)
	}
	t := time.Unix(epoch, 0) // Convert epoch to *time.Time, assuming epoch is in seconds.
	if useRFC1123Z {
		return t.Format(time.RFC1123Z), nil // Returns the time formatted using RFC1123Z layout.
	}
	return t.Format(time.RFC1123), nil // Returns the time formatted using RFC1123 layout.
}

// #######################################################################################
// ######################Conversion function for Timeout Policy Rule######################
// #######################################################################################
func parseHumanReadableTimeout(input string) (int, error) {
	var multiplier int
	var value int
	var unit string

	// Handle special case for "Never" or "never"
	if strings.ToLower(input) == "never" {
		return -1, nil
	}

	_, err := fmt.Sscanf(input, "%d %s", &value, &unit)
	if err != nil {
		return 0, fmt.Errorf("error parsing timeout value: %v", err)
	}

	unit = strings.ToLower(unit)
	switch unit {
	case "minute", "minutes":
		multiplier = 60
	case "hour", "hours":
		multiplier = 3600
	case "day", "days":
		multiplier = 86400
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}

	return value * multiplier, nil
}

// Convert seconds into a human-readable format. This function assumes `seconds` is a string that can be parsed into an int.
func secondsToHumanReadable(seconds string) string {
	sec, err := strconv.Atoi(seconds)
	if err != nil {
		log.Printf("[ERROR] Failed to parse seconds: %v", err)
		return ""
	}

	days := sec / 86400
	hours := (sec % 86400) / 3600
	minutes := (sec % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%d %s", days, pluralize(days, "Day", "Days"))
	} else if hours > 0 {
		return fmt.Sprintf("%d %s", hours, pluralize(hours, "Hour", "Hours"))
	} else if minutes > 0 {
		return fmt.Sprintf("%d %s", minutes, pluralize(minutes, "Minute", "Minutes"))
	}
	return fmt.Sprintf("%d %s", sec, pluralize(sec, "Second", "Seconds"))
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func detachSegmentGroup(client *Client, segmentID, segmentGroupID string) error {
	log.Printf("[INFO] Detaching application segment  %s from segment group: %s\n", segmentID, segmentGroupID)
	service := client.SegmentGroup

	segGroup, _, err := segmentgroup.Get(service, segmentGroupID)
	if err != nil {
		log.Printf("[error] Error while getting segment group id: %s", segmentGroupID)
		return err
	}
	adaptedApplications := []segmentgroup.Application{}
	for _, app := range segGroup.Applications {
		if app.ID != segmentID {
			adaptedApplications = append(adaptedApplications, app)
		}
	}
	segGroup.Applications = adaptedApplications
	_, err = segmentgroup.Update(service, segmentGroupID, segGroup)
	return err
}

var sensitiveFieldNames = []string{"password", "passphrase", "private_key", "debugMode.filePassword"}

func sanitizeFields(input interface{}) {
	val := reflect.ValueOf(input).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		if isSensitiveField(field.Name) {
			if val.Field(i).CanSet() && val.Field(i).Kind() == reflect.String {
				val.Field(i).SetString("***REDACTED***")
			}
		}
	}
}

func isSensitiveField(fieldName string) bool {
	fieldName = strings.ToLower(fieldName)
	for _, sensitiveField := range sensitiveFieldNames {
		if strings.Contains(fieldName, sensitiveField) {
			return true
		}
	}
	return false
}
