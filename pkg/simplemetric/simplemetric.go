// simplemetric helps with easy manupulation of simple telegraf metrics
//  without telegraf libraries dependency
//
// Author: Tesifonte Belda
// License: The MIT License (MIT)

package simplemetric

import (
	"strconv"
	"strings"
	"time"
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
//  (https://github.com/influxdata/telegraf/blob/master/docs/DATA_FORMATS_INPUT.md)
// Currently it only supports influx line protocol
func (m *SimpleMetric) String(format string) string {
	switch format {
	case "influx":
		return m.stringInfluxLP()
	}
	return ""
}

// stringInfluxLP returns a representation of the metric in influx line protocol
func (m *SimpleMetric) stringInfluxLP() string {
	const (
		comma = ","
		space = " "
		equal = "="
	)
	var firstf bool = true

	buf := new(strings.Builder)
	buf.WriteString(m.name)
	for k, v := range m.tags {
		buf.WriteString(comma)
		buf.WriteString(k)
		buf.WriteString(equal)
		v = strings.Replace(v, " ", "\\ ", -1)
		buf.WriteString(v)
	}
	buf.WriteString(space)

	for k, v := range m.fields {
		v := convertField(v)
		if len(v) > 0 {
			if !firstf {
				buf.WriteString(comma)
			} else {
				firstf = false
			}
			buf.WriteString(k)
			buf.WriteString(equal)
			buf.WriteString(v)
		}
	}
	buf.WriteString(space)

	buf.WriteString(strconv.FormatInt(m.tm.UnixNano(), 10))
	return buf.String()
}

func convertField(v interface{}) string {
	const i = "i"

	buf := new(strings.Builder)
	switch v := v.(type) {
	case bool:
		buf.WriteString(strconv.FormatBool(v))
	case float64:
		buf.WriteString(strconv.FormatFloat(v, 'f', -1, 32))
	case float32:
		buf.WriteString(strconv.FormatFloat(float64(v), 'f', -1, 32))
	case int64:
		buf.WriteString(strconv.FormatInt(v, 10))
		buf.WriteString(i)
	case int32:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
		buf.WriteString(i)
	case int16:
		buf.WriteString(strconv.FormatInt(int64(v), 10))
		buf.WriteString(i)
	case int:
		buf.WriteString(strconv.Itoa(v))
		buf.WriteString(i)
	case uint64:
		buf.WriteString(strconv.FormatUint(v, 10))
		buf.WriteString(i)
	case uint16:
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
		buf.WriteString(i)
	case uint:
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
		buf.WriteString(i)
	case string:
		if len(v) > 0 {
			v = strings.Replace(v, " ", "\\ ", -1)
			buf.WriteString(strconv.Quote(v))
		}
	}
	return buf.String()
}
