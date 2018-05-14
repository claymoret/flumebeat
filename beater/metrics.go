package beater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type AttrMap map[string]string
type MetricsMap map[string]AttrMap
type ParsedAttrMap map[string]interface{}
type ParsedMetrics map[string]ParsedAttrMap

func (pA *AttrMap) parse() ParsedAttrMap {
	parsed := make(ParsedAttrMap)
	for k, v := range *pA {
		i, err := strconv.Atoi(v)
		if err == nil {
			parsed[k] = i
			continue
		}
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			parsed[k] = f
			continue
		}
		parsed[k] = v
	}
	return parsed
}

func (p *MetricsMap) Fetch(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(p)
	if err != nil {
		return err
	}

	return nil
}

func (p *MetricsMap) Parse() ParsedMetrics {
	pm := ParsedMetrics{}
	for component, attrmap := range *p {
		parsed := attrmap.parse()
		index := strings.Index(component, ".")
		name := component[index+1:]
		parsed["name"] = name
		t, ok := parsed["Type"]
		if ok {
			parsed["type"] = t
			delete(parsed, "Type")
		}
		pm[component] = parsed
	}
	return pm
}

func (m ParsedAttrMap) String() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("Not valid json: %v", err)
	}
	return string(bytes)
}
