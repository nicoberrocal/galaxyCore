package ships

// FormationTreeCatalog contains all formation mastery trees.
// Key: FormationType (empty string "" = global tree)
var FormationTreeCatalog = map[FormationType]FormationTree{}

func init() {
	initGlobalTree()
	initLineTree()
	initBoxTree()
	initVanguardTree()
	initSkirmishTree()
	initEchelonTree()
	initPhalanxTree()
	initSwarmTree()
}
