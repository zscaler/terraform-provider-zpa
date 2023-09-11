package zpa

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"

	"github.com/fabiotavarespr/iso3166"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zscaler/zscaler-sdk-go/v2/zpa/services/common"
)

func ValidateLatitude(val interface{}, key string) (warns []string, errs []error) {
	v, _ := strconv.ParseFloat(val.(string), 64)
	if v < -90 || v > 90 {
		errs = append(errs, fmt.Errorf("latitude must be between -90 and 90"))
	}
	return
}

func ValidateLongitude(val interface{}, key string) (warns []string, errs []error) {
	v, _ := strconv.ParseFloat(val.(string), 64)
	if v < -180 || v > 180 {
		errs = append(errs, fmt.Errorf("longitude must be between -180 and 180"))
	}
	return
}

func DiffSuppressFuncCoordinate(k, old, new string, d *schema.ResourceData) bool {
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

func expandList(portRangeLst []interface{}) []string {
	portRanges := make([]string, len(portRangeLst))
	for i, port := range portRangeLst {
		portRanges[i] = port.(string)
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
		portSet, ok := portsInterface.(*schema.Set)
		if !ok {
			log.Printf("[ERROR] conversion failed, destUdpPortsInterface")
			return []string{}
		}
		ports = make([]string, len(portSet.List())*2)
		for i, val := range portSet.List() {
			portItem := val.(map[string]interface{})
			ports[2*i] = portItem["from"].(string)
			ports[2*i+1] = portItem["to"].(string)
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

// This function checks for ISO3166 Alpha2 Country codes
// It's used within the access policy RHS validation
func isValidAlpha2(code string) bool {
	return iso3166.ExistsIso3166ByAlpha2Code(code)
}
