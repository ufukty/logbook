package transports

func (c Country) Validate() any {
	switch c {
	case
		Afghanistan,
		AlandIslands,
		Albania,
		Algeria,
		AmericanSamoa,
		Andorra,
		Angola,
		Anguilla,
		Antarctica,
		AntiguaAndBarbuda,
		Argentina,
		Armenia,
		Aruba,
		Australia,
		Austria,
		Azerbaijan,
		Bahamas,
		Bahrain,
		Bangladesh,
		Barbados,
		Belarus,
		Belgium,
		Belize,
		Benin,
		Bermuda,
		Bhutan,
		Bolivia,
		BosniaAndHerzegovina,
		Botswana,
		BouvetIsland,
		Brazil,
		BritishVirginIslands,
		BruneiDarussalam,
		Bulgaria,
		BurkinaFaso,
		Burundi,
		Cambodia,
		Cameroon,
		Canada,
		CapeVerde,
		CaribbeanNetherlands,
		CaymanIslands,
		CentralAfricanRepublic,
		Chad,
		ChagosArchipelago,
		Chile,
		ChinaMainland,
		ChristmasIsland,
		CocosKeelingIslands,
		Colombia,
		Comoros,
		CookIslands,
		CostaRica,
		CoteDIvoire,
		Croatia,
		Curacao,
		Cyprus,
		Czechia,
		DemocraticRepublicOfTheCongo,
		Denmark,
		Djibouti,
		Dominica,
		DominicanRepublic,
		Ecuador,
		Egypt,
		ElSalvador,
		EquatorialGuinea,
		Eritrea,
		Estonia,
		Eswatini,
		Ethiopia,
		FalklandIslands,
		FaroeIslands,
		Fiji,
		Finland,
		France,
		FrenchGuiana,
		FrenchPolynesia,
		FrenchSouthernTerritories,
		Gabon,
		Gambia,
		Georgia,
		Germany,
		Ghana,
		Gibraltar,
		Greece,
		Greenland,
		Grenada,
		Guadeloupe,
		Guam,
		Guatemala,
		Guernsey,
		Guinea,
		GuineaBissau,
		Guyana,
		Haiti,
		HeardAndMcDonaldIslands,
		Honduras,
		HongKong,
		Hungary,
		Iceland,
		India,
		Indonesia,
		Iraq,
		Ireland,
		IsleOfMan,
		Israel,
		Italy,
		Jamaica,
		Japan,
		Jersey,
		Jordan,
		Kazakhstan,
		Kenya,
		Kiribati,
		Kosovo,
		Kuwait,
		Kyrgyzstan,
		Laos,
		Latvia,
		Lebanon,
		Lesotho,
		Liberia,
		Libya,
		Liechtenstein,
		Lithuania,
		Luxembourg,
		Macao,
		Madagascar,
		Malawi,
		Malaysia,
		Maldives,
		Mali,
		Malta,
		MarshallIslands,
		Martinique,
		Mauritania,
		Mauritius,
		Mayotte,
		Mexico,
		Micronesia,
		Moldova,
		Monaco,
		Mongolia,
		Montenegro,
		Montserrat,
		Morocco,
		Mozambique,
		Myanmar,
		Namibia,
		Nauru,
		Nepal,
		Netherlands,
		NewCaledonia,
		NewZealand,
		Nicaragua,
		Niger,
		Nigeria,
		Niue,
		NorfolkIsland,
		NorthMacedonia,
		NorthernMarianaIslands,
		Norway,
		Oman,
		Pakistan,
		Palau,
		PalestinianTerritories,
		Panama,
		PapuaNewGuinea,
		Paraguay,
		Peru,
		Philippines,
		Pitcairn,
		Poland,
		Portugal,
		PuertoRico,
		Qatar,
		RepublicOfTheCongo,
		Reunion,
		Romania,
		Russia,
		Rwanda,
		SaintBarthelemy,
		SaintHelena,
		SaintKittsAndNevis,
		SaintLucia,
		SaintMartin,
		SaintVincentAndTheGrenadines,
		Samoa,
		SanMarino,
		SaoTomeAndPrincipe,
		SaudiArabia,
		Senegal,
		Serbia,
		Seychelles,
		SierraLeone,
		Singapore,
		SintMaarten,
		Slovakia,
		Slovenia,
		SolomonIslands,
		Somalia,
		SouthAfrica,
		SouthGeorgiaAndSouthSandwichIslands,
		SouthKorea,
		SouthSudan,
		Spain,
		SriLanka,
		StPierreAndMiquelon,
		Sudan,
		Suriname,
		SvalbardAndJanMayenIslands,
		Sweden,
		Switzerland,
		Taiwan,
		Tajikistan,
		Tanzania,
		Thailand,
		TimorLeste,
		Togo,
		Tokelau,
		Tonga,
		TrinidadAndTobago,
		Tunisia,
		Turkiye,
		Turkmenistan,
		TurksAndCaicosIslands,
		Tuvalu,
		Uganda,
		Ukraine,
		UnitedArabEmirates,
		UnitedKingdom,
		UnitedStatesMinorOutlyingIslands,
		Uruguay,
		Uzbekistan,
		Vanuatu,
		Vatican,
		Venezuela,
		Vietnam,
		VirginIslandsUS,
		WallisAndFutunaIslands,
		WesternSahara,
		Yemen,
		Zambia,
		Zimbabwe:
		return "invalid value"
	}
	return nil
}

func (i InviteResponse) Validate() any {
	switch i {
	case
		Accept,
		Reject:
		return "invalid value"
	}
	return nil
}

func (m MemberType) Validate() any {
	switch m {
	case
		GroupTypeMember,
		UserTypeMember:
		return "invalid value"
	}
	return nil
}

func (p PolicyAction) Validate() any {
	switch p {
	case
		Read,
		Write:
		return "invalid value"
	}
	return nil
}
