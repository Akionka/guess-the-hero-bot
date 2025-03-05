package main

// var items = []Item{
// 	{
// 		ID:          1,
// 		DisplayName: "Blink Dagger",
// 		ShortName:   "blink",
// 	}, {
// 		ID:          2,
// 		DisplayName: "Blades of Attack",
// 		ShortName:   "blades_of_attack",
// 	}, {
// 		ID:          3,
// 		DisplayName: "Broadsword",
// 		ShortName:   "broadsword",
// 	}, {
// 		ID:          4,
// 		DisplayName: "Chainmail",
// 		ShortName:   "chainmail",
// 	}, {
// 		ID:          5,
// 		DisplayName: "Claymore",
// 		ShortName:   "claymore",
// 	}, {
// 		ID:          6,
// 		DisplayName: "Helm of Iron Will",
// 		ShortName:   "helm_of_iron_will",
// 	}, {
// 		ID:          7,
// 		DisplayName: "Javelin",
// 		ShortName:   "javelin",
// 	}, {
// 		ID:          8,
// 		DisplayName: "Mithril Hammer",
// 		ShortName:   "mithril_hammer",
// 	}, {
// 		ID:          9,
// 		DisplayName: "Platemail",
// 		ShortName:   "platemail",
// 	}, {
// 		ID:          10,
// 		DisplayName: "Quarterstaff",
// 		ShortName:   "quarterstaff",
// 	}, {
// 		ID:          11,
// 		DisplayName: "Quelling Blade",
// 		ShortName:   "quelling_blade",
// 	}, {
// 		ID:          12,
// 		DisplayName: "Ring of Protection",
// 		ShortName:   "ring_of_protection",
// 	}, {
// 		ID:          13,
// 		DisplayName: "Gauntlets of Strength",
// 		ShortName:   "gauntlets",
// 	}, {
// 		ID:          14,
// 		DisplayName: "Slippers of Agility",
// 		ShortName:   "slippers",
// 	}, {
// 		ID:          15,
// 		DisplayName: "Mantle of Intelligence",
// 		ShortName:   "mantle",
// 	}, {
// 		ID:          16,
// 		DisplayName: "Iron Branch",
// 		ShortName:   "branches",
// 	}, {
// 		ID:          17,
// 		DisplayName: "Belt of Strength",
// 		ShortName:   "belt_of_strength",
// 	}, {
// 		ID:          18,
// 		DisplayName: "Band of Elvenskin",
// 		ShortName:   "boots_of_elves",
// 	}, {
// 		ID:          19,
// 		DisplayName: "Robe of the Magi",
// 		ShortName:   "robe",
// 	}, {
// 		ID:          20,
// 		DisplayName: "Circlet",
// 		ShortName:   "circlet",
// 	}, {
// 		ID:          21,
// 		DisplayName: "Ogre Axe",
// 		ShortName:   "ogre_axe",
// 	}, {
// 		ID:          22,
// 		DisplayName: "Blade of Alacrity",
// 		ShortName:   "blade_of_alacrity",
// 	}, {
// 		ID:          23,
// 		DisplayName: "Staff of Wizardry",
// 		ShortName:   "staff_of_wizardry",
// 	}, {
// 		ID:          24,
// 		DisplayName: "Ultimate Orb",
// 		ShortName:   "ultimate_orb",
// 	}, {
// 		ID:          25,
// 		DisplayName: "Gloves of Haste",
// 		ShortName:   "gloves",
// 	}, {
// 		ID:          26,
// 		DisplayName: "Morbid Mask",
// 		ShortName:   "lifesteal",
// 	}, {
// 		ID:          27,
// 		DisplayName: "Ring of Regen",
// 		ShortName:   "ring_of_regen",
// 	}, {
// 		ID:          28,
// 		DisplayName: "Sage's Mask",
// 		ShortName:   "sobi_mask",
// 	}, {
// 		ID:          29,
// 		DisplayName: "Boots of Speed",
// 		ShortName:   "boots",
// 	}, {
// 		ID:          30,
// 		DisplayName: "Gem of True Sight",
// 		ShortName:   "gem",
// 	}, {
// 		ID:          31,
// 		DisplayName: "Cloak",
// 		ShortName:   "cloak",
// 	}, {
// 		ID:          32,
// 		DisplayName: "Talisman of Evasion",
// 		ShortName:   "talisman_of_evasion",
// 	}, {
// 		ID:          33,
// 		DisplayName: "Cheese",
// 		ShortName:   "cheese",
// 	}, {
// 		ID:          34,
// 		DisplayName: "Magic Stick",
// 		ShortName:   "magic_stick",
// 	}, {
// 		ID:          35,
// 		DisplayName: "Magic Wand Recipe",
// 		ShortName:   "recipe_magic_wand",
// 	}, {
// 		ID:          36,
// 		DisplayName: "Magic Wand",
// 		ShortName:   "magic_wand",
// 	}, {
// 		ID:          37,
// 		DisplayName: "Ghost Scepter",
// 		ShortName:   "ghost",
// 	}, {
// 		ID:          38,
// 		DisplayName: "Clarity",
// 		ShortName:   "clarity",
// 	}, {
// 		ID:          39,
// 		DisplayName: "Healing Salve",
// 		ShortName:   "flask",
// 	}, {
// 		ID:          40,
// 		DisplayName: "Dust of Appearance",
// 		ShortName:   "dust",
// 	}, {
// 		ID:          41,
// 		DisplayName: "Bottle",
// 		ShortName:   "bottle",
// 	}, {
// 		ID:          42,
// 		DisplayName: "Observer Ward",
// 		ShortName:   "ward_observer",
// 	}, {
// 		ID:          43,
// 		DisplayName: "Sentry Ward",
// 		ShortName:   "ward_sentry",
// 	}, {
// 		ID:          44,
// 		DisplayName: "Tango",
// 		ShortName:   "tango",
// 	}, {
// 		ID:          45,
// 		DisplayName: "Animal Courier",
// 		ShortName:   "courier",
// 	}, {
// 		ID:          46,
// 		DisplayName: "Town Portal Scroll",
// 		ShortName:   "tpscroll",
// 	}, {
// 		ID:          47,
// 		DisplayName: "Boots of Travel Recipe",
// 		ShortName:   "recipe_travel_boots",
// 	}, {
// 		ID:          48,
// 		DisplayName: "Boots of Travel",
// 		ShortName:   "travel_boots",
// 	}, {
// 		ID:          49,
// 		DisplayName: "",
// 		ShortName:   "recipe_phase_boots",
// 	}, {
// 		ID:          50,
// 		DisplayName: "Phase Boots",
// 		ShortName:   "phase_boots",
// 	}, {
// 		ID:          51,
// 		DisplayName: "Demon Edge",
// 		ShortName:   "demon_edge",
// 	}, {
// 		ID:          52,
// 		DisplayName: "Eaglesong",
// 		ShortName:   "eagle",
// 	}, {
// 		ID:          53,
// 		DisplayName: "Reaver",
// 		ShortName:   "reaver",
// 	}, {
// 		ID:          54,
// 		DisplayName: "Sacred Relic",
// 		ShortName:   "relic",
// 	}, {
// 		ID:          55,
// 		DisplayName: "Hyperstone",
// 		ShortName:   "hyperstone",
// 	}, {
// 		ID:          56,
// 		DisplayName: "Ring of Health",
// 		ShortName:   "ring_of_health",
// 	}, {
// 		ID:          57,
// 		DisplayName: "Void Stone",
// 		ShortName:   "void_stone",
// 	}, {
// 		ID:          58,
// 		DisplayName: "Mystic Staff",
// 		ShortName:   "mystic_staff",
// 	}, {
// 		ID:          59,
// 		DisplayName: "Energy Booster",
// 		ShortName:   "energy_booster",
// 	}, {
// 		ID:          60,
// 		DisplayName: "Point Booster",
// 		ShortName:   "point_booster",
// 	}, {
// 		ID:          61,
// 		DisplayName: "Vitality Booster",
// 		ShortName:   "vitality_booster",
// 	}, {
// 		ID:          62,
// 		DisplayName: "",
// 		ShortName:   "recipe_power_treads",
// 	}, {
// 		ID:          63,
// 		DisplayName: "Power Treads",
// 		ShortName:   "power_treads",
// 	}, {
// 		ID:          64,
// 		DisplayName: "Hand of Midas Recipe",
// 		ShortName:   "recipe_hand_of_midas",
// 	}, {
// 		ID:          65,
// 		DisplayName: "Hand of Midas",
// 		ShortName:   "hand_of_midas",
// 	}, {
// 		ID:          66,
// 		DisplayName: "",
// 		ShortName:   "recipe_oblivion_staff",
// 	}, {
// 		ID:          67,
// 		DisplayName: "Oblivion Staff",
// 		ShortName:   "oblivion_staff",
// 	}, {
// 		ID:          68,
// 		DisplayName: "",
// 		ShortName:   "recipe_pers",
// 	}, {
// 		ID:          69,
// 		DisplayName: "Perseverance",
// 		ShortName:   "pers",
// 	}, {
// 		ID:          70,
// 		DisplayName: "",
// 		ShortName:   "recipe_poor_mans_shield",
// 	}, {
// 		ID:          71,
// 		DisplayName: "Poor Man's Shield",
// 		ShortName:   "poor_mans_shield",
// 	}, {
// 		ID:          72,
// 		DisplayName: "Bracer Recipe",
// 		ShortName:   "recipe_bracer",
// 	}, {
// 		ID:          73,
// 		DisplayName: "Bracer",
// 		ShortName:   "bracer",
// 	}, {
// 		ID:          74,
// 		DisplayName: "Wraith Band Recipe",
// 		ShortName:   "recipe_wraith_band",
// 	}, {
// 		ID:          75,
// 		DisplayName: "Wraith Band",
// 		ShortName:   "wraith_band",
// 	}, {
// 		ID:          76,
// 		DisplayName: "Null Talisman Recipe",
// 		ShortName:   "recipe_null_talisman",
// 	}, {
// 		ID:          77,
// 		DisplayName: "Null Talisman",
// 		ShortName:   "null_talisman",
// 	}, {
// 		ID:          78,
// 		DisplayName: "Mekansm Recipe",
// 		ShortName:   "recipe_mekansm",
// 	}, {
// 		ID:          79,
// 		DisplayName: "Mekansm",
// 		ShortName:   "mekansm",
// 	}, {
// 		ID:          80,
// 		DisplayName: "Vladmir's Offering Recipe",
// 		ShortName:   "recipe_vladmir",
// 	}, {
// 		ID:          81,
// 		DisplayName: "Vladmir's Offering",
// 		ShortName:   "vladmir",
// 	}, {
// 		ID:          84,
// 		DisplayName: "Flying Courier",
// 		ShortName:   "flying_courier",
// 	}, {
// 		ID:          85,
// 		DisplayName: "Buckler Recipe",
// 		ShortName:   "recipe_buckler",
// 	}, {
// 		ID:          86,
// 		DisplayName: "Buckler",
// 		ShortName:   "buckler",
// 	}, {
// 		ID:          87,
// 		DisplayName: "Ring of Basilius Recipe",
// 		ShortName:   "recipe_ring_of_basilius",
// 	}, {
// 		ID:          88,
// 		DisplayName: "Ring of Basilius",
// 		ShortName:   "ring_of_basilius",
// 	}, {
// 		ID:          89,
// 		DisplayName: "Pipe of Insight Recipe",
// 		ShortName:   "recipe_pipe",
// 	}, {
// 		ID:          90,
// 		DisplayName: "Pipe of Insight",
// 		ShortName:   "pipe",
// 	}, {
// 		ID:          91,
// 		DisplayName: "Urn of Shadows Recipe",
// 		ShortName:   "recipe_urn_of_shadows",
// 	}, {
// 		ID:          92,
// 		DisplayName: "Urn of Shadows",
// 		ShortName:   "urn_of_shadows",
// 	}, {
// 		ID:          93,
// 		DisplayName: "Headdress Recipe",
// 		ShortName:   "recipe_headdress",
// 	}, {
// 		ID:          94,
// 		DisplayName: "Headdress",
// 		ShortName:   "headdress",
// 	}, {
// 		ID:          95,
// 		DisplayName: "Scythe of Vyse Recipe",
// 		ShortName:   "recipe_sheepstick",
// 	}, {
// 		ID:          96,
// 		DisplayName: "Scythe of Vyse",
// 		ShortName:   "sheepstick",
// 	}, {
// 		ID:          97,
// 		DisplayName: "Orchid Malevolence Recipe",
// 		ShortName:   "recipe_orchid",
// 	}, {
// 		ID:          98,
// 		DisplayName: "Orchid Malevolence",
// 		ShortName:   "orchid",
// 	}, {
// 		ID:          99,
// 		DisplayName: "Eul's Scepter Recipe",
// 		ShortName:   "recipe_cyclone",
// 	}, {
// 		ID:          100,
// 		DisplayName: "Eul's Scepter of Divinity",
// 		ShortName:   "cyclone",
// 	}, {
// 		ID:          101,
// 		DisplayName: "Force Staff Recipe",
// 		ShortName:   "recipe_force_staff",
// 	}, {
// 		ID:          102,
// 		DisplayName: "Force Staff",
// 		ShortName:   "force_staff",
// 	}, {
// 		ID:          103,
// 		DisplayName: "Dagon Recipe",
// 		ShortName:   "recipe_dagon",
// 	}, {
// 		ID:          104,
// 		DisplayName: "Dagon",
// 		ShortName:   "dagon",
// 	}, {
// 		ID:          105,
// 		DisplayName: "Necronomicon Recipe",
// 		ShortName:   "recipe_necronomicon",
// 	}, {
// 		ID:          106,
// 		DisplayName: "Necronomicon",
// 		ShortName:   "necronomicon",
// 	}, {
// 		ID:          107,
// 		DisplayName: "",
// 		ShortName:   "recipe_ultimate_scepter",
// 	}, {
// 		ID:          108,
// 		DisplayName: "Aghanim's Scepter",
// 		ShortName:   "ultimate_scepter",
// 	}, {
// 		ID:          109,
// 		DisplayName: "Refresher Orb Recipe",
// 		ShortName:   "recipe_refresher",
// 	}, {
// 		ID:          110,
// 		DisplayName: "Refresher Orb",
// 		ShortName:   "refresher",
// 	}, {
// 		ID:          111,
// 		DisplayName: "Assault Cuirass Recipe",
// 		ShortName:   "recipe_assault",
// 	}, {
// 		ID:          112,
// 		DisplayName: "Assault Cuirass",
// 		ShortName:   "assault",
// 	}, {
// 		ID:          113,
// 		DisplayName: "Heart of Tarrasque Recipe",
// 		ShortName:   "recipe_heart",
// 	}, {
// 		ID:          114,
// 		DisplayName: "Heart of Tarrasque",
// 		ShortName:   "heart",
// 	}, {
// 		ID:          115,
// 		DisplayName: "Black King Bar Recipe",
// 		ShortName:   "recipe_black_king_bar",
// 	}, {
// 		ID:          116,
// 		DisplayName: "Black King Bar",
// 		ShortName:   "black_king_bar",
// 	}, {
// 		ID:          117,
// 		DisplayName: "Aegis of the Immortal",
// 		ShortName:   "aegis",
// 	}, {
// 		ID:          118,
// 		DisplayName: "Shiva's Guard Recipe",
// 		ShortName:   "recipe_shivas_guard",
// 	}, {
// 		ID:          119,
// 		DisplayName: "Shiva's Guard",
// 		ShortName:   "shivas_guard",
// 	}, {
// 		ID:          120,
// 		DisplayName: "Bloodstone Recipe",
// 		ShortName:   "recipe_bloodstone",
// 	}, {
// 		ID:          121,
// 		DisplayName: "Bloodstone",
// 		ShortName:   "bloodstone",
// 	}, {
// 		ID:          122,
// 		DisplayName: "Linken's Sphere Recipe",
// 		ShortName:   "recipe_sphere",
// 	}, {
// 		ID:          123,
// 		DisplayName: "Linken's Sphere",
// 		ShortName:   "sphere",
// 	}, {
// 		ID:          124,
// 		DisplayName: "",
// 		ShortName:   "recipe_vanguard",
// 	}, {
// 		ID:          125,
// 		DisplayName: "Vanguard",
// 		ShortName:   "vanguard",
// 	}, {
// 		ID:          126,
// 		DisplayName: "Blade Mail Recipe",
// 		ShortName:   "recipe_blade_mail",
// 	}, {
// 		ID:          127,
// 		DisplayName: "Blade Mail",
// 		ShortName:   "blade_mail",
// 	}, {
// 		ID:          128,
// 		DisplayName: "",
// 		ShortName:   "recipe_soul_booster",
// 	}, {
// 		ID:          129,
// 		DisplayName: "Soul Booster",
// 		ShortName:   "soul_booster",
// 	}, {
// 		ID:          130,
// 		DisplayName: "",
// 		ShortName:   "recipe_hood_of_defiance",
// 	}, {
// 		ID:          131,
// 		DisplayName: "Hood of Defiance",
// 		ShortName:   "hood_of_defiance",
// 	}, {
// 		ID:          132,
// 		DisplayName: "Divine Rapier Recipe",
// 		ShortName:   "recipe_rapier",
// 	}, {
// 		ID:          133,
// 		DisplayName: "Divine Rapier",
// 		ShortName:   "rapier",
// 	}, {
// 		ID:          134,
// 		DisplayName: "Monkey King Bar Recipe",
// 		ShortName:   "recipe_monkey_king_bar",
// 	}, {
// 		ID:          135,
// 		DisplayName: "Monkey King Bar",
// 		ShortName:   "monkey_king_bar",
// 	}, {
// 		ID:          136,
// 		DisplayName: "Radiance Recipe",
// 		ShortName:   "recipe_radiance",
// 	}, {
// 		ID:          137,
// 		DisplayName: "Radiance",
// 		ShortName:   "radiance",
// 	}, {
// 		ID:          138,
// 		DisplayName: "",
// 		ShortName:   "recipe_butterfly",
// 	}, {
// 		ID:          139,
// 		DisplayName: "Butterfly",
// 		ShortName:   "butterfly",
// 	}, {
// 		ID:          140,
// 		DisplayName: "Daedalus Recipe",
// 		ShortName:   "recipe_greater_crit",
// 	}, {
// 		ID:          141,
// 		DisplayName: "Daedalus",
// 		ShortName:   "greater_crit",
// 	}, {
// 		ID:          142,
// 		DisplayName: "Skull Basher Recipe",
// 		ShortName:   "recipe_basher",
// 	}, {
// 		ID:          143,
// 		DisplayName: "Skull Basher",
// 		ShortName:   "basher",
// 	}, {
// 		ID:          144,
// 		DisplayName: "Battle Fury Recipe",
// 		ShortName:   "recipe_bfury",
// 	}, {
// 		ID:          145,
// 		DisplayName: "Battle Fury",
// 		ShortName:   "bfury",
// 	}, {
// 		ID:          146,
// 		DisplayName: "Manta Style Recipe",
// 		ShortName:   "recipe_manta",
// 	}, {
// 		ID:          147,
// 		DisplayName: "Manta Style",
// 		ShortName:   "manta",
// 	}, {
// 		ID:          148,
// 		DisplayName: "Crystalys Recipe",
// 		ShortName:   "recipe_lesser_crit",
// 	}, {
// 		ID:          149,
// 		DisplayName: "Crystalys",
// 		ShortName:   "lesser_crit",
// 	}, {
// 		ID:          150,
// 		DisplayName: "Armlet of Mordiggian Recipe",
// 		ShortName:   "recipe_armlet",
// 	}, {
// 		ID:          151,
// 		DisplayName: "Armlet of Mordiggian",
// 		ShortName:   "armlet",
// 	}, {
// 		ID:          152,
// 		DisplayName: "Shadow Blade",
// 		ShortName:   "invis_sword",
// 	}, {
// 		ID:          153,
// 		DisplayName: "",
// 		ShortName:   "recipe_sange_and_yasha",
// 	}, {
// 		ID:          154,
// 		DisplayName: "Sange and Yasha",
// 		ShortName:   "sange_and_yasha",
// 	}, {
// 		ID:          155,
// 		DisplayName: "Satanic Recipe",
// 		ShortName:   "recipe_satanic",
// 	}, {
// 		ID:          156,
// 		DisplayName: "Satanic",
// 		ShortName:   "satanic",
// 	}, {
// 		ID:          157,
// 		DisplayName: "Mjollnir Recipe",
// 		ShortName:   "recipe_mjollnir",
// 	}, {
// 		ID:          158,
// 		DisplayName: "Mjollnir",
// 		ShortName:   "mjollnir",
// 	}, {
// 		ID:          159,
// 		DisplayName: "Eye of Skadi Recipe",
// 		ShortName:   "recipe_skadi",
// 	}, {
// 		ID:          160,
// 		DisplayName: "Eye of Skadi",
// 		ShortName:   "skadi",
// 	}, {
// 		ID:          161,
// 		DisplayName: "Sange Recipe",
// 		ShortName:   "recipe_sange",
// 	}, {
// 		ID:          162,
// 		DisplayName: "Sange",
// 		ShortName:   "sange",
// 	}, {
// 		ID:          163,
// 		DisplayName: "Helm of the Dominator Recipe",
// 		ShortName:   "recipe_helm_of_the_dominator",
// 	}, {
// 		ID:          164,
// 		DisplayName: "Helm of the Dominator",
// 		ShortName:   "helm_of_the_dominator",
// 	}, {
// 		ID:          165,
// 		DisplayName: "Maelstrom Recipe",
// 		ShortName:   "recipe_maelstrom",
// 	}, {
// 		ID:          166,
// 		DisplayName: "Maelstrom",
// 		ShortName:   "maelstrom",
// 	}, {
// 		ID:          167,
// 		DisplayName: "",
// 		ShortName:   "recipe_desolator",
// 	}, {
// 		ID:          168,
// 		DisplayName: "Desolator",
// 		ShortName:   "desolator",
// 	}, {
// 		ID:          169,
// 		DisplayName: "Yasha Recipe",
// 		ShortName:   "recipe_yasha",
// 	}, {
// 		ID:          170,
// 		DisplayName: "Yasha",
// 		ShortName:   "yasha",
// 	}, {
// 		ID:          171,
// 		DisplayName: "Mask of Madness Recipe",
// 		ShortName:   "recipe_mask_of_madness",
// 	}, {
// 		ID:          172,
// 		DisplayName: "Mask of Madness",
// 		ShortName:   "mask_of_madness",
// 	}, {
// 		ID:          173,
// 		DisplayName: "Diffusal Blade Recipe",
// 		ShortName:   "recipe_diffusal_blade",
// 	}, {
// 		ID:          174,
// 		DisplayName: "Diffusal Blade",
// 		ShortName:   "diffusal_blade",
// 	}, {
// 		ID:          175,
// 		DisplayName: "Ethereal Blade Recipe",
// 		ShortName:   "recipe_ethereal_blade",
// 	}, {
// 		ID:          176,
// 		DisplayName: "Ethereal Blade",
// 		ShortName:   "ethereal_blade",
// 	}, {
// 		ID:          177,
// 		DisplayName: "Soul Ring Recipe",
// 		ShortName:   "recipe_soul_ring",
// 	}, {
// 		ID:          178,
// 		DisplayName: "Soul Ring",
// 		ShortName:   "soul_ring",
// 	}, {
// 		ID:          179,
// 		DisplayName: "Arcane Boots Recipe",
// 		ShortName:   "recipe_arcane_boots",
// 	}, {
// 		ID:          180,
// 		DisplayName: "Arcane Boots",
// 		ShortName:   "arcane_boots",
// 	}, {
// 		ID:          181,
// 		DisplayName: "Orb of Venom",
// 		ShortName:   "orb_of_venom",
// 	}, {
// 		ID:          182,
// 		DisplayName: "Stout Shield",
// 		ShortName:   "stout_shield",
// 	}, {
// 		ID:          183,
// 		DisplayName: "",
// 		ShortName:   "recipe_invis_sword",
// 	}, {
// 		ID:          184,
// 		DisplayName: "Drum of Endurance Recipe",
// 		ShortName:   "recipe_ancient_janggo",
// 	}, {
// 		ID:          185,
// 		DisplayName: "Drum of Endurance",
// 		ShortName:   "ancient_janggo",
// 	}, {
// 		ID:          186,
// 		DisplayName: "",
// 		ShortName:   "recipe_medallion_of_courage",
// 	}, {
// 		ID:          187,
// 		DisplayName: "Medallion of Courage",
// 		ShortName:   "medallion_of_courage",
// 	}, {
// 		ID:          188,
// 		DisplayName: "Smoke of Deceit",
// 		ShortName:   "smoke_of_deceit",
// 	}, {
// 		ID:          189,
// 		DisplayName: "Veil of Discord Recipe",
// 		ShortName:   "recipe_veil_of_discord",
// 	}, {
// 		ID:          190,
// 		DisplayName: "Veil of Discord",
// 		ShortName:   "veil_of_discord",
// 	}, {
// 		ID:          191,
// 		DisplayName: "",
// 		ShortName:   "recipe_necronomicon_2",
// 	}, {
// 		ID:          192,
// 		DisplayName: "",
// 		ShortName:   "recipe_necronomicon_3",
// 	}, {
// 		ID:          193,
// 		DisplayName: "Necronomicon",
// 		ShortName:   "necronomicon_2",
// 	}, {
// 		ID:          194,
// 		DisplayName: "Necronomicon",
// 		ShortName:   "necronomicon_3",
// 	}, {
// 		ID:          195,
// 		DisplayName: "Recipe: Diffusal Blade (level 2)",
// 		ShortName:   "recipe_diffusal_blade_2",
// 	}, {
// 		ID:          196,
// 		DisplayName: "Diffusal Blade (level 2)",
// 		ShortName:   "diffusal_blade_2",
// 	}, {
// 		ID:          197,
// 		DisplayName: "",
// 		ShortName:   "recipe_dagon_2",
// 	}, {
// 		ID:          198,
// 		DisplayName: "",
// 		ShortName:   "recipe_dagon_3",
// 	}, {
// 		ID:          199,
// 		DisplayName: "",
// 		ShortName:   "recipe_dagon_4",
// 	}, {
// 		ID:          200,
// 		DisplayName: "",
// 		ShortName:   "recipe_dagon_5",
// 	}, {
// 		ID:          201,
// 		DisplayName: "Dagon",
// 		ShortName:   "dagon_2",
// 	}, {
// 		ID:          202,
// 		DisplayName: "Dagon",
// 		ShortName:   "dagon_3",
// 	}, {
// 		ID:          203,
// 		DisplayName: "Dagon",
// 		ShortName:   "dagon_4",
// 	}, {
// 		ID:          204,
// 		DisplayName: "Dagon",
// 		ShortName:   "dagon_5",
// 	}, {
// 		ID:          205,
// 		DisplayName: "Rod of Atos Recipe",
// 		ShortName:   "recipe_rod_of_atos",
// 	}, {
// 		ID:          206,
// 		DisplayName: "Rod of Atos",
// 		ShortName:   "rod_of_atos",
// 	}, {
// 		ID:          207,
// 		DisplayName: "Abyssal Blade Recipe",
// 		ShortName:   "recipe_abyssal_blade",
// 	}, {
// 		ID:          208,
// 		DisplayName: "Abyssal Blade",
// 		ShortName:   "abyssal_blade",
// 	}, {
// 		ID:          209,
// 		DisplayName: "Heaven's Halberd Recipe",
// 		ShortName:   "recipe_heavens_halberd",
// 	}, {
// 		ID:          210,
// 		DisplayName: "Heaven's Halberd",
// 		ShortName:   "heavens_halberd",
// 	}, {
// 		ID:          211,
// 		DisplayName: "",
// 		ShortName:   "recipe_ring_of_aquila",
// 	}, {
// 		ID:          212,
// 		DisplayName: "Ring of Aquila",
// 		ShortName:   "ring_of_aquila",
// 	}, {
// 		ID:          213,
// 		DisplayName: "Tranquil Boots Recipe",
// 		ShortName:   "recipe_tranquil_boots",
// 	}, {
// 		ID:          214,
// 		DisplayName: "Tranquil Boots",
// 		ShortName:   "tranquil_boots",
// 	}, {
// 		ID:          215,
// 		DisplayName: "Shadow Amulet",
// 		ShortName:   "shadow_amulet",
// 	}, {
// 		ID:          216,
// 		DisplayName: "Enchanted Mango",
// 		ShortName:   "enchanted_mango",
// 	}, {
// 		ID:          217,
// 		DisplayName: "",
// 		ShortName:   "recipe_ward_dispenser",
// 	}, {
// 		ID:          218,
// 		DisplayName: "Observer and Sentry Wards",
// 		ShortName:   "ward_dispenser",
// 	}, {
// 		ID:          219,
// 		DisplayName: "",
// 		ShortName:   "recipe_travel_boots_2",
// 	}, {
// 		ID:          220,
// 		DisplayName: "Boots of Travel 2",
// 		ShortName:   "travel_boots_2",
// 	}, {
// 		ID:          221,
// 		DisplayName: "Lotus Orb Recipe",
// 		ShortName:   "recipe_lotus_orb",
// 	}, {
// 		ID:          222,
// 		DisplayName: "Meteor Hammer Recipe",
// 		ShortName:   "recipe_meteor_hammer",
// 	}, {
// 		ID:          223,
// 		DisplayName: "Meteor Hammer",
// 		ShortName:   "meteor_hammer",
// 	}, {
// 		ID:          224,
// 		DisplayName: "Nullifier Recipe",
// 		ShortName:   "recipe_nullifier",
// 	}, {
// 		ID:          225,
// 		DisplayName: "Nullifier",
// 		ShortName:   "nullifier",
// 	}, {
// 		ID:          226,
// 		DisplayName: "Lotus Orb",
// 		ShortName:   "lotus_orb",
// 	}, {
// 		ID:          227,
// 		DisplayName: "Solar Crest Recipe",
// 		ShortName:   "recipe_solar_crest",
// 	}, {
// 		ID:          228,
// 		DisplayName: "Octarine Core Recipe",
// 		ShortName:   "recipe_octarine_core",
// 	}, {
// 		ID:          229,
// 		DisplayName: "Solar Crest",
// 		ShortName:   "solar_crest",
// 	}, {
// 		ID:          230,
// 		DisplayName: "Guardian Greaves Recipe",
// 		ShortName:   "recipe_guardian_greaves",
// 	}, {
// 		ID:          231,
// 		DisplayName: "Guardian Greaves",
// 		ShortName:   "guardian_greaves",
// 	}, {
// 		ID:          232,
// 		DisplayName: "Aether Lens",
// 		ShortName:   "aether_lens",
// 	}, {
// 		ID:          233,
// 		DisplayName: "Aether Lens Recipe",
// 		ShortName:   "recipe_aether_lens",
// 	}, {
// 		ID:          234,
// 		DisplayName: "Dragon Lance Recipe",
// 		ShortName:   "recipe_dragon_lance",
// 	}, {
// 		ID:          235,
// 		DisplayName: "Octarine Core",
// 		ShortName:   "octarine_core",
// 	}, {
// 		ID:          236,
// 		DisplayName: "Dragon Lance",
// 		ShortName:   "dragon_lance",
// 	}, {
// 		ID:          237,
// 		DisplayName: "Faerie Fire",
// 		ShortName:   "faerie_fire",
// 	}, {
// 		ID:          238,
// 		DisplayName: "Iron Talon Recipe",
// 		ShortName:   "recipe_iron_talon",
// 	}, {
// 		ID:          239,
// 		DisplayName: "Iron Talon",
// 		ShortName:   "iron_talon",
// 	}, {
// 		ID:          240,
// 		DisplayName: "Orb of Blight",
// 		ShortName:   "blight_stone",
// 	}, {
// 		ID:          241,
// 		DisplayName: "Tango (Shared)",
// 		ShortName:   "tango_single",
// 	}, {
// 		ID:          242,
// 		DisplayName: "Crimson Guard",
// 		ShortName:   "crimson_guard",
// 	}, {
// 		ID:          243,
// 		DisplayName: "Crimson Guard Recipe",
// 		ShortName:   "recipe_crimson_guard",
// 	}, {
// 		ID:          244,
// 		DisplayName: "Wind Lace",
// 		ShortName:   "wind_lace",
// 	}, {
// 		ID:          245,
// 		DisplayName: "Bloodthorn Recipe",
// 		ShortName:   "recipe_bloodthorn",
// 	}, {
// 		ID:          246,
// 		DisplayName: "",
// 		ShortName:   "recipe_moon_shard",
// 	}, {
// 		ID:          247,
// 		DisplayName: "Moon Shard",
// 		ShortName:   "moon_shard",
// 	}, {
// 		ID:          248,
// 		DisplayName: "Silver Edge Recipe",
// 		ShortName:   "recipe_silver_edge",
// 	}, {
// 		ID:          249,
// 		DisplayName: "Silver Edge",
// 		ShortName:   "silver_edge",
// 	}, {
// 		ID:          250,
// 		DisplayName: "Bloodthorn",
// 		ShortName:   "bloodthorn",
// 	}, {
// 		ID:          251,
// 		DisplayName: "",
// 		ShortName:   "recipe_echo_sabre",
// 	}, {
// 		ID:          252,
// 		DisplayName: "Echo Sabre",
// 		ShortName:   "echo_sabre",
// 	}, {
// 		ID:          253,
// 		DisplayName: "Glimmer Cape Recipe",
// 		ShortName:   "recipe_glimmer_cape",
// 	}, {
// 		ID:          254,
// 		DisplayName: "Glimmer Cape",
// 		ShortName:   "glimmer_cape",
// 	}, {
// 		ID:          255,
// 		DisplayName: "Aeon Disk Recipe",
// 		ShortName:   "recipe_aeon_disk",
// 	}, {
// 		ID:          256,
// 		DisplayName: "Aeon Disk",
// 		ShortName:   "aeon_disk",
// 	}, {
// 		ID:          257,
// 		DisplayName: "Tome of Knowledge",
// 		ShortName:   "tome_of_knowledge",
// 	}, {
// 		ID:          258,
// 		DisplayName: "Kaya Recipe",
// 		ShortName:   "recipe_kaya",
// 	}, {
// 		ID:          259,
// 		DisplayName: "Kaya",
// 		ShortName:   "kaya",
// 	}, {
// 		ID:          260,
// 		DisplayName: "Refresher Shard",
// 		ShortName:   "refresher_shard",
// 	}, {
// 		ID:          261,
// 		DisplayName: "Crown",
// 		ShortName:   "crown",
// 	}, {
// 		ID:          262,
// 		DisplayName: "Hurricane Pike Recipe",
// 		ShortName:   "recipe_hurricane_pike",
// 	}, {
// 		ID:          263,
// 		DisplayName: "Hurricane Pike",
// 		ShortName:   "hurricane_pike",
// 	}, {
// 		ID:          265,
// 		DisplayName: "Infused Raindrops",
// 		ShortName:   "infused_raindrop",
// 	}, {
// 		ID:          266,
// 		DisplayName: "Spirit Vessel Recipe",
// 		ShortName:   "recipe_spirit_vessel",
// 	}, {
// 		ID:          267,
// 		DisplayName: "Spirit Vessel",
// 		ShortName:   "spirit_vessel",
// 	}, {
// 		ID:          268,
// 		DisplayName: "Holy Locket Recipe",
// 		ShortName:   "recipe_holy_locket",
// 	}, {
// 		ID:          269,
// 		DisplayName: "Holy Locket",
// 		ShortName:   "holy_locket",
// 	}, {
// 		ID:          270,
// 		DisplayName: "Aghanim's Blessing Recipe",
// 		ShortName:   "recipe_ultimate_scepter_2",
// 	}, {
// 		ID:          271,
// 		DisplayName: "Aghanim's Blessing",
// 		ShortName:   "ultimate_scepter_2",
// 	}, {
// 		ID:          272,
// 		DisplayName: "Kaya and Sange Recipe",
// 		ShortName:   "recipe_kaya_and_sange",
// 	}, {
// 		ID:          273,
// 		DisplayName: "Kaya and Sange",
// 		ShortName:   "kaya_and_sange",
// 	}, {
// 		ID:          274,
// 		DisplayName: "Yasha and Kaya Recipe",
// 		ShortName:   "recipe_yasha_and_kaya",
// 	}, {
// 		ID:          275,
// 		DisplayName: "Trident Recipe",
// 		ShortName:   "recipe_trident",
// 	}, {
// 		ID:          276,
// 		DisplayName: "",
// 		ShortName:   "combo_breaker",
// 	}, {
// 		ID:          277,
// 		DisplayName: "Yasha and Kaya",
// 		ShortName:   "yasha_and_kaya",
// 	}, {
// 		ID:          279,
// 		DisplayName: "Ring of Tarrasque",
// 		ShortName:   "ring_of_tarrasque",
// 	}, {
// 		ID:          286,
// 		DisplayName: "Flying Courier",
// 		ShortName:   "flying_courier",
// 	}, {
// 		ID:          287,
// 		DisplayName: "Keen Optic",
// 		ShortName:   "keen_optic",
// 	}, {
// 		ID:          288,
// 		DisplayName: "Grove Bow",
// 		ShortName:   "grove_bow",
// 	}, {
// 		ID:          289,
// 		DisplayName: "Quickening Charm",
// 		ShortName:   "quickening_charm",
// 	}, {
// 		ID:          290,
// 		DisplayName: "Philosopher's Stone",
// 		ShortName:   "philosophers_stone",
// 	}, {
// 		ID:          291,
// 		DisplayName: "Force Boots",
// 		ShortName:   "force_boots",
// 	}, {
// 		ID:          292,
// 		DisplayName: "Stygian Desolator",
// 		ShortName:   "desolator_2",
// 	}, {
// 		ID:          293,
// 		DisplayName: "Phoenix Ash",
// 		ShortName:   "phoenix_ash",
// 	}, {
// 		ID:          294,
// 		DisplayName: "Seer Stone",
// 		ShortName:   "seer_stone",
// 	}, {
// 		ID:          295,
// 		DisplayName: "Greater Mango",
// 		ShortName:   "greater_mango",
// 	}, {
// 		ID:          297,
// 		DisplayName: "Vampire Fangs",
// 		ShortName:   "vampire_fangs",
// 	}, {
// 		ID:          298,
// 		DisplayName: "Craggy Coat",
// 		ShortName:   "craggy_coat",
// 	}, {
// 		ID:          299,
// 		DisplayName: "Greater Faerie Fire",
// 		ShortName:   "greater_faerie_fire",
// 	}, {
// 		ID:          300,
// 		DisplayName: "Timeless Relic",
// 		ShortName:   "timeless_relic",
// 	}, {
// 		ID:          301,
// 		DisplayName: "Mirror Shield",
// 		ShortName:   "mirror_shield",
// 	}, {
// 		ID:          302,
// 		DisplayName: "Elixir",
// 		ShortName:   "elixer",
// 	}, {
// 		ID:          303,
// 		DisplayName: "Ironwood Tree Recipe",
// 		ShortName:   "recipe_ironwood_tree",
// 	}, {
// 		ID:          304,
// 		DisplayName: "Ironwood Tree",
// 		ShortName:   "ironwood_tree",
// 	}, {
// 		ID:          305,
// 		DisplayName: "Royal Jelly",
// 		ShortName:   "royal_jelly",
// 	}, {
// 		ID:          306,
// 		DisplayName: "Pupil's Gift",
// 		ShortName:   "pupils_gift",
// 	}, {
// 		ID:          307,
// 		DisplayName: "Tome of Aghanim",
// 		ShortName:   "tome_of_aghanim",
// 	}, {
// 		ID:          308,
// 		DisplayName: "Repair Kit",
// 		ShortName:   "repair_kit",
// 	}, {
// 		ID:          309,
// 		DisplayName: "Mind Breaker",
// 		ShortName:   "mind_breaker",
// 	}, {
// 		ID:          310,
// 		DisplayName: "Third Eye",
// 		ShortName:   "third_eye",
// 	}, {
// 		ID:          311,
// 		DisplayName: "Spell Prism",
// 		ShortName:   "spell_prism",
// 	}, {
// 		ID:          312,
// 		DisplayName: "Horizon",
// 		ShortName:   "horizon",
// 	}, {
// 		ID:          313,
// 		DisplayName: "Fusion Rune",
// 		ShortName:   "fusion_rune",
// 	}, {
// 		ID:          317,
// 		DisplayName: "Recipe: Fallen Sky",
// 		ShortName:   "recipe_fallen_sky",
// 	}, {
// 		ID:          325,
// 		DisplayName: "Prince's Knife",
// 		ShortName:   "princes_knife",
// 	}, {
// 		ID:          326,
// 		DisplayName: "Spider Legs",
// 		ShortName:   "spider_legs",
// 	}, {
// 		ID:          327,
// 		DisplayName: "Helm of the Undying",
// 		ShortName:   "helm_of_the_undying",
// 	}, {
// 		ID:          328,
// 		DisplayName: "Mango Tree",
// 		ShortName:   "mango_tree",
// 	}, {
// 		ID:          329,
// 		DisplayName: "Vambrace Recipe",
// 		ShortName:   "recipe_vambrace",
// 	}, {
// 		ID:          330,
// 		DisplayName: "Witless Shako",
// 		ShortName:   "witless_shako",
// 	}, {
// 		ID:          331,
// 		DisplayName: "Vambrace",
// 		ShortName:   "vambrace",
// 	}, {
// 		ID:          334,
// 		DisplayName: "Imp Claw",
// 		ShortName:   "imp_claw",
// 	}, {
// 		ID:          335,
// 		DisplayName: "Flicker",
// 		ShortName:   "flicker",
// 	}, {
// 		ID:          336,
// 		DisplayName: "Telescope",
// 		ShortName:   "spy_gadget",
// 	}, {
// 		ID:          349,
// 		DisplayName: "Arcane Ring",
// 		ShortName:   "arcane_ring",
// 	}, {
// 		ID:          354,
// 		DisplayName: "Ocean Heart",
// 		ShortName:   "ocean_heart",
// 	}, {
// 		ID:          355,
// 		DisplayName: "Broom Handle",
// 		ShortName:   "broom_handle",
// 	}, {
// 		ID:          356,
// 		DisplayName: "Trusty Shovel",
// 		ShortName:   "trusty_shovel",
// 	}, {
// 		ID:          357,
// 		DisplayName: "Nether Shawl",
// 		ShortName:   "nether_shawl",
// 	}, {
// 		ID:          358,
// 		DisplayName: "Dragon Scale",
// 		ShortName:   "dragon_scale",
// 	}, {
// 		ID:          359,
// 		DisplayName: "Essence Ring",
// 		ShortName:   "essence_ring",
// 	}, {
// 		ID:          360,
// 		DisplayName: "Clumsy Net",
// 		ShortName:   "clumsy_net",
// 	}, {
// 		ID:          361,
// 		DisplayName: "Enchanted Quiver",
// 		ShortName:   "enchanted_quiver",
// 	}, {
// 		ID:          362,
// 		DisplayName: "Ninja Gear",
// 		ShortName:   "ninja_gear",
// 	}, {
// 		ID:          363,
// 		DisplayName: "Illusionist's Cape",
// 		ShortName:   "illusionsts_cape",
// 	}, {
// 		ID:          364,
// 		DisplayName: "Havoc Hammer",
// 		ShortName:   "havoc_hammer",
// 	}, {
// 		ID:          365,
// 		DisplayName: "Magic Lamp",
// 		ShortName:   "panic_button",
// 	}, {
// 		ID:          366,
// 		DisplayName: "Apex",
// 		ShortName:   "apex",
// 	}, {
// 		ID:          367,
// 		DisplayName: "Ballista",
// 		ShortName:   "ballista",
// 	}, {
// 		ID:          368,
// 		DisplayName: "Woodland Striders",
// 		ShortName:   "woodland_striders",
// 	}, {
// 		ID:          369,
// 		DisplayName: "Trident",
// 		ShortName:   "trident",
// 	}, {
// 		ID:          370,
// 		DisplayName: "Book of the Dead",
// 		ShortName:   "demonicon",
// 	}, {
// 		ID:          371,
// 		DisplayName: "Fallen Sky",
// 		ShortName:   "fallen_sky",
// 	}, {
// 		ID:          372,
// 		DisplayName: "Pirate Hat",
// 		ShortName:   "pirate_hat",
// 	}, {
// 		ID:          373,
// 		DisplayName: "Dimensional Doorway",
// 		ShortName:   "dimensional_doorway",
// 	}, {
// 		ID:          374,
// 		DisplayName: "Ex Machina",
// 		ShortName:   "ex_machina",
// 	}, {
// 		ID:          375,
// 		DisplayName: "Faded Broach",
// 		ShortName:   "faded_broach",
// 	}, {
// 		ID:          376,
// 		DisplayName: "Paladin Sword",
// 		ShortName:   "paladin_sword",
// 	}, {
// 		ID:          377,
// 		DisplayName: "Minotaur Horn",
// 		ShortName:   "minotaur_horn",
// 	}, {
// 		ID:          378,
// 		DisplayName: "Orb of Destruction",
// 		ShortName:   "orb_of_destruction",
// 	}, {
// 		ID:          379,
// 		DisplayName: "The Leveller",
// 		ShortName:   "the_leveller",
// 	}, {
// 		ID:          381,
// 		DisplayName: "Titan Sliver",
// 		ShortName:   "titan_sliver",
// 	}, {
// 		ID:          473,
// 		DisplayName: "Voodoo Mask",
// 		ShortName:   "voodoo_mask",
// 	}, {
// 		ID:          485,
// 		DisplayName: "Blitz Knuckles",
// 		ShortName:   "blitz_knuckles",
// 	}, {
// 		ID:          533,
// 		DisplayName: "Witch Blade Recipe",
// 		ShortName:   "recipe_witch_blade",
// 	}, {
// 		ID:          534,
// 		DisplayName: "Witch Blade",
// 		ShortName:   "witch_blade",
// 	}, {
// 		ID:          565,
// 		DisplayName: "Chipped Vest",
// 		ShortName:   "chipped_vest",
// 	}, {
// 		ID:          566,
// 		DisplayName: "Wizard Glass",
// 		ShortName:   "wizard_glass",
// 	}, {
// 		ID:          569,
// 		DisplayName: "Orb of Corrosion",
// 		ShortName:   "orb_of_corrosion",
// 	}, {
// 		ID:          570,
// 		DisplayName: "Gloves of Travel",
// 		ShortName:   "gloves_of_travel",
// 	}, {
// 		ID:          571,
// 		DisplayName: "Trickster Cloak",
// 		ShortName:   "trickster_cloak",
// 	}, {
// 		ID:          573,
// 		DisplayName: "Elven Tunic",
// 		ShortName:   "elven_tunic",
// 	}, {
// 		ID:          574,
// 		DisplayName: "Cloak of Flames",
// 		ShortName:   "cloak_of_flames",
// 	}, {
// 		ID:          575,
// 		DisplayName: "Venom Gland",
// 		ShortName:   "venom_gland",
// 	}, {
// 		ID:          576,
// 		DisplayName: "Helm of the Gladiator",
// 		ShortName:   "gladiator_helm",
// 	}, {
// 		ID:          577,
// 		DisplayName: "Possessed Mask",
// 		ShortName:   "possessed_mask",
// 	}, {
// 		ID:          578,
// 		DisplayName: "Ancient Perseverance",
// 		ShortName:   "ancient_perseverance",
// 	}, {
// 		ID:          582,
// 		DisplayName: "Oakheart",
// 		ShortName:   "oakheart",
// 	}, {
// 		ID:          585,
// 		DisplayName: "Stormcrafter",
// 		ShortName:   "stormcrafter",
// 	}, {
// 		ID:          588,
// 		DisplayName: "Overflowing Elixir",
// 		ShortName:   "overflowing_elixir",
// 	}, {
// 		ID:          589,
// 		DisplayName: "Fairy's Trinket",
// 		ShortName:   "mysterious_hat",
// 	}, {
// 		ID:          593,
// 		DisplayName: "Fluffy Hat",
// 		ShortName:   "fluffy_hat",
// 	}, {
// 		ID:          596,
// 		DisplayName: "Falcon Blade",
// 		ShortName:   "falcon_blade",
// 	}, {
// 		ID:          597,
// 		DisplayName: "Mage Slayer Recipe",
// 		ShortName:   "recipe_mage_slayer",
// 	}, {
// 		ID:          598,
// 		DisplayName: "Mage Slayer",
// 		ShortName:   "mage_slayer",
// 	}, {
// 		ID:          599,
// 		DisplayName: "Falcon Blade Recipe",
// 		ShortName:   "recipe_falcon_blade",
// 	}, {
// 		ID:          600,
// 		DisplayName: "Overwhelming Blink",
// 		ShortName:   "overwhelming_blink",
// 	}, {
// 		ID:          603,
// 		DisplayName: "Swift Blink",
// 		ShortName:   "swift_blink",
// 	}, {
// 		ID:          604,
// 		DisplayName: "Arcane Blink",
// 		ShortName:   "arcane_blink",
// 	}, {
// 		ID:          606,
// 		DisplayName: "Arcane Blink Recipe",
// 		ShortName:   "recipe_arcane_blink",
// 	}, {
// 		ID:          607,
// 		DisplayName: "Swift Blink Recipe",
// 		ShortName:   "recipe_swift_blink",
// 	}, {
// 		ID:          608,
// 		DisplayName: "Overwhelming Blink Recipe",
// 		ShortName:   "recipe_overwhelming_blink",
// 	}, {
// 		ID:          609,
// 		DisplayName: "Aghanim's Shard",
// 		ShortName:   "aghanims_shard",
// 	}, {
// 		ID:          610,
// 		DisplayName: "Wind Waker",
// 		ShortName:   "wind_waker",
// 	}, {
// 		ID:          612,
// 		DisplayName: "Wind Waker Recipe",
// 		ShortName:   "recipe_wind_waker",
// 	}, {
// 		ID:          633,
// 		DisplayName: "Helm of the Overlord Recipe",
// 		ShortName:   "recipe_helm_of_the_overlord",
// 	}, {
// 		ID:          635,
// 		DisplayName: "Helm of the Overlord",
// 		ShortName:   "helm_of_the_overlord",
// 	}, {
// 		ID:          637,
// 		DisplayName: "Star Mace",
// 		ShortName:   "star_mace",
// 	}, {
// 		ID:          638,
// 		DisplayName: "Penta-Edged Sword",
// 		ShortName:   "penta_edged_sword",
// 	}, {
// 		ID:          640,
// 		DisplayName: "Orb of Corrosion Recipe",
// 		ShortName:   "recipe_orb_of_corrosion",
// 	}, {
// 		ID:          653,
// 		DisplayName: "",
// 		ShortName:   "recipe_grandmasters_glaive",
// 	}, {
// 		ID:          655,
// 		DisplayName: "Grandmaster's Glaive",
// 		ShortName:   "grandmasters_glaive",
// 	}, {
// 		ID:          674,
// 		DisplayName: "Warhammer",
// 		ShortName:   "warhammer",
// 	}, {
// 		ID:          675,
// 		DisplayName: "Psychic Headband",
// 		ShortName:   "psychic_headband",
// 	}, {
// 		ID:          676,
// 		DisplayName: "Ceremonial Robe",
// 		ShortName:   "ceremonial_robe",
// 	}, {
// 		ID:          677,
// 		DisplayName: "Book of Shadows",
// 		ShortName:   "book_of_shadows",
// 	}, {
// 		ID:          678,
// 		DisplayName: "Giant's Ring",
// 		ShortName:   "giants_ring",
// 	}, {
// 		ID:          679,
// 		DisplayName: "Shadow of Vengeance",
// 		ShortName:   "vengeances_shadow",
// 	}, {
// 		ID:          680,
// 		DisplayName: "Bullwhip",
// 		ShortName:   "bullwhip",
// 	}, {
// 		ID:          686,
// 		DisplayName: "Quicksilver Amulet",
// 		ShortName:   "quicksilver_amulet",
// 	}, {
// 		ID:          691,
// 		DisplayName: "Eternal Shroud Recipe",
// 		ShortName:   "recipe_eternal_shroud",
// 	}, {
// 		ID:          692,
// 		DisplayName: "Eternal Shroud",
// 		ShortName:   "eternal_shroud",
// 	}, {
// 		ID:          725,
// 		DisplayName: "Aghanim's Shard - Consumable",
// 		ShortName:   "aghanims_shard_roshan",
// 	}, {
// 		ID:          727,
// 		DisplayName: "Aghanim's Blessing - Roshan",
// 		ShortName:   "ultimate_scepter_roshan",
// 	}, {
// 		ID:          731,
// 		DisplayName: "Satchel",
// 		ShortName:   "satchel",
// 	}, {
// 		ID:          824,
// 		DisplayName: "Assassin's Dagger",
// 		ShortName:   "assassins_dagger",
// 	}, {
// 		ID:          825,
// 		DisplayName: "Ascetic's Cap",
// 		ShortName:   "ascetic_cap",
// 	}, {
// 		ID:          826,
// 		DisplayName: "Assassin's Contract",
// 		ShortName:   "sample_picker",
// 	}, {
// 		ID:          827,
// 		DisplayName: "Icarus Wings",
// 		ShortName:   "icarus_wings",
// 	}, {
// 		ID:          828,
// 		DisplayName: "Brigand's Blade",
// 		ShortName:   "misericorde",
// 	}, {
// 		ID:          829,
// 		DisplayName: "Arcanist's Armor",
// 		ShortName:   "force_field",
// 	}, {
// 		ID:          833,
// 		DisplayName: "Bruiser's Maul Recipe",
// 		ShortName:   "recipe_tenderizer",
// 	}, {
// 		ID:          834,
// 		DisplayName: "Blast Rig",
// 		ShortName:   "black_powder_bag",
// 	}, {
// 		ID:          835,
// 		DisplayName: "Fae Grenade",
// 		ShortName:   "paintball",
// 	}, {
// 		ID:          836,
// 		DisplayName: "Light Robes",
// 		ShortName:   "light_robes",
// 	}, {
// 		ID:          837,
// 		DisplayName: "Witchbane",
// 		ShortName:   "heavy_blade",
// 	}, {
// 		ID:          838,
// 		DisplayName: "Pig Pole",
// 		ShortName:   "unstable_wand",
// 	}, {
// 		ID:          839,
// 		DisplayName: "Ring of Fortitude",
// 		ShortName:   "fortitude_ring",
// 	}, {
// 		ID:          840,
// 		DisplayName: "Tumbler's Toy",
// 		ShortName:   "pogo_stick",
// 	}, {
// 		ID:          849,
// 		DisplayName: "Mechanical Arm",
// 		ShortName:   "mechanical_arm",
// 	}, {
// 		ID:          859,
// 		DisplayName: "Voidwalker Scythe Recipe",
// 		ShortName:   "recipe_voidwalker_scythe",
// 	}, {
// 		ID:          904,
// 		DisplayName: "Voidwalker Scythe",
// 		ShortName:   "voidwalker_scythe",
// 	}, {
// 		ID:          906,
// 		DisplayName: "Bruiser's Maul",
// 		ShortName:   "tenderizer",
// 	}, {
// 		ID:          907,
// 		DisplayName: "Wraith Pact Recipe",
// 		ShortName:   "recipe_wraith_pact",
// 	}, {
// 		ID:          908,
// 		DisplayName: "Wraith Pact",
// 		ShortName:   "wraith_pact",
// 	}, {
// 		ID:          910,
// 		DisplayName: "Revenant's Brooch Recipe",
// 		ShortName:   "recipe_revenants_brooch",
// 	}, {
// 		ID:          911,
// 		DisplayName: "Revenant's Brooch",
// 		ShortName:   "revenants_brooch",
// 	}, {
// 		ID:          928,
// 		DisplayName: "",
// 		ShortName:   "recipe_eagle_eye",
// 	}, {
// 		ID:          929,
// 		DisplayName: "Eagle Eye",
// 		ShortName:   "eagle_eye",
// 	}, {
// 		ID:          930,
// 		DisplayName: "Boots of Bearing Recipe",
// 		ShortName:   "recipe_boots_of_bearing",
// 	}, {
// 		ID:          931,
// 		DisplayName: "Boots of Bearing",
// 		ShortName:   "boots_of_bearing",
// 	}, {
// 		ID:          938,
// 		DisplayName: "",
// 		ShortName:   "slime_vial",
// 	}, {
// 		ID:          939,
// 		DisplayName: "Harpoon",
// 		ShortName:   "harpoon",
// 	}, {
// 		ID:          940,
// 		DisplayName: "Wand of the Brine",
// 		ShortName:   "wand_of_the_brine",
// 	}, {
// 		ID:          945,
// 		DisplayName: "Seeds of Serenity",
// 		ShortName:   "seeds_of_serenity",
// 	}, {
// 		ID:          946,
// 		DisplayName: "Lance of Pursuit",
// 		ShortName:   "lance_of_pursuit",
// 	}, {
// 		ID:          947,
// 		DisplayName: "Occult Bracelet",
// 		ShortName:   "occult_bracelet",
// 	}, {
// 		ID:          948,
// 		DisplayName: "",
// 		ShortName:   "tome_of_omniscience",
// 	}, {
// 		ID:          949,
// 		DisplayName: "Ogre Seal Totem",
// 		ShortName:   "ogre_seal_totem",
// 	}, {
// 		ID:          950,
// 		DisplayName: "Defiant Shell",
// 		ShortName:   "defiant_shell",
// 	}, {
// 		ID:          964,
// 		DisplayName: "Diffusal Blade",
// 		ShortName:   "diffusal_blade_2",
// 	}, {
// 		ID:          965,
// 		DisplayName: "",
// 		ShortName:   "recipe_diffusal_blade_2",
// 	}, {
// 		ID:          968,
// 		DisplayName: "",
// 		ShortName:   "arcane_scout",
// 	}, {
// 		ID:          969,
// 		DisplayName: "",
// 		ShortName:   "barricade",
// 	}, {
// 		ID:          990,
// 		DisplayName: "Eye of the Vizier",
// 		ShortName:   "eye_of_the_vizier",
// 	}, {
// 		ID:          998,
// 		DisplayName: "",
// 		ShortName:   "manacles_of_power",
// 	}, {
// 		ID:          1000,
// 		DisplayName: "",
// 		ShortName:   "bottomless_chalice",
// 	}, {
// 		ID:          1017,
// 		DisplayName: "",
// 		ShortName:   "wand_of_sanctitude",
// 	}, {
// 		ID:          1021,
// 		DisplayName: "River Vial: Chrome",
// 		ShortName:   "river_painter",
// 	}, {
// 		ID:          1022,
// 		DisplayName: "River Vial: Dry",
// 		ShortName:   "river_painter2",
// 	}, {
// 		ID:          1023,
// 		DisplayName: "River Vial: Slime",
// 		ShortName:   "river_painter3",
// 	}, {
// 		ID:          1024,
// 		DisplayName: "River Vial: Oil",
// 		ShortName:   "river_painter4",
// 	}, {
// 		ID:          1025,
// 		DisplayName: "River Vial: Electrified",
// 		ShortName:   "river_painter5",
// 	}, {
// 		ID:          1026,
// 		DisplayName: "River Vial: Potion",
// 		ShortName:   "river_painter6",
// 	}, {
// 		ID:          1027,
// 		DisplayName: "River Vial: Blood",
// 		ShortName:   "river_painter7",
// 	}, {
// 		ID:          1028,
// 		DisplayName: "Tombstone",
// 		ShortName:   "mutation_tombstone",
// 	}, {
// 		ID:          1029,
// 		DisplayName: "Super Blink Dagger",
// 		ShortName:   "super_blink",
// 	}, {
// 		ID:          1030,
// 		DisplayName: "Pocket Tower",
// 		ShortName:   "pocket_tower",
// 	}, {
// 		ID:          1032,
// 		DisplayName: "Pocket Roshan",
// 		ShortName:   "pocket_roshan",
// 	}, {
// 		ID:          1076,
// 		DisplayName: "Specialist's Array",
// 		ShortName:   "specialists_array",
// 	}, {
// 		ID:          1077,
// 		DisplayName: "Dagger of Ristul",
// 		ShortName:   "dagger_of_ristul",
// 	}, {
// 		ID:          1090,
// 		DisplayName: "Mercy & Grace",
// 		ShortName:   "muertas_gun",
// 	}, {
// 		ID:          1091,
// 		DisplayName: "Samurai Tabi",
// 		ShortName:   "samurai_tabi",
// 	}, {
// 		ID:          1092,
// 		DisplayName: "Hermes Sandals Recipe",
// 		ShortName:   "recipe_hermes_sandals",
// 	}, {
// 		ID:          1093,
// 		DisplayName: "Hermes Sandals",
// 		ShortName:   "hermes_sandals",
// 	}, {
// 		ID:          1094,
// 		DisplayName: "Lunar Crest Recipe",
// 		ShortName:   "recipe_lunar_crest",
// 	}, {
// 		ID:          1095,
// 		DisplayName: "Lunar Crest",
// 		ShortName:   "lunar_crest",
// 	}, {
// 		ID:          1096,
// 		DisplayName: "Disperser Recipe",
// 		ShortName:   "recipe_disperser",
// 	}, {
// 		ID:          1097,
// 		DisplayName: "Disperser",
// 		ShortName:   "disperser",
// 	}, {
// 		ID:          1098,
// 		DisplayName: "Samurai Tabi Recipe",
// 		ShortName:   "recipe_samurai_tabi",
// 	}, {
// 		ID:          1099,
// 		DisplayName: "Witches Switch Recipe",
// 		ShortName:   "recipe_witches_switch",
// 	}, {
// 		ID:          1100,
// 		DisplayName: "Witches Switch",
// 		ShortName:   "witches_switch",
// 	}, {
// 		ID:          1101,
// 		DisplayName: "Harpoon Recipe",
// 		ShortName:   "recipe_harpoon",
// 	}, {
// 		ID:          1106,
// 		DisplayName: "Phylactery Recipe",
// 		ShortName:   "recipe_phylactery",
// 	}, {
// 		ID:          1107,
// 		DisplayName: "Phylactery",
// 		ShortName:   "phylactery",
// 	}, {
// 		ID:          1122,
// 		DisplayName: "Diadem",
// 		ShortName:   "diadem",
// 	}, {
// 		ID:          1123,
// 		DisplayName: "Blood Grenade",
// 		ShortName:   "blood_grenade",
// 	}, {
// 		ID:          1124,
// 		DisplayName: "Spark of Courage",
// 		ShortName:   "spark_of_courage",
// 	}, {
// 		ID:          1125,
// 		DisplayName: "Cornucopia",
// 		ShortName:   "cornucopia",
// 	}, {
// 		ID:          1127,
// 		DisplayName: "Pavise Recipe",
// 		ShortName:   "recipe_pavise",
// 	}, {
// 		ID:          1128,
// 		DisplayName: "Pavise",
// 		ShortName:   "pavise",
// 	}, {
// 		ID:          1154,
// 		DisplayName: "Block of Cheese",
// 		ShortName:   "royale_with_cheese",
// 	}, {
// 		ID:          1156,
// 		DisplayName: "Ancient Guardian",
// 		ShortName:   "ancient_guardian",
// 	}, {
// 		ID:          1157,
// 		DisplayName: "Safety Bubble",
// 		ShortName:   "safety_bubble",
// 	}, {
// 		ID:          1158,
// 		DisplayName: "Whisper of the Dread",
// 		ShortName:   "whisper_of_the_dread",
// 	}, {
// 		ID:          1159,
// 		DisplayName: "Nemesis Curse",
// 		ShortName:   "nemesis_curse",
// 	}, {
// 		ID:          1160,
// 		DisplayName: "Aviana's Feather",
// 		ShortName:   "avianas_feather",
// 	}, {
// 		ID:          1161,
// 		DisplayName: "Unwavering Condition",
// 		ShortName:   "unwavering_condition",
// 	}, {
// 		ID:          1164,
// 		DisplayName: "Aetherial Hammer",
// 		ShortName:   "aetherial_halo",
// 	}, {
// 		ID:          1167,
// 		DisplayName: "Light Collector",
// 		ShortName:   "light_collector",
// 	}, {
// 		ID:          1168,
// 		DisplayName: "Rattlecage",
// 		ShortName:   "rattlecage",
// 	}, {
// 		ID:          1440,
// 		DisplayName: "Black Grimoire",
// 		ShortName:   "black_grimoire",
// 	}, {
// 		ID:          1441,
// 		DisplayName: "Gris-Gris",
// 		ShortName:   "grisgris",
// 	}, {
// 		ID:          1466,
// 		DisplayName: "Gleipnir",
// 		ShortName:   "gungir",
// 	}, {
// 		ID:          1565,
// 		DisplayName: "Gleipnir Recipe",
// 		ShortName:   "recipe_gungir",
// 	}, {
// 		ID:          1575,
// 		DisplayName: "Orb of Frost",
// 		ShortName:   "orb_of_frost",
// 	}, {
// 		ID:          1576,
// 		DisplayName: "Vast",
// 		ShortName:   "enhancement_vast",
// 	}, {
// 		ID:          1577,
// 		DisplayName: "Quickened",
// 		ShortName:   "enhancement_quickened",
// 	}, {
// 		ID:          1578,
// 		DisplayName: "Accursed",
// 		ShortName:   "cursed_circlet",
// 	}, {
// 		ID:          1579,
// 		DisplayName: "Restorative",
// 		ShortName:   "ogre_heart",
// 	}, {
// 		ID:          1580,
// 		DisplayName: "Elusive",
// 		ShortName:   "neutral_tabi",
// 	}, {
// 		ID:          1581,
// 		DisplayName: "Audacious",
// 		ShortName:   "enhancement_audacious",
// 	}, {
// 		ID:          1582,
// 		DisplayName: "",
// 		ShortName:   "hellbear_totem",
// 	}, {
// 		ID:          1583,
// 		DisplayName: "Mystical",
// 		ShortName:   "enhancement_mystical",
// 	}, {
// 		ID:          1584,
// 		DisplayName: "Alert",
// 		ShortName:   "enhancement_alert",
// 	}, {
// 		ID:          1585,
// 		DisplayName: "Brawny",
// 		ShortName:   "enhancement_brawny",
// 	}, {
// 		ID:          1586,
// 		DisplayName: "Tough",
// 		ShortName:   "enhancement_tough",
// 	}, {
// 		ID:          1587,
// 		DisplayName: "Feverish",
// 		ShortName:   "enhancement_feverish",
// 	}, {
// 		ID:          1588,
// 		DisplayName: "Fleetfooted",
// 		ShortName:   "enhancement_fleetfooted",
// 	}, {
// 		ID:          1589,
// 		DisplayName: "Crude",
// 		ShortName:   "enhancement_crude",
// 	}, {
// 		ID:          1590,
// 		DisplayName: "Boundless",
// 		ShortName:   "enhancement_boundless",
// 	}, {
// 		ID:          1591,
// 		DisplayName: "Wise",
// 		ShortName:   "enhancement_wise",
// 	}, {
// 		ID:          1592,
// 		DisplayName: "Timeless",
// 		ShortName:   "enhancement_timeless",
// 	}, {
// 		ID:          1593,
// 		DisplayName: "Greedy",
// 		ShortName:   "enhancement_greedy",
// 	}, {
// 		ID:          1594,
// 		DisplayName: "Vampiric",
// 		ShortName:   "enhancement_vampiric",
// 	}, {
// 		ID:          1595,
// 		DisplayName: "Keen-eyed",
// 		ShortName:   "enhancement_keen_eyed",
// 	}, {
// 		ID:          1596,
// 		DisplayName: "Evolved",
// 		ShortName:   "enhancement_evolved",
// 	}, {
// 		ID:          1597,
// 		DisplayName: "Titanic",
// 		ShortName:   "enhancement_titanic",
// 	}, {
// 		ID:          1598,
// 		DisplayName: "Unrelenting Eye",
// 		ShortName:   "unrelenting_eye",
// 	}, {
// 		ID:          1599,
// 		DisplayName: "Mana Draught",
// 		ShortName:   "mana_draught",
// 	}, {
// 		ID:          1600,
// 		DisplayName: "Ripper's Lash",
// 		ShortName:   "rippers_lash",
// 	}, {
// 		ID:          1601,
// 		DisplayName: "Crippling Crossbow",
// 		ShortName:   "crippling_crossbow",
// 	}, {
// 		ID:          1602,
// 		DisplayName: "Gale Guard",
// 		ShortName:   "gale_guard",
// 	}, {
// 		ID:          1603,
// 		DisplayName: "Gunpowder Gauntlet",
// 		ShortName:   "gunpowder_gauntlets",
// 	}, {
// 		ID:          1604,
// 		DisplayName: "Searing Signet",
// 		ShortName:   "searing_signet",
// 	}, {
// 		ID:          1605,
// 		DisplayName: "Serrated Shiv",
// 		ShortName:   "serrated_shiv",
// 	}, {
// 		ID:          1606,
// 		DisplayName: "Pollywog Charm",
// 		ShortName:   "polliwog_charm",
// 	}, {
// 		ID:          1607,
// 		DisplayName: "Magnifying Monocle",
// 		ShortName:   "magnifying_monocle",
// 	}, {
// 		ID:          1608,
// 		DisplayName: "Pyrrhic Cloak",
// 		ShortName:   "pyrrhic_cloak",
// 	}, {
// 		ID:          1609,
// 		DisplayName: "Madstone Bundle",
// 		ShortName:   "madstone_bundle",
// 	}, {
// 		ID:          1610,
// 		DisplayName: "",
// 		ShortName:   "miniboss_minion_summoner",
// 	}, {
// 		ID:          1801,
// 		DisplayName: "Caster Rapier",
// 		ShortName:   "caster_rapier",
// 	}, {
// 		ID:          1802,
// 		DisplayName: "Tiara of Selemene",
// 		ShortName:   "tiara_of_selemene",
// 	}, {
// 		ID:          1803,
// 		DisplayName: "Doubloon",
// 		ShortName:   "doubloon",
// 	}, {
// 		ID:          1804,
// 		DisplayName: "Roshan's Banner",
// 		ShortName:   "roshans_banner",
// 	}, {
// 		ID:          1805,
// 		DisplayName: "Parasma Recipe",
// 		ShortName:   "recipe_devastator",
// 	}, {
// 		ID:          1806,
// 		DisplayName: "Parasma",
// 		ShortName:   "devastator",
// 	}, {
// 		ID:          1807,
// 		DisplayName: "Khanda Recipe",
// 		ShortName:   "recipe_angels_demise",
// 	}, {
// 		ID:          1808,
// 		DisplayName: "Khanda",
// 		ShortName:   "angels_demise",
// 	}, {
// 		ID:          2091,
// 		DisplayName: "Tier 1 Token",
// 		ShortName:   "tier1_token",
// 	}, {
// 		ID:          2092,
// 		DisplayName: "Tier 2 Token",
// 		ShortName:   "tier2_token",
// 	}, {
// 		ID:          2093,
// 		DisplayName: "Tier 3 Token",
// 		ShortName:   "tier3_token",
// 	}, {
// 		ID:          2094,
// 		DisplayName: "Tier 4 Token",
// 		ShortName:   "tier4_token",
// 	}, {
// 		ID:          2095,
// 		DisplayName: "Tier 5 Token",
// 		ShortName:   "tier5_token",
// 	}, {
// 		ID:          2096,
// 		DisplayName: "Vindicator's Axe",
// 		ShortName:   "vindicators_axe",
// 	}, {
// 		ID:          2097,
// 		DisplayName: "Duelist Gloves",
// 		ShortName:   "duelist_gloves",
// 	}, {
// 		ID:          2098,
// 		DisplayName: "Horizon's Equilibrium",
// 		ShortName:   "horizons_equilibrium",
// 	}, {
// 		ID:          2099,
// 		DisplayName: "Blighted Spirit",
// 		ShortName:   "blighted_spirit",
// 	}, {
// 		ID:          2190,
// 		DisplayName: "Dandelion Amulet",
// 		ShortName:   "dandelion_amulet",
// 	}, {
// 		ID:          2191,
// 		DisplayName: "Turtle Shell",
// 		ShortName:   "turtle_shell",
// 	}, {
// 		ID:          2192,
// 		DisplayName: "Martyr's Plate",
// 		ShortName:   "martyrs_plate",
// 	}, {
// 		ID:          2193,
// 		DisplayName: "Gossamer Cape",
// 		ShortName:   "gossamer_cape",
// 	}, {
// 		ID:          4204,
// 		DisplayName: "Healing Lotus",
// 		ShortName:   "famango",
// 	}, {
// 		ID:          4205,
// 		DisplayName: "Great Healing Lotus",
// 		ShortName:   "great_famango",
// 	}, {
// 		ID:          4206,
// 		DisplayName: "Greater Healing Lotus",
// 		ShortName:   "greater_famango",
// 	}, {
// 		ID:          4207,
// 		DisplayName: "",
// 		ShortName:   "recipe_great_famango",
// 	}, {
// 		ID:          4208,
// 		DisplayName: "",
// 		ShortName:   "recipe_greater_famango",
// 	}, {
// 		ID:          4300,
// 		DisplayName: "Beloved Memory",
// 		ShortName:   "ofrenda",
// 	}, {
// 		ID:          4301,
// 		DisplayName: "Scrying Shovel",
// 		ShortName:   "ofrenda_shovel",
// 	}, {
// 		ID:          4302,
// 		DisplayName: "Forebearer's Fortune",
// 		ShortName:   "ofrenda_pledge",
// 	},
// }
