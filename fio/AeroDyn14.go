package fio

var AeroDyn14 = NewFile("AeroDyn14", []Field{
	{Type: Heading, Desc: "AeroDyn14 Input File"},
	{Name: "Title", Type: Title},
	{Name: "StallMod", Type: String, Desc: `Dynamic stall included [BEDDOES or STEADY]`, Unit: "unquoted string"},
	{Name: "UseCm", Type: String, Desc: `Use aerodynamic pitching moment model? [USE_CM or NO_CM]`, Unit: "unquoted string"},
	{Name: "InfModel", Type: String, Desc: `Inflow model [DYNIN or EQUIL]`, Unit: "unquoted string"},
	{Name: "IndModel", Type: String, Desc: `Induction-factor model [NONE or WAKE or SWIRL]`, Unit: "unquoted string"},
	{Name: "AToler", Type: Float, Desc: `Induction-factor tolerance (convergence criteria)`},
	{Name: "TLModel", Type: String, Desc: `Tip-loss model (EQUIL only) [PRANDtl, GTECH, or NONE]`, Unit: "unquoted string"},
	{Name: "HLModel", Type: String, Desc: `Hub-loss model (EQUIL only) [PRANdtl or NONE]`, Unit: "unquoted string"},
	{Name: "TwrShad", Type: Float, Desc: `Tower-shadow velocity deficit`},
	{Name: "ShadHWid", Type: Float, Desc: `Tower-shadow half width`, Unit: "m"},
	{Name: "T_Shad_Refpt", Type: Float, Desc: `Tower-shadow reference point`, Unit: "m"},
	{Name: "AirDens", Type: Float, Desc: `Air density`, Unit: "kg/m^3"},
	{Name: "KinVisc", Type: Float, Desc: `Kinematic air viscosity [CURRENTLY IGNORED]`, Unit: "m^2/sec"},
	{Name: "DTAero", Type: FloatDefault, Desc: `Time interval for aerodynamic calculations`, Unit: "sec"},
	{Name: "NumFoil", Type: Int, Desc: `Number of airfoil files`},
	{Name: "FoilNm", Type: Paths, Num: "NumFoil", Desc: `Names of the airfoil files [NumFoil lines]`, Unit: "quoted string",
		PathFileType: "Foil"},
	{Name: "BldNodes", Type: Int, Desc: `Number of blade nodes used for analysis`},
	{Name: "BldNodeData", Type: Table, Num: "BldNodes",
		TableHeaderSize: 1,
		TableColumns: []Column{
			{Name: "RNodes", Type: Float},
			{Name: "AeroTwst", Type: Float},
			{Name: "DRNodes", Type: Float},
			{Name: "Chord", Type: Float},
			{Name: "NFoil", Type: Int},
			{Name: "PrnElm", Type: String},
		}},
})
