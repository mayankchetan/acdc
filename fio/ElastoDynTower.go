package fio

var ElastoDynTower = NewFile("ElastoDynTower", []Field{
	{Type: Heading, Desc: "ElastoDynTower Input File"},
	{Name: "Title", Type: Title},
	{Type: Heading, Desc: "Tower Parameters"},
	{Name: "NTwInpSt", Type: Int, Desc: `Number of input stations to specify tower geometry`},
	{Name: "TwrFADmp(1)", Type: Float, Desc: `Tower 1st fore-aft mode structural damping ratio`, Unit: "%"},
	{Name: "TwrFADmp(2)", Type: Float, Desc: `Tower 2nd fore-aft mode structural damping ratio`, Unit: "%"},
	{Name: "TwrSSDmp(1)", Type: Float, Desc: `Tower 1st side-to-side mode structural damping ratio`, Unit: "%"},
	{Name: "TwrSSDmp(2)", Type: Float, Desc: `Tower 2nd side-to-side mode structural damping ratio`, Unit: "%"},
	{Type: Heading, Desc: "Tower Adjustment Factors"},
	{Name: "FAStTunr(1)", Type: Float, Desc: `Tower fore-aft modal stiffness tuner, 1st mode`},
	{Name: "FAStTunr(2)", Type: Float, Desc: `Tower fore-aft modal stiffness tuner, 2nd mode`},
	{Name: "SSStTunr(1)", Type: Float, Desc: `Tower side-to-side stiffness tuner, 1st mode`},
	{Name: "SSStTunr(2)", Type: Float, Desc: `Tower side-to-side stiffness tuner, 2nd mode`},
	{Name: "AdjTwMa", Type: Float, Desc: `Factor to adjust tower mass density`},
	{Name: "AdjFASt", Type: Float, Desc: `Factor to adjust tower fore-aft stiffness`},
	{Name: "AdjSSSt", Type: Float, Desc: `Factor to adjust tower side-to-side stiffness`},
	{Type: Heading, Desc: "Distributed Tower Properties"},
	{Name: "TwInpSt", Type: Table, Num: "NTwInpSt",
		TableHeaderSize: 2,
		TableColumns: []Column{
			{Name: "HtFract", Type: Float, Unit: "-"},
			{Name: "TMassDen", Type: Float, Unit: "kg/m"},
			{Name: "TwFAStif", Type: Float, Unit: "Nm^2"},
			{Name: "TwSSStif", Type: Float, Unit: "Nm^2"},
		}},
	{Type: Heading, Desc: "Tower Fore-Aft Mode Shapes"},
	{Name: "TwFAM1Sh(2)", Type: Float, Desc: `Mode 1, coefficient of x^2 term`},
	{Name: "TwFAM1Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Name: "TwFAM1Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Name: "TwFAM1Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Name: "TwFAM1Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
	{Name: "TwFAM2Sh(2)", Type: Float, Desc: `Mode 2, coefficient of x^2 term`},
	{Name: "TwFAM2Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Name: "TwFAM2Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Name: "TwFAM2Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Name: "TwFAM2Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
	{Type: Heading, Desc: "Tower Side-To-Side Mode Shapes"},
	{Name: "TwSSM1Sh(2)", Type: Float, Desc: `Mode 1, coefficient of x^2 term`},
	{Name: "TwSSM1Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Name: "TwSSM1Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Name: "TwSSM1Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Name: "TwSSM1Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
	{Name: "TwSSM2Sh(2)", Type: Float, Desc: `Mode 2, coefficient of x^2 term`},
	{Name: "TwSSM2Sh(3)", Type: Float, Desc: `      , coefficient of x^3 term`},
	{Name: "TwSSM2Sh(4)", Type: Float, Desc: `      , coefficient of x^4 term`},
	{Name: "TwSSM2Sh(5)", Type: Float, Desc: `      , coefficient of x^5 term`},
	{Name: "TwSSM2Sh(6)", Type: Float, Desc: `      , coefficient of x^6 term`},
})
