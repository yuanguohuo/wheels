package api

import "encoding/json"

type RaftMember struct {
	Id    uint64            `json:"id"`
	Ip    string            `json:"ip"`
	Port  int32             `json:"port"`
	Attrs map[string]string `json:"attrs"`
}

func (m *RaftMember) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *RaftMember) Unmarshal(data []byte) error {
	return json.Unmarshal(data, m)
}
