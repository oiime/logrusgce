package logrusgce

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type gceSeverity string

const (
	gceSeverityDEFAULT  gceSeverity = "DEFAULT"
	gceSeverityDEBUG    gceSeverity = "DEBUG"
	gceSeverityINFO     gceSeverity = "INFO"
	gceSeverityWARNING  gceSeverity = "WARNING"
	gceSeverityERROR    gceSeverity = "ERROR"
	gceSeverityCRITICAL gceSeverity = "CRITICAL"
	gceSeverityALERT    gceSeverity = "ALERT"
)

var (
	severityMapping = map[logrus.Level]gceSeverity{
		logrus.TraceLevel: gceSeverityDEFAULT,
		logrus.DebugLevel: gceSeverityDEBUG,
		logrus.InfoLevel:  gceSeverityINFO,
		logrus.WarnLevel:  gceSeverityWARNING,
		logrus.ErrorLevel: gceSeverityERROR,
		logrus.FatalLevel: gceSeverityCRITICAL,
		logrus.PanicLevel: gceSeverityALERT,
	}
)

type GCELogEntry struct {
	Severity gceSeverity            `json:"severity"`
	Message  map[string]interface{} `json:"message"`
	Time     string                 `json:"time"`
}

// PostProcessFunc add additional fields
type PostProcessFunc func(*logrus.Entry, logrus.Fields)

// GCEFormatter google cloud logging format
type GCEFormatter struct {
	postprocess []PostProcessFunc
}

// NewGCEFormatter retrurn new formatter
func NewGCEFormatter() *GCEFormatter {
	return &GCEFormatter{}
}

// PostProcess add a post process
func (f *GCEFormatter) WithPostProcess(fn PostProcessFunc) *GCEFormatter {
	nf := *f
	if f.postprocess != nil {
		nf.postprocess = append(f.postprocess, fn)
	} else {
		nf.postprocess = []PostProcessFunc{fn}
	}
	return &nf
}

// Format format logrus entry
func (f *GCEFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	serialized, err := json.Marshal(f.process(entry))
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields into JSON: %v", err)
	}
	return append(serialized, '\n'), nil
}

func (f *GCEFormatter) process(logrusEntry *logrus.Entry) logrus.Fields {
	gceEntry := make(logrus.Fields, 3)

	gceEntry["timestamp"] = logrusEntry.Time.Format(time.RFC3339Nano)
	gceEntry["severity"] = severityMapping[logrusEntry.Level]
	gceEntry["message"] = logrusEntry.Message

	if f.postprocess != nil {
		for _, p := range f.postprocess {
			p(logrusEntry, gceEntry)
		}
	}
	return gceEntry
}
