package names

import (
	"embed"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"

	"golang.org/x/exp/rand"
)

const DefaultName = "placeholder"

//go:embed names.json
var embeddedNames embed.FS

type Names struct {
	namesMap map[string][]string
}

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func New() (*Names, error) {
	names := Names{}

	if err := names.readNames(); err != nil {
		return nil, fmt.Errorf("names.readNames: %w", err)
	}

	return &names, nil
}

func (n *Names) readNames() error {
	data, err := embeddedNames.ReadFile("names.json")
	if err != nil {
		return fmt.Errorf("embeddedNames.ReadFile: %w", err)
	}

	var namesMap map[string][]string

	if err = json.Unmarshal(data, &namesMap); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	n.namesMap = namesMap

	return nil
}

func (n *Names) GetRandomName(dictionary string) string {
	namesList, ok := n.namesMap[dictionary]
	if !ok {
		zap.S().Warn("Names/GetRandomName/n.namesMap[dictionary]: !ok")
		return DefaultName
	}

	return namesList[rand.Intn(len(namesList))]
}
