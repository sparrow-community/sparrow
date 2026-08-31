package deploy

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/sparrow-community/sparrow/bpmn/element"
)

const (
	DefaultIncidentThreshold = 3
	sparrowNS                = "http://sparrow.example/bpmn"
)

// IncidentThreshold returns the configured fail count before opening an incident.
func (d *Deployment) IncidentThreshold(elementID string) int {
	if d == nil || d.incidentThresholds == nil {
		return DefaultIncidentThreshold
	}
	if n, ok := d.incidentThresholds[elementID]; ok {
		return n
	}
	return DefaultIncidentThreshold
}

func indexIncidentThreshold(d *Deployment, elementID string, ext element.ExtensionElements) error {
	n, ok, err := parseIncidentThreshold(ext)
	if err != nil {
		return fmt.Errorf("service task %q: %w", elementID, err)
	}
	if !ok {
		return nil
	}
	if d.incidentThresholds == nil {
		d.incidentThresholds = make(map[string]int)
	}
	d.incidentThresholds[elementID] = n
	return nil
}

func parseIncidentThreshold(ext element.ExtensionElements) (int, bool, error) {
	for _, tok := range ext.Any {
		se, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		if se.Name.Local != "failedJobIncidentThreshold" {
			continue
		}
		for _, attr := range se.Attr {
			if attr.Name.Local == "failedJobIncidentThreshold" && attr.Name.Space == sparrowNS {
				return parseThresholdAttr(attr.Value)
			}
		}
		if len(se.Attr) > 0 {
			for _, attr := range se.Attr {
				if attr.Name.Local == "value" || attr.Name.Local == "failedJobIncidentThreshold" {
					return parseThresholdAttr(attr.Value)
				}
			}
		}
	}
	return 0, false, nil
}

func parseThresholdAttr(raw string) (int, bool, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false, fmt.Errorf("empty failedJobIncidentThreshold")
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, false, fmt.Errorf("invalid failedJobIncidentThreshold %q", raw)
	}
	return n, true, nil
}
