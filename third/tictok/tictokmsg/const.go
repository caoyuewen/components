package tictokmsg

const (
	TaskStartMsgTypeComment  = "live_comment"  // 评论
	TaskStartMsgTypeGift     = "live_gift"     // 礼物
	TaskStartMsgTypeLike     = "live_like"     // 点赞
	TaskStartMsgTypeFansclub = "live_fansclub" // 粉丝团
)

/*
	sid		礼物				钻石价值	增加足球数	增加积分	gid
	1		仙女棒			1		50			10		n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs=
	2		仙女棒1组（10个）	10		500			110		n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs=
	3		能力药丸			10		500			110		28rYzVFNyXEXFC8HI+f/WG+I7a6lfl3OyZZjUS+CVuwCgYZrPrUdytGHu0c=
	4		魔法镜			19		500			200		fJs8HKQ0xlPRixn8JAUiL2gFRiLD9S6IFCFdvZODSnhyo9YN8q7xUuVVyZI=
	5		甜甜圈			52		1000		5500	PJ0FFeaDzXUreuUBZH6Hs+b56Jh0tQjrq0bIrrlZmv13GSAL9Q1hf59fjGk=
	6		能量电池			99		2000		1200	IkkadLfz7O/a5UR45p/OOCCG6ewAWVbsuzR/Z+v1v76CBU+mTG/wPjqdpfg=
	7		爱的爆炸(恶魔炸弹)	199		4000		2500	gx7pmjQfhBaDOG2XkWI2peZ66YFWkCWRjZXpTqb23O/epru+sxWyTV/3Ufs=
	8		神秘空投			520		10000		6500	pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI=
*/

const (
	GiftFairySticksGroupNum = 10
)

/*
   "NzWjErGItCd1JbTkcZPXZ+McUN9QRhv+JoG43QjAGnyzhCsNmjC9hZ5DNSs=",
    "I+pWrWQoF67euZyPuNMG9TJYOVbjQqBPHR4PZOkSXH+u3gq3MN5AxYzuIIs=",
    "cTrIsyYM88TqgVN3+7Ix8ZCy503ZAGcFpvnSYBKMgGcr8bO9cnjsKfewaPw=",
    "k40//24g4+xpRxN2+tVL4ixvFIAl88Un0lasphvr5e7uUj0WQFG7iyiYpWA=",
    "HjvyGu7gcvbHRi+cFRlizUoIhVbFvSrSyaummvbnZISMUm7fJwTjosE6+aM=",
    "yCUpgk2ewIRIG9iVTXOSD6gGqvlBZY1fnNeoIiEJ/7oUGMLLSpCQ9WQlC+U=",
    "4/4bykWAIbgxoJ+bq//eHpf9U0FcRpiLjJ6LWIYrpyVNDjjkDgwMrwGmGZc=",
    "ZaSHqPZw2yhMcMcTFteJePlAkbBXjGZNP89Tz9ytoE9rnb1R0C01kvsYSQ4=",
    "6meQjmMnlskkYKBc5Up9F1W4FTdDmsTwlMq6xqqoP6c0EM0I6Z26OmEeWOA=",
    "ueaBKILcj49t8VyUWARYHFm8Ts7fRMbUVBU3t6idxJn1ZFhLaCLelC2qJGY=",
    "0lxm9m9bzdZJc9aCcrbnoXMlpJ4mSBgGeIYaZgp2gwXOJUTF6zHKOKj9gGw=",
    "iJVlyMdf9XDI9WelF7lwI/Jsll68G4syWqgP5pdYUK11A9EVF5OB1RW9u+M=",
    "wknBHep7wqRlDtRsW2t7QeykNg3qCWhKfAVpoPitioA/DcXmBi6zIgUV34Q=",
    "fkBeJecU782reNzqpi8iZRJILYj73wkMg7SnFzTsQrkY8chji4vL2jv41wM=",
    "94JsIHHrd0y8T5Htq68/FOcXs0bgrK3vOMKjY04uy7S+kLQa+KHeBN1EgE0=",
    "PI0TosNkq2eSpIlnTRKfS/KiFlJkLIU7SHVl7rzmzXDXivbeyJzLoK5icFM=",
    "okIrEUvWBV7dEa1Bp5DOsGxA2dMtPuJty8I1kq4SfxKeX5v72wdzC/fa0WM=",
    "BgPiF44wXUV9SFs1vceqEH0CTuEbZIBZRlOpL7fLsEThNe0TKPADNMKYzV8="
*/

const (
	GiftFairySticksId        = "n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs="        //
	GiftFairySticksGroupId   = "n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs=group"   //  仙女棒组刷 => 双管炮
	GiftFairySticksGroup66Id = "n1/Dg1905sj1FyoBlQBvmbaDZFBNaKuKZH6zxHkv8Lg5x2cRfrKUTb8gzMs=group66" //  仙女棒66组刷
	GiftAbilityPillsId       = "28rYzVFNyXEXFC8HI+f/WG+I7a6lfl3OyZZjUS+CVuwCgYZrPrUdytGHu0c="        // 能量药丸 => 能力突破(速射炮)
	GiftMagicMirrorId        = "fJs8HKQ0xlPRixn8JAUiL2gFRiLD9S6IFCFdvZODSnhyo9YN8q7xUuVVyZI="        // 魔法镜 => 高位压迫+双管炮
	GiftDonutId              = "PJ0FFeaDzXUreuUBZH6Hs+b56Jh0tQjrq0bIrrlZmv13GSAL9Q1hf59fjGk="        // 甜甜圈 => 大力抽射+(速射炮)
	GiftAbilityBatteryId     = "IkkadLfz7O/a5UR45p/OOCCG6ewAWVbsuzR/Z+v1v76CBU+mTG/wPjqdpfg="
	GiftEvilBombId           = "gx7pmjQfhBaDOG2XkWI2peZ66YFWkCWRjZXpTqb23O/epru+sxWyTV/3Ufs="
	GiftMysteryAirdropsId    = "pGLo7HKNk1i4djkicmJXf6iWEyd+pfPBjbsHmd3WcX0Ierm2UdnRR7UINvI="
)

var EffectiveGiftIds = []string{
	GiftFairySticksId,
	GiftFairySticksGroupId,
	GiftAbilityPillsId,
	GiftMagicMirrorId,
	GiftDonutId,
	GiftAbilityBatteryId,
	GiftEvilBombId,
	GiftMysteryAirdropsId,
}

var TopGiftIds = []string{
	GiftFairySticksId,
	GiftAbilityPillsId,
	GiftMagicMirrorId,
	GiftDonutId,
	GiftAbilityBatteryId,
	GiftEvilBombId,
	//GiftMysteryAirdropsId,
}
