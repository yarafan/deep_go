package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		copy(person.name[:], name)
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = person.manaHealthFlags | uint32(mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = person.manaHealthFlags | (uint32(health) << 10)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attrs = person.attrs | uint16(respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attrs = person.attrs | (uint16(strength) << 4)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attrs = person.attrs | (uint16(experience) << 8)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.attrs = person.attrs | (uint16(level) << 12)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = person.manaHealthFlags | (uint32(1) << 20)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = person.manaHealthFlags | (uint32(1) << 22)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = person.manaHealthFlags | (uint32(1) << 21)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.manaHealthFlags = person.manaHealthFlags | (uint32(personType) << 23)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x               int32
	y               int32
	z               int32
	gold            uint32
	manaHealthFlags uint32
	attrs           uint16
	name            [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}

	for _, option := range options {
		option(&person)
	}

	return person
}

func (p *GamePerson) Name() string {
	n := 0
	for n < len(p.name) && p.name[n] != 0 {
		n++
	}

	return string(p.name[:n])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.manaHealthFlags & 0x3ff)
}

func (p *GamePerson) Health() int {
	return int((p.manaHealthFlags & (0x3ff << 10)) >> 10)
}

func (p *GamePerson) Respect() int {
	return int(p.attrs & 0x000f)
}

func (p *GamePerson) Strength() int {
	return int((p.attrs & 0x00f0) >> 4)
}

func (p *GamePerson) Experience() int {
	return int((p.attrs & 0x0f00) >> 8)
}

func (p *GamePerson) Level() int {
	return int((p.attrs & 0xf000) >> 12)
}

func (p *GamePerson) HasHouse() bool {
	return (p.manaHealthFlags & 0x00100000) != 0
}

func (p *GamePerson) HasGun() bool {
	return (p.manaHealthFlags & 0x00400000) != 0
}

func (p *GamePerson) HasFamily() bool {
	return (p.manaHealthFlags & 0x00200000) != 0
}

func (p *GamePerson) Type() int {
	return int((p.manaHealthFlags & 0x01800000) >> 23)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}
