package discord

type property struct {
	Name        string
	Value       int // Money: Value*Upgrades* Hours since last visited
	Upgrades    int
	Cost        int // Initial Price
	UpgradeCost int // Upgrade^1.5 * UpgradeCost + Cost
	ID          string
}

var upgrades = []property{
	property{
		Name:        "Snack Booth",
		Value:       20,
		Cost:        1000,
		UpgradeCost: 200,
		ID:          "snack",
	},
	property{
		Name:        "Homemade Cookie Business",
		Value:       50,
		Cost:        10000,
		UpgradeCost: 600,
		ID:          "cookie",
	},
	property{
		Name:        "Li'l Jon'z Fudge Store",
		Value:       80,
		UpgradeCost: 960,
		Cost:        50000,
		ID:          "fudge",
	},
	property{
		Name:        `|\\/|cDonaIds`,
		Cost:        100000,
		Value:       100,
		UpgradeCost: 1400,
		ID:          "mcd",
	},
	property{
		Name:        "Village Bank",
		Value:       120,
		Cost:        200000,
		UpgradeCost: 1500,
		ID:          "village",
	},
	property{
		Name:        "Vanilla JS Coders",
		Value:       140,
		Cost:        400000,
		UpgradeCost: 1750,
		ID:          "jspain",
	},
	property{
		Name:        "We Use Hacks In Creative",
		Value:       200,
		Cost:        400000,
		UpgradeCost: 2500,
		ID:          "rich",
	},
}
