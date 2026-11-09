package internal

// Currency metadata. History data comes from the ECB reference set (frankfurter);
// these 30 codes are the complete supported set.
type Currency struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Zh       string `json:"zh"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
	Popular  int    `json:"popular"`
}

var Currencies = []Currency{
	{"USD", "US Dollar", "美元", "$", 2, 1},
	{"CNY", "Chinese Yuan", "人民币", "¥", 2, 1},
	{"EUR", "Euro", "欧元", "€", 2, 1},
	{"JPY", "Japanese Yen", "日元", "¥", 0, 1},
	{"GBP", "British Pound", "英镑", "£", 2, 1},
	{"HKD", "Hong Kong Dollar", "港币", "HK$", 2, 1},
	{"KRW", "South Korean Won", "韩元", "₩", 0, 1},
	{"AUD", "Australian Dollar", "澳元", "A$", 2, 1},
	{"CHF", "Swiss Franc", "瑞士法郎", "Fr", 2, 1},
	{"SGD", "Singapore Dollar", "新加坡元", "S$", 2, 1},
	{"CAD", "Canadian Dollar", "加元", "C$", 2, 0},
	{"NZD", "New Zealand Dollar", "新西兰元", "NZ$", 2, 0},
	{"THB", "Thai Baht", "泰铢", "฿", 2, 0},
	{"MYR", "Malaysian Ringgit", "林吉特", "RM", 2, 0},
	{"IDR", "Indonesian Rupiah", "印尼盾", "Rp", 0, 0},
	{"PHP", "Philippine Peso", "菲律宾比索", "₱", 2, 0},
	{"INR", "Indian Rupee", "印度卢比", "₹", 2, 0},
	{"ILS", "Israeli Shekel", "以色列谢克尔", "₪", 2, 0},
	{"TRY", "Turkish Lira", "土耳其里拉", "₺", 2, 0},
	{"MXN", "Mexican Peso", "墨西哥比索", "$", 2, 0},
	{"BRL", "Brazilian Real", "巴西雷亚尔", "R$", 2, 0},
	{"ZAR", "South African Rand", "南非兰特", "R", 2, 0},
	{"SEK", "Swedish Krona", "瑞典克朗", "kr", 2, 0},
	{"NOK", "Norwegian Krone", "挪威克朗", "kr", 2, 0},
	{"DKK", "Danish Krone", "丹麦克朗", "kr", 2, 0},
	{"PLN", "Polish Złoty", "兹罗提", "zł", 2, 0},
	{"CZK", "Czech Koruna", "捷克克朗", "Kč", 2, 0},
	{"HUF", "Hungarian Forint", "福林", "Ft", 2, 0},
	{"ISK", "Icelandic Króna", "冰岛克朗", "kr", 0, 0},
}

var CurrencyByCode = func() map[string]Currency {
	m := make(map[string]Currency, len(Currencies))
	for _, c := range Currencies {
		m[c.Code] = c
	}
	return m
}()
