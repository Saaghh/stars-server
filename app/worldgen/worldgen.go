package worldgen

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/exp/rand"
	"stars-server/app/models"
)

const (
	DefaultOwnerID = 1

	PlanetsPerSystem  = 7
	SystemsToGenerate = 3

	DictionaryPlanets = "planets"
	DictionaryStars   = "stars"

	StarMaxMass     = 100000
	StarMaxDiameter = 1000

	PlanetMaxMass          = 1000
	PlanetMaxDiameter      = 100
	PlanetMaxOrbitalRadius = 100
	PlanetMaxAngleSpeedMod = 5

	StellarBodyTypeStarID   = 1
	StellarBodyTypePlanetID = 1
)

type store interface {
	WithTx(ctx context.Context, f func(context.Context) error) error

	TxCreateGame(ctx context.Context, game models.Game) (models.Game, error)
	TxCreateSystem(ctx context.Context, starSystem models.System) (models.System, error)
	TxCreateBody(ctx context.Context, body models.StellarBody) (models.StellarBody, error)
}

type namesGenerator interface {
	GetRandomName(dictionary string) string
}

type WorldGenerator struct {
	db    store
	names namesGenerator
}

func init() {
	rand.Seed(uint64(time.Now().UnixNano()))
}

func New(db store, names namesGenerator) *WorldGenerator {
	return &WorldGenerator{
		db:    db,
		names: names,
	}
}

func (w *WorldGenerator) AutoGenerateGame(ctx context.Context) (models.Game, error) {
	var err error
	var game = models.Game{
		Owner:     models.User{ID: DefaultOwnerID},
		WorldTime: time.Now(),
		Name:      time.Now().Format("2006/01/02 15:04:05") + " - Game",
	}

	if err = w.db.WithTx(ctx, func(ctx context.Context) error {
		if game, err = w.txGenerateGame(ctx, game); err != nil {
			return fmt.Errorf("w.db.TxCreateGame: %w", err)
		}

		return nil
	}); err != nil {
		return models.Game{}, fmt.Errorf("w.db.WithTx: %w", err)
	}

	return game, nil
}

func (w *WorldGenerator) txGenerateGame(ctx context.Context, game models.Game) (models.Game, error) {
	var err error

	// save game
	if game, err = w.db.TxCreateGame(ctx, game); err != nil {
		return models.Game{}, fmt.Errorf("w.db.TxCreateGame: %w", err)
	}

	// generate systems
	for range SystemsToGenerate {
		systemName := w.names.GetRandomName(DictionaryStars)

		var starSystem = models.System{
			GameID: game.ID,
			Name:   systemName,
		}

		starSystem, err = w.txGenerateSystem(ctx, starSystem)
		if err != nil {
			return models.Game{}, fmt.Errorf("w.txGenerateSystem: %w", err)
		}

		game.Systems = append(game.Systems, starSystem)
	}

	return game, nil
}

func (w *WorldGenerator) txGenerateSystem(ctx context.Context, starSystem models.System) (models.System, error) {
	var err error

	// Creating Star systems
	starSystem, err = w.db.TxCreateSystem(ctx, starSystem)
	if err != nil {
		return models.System{}, fmt.Errorf("w.db.TxCreateSystem: %w", err)
	}

	// Creating Star
	star := w.getRandomStar()
	star.SystemID = starSystem.ID
	star.Name = starSystem.Name + " Star"

	if star, err = w.db.TxCreateBody(ctx, star); err != nil {
		return models.System{}, fmt.Errorf("w.db.TxCreateStar: %w", err)
	}

	starSystem.Bodies = append(starSystem.Bodies, star)

	//Creating planets
	var planetsCount = 0

	for range PlanetsPerSystem {
		var planet = w.getRandomPlanet()
		planet.SystemID = starSystem.ID
		planetsCount++
		planet.Name = starSystem.Name + " " + strconv.Itoa(planetsCount)
		planet.ParentBodyID = &star.ID

		planet, err = w.db.TxCreateBody(ctx, planet)
		if err != nil {
			return models.System{}, fmt.Errorf("w.db.TxCreatePlanet: %w", err)
		}

		starSystem.Bodies = append(starSystem.Bodies)

		// TODO: add moons generation
	}

	return starSystem, nil
}

func (w *WorldGenerator) getRandomStar() models.StellarBody {
	return models.StellarBody{
		TypeID:   StellarBodyTypeStarID,
		Mass:     float64(rand.Intn(StarMaxMass)),
		Diameter: float64(rand.Intn(StarMaxDiameter)),
	}
}

func (w *WorldGenerator) getRandomPlanet() models.StellarBody {
	var orbitalRadius = float64(rand.Intn(PlanetMaxOrbitalRadius)) + rand.Float64()
	var angle = float64(rand.Intn(359)) + rand.Float64()
	var angleSpeed = rand.Float64() * PlanetMaxAngleSpeedMod

	return models.StellarBody{
		TypeID:        StellarBodyTypePlanetID,
		Mass:          float64(rand.Intn(PlanetMaxMass)),
		Diameter:      float64(rand.Intn(PlanetMaxDiameter)),
		OrbitalRadius: &orbitalRadius,
		Angle:         &angle,
		AngleSpeed:    &angleSpeed,
	}

	// TODO: add stockpiles generation
}
