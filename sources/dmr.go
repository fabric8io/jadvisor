package sources

type DmrContainer struct {
    Name        string        `json:"name,omitempty"`
    Host        string        `json:"host"`
    DmrPort int               `json:"dmrPort"`
    Stats       *StatsEntry   `json:"stats,omitempty"`
}

func (self *DmrContainer) GetStats() (*StatsEntry, error) {
    return &StatsEntry{}, nil
}
