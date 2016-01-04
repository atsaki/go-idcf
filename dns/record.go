package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Record struct {
	UUID      string
	Name      string
	Type      string
	TTL       int
	Content   interface{}
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *Record) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return err
	}
	m := v.(map[string]interface{})
	r.UUID = m["uuid"].(string)
	r.Name = m["name"].(string)
	r.Type = m["type"].(string)
	r.TTL = int(m["ttl"].(float64))
	r.Content = m["content"]

	if m["priority"] != nil {
		r.Priority = int(m["priority"].(float64))
	}

	if m["created_at"] != nil {
		t, err := parseTime(m["created_at"].(string))
		if err != nil {
			return err
		}
		r.CreatedAt = t
	}

	if m["updated_at"] != nil {
		t, err := parseTime(m["updated_at"].(string))
		if err != nil {
			return err
		}
		r.UpdatedAt = t
	} else {
		r.UpdatedAt = r.CreatedAt
	}
	return nil
}

func (c *Client) Records(zoneid string) (rs []*Record, err error) {
	b, err := c.Request(
		"GET", fmt.Sprintf("/api/v1/zones/%s/records", zoneid), nil)
	if err != nil {
		return
	}
	err = unmarshal(b, &rs)
	return
}

func (c *Client) Record(zoneid, recordid string) (r *Record, err error) {
	b, err := c.Request(
		"GET", fmt.Sprintf("/api/v1/zones/%s/records/%s", zoneid, recordid), nil)
	if err != nil {
		return
	}
	r = new(Record)
	err = unmarshal(b, r)
	return
}

type CreateRecordParameter map[string]interface{}

func NewCreateRecordParamter(
	zoneid, name, recordType, content string) CreateRecordParameter {
	p := make(map[string]interface{})
	p["zoneid"] = zoneid
	p["name"] = name
	p["type"] = recordType
	p["content"] = content
	return p
}

func (p CreateRecordParameter) SetTTL(ttl int) {
	p["ttl"] = ttl
}

func (p CreateRecordParameter) SetPriority(priority int) {
	p["priority"] = priority
}

func (c *Client) CreateRecord(p CreateRecordParameter) (r *Record, err error) {

	zoneid := p["zoneid"].(string)
	delete(p, "zoneid")

	if _, ok := p["ttl"]; !ok {
		z, err := c.Zone(zoneid)
		if err != nil {
			return r, err
		}
		p["ttl"] = z.DefaultTTL
	}

	b, err := c.Request(
		"POST", fmt.Sprintf("/api/v1/zones/%s/records", zoneid), p)
	if err != nil {
		return
	}
	r = new(Record)
	err = unmarshal(b, r)
	return
}

type UpdateRecordParameter map[string]interface{}

func NewUpdateRecordParamter(zoneid, recordid string) UpdateRecordParameter {
	p := make(map[string]interface{})
	p["zoneid"] = zoneid
	p["recordid"] = recordid
	return p
}

func (p UpdateRecordParameter) SetName(name string) {
	p["name"] = name
}

func (p UpdateRecordParameter) SetContent(content interface{}) {
	p["content"] = content
}

func (p UpdateRecordParameter) SetType(recordType string) {
	p["type"] = recordType
}

func (p UpdateRecordParameter) SetTTL(ttl int) {
	p["ttl"] = ttl
}

func (p UpdateRecordParameter) SetPriority(priority int) {
	p["priority"] = priority
}

func (c *Client) UpdateRecord(p UpdateRecordParameter) (r *Record, err error) {
	zoneid := p["zoneid"].(string)
	delete(p, "zoneid")

	recordid := p["recordid"].(string)
	delete(p, "recordid")

	b, err := c.Request(
		"PUT", fmt.Sprintf("/api/v1/zones/%s/records/%s", zoneid, recordid), p)
	if err != nil {
		return
	}
	r = new(Record)
	err = unmarshal(b, r)
	return
}

func (c *Client) DeleteRecord(zoneid, recordid string) (err error) {
	b, err := c.Request(
		"DELETE", fmt.Sprintf("/api/v1/zones/%s/records/%s", zoneid, recordid), nil)
	if err != nil {
		return
	}
	if string(b) != "{}" {
		return errors.New(string(b))
	}
	return
}
