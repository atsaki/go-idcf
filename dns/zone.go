package dns

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Zone struct {
	UUID          string
	Name          string
	Description   string
	DefaultTTL    int
	Authenticated bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (z *Zone) UnmarshalJSON(data []byte) error {
	var v interface{}
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return err
	}
	m := v.(map[string]interface{})
	z.UUID = m["uuid"].(string)
	z.Name = m["name"].(string)
	z.Description = m["description"].(string)
	z.DefaultTTL = int(m["default_ttl"].(float64))
	z.Authenticated = m["authenticated"].(bool)

	if m["created_at"] != nil {
		t, err := parseTime(m["created_at"].(string))
		if err != nil {
			return err
		}
		z.CreatedAt = t
	}

	if m["updated_at"] != nil {
		t, err := parseTime(m["updated_at"].(string))
		if err != nil {
			return err
		}
		z.UpdatedAt = t
	} else {
		z.UpdatedAt = z.CreatedAt
	}
	return nil
}

func (c *Client) Zones() (zs []*Zone, err error) {
	b, err := c.Request("GET", "/api/v1/zones", nil)
	if err != nil {
		return
	}
	err = unmarshal(b, &zs)
	return
}

func (c *Client) Zone(zoneid string) (z *Zone, err error) {
	b, err := c.Request("GET", fmt.Sprintf("/api/v1/zones/%s", zoneid), nil)
	if err != nil {
		return
	}
	z = new(Zone)
	err = unmarshal(b, z)
	return
}

type CreateZoneParameter map[string]interface{}

func NewCreateZoneParamter(name, email string) CreateZoneParameter {
	p := make(map[string]interface{})
	p["name"] = name
	p["email"] = email
	return p
}

func (p CreateZoneParameter) SetDescription(description string) {
	p["description"] = description
}

func (p CreateZoneParameter) SetDefaultTTL(ttl int) {
	p["default_ttl"] = ttl
}

func (c *Client) CreateZone(p CreateZoneParameter) (z *Zone, err error) {
	b, err := c.Request("POST", "/api/v1/zones", p)
	if err != nil {
		return
	}
	z = new(Zone)
	err = unmarshal(b, z)
	return
}

type UpdateZoneParameter map[string]interface{}

func NewUpdateZoneParamter(zoneid string) UpdateZoneParameter {
	p := make(map[string]interface{})
	p["zoneid"] = zoneid
	return p
}

func (p UpdateZoneParameter) SetDescription(description string) {
	p["description"] = description
}

func (p UpdateZoneParameter) SetDefaultTTL(ttl int) {
	p["default_ttl"] = ttl
}

func (c *Client) UpdateZone(p UpdateZoneParameter) (z *Zone, err error) {
	zoneid := p["zoneid"].(string)
	delete(p, "zoneid")

	b, err := c.Request("PUT", fmt.Sprintf("/api/v1/zones/%s", zoneid), p)
	if err != nil {
		return
	}
	z = new(Zone)
	err = unmarshal(b, &z)
	return
}

func (c *Client) DeleteZone(zoneid string) (err error) {
	b, err := c.Request("DELETE", fmt.Sprintf("/api/v1/zones/%s", zoneid), nil)
	if err != nil {
		return
	}
	if string(b) != "{}" {
		return errors.New(string(b))
	}
	return
}

func (c *Client) VerifyZone(zoneid string) (err error) {
	b, err := c.Request(
		"POST", fmt.Sprintf("/api/v1/zones/%s/verify", zoneid), nil)
	if err != nil {
		return
	}
	if string(b) != "{}" {
		return errors.New(string(b))
	}
	return
}
