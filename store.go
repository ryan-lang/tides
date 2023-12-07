package tides

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type (
	StationStore struct {
		Stations map[string]*Station
	}
)

func (s *StationStore) GetStation(id string) *Station {
	return s.Stations[id]
}

func LoadStationData(dataDir string) *StationStore {
	store := &StationStore{
		Stations: make(map[string]*Station),
	}

	filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var station Station
			if err := json.Unmarshal(data, &station); err != nil {
				return err
			}

			store.Stations[station.ID] = &station
		}

		return nil
	})

	return store
}
