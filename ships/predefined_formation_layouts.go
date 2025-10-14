package ships

type FormationLayoutPosition struct {
	Filled      bool    `json:"filled"`
	BucketIndex *int    `json:"bucket_index"`
	Quantity    int     `json:"quantity"`
	Position    Vector2 `json:"position"`
	Order       int     `json:"order"`
}

type Vector2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

var LineFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 0}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 0}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 2}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 2}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 3}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 3}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 4}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 4}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 5}, Order: 13},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 5}, Order: 14},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -1}, Order: 15},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -1}, Order: 16},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -2}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -2}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -2}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -3}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -3}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -3}, Order: 8},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -2}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -2}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -2}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -2}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -2}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -3}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -3}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -3}, Order: 10},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -4}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -4}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -5}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -5}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -6}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -6}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -7}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -7}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -8}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -8}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -9}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -9}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -4}, Order: 13},
	},
}

var VanguardFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 0}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 2}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 4}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 5}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 6}, Order: 13},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 7}, Order: 15},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 8}, Order: 17},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 9}, Order: 19},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 10}, Order: 21},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 0}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 2}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 3}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 4}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 5}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 6}, Order: 14},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 7}, Order: 16},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 8}, Order: 18},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 9}, Order: 20},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 10}, Order: 22},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, 1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, 0}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, -2}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 0}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -2}, Order: 8},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -1}, Order: 5},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -2}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -2}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -2}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -3}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -4}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -5}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -3}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -4}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -5}, Order: 10},
	},
}

var EchelonFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 6}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 5}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 4}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, 3}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, 0}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -1}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -2}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -3}, Order: 10},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -2}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 0}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, 1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 2}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 3}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 4}, Order: 6},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, 4}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, 3}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 2}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 0}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -2}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -3}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -4}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -5}, Order: 10},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, 6}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -4}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{7, 7}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -5}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{7, 5}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -6}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{8, 6}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, -5}, Order: 8},
	},
}

var BoxFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, 1}, Order: 15},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, 1}, Order: 13},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, 1}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, 1}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 1}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, 1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 1}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 1}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 1}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, 1}, Order: 14},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, 0}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -2}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -3}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -4}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, 0}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -2}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -4}, Order: 9},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -3}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -3}, Order: 4},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -5}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -5}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -5}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -5}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -5}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -5}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -5}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -5}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -5}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -5}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -5}, Order: 10},
	},
}

var SkirmishFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 0}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, 0}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, 1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, 1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, 1}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 1}, Order: 8},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -2}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -2}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -2}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -2}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -3}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -3}, Order: 8},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, -1}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, -1}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-7, -1}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{7, -1}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-8, 0}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{8, 0}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-8, -2}, Order: 13},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{8, -2}, Order: 14},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-9, 0}, Order: 15},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{9, 0}, Order: 16},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-9, -2}, Order: 17},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{9, -2}, Order: 18},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-10, -1}, Order: 19},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{10, -1}, Order: 20},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -5}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -6}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -5}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -5}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -4}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -4}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -3}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -3}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -3}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -4}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -4}, Order: 12},
	},
}

var SwarmFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 0}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 2}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 3}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 4}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 0}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 2}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 3}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 4}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -1}, Order: 2},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, 0}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, 1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, 2}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, 3}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-7, 4}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 0}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 2}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, 3}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, 4}, Order: 12},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -2}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -2}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -2}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -2}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -2}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -3}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -3}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -3}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -3}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -3}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -3}, Order: 11},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -4}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -5}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, -6}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, -7}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-7, -8}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-8, -9}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -4}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -5}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -6}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, -7}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, -8}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{7, -9}, Order: 12},
	},
}

var PhalanxFormation = map[string][]FormationLayoutPosition{
	"front": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-7, 1}, Order: 15},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, 1}, Order: 13},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-5, 1}, Order: 11},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, 1}, Order: 9},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, 1}, Order: 7},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, 1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, 1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, 1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, 1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, 1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, 1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, 1}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{5, 1}, Order: 10},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, 1}, Order: 12},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{7, 1}, Order: 14},
	},

	"flank": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-6, -1}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-7, -1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-8, -1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-9, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{6, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{7, -1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{8, -1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{9, -1}, Order: 7},
	},

	"support": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -1}, Order: 8},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -1}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-2, -1}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -1}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -1}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{2, -1}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -1}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -1}, Order: 7},
	},

	"back": {
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{0, -2}, Order: 1},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-1, -3}, Order: 2},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{1, -3}, Order: 3},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{3, -3}, Order: 4},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-3, -3}, Order: 5},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{-4, -4}, Order: 6},
		{Filled: false, BucketIndex: nil, Quantity: 0, Position: Vector2{4, -4}, Order: 7},
	},
}
