package models

import (
	"github.com/google/uuid"
	"time"
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
}

type StellarBodyFilter struct {
	Systems         []int
	Game            int
	NameSearch      string
	StellarBodyType int
}
