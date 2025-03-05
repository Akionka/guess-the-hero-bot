package main

// var heroes = []Hero{
// 	{
// 		ID:          1,
// 		DisplayName: "Anti-Mage",
// 		ShortName:   "antimage",
// 	},
// 	{
// 		ID:          2,
// 		DisplayName: "Axe",
// 		ShortName:   "axe",
// 	}, {
// 		ID:          3,
// 		DisplayName: "Bane",
// 		ShortName:   "bane",
// 	}, {
// 		ID:          4,
// 		DisplayName: "Bloodseeker",
// 		ShortName:   "bloodseeker",
// 	}, {
// 		ID:          5,
// 		DisplayName: "Crystal Maiden",
// 		ShortName:   "crystal_maiden",
// 	}, {
// 		ID:          6,
// 		DisplayName: "Drow Ranger",
// 		ShortName:   "drow_ranger",
// 	}, {
// 		ID:          7,
// 		DisplayName: "Earthshaker",
// 		ShortName:   "earthshaker",
// 	}, {
// 		ID:          8,
// 		DisplayName: "Juggernaut",
// 		ShortName:   "juggernaut",
// 	}, {
// 		ID:          9,
// 		DisplayName: "Mirana",
// 		ShortName:   "mirana",
// 	}, {
// 		ID:          10,
// 		DisplayName: "Morphling",
// 		ShortName:   "morphling",
// 	}, {
// 		ID:          11,
// 		DisplayName: "Shadow Fiend",
// 		ShortName:   "nevermore",
// 	}, {
// 		ID:          12,
// 		DisplayName: "Phantom Lancer",
// 		ShortName:   "phantom_lancer",
// 	}, {
// 		ID:          13,
// 		DisplayName: "Puck",
// 		ShortName:   "puck",
// 	}, {
// 		ID:          14,
// 		DisplayName: "Pudge",
// 		ShortName:   "pudge",
// 	}, {
// 		ID:          15,
// 		DisplayName: "Razor",
// 		ShortName:   "razor",
// 	}, {
// 		ID:          16,
// 		DisplayName: "Sand King",
// 		ShortName:   "sand_king",
// 	}, {
// 		ID:          17,
// 		DisplayName: "Storm Spirit",
// 		ShortName:   "storm_spirit",
// 	}, {
// 		ID:          18,
// 		DisplayName: "Sven",
// 		ShortName:   "sven",
// 	}, {
// 		ID:          19,
// 		DisplayName: "Tiny",
// 		ShortName:   "tiny",
// 	}, {
// 		ID:          20,
// 		DisplayName: "Vengeful Spirit",
// 		ShortName:   "vengefulspirit",
// 	}, {
// 		ID:          21,
// 		DisplayName: "Windranger",
// 		ShortName:   "windrunner",
// 	}, {
// 		ID:          22,
// 		DisplayName: "Zeus",
// 		ShortName:   "zeus",
// 	}, {
// 		ID:          23,
// 		DisplayName: "Kunkka",
// 		ShortName:   "kunkka",
// 	}, {
// 		ID:          25,
// 		DisplayName: "Lina",
// 		ShortName:   "lina",
// 	}, {
// 		ID:          26,
// 		DisplayName: "Lion",
// 		ShortName:   "lion",
// 	}, {
// 		ID:          27,
// 		DisplayName: "Shadow Shaman",
// 		ShortName:   "shadow_shaman",
// 	}, {
// 		ID:          28,
// 		DisplayName: "Slardar",
// 		ShortName:   "slardar",
// 	}, {
// 		ID:          29,
// 		DisplayName: "Tidehunter",
// 		ShortName:   "tidehunter",
// 	}, {
// 		ID:          30,
// 		DisplayName: "Witch Doctor",
// 		ShortName:   "witch_doctor",
// 	}, {
// 		ID:          31,
// 		DisplayName: "Lich",
// 		ShortName:   "lich",
// 	}, {
// 		ID:          32,
// 		DisplayName: "Riki",
// 		ShortName:   "riki",
// 	}, {
// 		ID:          33,
// 		DisplayName: "Enigma",
// 		ShortName:   "enigma",
// 	}, {
// 		ID:          34,
// 		DisplayName: "Tinker",
// 		ShortName:   "tinker",
// 	}, {
// 		ID:          35,
// 		DisplayName: "Sniper",
// 		ShortName:   "sniper",
// 	}, {
// 		ID:          36,
// 		DisplayName: "Necrophos",
// 		ShortName:   "necrolyte",
// 	}, {
// 		ID:          37,
// 		DisplayName: "Warlock",
// 		ShortName:   "warlock",
// 	}, {
// 		ID:          38,
// 		DisplayName: "Beastmaster",
// 		ShortName:   "beastmaster",
// 	}, {
// 		ID:          39,
// 		DisplayName: "Queen of Pain",
// 		ShortName:   "queenofpain",
// 	}, {
// 		ID:          40,
// 		DisplayName: "Venomancer",
// 		ShortName:   "venomancer",
// 	}, {
// 		ID:          41,
// 		DisplayName: "Faceless Void",
// 		ShortName:   "faceless_void",
// 	}, {
// 		ID:          42,
// 		DisplayName: "Wraith King",
// 		ShortName:   "skeleton_king",
// 	}, {
// 		ID:          43,
// 		DisplayName: "Death Prophet",
// 		ShortName:   "death_prophet",
// 	}, {
// 		ID:          44,
// 		DisplayName: "Phantom Assassin",
// 		ShortName:   "phantom_assassin",
// 	}, {
// 		ID:          45,
// 		DisplayName: "Pugna",
// 		ShortName:   "pugna",
// 	}, {
// 		ID:          46,
// 		DisplayName: "Templar Assassin",
// 		ShortName:   "templar_assassin",
// 	}, {
// 		ID:          47,
// 		DisplayName: "Viper",
// 		ShortName:   "viper",
// 	}, {
// 		ID:          48,
// 		DisplayName: "Luna",
// 		ShortName:   "luna",
// 	}, {
// 		ID:          49,
// 		DisplayName: "Dragon Knight",
// 		ShortName:   "dragon_knight",
// 	}, {
// 		ID:          50,
// 		DisplayName: "Dazzle",
// 		ShortName:   "dazzle",
// 	}, {
// 		ID:          51,
// 		DisplayName: "Clockwerk",
// 		ShortName:   "rattletrap",
// 	}, {
// 		ID:          52,
// 		DisplayName: "Leshrac",
// 		ShortName:   "leshrac",
// 	}, {
// 		ID:          53,
// 		DisplayName: "Nature's Prophet",
// 		ShortName:   "furion",
// 	}, {
// 		ID:          54,
// 		DisplayName: "Lifestealer",
// 		ShortName:   "life_stealer",
// 	}, {
// 		ID:          55,
// 		DisplayName: "Dark Seer",
// 		ShortName:   "dark_seer",
// 	}, {
// 		ID:          56,
// 		DisplayName: "Clinkz",
// 		ShortName:   "clinkz",
// 	}, {
// 		ID:          57,
// 		DisplayName: "Omniknight",
// 		ShortName:   "omniknight",
// 	}, {
// 		ID:          58,
// 		DisplayName: "Enchantress",
// 		ShortName:   "enchantress",
// 	}, {
// 		ID:          59,
// 		DisplayName: "Huskar",
// 		ShortName:   "huskar",
// 	}, {
// 		ID:          60,
// 		DisplayName: "Night Stalker",
// 		ShortName:   "night_stalker",
// 	}, {
// 		ID:          61,
// 		DisplayName: "Broodmother",
// 		ShortName:   "broodmother",
// 	}, {
// 		ID:          62,
// 		DisplayName: "Bounty Hunter",
// 		ShortName:   "bounty_hunter",
// 	}, {
// 		ID:          63,
// 		DisplayName: "Weaver",
// 		ShortName:   "weaver",
// 	}, {
// 		ID:          64,
// 		DisplayName: "Jakiro",
// 		ShortName:   "jakiro",
// 	}, {
// 		ID:          65,
// 		DisplayName: "Batrider",
// 		ShortName:   "batrider",
// 	}, {
// 		ID:          66,
// 		DisplayName: "Chen",
// 		ShortName:   "chen",
// 	}, {
// 		ID:          67,
// 		DisplayName: "Spectre",
// 		ShortName:   "spectre",
// 	}, {
// 		ID:          68,
// 		DisplayName: "Ancient Apparition",
// 		ShortName:   "ancient_apparition",
// 	}, {
// 		ID:          69,
// 		DisplayName: "Doom",
// 		ShortName:   "doom_bringer",
// 	}, {
// 		ID:          70,
// 		DisplayName: "Ursa",
// 		ShortName:   "ursa",
// 	}, {
// 		ID:          71,
// 		DisplayName: "Spirit Breaker",
// 		ShortName:   "spirit_breaker",
// 	}, {
// 		ID:          72,
// 		DisplayName: "Gyrocopter",
// 		ShortName:   "gyrocopter",
// 	}, {
// 		ID:          73,
// 		DisplayName: "Alchemist",
// 		ShortName:   "alchemist",
// 	}, {
// 		ID:          74,
// 		DisplayName: "Invoker",
// 		ShortName:   "invoker",
// 	}, {
// 		ID:          75,
// 		DisplayName: "Silencer",
// 		ShortName:   "silencer",
// 	}, {
// 		ID:          76,
// 		DisplayName: "Outworld Destroyer",
// 		ShortName:   "obsidian_destroyer",
// 	}, {
// 		ID:          77,
// 		DisplayName: "Lycan",
// 		ShortName:   "lycan",
// 	}, {
// 		ID:          78,
// 		DisplayName: "Brewmaster",
// 		ShortName:   "brewmaster",
// 	}, {
// 		ID:          79,
// 		DisplayName: "Shadow Demon",
// 		ShortName:   "shadow_demon",
// 	}, {
// 		ID:          80,
// 		DisplayName: "Lone Druid",
// 		ShortName:   "lone_druid",
// 	}, {
// 		ID:          81,
// 		DisplayName: "Chaos Knight",
// 		ShortName:   "chaos_knight",
// 	}, {
// 		ID:          82,
// 		DisplayName: "Meepo",
// 		ShortName:   "meepo",
// 	}, {
// 		ID:          83,
// 		DisplayName: "Treant Protector",
// 		ShortName:   "treant",
// 	}, {
// 		ID:          84,
// 		DisplayName: "Ogre Magi",
// 		ShortName:   "ogre_magi",
// 	}, {
// 		ID:          85,
// 		DisplayName: "Undying",
// 		ShortName:   "undying",
// 	}, {
// 		ID:          86,
// 		DisplayName: "Rubick",
// 		ShortName:   "rubick",
// 	}, {
// 		ID:          87,
// 		DisplayName: "Disruptor",
// 		ShortName:   "disruptor",
// 	}, {
// 		ID:          88,
// 		DisplayName: "Nyx Assassin",
// 		ShortName:   "nyx_assassin",
// 	}, {
// 		ID:          89,
// 		DisplayName: "Naga Siren",
// 		ShortName:   "naga_siren",
// 	}, {
// 		ID:          90,
// 		DisplayName: "Keeper of the Light",
// 		ShortName:   "keeper_of_the_light",
// 	}, {
// 		ID:          91,
// 		DisplayName: "Io",
// 		ShortName:   "wisp",
// 	}, {
// 		ID:          92,
// 		DisplayName: "Visage",
// 		ShortName:   "visage",
// 	}, {
// 		ID:          93,
// 		DisplayName: "Slark",
// 		ShortName:   "slark",
// 	}, {
// 		ID:          94,
// 		DisplayName: "Medusa",
// 		ShortName:   "medusa",
// 	}, {
// 		ID:          95,
// 		DisplayName: "Troll Warlord",
// 		ShortName:   "troll_warlord",
// 	}, {
// 		ID:          96,
// 		DisplayName: "Centaur Warrunner",
// 		ShortName:   "centaur",
// 	}, {
// 		ID:          97,
// 		DisplayName: "Magnus",
// 		ShortName:   "magnataur",
// 	}, {
// 		ID:          98,
// 		DisplayName: "Timbersaw",
// 		ShortName:   "shredder",
// 	}, {
// 		ID:          99,
// 		DisplayName: "Bristleback",
// 		ShortName:   "bristleback",
// 	}, {
// 		ID:          100,
// 		DisplayName: "Tusk",
// 		ShortName:   "tusk",
// 	}, {
// 		ID:          101,
// 		DisplayName: "Skywrath Mage",
// 		ShortName:   "skywrath_mage",
// 	}, {
// 		ID:          102,
// 		DisplayName: "Abaddon",
// 		ShortName:   "abaddon",
// 	}, {
// 		ID:          103,
// 		DisplayName: "Elder Titan",
// 		ShortName:   "elder_titan",
// 	}, {
// 		ID:          104,
// 		DisplayName: "Legion Commander",
// 		ShortName:   "legion_commander",
// 	}, {
// 		ID:          105,
// 		DisplayName: "Techies",
// 		ShortName:   "techies",
// 	}, {
// 		ID:          106,
// 		DisplayName: "Ember Spirit",
// 		ShortName:   "ember_spirit",
// 	}, {
// 		ID:          107,
// 		DisplayName: "Earth Spirit",
// 		ShortName:   "earth_spirit",
// 	}, {
// 		ID:          108,
// 		DisplayName: "Underlord",
// 		ShortName:   "abyssal_underlord",
// 	}, {
// 		ID:          109,
// 		DisplayName: "Terrorblade",
// 		ShortName:   "terrorblade",
// 	}, {
// 		ID:          110,
// 		DisplayName: "Phoenix",
// 		ShortName:   "phoenix",
// 	}, {
// 		ID:          111,
// 		DisplayName: "Oracle",
// 		ShortName:   "oracle",
// 	}, {
// 		ID:          112,
// 		DisplayName: "Winter Wyvern",
// 		ShortName:   "winter_wyvern",
// 	}, {
// 		ID:          113,
// 		DisplayName: "Arc Warden",
// 		ShortName:   "arc_warden",
// 	}, {
// 		ID:          114,
// 		DisplayName: "Monkey King",
// 		ShortName:   "monkey_king",
// 	}, {
// 		ID:          119,
// 		DisplayName: "Dark Willow",
// 		ShortName:   "dark_willow",
// 	}, {
// 		ID:          120,
// 		DisplayName: "Pangolier",
// 		ShortName:   "pangolier",
// 	}, {
// 		ID:          121,
// 		DisplayName: "Grimstroke",
// 		ShortName:   "grimstroke",
// 	}, {
// 		ID:          123,
// 		DisplayName: "Hoodwink",
// 		ShortName:   "hoodwink",
// 	}, {
// 		ID:          126,
// 		DisplayName: "Void Spirit",
// 		ShortName:   "void_spirit",
// 	}, {
// 		ID:          128,
// 		DisplayName: "Snapfire",
// 		ShortName:   "snapfire",
// 	}, {
// 		ID:          129,
// 		DisplayName: "Mars",
// 		ShortName:   "mars",
// 	}, {
// 		ID:          131,
// 		DisplayName: "Ringmaster",
// 		ShortName:   "ringmaster",
// 	}, {
// 		ID:          135,
// 		DisplayName: "Dawnbreaker",
// 		ShortName:   "dawnbreaker",
// 	}, {
// 		ID:          136,
// 		DisplayName: "Marci",
// 		ShortName:   "marci",
// 	}, {
// 		ID:          137,
// 		DisplayName: "Primal Beast",
// 		ShortName:   "primal_beast",
// 	}, {
// 		ID:          138,
// 		DisplayName: "Muerta",
// 		ShortName:   "muerta",
// 	}, {
// 		ID:          145,
// 		DisplayName: "Kez",
// 		ShortName:   "kez",
// 	},
// }
