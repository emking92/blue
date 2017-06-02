package entities

import (
	"strconv"
)

type IEntity interface {
	Init()
}

type Entity struct {
	Name        string      `json:"name"`
	EnityNumber int         `json:"entity_number"`
	Position    Vector2     `json:"position"`
	Direction   direction   `json:"direction"`
	Connections Connections `json:"connections"`
}

type Vector2 struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Circuit struct {
	CircuitId int `json:"circuit_id,omitempty"`
	EntityId  int `json:"entity_id,omitempty"`
}

type Connections map[string]map[wire][]Circuit

type wire string

const (
	WireGreen = wire("green")
	WireRed   = wire("red")
)

type direction int

const (
	DirectionUp    = direction(0)
	DirectionRight = direction(2)
	DirectionDown  = direction(4)
	DirectionLeft  = direction(6)
)

var (
	idTracker int
)

func (entity *Entity) initEntity() {
	entity.Connections = make(Connections)
}

func nextID() int {
	idTracker++
	return idTracker
}

func (point *Vector2) Set(x, y float32) {
	point.X = x
	point.Y = y
}

func ConnectEntities(first *Entity, firstConnector int, second *Entity, secondConnector int, w wire) {
	first.addConnection(second.EnityNumber, strconv.Itoa(firstConnector), w, secondConnector)
	second.addConnection(first.EnityNumber, strconv.Itoa(secondConnector), w, firstConnector)
}

func (entity *Entity) addConnection(otherEntityId int, connectionId string, w wire, circuitId int) {

	connection, ok := entity.Connections[connectionId]
	if !ok {
		connection = make(map[wire][]Circuit)
		entity.Connections[connectionId] = connection
	}

	circuit, ok := connection[w]
	if !ok {
		circuit = []Circuit{}
		connection[w] = circuit
	}

	c := Circuit{EntityId: otherEntityId, CircuitId: circuitId}
	connection[w] = append(circuit, c)
}
