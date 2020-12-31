package discord

type property struct {
	Name        string
	Value       int // Money: Value*Upgrades* Hours since last visited
	Cost        int // Initial Price
	UpgradeCost int // Upgrade^1.5 * UpgradeCost
	ID          string
	Credit      int
}

var upgrades = []property{
	property{
		Name:        "Snack Booth",
		Value:       20,
		Cost:        1000,
		UpgradeCost: 200,
		ID:          "snack",
		Credit:      50,
	},
	property{
		Name:        "Homemade Cookie Business",
		Value:       50,
		Cost:        10000,
		UpgradeCost: 600,
		ID:          "cookie",
		Credit:      100,
	},
	property{
		Name:        "Li'l Jon'z Fudge Store",
		Value:       80,
		UpgradeCost: 960,
		Cost:        50000,
		ID:          "fudge",
		Credit:      200,
	},
	property{
		Name:        `|\\/|cDonaIds`,
		Cost:        100000,
		Value:       100,
		UpgradeCost: 1400,
		ID:          "mcd",
		Credit:      250,
	},
	property{
		Name:        "Village Bank",
		Value:       120,
		Cost:        150000,
		UpgradeCost: 1500,
		ID:          "village",
		Credit:      300,
	},
	property{
		Name:        "Vanilla JS Coders",
		Value:       140,
		Cost:        200000,
		UpgradeCost: 1750,
		ID:          "jspain",
		Credit:      350,
	},
	property{
		Name:        "We Use Hacks In Creative",
		Value:       300,
		Cost:        400000,
		UpgradeCost: 2500,
		ID:          "rich",
		Credit:      400,
	},
	property{
		Name:        "Scammers Inc.",
		Value:       0,
		Cost:        10000,
		UpgradeCost: 600,
		ID:          "rob",
		Credit:      100,
	},
}
