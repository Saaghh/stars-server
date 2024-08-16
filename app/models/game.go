package models

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID        int
	Owner     User
	WorldTime time.Time
	Name      string
	IsArchive bool
	CreatedAt time.Time
	Systems   []System
	Bodies    []StellarBody
}

type DBGame struct {
	ID        int
	OwnerID   int
	WorldTime time.Time
	Name      string
	IsArchive bool
	CreatedAt time.Time
}

type System struct {
	ID     int
	GameID int
	Name   string
}

type StellarBodyType struct {
	ID   int
	Name string
}

type StellarBody struct {
	ID            uuid.UUID
	SystemID      int
	Name          string
	TypeID        int
	Mass          float64
	Diameter      float64
	ParentBodyID  *uuid.UUID
	OrbitalRadius *float64
	Angle         *float64
	AngleSpeed    *float64
	LinearSpeed   *float64
	CoordinateX   *float64
	CoordinateY   *float64
	Stockpiles    []Stockpile
}

type StellarBodyFilter struct {
	Systems         []int
	Game            int
	NameSearch      string
	StellarBodyType int
}

type ResourceType struct {
	ID      int
	Density float64
	Name    string
}

type Stockpile struct {
	ID           int
	Quantity     float64
	ResourceType ResourceType
}
