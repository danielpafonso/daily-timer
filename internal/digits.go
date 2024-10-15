package internal

var (
	Digits        = make(map[int]string)
	Dots   string = `
 db
 VP

 db
 VP`
)

func init() {
	Digits[0] = ` .d88b.
.8P  88.
88  d'88
88 d' 88
'88  d8'
 'Y88P'`

	Digits[1] = `   db
  o88
   88
   88
   88
   VP`

	Digits[2] = `.d888b.
VP  '8D
   odD'
 .88'
j88.
888888D`

	Digits[3] = `d8888b.
VP  '8D
  oooY'
     b.
db   8D
Y8888P'`

	Digits[4] = `  j88D
 j8~88
j8' 88
V88888D
    88
    VP`

	Digits[5] = `oooooo
8P
8P
V8888b.
    '8D
88oobY'`

	Digits[6] = `   dD
  d8'
 d8'
d8888b.
88' '8D
'8888P`

	Digits[7] = `d88888D
VP  d8'
   d8'
  d8'
 d8'
d8'`

	Digits[8] = `.d888b.
88   8D
'VoooY'
.d   b.
88   8D
'Y888P'`

	Digits[9] = `.d888b.
88' '8D
'V8o88'
   d8'
  d8'
 d8'`
}
