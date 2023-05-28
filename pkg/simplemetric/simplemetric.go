// simplemetric helps with easy manipulation of simple telegraf metrics
//  without telegraf libraries dependency
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)

package simplemetric

import (
	"time"

	"github.com/influxdata/line-protocol/v2/lineprotocol"
)

type SimpleMetric struct {
	name   string
	tags   map[string]string
	fields map[string]interface{}
	tm     time.Time
}

func New(n string, t map[string]string, f map[string]interface{}) SimpleMetric {
	m := SimpleMetric{
		name:   n,
		tags:   t,
		fields: f,
		tm:     time.Now().Round(time.Second),
	}
	return m
}

func (m *SimpleMetric) Name() string {
	return m.name
}

func (m *SimpleMetric) SetTime(t time.Time) {
	m.tm = t
}

// String returns a representation of the metric in the given format string
//
//	(https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md)
//
// Currently it only supports influx line protocol
func (m *SimpleMetric) String(format string) string {
	if format == "influx" {
		return m.stringInfluxLP()
	}
	return ""
}

// stringInfluxLP returns a representation of the metric in influx line protocol
func (m *SimpleMetric) stringInfluxLP() string {
	var (
		enc lineprotocol.Encoder
		val lineprotocol.Value
		ok  bool
	)

	enc.StartLine(m.name)
	for k, v := range m.tags {
		enc.AddTag(k, v)
	}
	for k, v := range m.fields {
		if val, ok = lineprotocol.NewValue(v); ok {
			enc.AddField(k, val)
		}
	}
	enc.EndLine(m.tm)

	return string(enc.Bytes())
}
