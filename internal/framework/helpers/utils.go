package helpers

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clientpkg "github.com/zscaler/terraform-provider-zpa/v4/internal/framework/client"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/common"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/policysetcontroller"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zpa/services/segmentgroup"
)

var PolicyRulesDetachLock sync.Mutex

func BoolValue(value types.Bool, defaultValue bool) bool {
	if value.IsUnknown() || value.IsNull() {
		return defaultValue
	}
	return value.ValueBool()
}

func BoolValueDefaultFalse(value types.Bool) bool {
	return BoolValue(value, false)
}

func StringValue(value types.String) string {
	if value.IsNull() || value.IsUnknown() {
		return ""
	}
	return strings.TrimSpace(value.ValueString())
}

func StringValueOrNull(value string) types.String {
	if strings.TrimSpace(value) == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

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

func SetToIntList(d *schema.ResourceData, key string) []int {
	setObj, ok := d.GetOk(key)
	if !ok {
		return []int{}
	}
	set, ok := setObj.(*schema.Set)
	if !ok {
		return []int{}
	}
	return SetToIntSlice(set)
}

func SetToIntSlice(d *schema.Set) []int {
	if d == nil || d.Len() == 0 {
		return []int{}
	}
	list := d.List()
	ans := make([]int, 0, len(list))
	for _, v := range list {
		switch x := v.(type) {
		case string:
			// Parse string to int
			if intVal, err := strconv.Atoi(x); err == nil {
				ans = append(ans, intVal)
			}
		case int:
			ans = append(ans, x)
		}
	}
	return ans
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

func convertPortsToListString(portRangeLst []common.NetworkPorts) []string {
	portRanges := make([]string, len(portRangeLst)*2)
	for i := range portRangeLst {
		portRanges[2*i] = portRangeLst[i].From
		portRanges[2*i+1] = portRangeLst[i].To
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

func expandAppSegmentNetwokPorts(d *schema.ResourceData, key string) []string {
	var ports []string
	if portsInterface, ok := d.GetOk(key); ok {
		portList, ok := portsInterface.([]interface{})
		if !ok {
			log.Printf("[ERROR] conversion failed, destUdpPortsInterface")
			return []string{}
		}
		ports = make([]string, len(portList)*2)
		for i, val := range portList {
			portItem := val.(map[string]interface{})
			ports[2*i] = portItem["from"].(string)
			ports[2*i+1] = portItem["to"].(string)
		}
	}
	return ports
}

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

// DetachAppFromAllPolicyRules removes the supplied application segment from every policy rule
// that references it by ID.
func DetachAppFromAllPolicyRules(ctx context.Context, service *zscaler.Service, appID string) {
	PolicyRulesDetachLock.Lock()
	defer PolicyRulesDetachLock.Unlock()

	types := []string{"ACCESS_POLICY", "TIMEOUT_POLICY", "SIEM_POLICY", "CLIENT_FORWARDING_POLICY", "INSPECTION_POLICY"}
	var rules []policysetcontroller.PolicyRule

	for _, policyType := range types {
		policySet, _, err := policysetcontroller.GetByPolicyType(ctx, service, policyType)
		if err != nil {
			continue
		}

		policyRules, _, err := policysetcontroller.GetAllByType(ctx, service, policyType)
		if err != nil {
			continue
		}

		for _, rule := range policyRules {
			rule.PolicySetID = policySet.ID
			rules = append(rules, rule)
		}
	}

	for _, rr := range rules {
		rule := rr
		changed := false
		for i, condition := range rule.Conditions {
			operands := make([]policysetcontroller.Operands, 0, len(condition.Operands))
			for _, operand := range condition.Operands {
				if operand.ObjectType == "APP" && strings.EqualFold(operand.LHS, "id") && operand.RHS == appID {
					changed = true
					continue
				}
				operands = append(operands, operand)
			}
			rule.Conditions[i].Operands = operands
		}

		if len(rule.Conditions) == 0 {
			rule.Conditions = []policysetcontroller.Conditions{}
		}

		if changed {
			if _, err := policysetcontroller.UpdateRule(ctx, service, rule.PolicySetID, rule.ID, &rule); err != nil {
				log.Printf("[WARN] Failed to detach application segment %s from policy rule %s: %v", appID, rule.ID, err)
			}
		}
	}
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
// func GetBool(v interface{}) bool {
// 	if b, ok := v.(bool); ok {
// 		return b
// 	}
// 	return false
// }

func GetBool(input interface{}) bool {
	if input == nil {
		return false
	}
	return input.(bool)
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

// RFC1123FromEpoch converts an epoch string (seconds) to an RFC1123 formatted string.
func RFC1123FromEpoch(epochStr string) (string, error) {
	return epochToRFC1123(epochStr, false)
}

// RFC1123ZFromEpoch converts an epoch string (seconds) to an RFC1123Z formatted string.
func RFC1123ZFromEpoch(epochStr string) (string, error) {
	return epochToRFC1123(epochStr, true)
}

// ConvertRFC1123ToEpoch converts a RFC1123 or RFC1123Z formatted timestamp into epoch seconds.
func ConvertRFC1123ToEpoch(dateStr string) (int64, error) {
	layouts := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t.Unix(), nil
		}
	}

	return 0, fmt.Errorf("unable to parse date: %s", dateStr)
}

// ValidatePRATimeRange validates start and end times for PRA approvals.
func ValidatePRATimeRange(startTimeStr, endTimeStr string) error {
	startTimeEpoch, err := ConvertRFC1123ToEpoch(startTimeStr)
	if err != nil {
		return fmt.Errorf("start time conversion error: %s", err)
	}

	currentTimeEpoch := time.Now().Unix()
	if startTimeEpoch < currentTimeEpoch-int64(3600) {
		return fmt.Errorf("the approval start time cannot be more than 1 hour in the past")
	}

	endTimeEpoch, err := ConvertRFC1123ToEpoch(endTimeStr)
	if err != nil {
		return fmt.Errorf("end time conversion error: %s", err)
	}

	const oneYear = int64(31536000) // 365 * 24 * 60 * 60
	if endTimeEpoch > startTimeEpoch+oneYear {
		return fmt.Errorf("the start time should be less than the future end time with a max range of 1 year")
	}

	return nil
}

// Validate24HourTimeFormat ensures the given string uses HH:MM 24-hour format.
func Validate24HourTimeFormat(value string) error {
	if _, err := time.Parse("15:04", value); err != nil {
		return fmt.Errorf("%q is not a valid time format (expected HH:MM)", value)
	}
	return nil
}

// ValidateTimeZone verifies that the provided string is a valid IANA timezone name.
func ValidateTimeZone(value string) error {
	if value == "" {
		return nil
	}
	if _, err := time.LoadLocation(value); err != nil {
		return fmt.Errorf("%q is not a valid timezone", value)
	}
	return nil
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

// ParseHumanReadableTimeout exposes the timeout parser for other packages.
func ParseHumanReadableTimeout(input string) (int, error) {
	return parseHumanReadableTimeout(input)
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

// SecondsToHumanReadable exposes the internal conversion helper for other packages.
func SecondsToHumanReadable(seconds string) string {
	return secondsToHumanReadable(seconds)
}

func pluralize(count int, singular, plural string) string {
	if count == 1 {
		return singular
	}
	return plural
}

func DetachSegmentGroup(ctx context.Context, zClient *clientpkg.Client, segmentID, segmentGroupID string) error {
	log.Printf("[INFO] Detaching application segment  %s from segment group: %s\n", segmentID, segmentGroupID)
	service := zClient.Service

	segGroup, _, err := segmentgroup.Get(ctx, service, segmentGroupID)
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
	_, err = segmentgroup.Update(ctx, service, segmentGroupID, segGroup)
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

// GenerateShortID creates a short, deterministic hash identifier for a given string.
func GenerateShortID(input string) string {
	hash := md5.Sum([]byte(input))
	return fmt.Sprintf("%x", hash)[:8] // Use first 8 characters of MD5 hash
}
