package types

type Location struct {
	Name   string
	Delays map[string]int64
}

var (
	LocCN = Location{
		Name: "CN",
		Delays: map[string]int64{
			"CN":  50,
			"NA":  300,
			"EU":  300,
			"RU":  150,
			"SEA": 100,
			"JP":  100,
		},
	}
	LocNA = Location{
		Name: "NA",
		Delays: map[string]int64{
			"CN":  300,
			"NA":  50,
			"EU":  200,
			"RU":  400,
			"SEA": 300,
			"JP":  200,
		},
	}
	LocEU = Location{
		Name: "EU",
		Delays: map[string]int64{
			"CN":  300,
			"NA":  200,
			"EU":  50,
			"RU":  200,
			"SEA": 200,
			"JP":  200,
		},
	}
	LocRU = Location{
		Name: "RU",
		Delays: map[string]int64{
			"CN":  150,
			"NA":  400,
			"EU":  200,
			"RU":  50,
			"SEA": 150,
			"JP":  150,
		},
	}
	LocSEA = Location{
		Name: "SEA",
		Delays: map[string]int64{
			"CN":  100,
			"NA":  300,
			"EU":  200,
			"RU":  150,
			"SEA": 50,
			"JP":  100,
		},
	}
	LocJP = Location{
		Name: "JP",
		Delays: map[string]int64{
			"CN":  100,
			"NA":  200,
			"EU":  200,
			"RU":  150,
			"SEA": 100,
			"JP":  50,
		},
	}
)

var ConstLocations = []Location{
	LocCN, LocNA, LocEU, LocRU, LocSEA, LocJP,
}
