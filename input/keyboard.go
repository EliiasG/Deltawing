package input

type Keyboard interface {
	SetKeyPressedHandler(handler func(Key))
	SetKeyReleasedHandler(handler func(Key))
	SetKeyHeldHandler(handler func(Key))
	SetKeyTypedHandler(handler func(rune))
}

type Key int32

// Maybe i should just une constants fronm glfw, but i did not want this to be related to glfw
// Generated from glfw keys, by python and aligned by go fmt

const (
	KeyUnknown      Key = -1
	KeySpace        Key = 32
	KeyApostrophe   Key = 39 // '
	KeyComma        Key = 44 // ,
	KeyMinus        Key = 45 // -
	KeyPeriod       Key = 46 // .
	KeySlash        Key = 47 // /
	Key0            Key = 48
	Key1            Key = 49
	Key2            Key = 50
	Key3            Key = 51
	Key4            Key = 52
	Key5            Key = 53
	Key6            Key = 54
	Key7            Key = 55
	Key8            Key = 56
	Key9            Key = 57
	KeySemicolon    Key = 59 // ;
	KeyEqual        Key = 61 // =
	KeyA            Key = 65
	KeyB            Key = 66
	KeyC            Key = 67
	KeyD            Key = 68
	KeyE            Key = 69
	KeyF            Key = 70
	KeyG            Key = 71
	KeyH            Key = 72
	KeyI            Key = 73
	KeyJ            Key = 74
	KeyK            Key = 75
	KeyL            Key = 76
	KeyM            Key = 77
	KeyN            Key = 78
	KeyO            Key = 79
	KeyP            Key = 80
	KeyQ            Key = 81
	KeyR            Key = 82
	KeyS            Key = 83
	KeyT            Key = 84
	KeyU            Key = 85
	KeyV            Key = 86
	KeyW            Key = 87
	KeyX            Key = 88
	KeyY            Key = 89
	KeyZ            Key = 90
	KeyLeftBracket  Key = 91  // [
	KeyBackslash    Key = 92  // \
	KeyRightBracket Key = 93  // ]
	KeyGraveAccent  Key = 96  // `
	KeyWorld1       Key = 161 // non-US #1
	KeyWorld2       Key = 162 // non-US #2
	KeyEscape       Key = 256
	KeyEnter        Key = 257
	KeyTab          Key = 258
	KeyBackspace    Key = 259
	KeyInsert       Key = 260
	KeyDelete       Key = 261
	KeyRight        Key = 262
	KeyLeft         Key = 263
	KeyDown         Key = 264
	KeyUp           Key = 265
	KeyPageUp       Key = 266
	KeyPageDown     Key = 267
	KeyHome         Key = 268
	KeyEnd          Key = 269
	KeyCapsLock     Key = 280
	KeyScrollLock   Key = 281
	KeyNumLock      Key = 282
	KeyPrintScreen  Key = 283
	KeyPause        Key = 284
	KeyF1           Key = 290
	KeyF2           Key = 291
	KeyF3           Key = 292
	KeyF4           Key = 293
	KeyF5           Key = 294
	KeyF6           Key = 295
	KeyF7           Key = 296
	KeyF8           Key = 297
	KeyF9           Key = 298
	KeyF10          Key = 299
	KeyF11          Key = 300
	KeyF12          Key = 301
	KeyF13          Key = 302
	KeyF14          Key = 303
	KeyF15          Key = 304
	KeyF16          Key = 305
	KeyF17          Key = 306
	KeyF18          Key = 307
	KeyF19          Key = 308
	KeyF20          Key = 309
	KeyF21          Key = 310
	KeyF22          Key = 311
	KeyF23          Key = 312
	KeyF24          Key = 313
	KeyF25          Key = 314
	KeyKp0          Key = 320
	KeyKp1          Key = 321
	KeyKp2          Key = 322
	KeyKp3          Key = 323
	KeyKp4          Key = 324
	KeyKp5          Key = 325
	KeyKp6          Key = 326
	KeyKp7          Key = 327
	KeyKp8          Key = 328
	KeyKp9          Key = 329
	KeyKpDecimal    Key = 330
	KeyKpDivide     Key = 331
	KeyKpMultiply   Key = 332
	KeyKpSubtract   Key = 333
	KeyKpAdd        Key = 334
	KeyKpEnter      Key = 335
	KeyKpEqual      Key = 336
	KeyLeftShift    Key = 340
	KeyLeftControl  Key = 341
	KeyLeftAlt      Key = 342
	KeyLeftSuper    Key = 343
	KeyRightShift   Key = 344
	KeyRightControl Key = 345
	KeyRightAlt     Key = 346
	KeyRightSuper   Key = 347
)
