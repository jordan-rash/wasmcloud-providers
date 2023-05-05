### LCD Display

This provider uses a "display" contract that can be found [here](https://github.com/jordan-rash/wasmcloud-interfaces/tree/main/lcd-display).    

This project takes advantage of a pre-existing LCM1602 driver library that can be found [here](https://github.com/wjessop/lcm1602_lcd), so you will need an I2C compatible display.

### Link Settings

None

### API

#### Actor to Provider

##### `Display.DisplayLine`

This RPC call will print a message on the given line.  They can be used in succession to fill both lines.  The only way to get rid of something on the screen is to run clear.  So, in order to print 
```bash
hello,
wasmcloud!
```

You will need to send two payloads
`{"msg": "hello", "line": 1}`
`{"msg": "wasmcloud!", "line": 2}`


##### `Display.Clear`

This RPC call will clear the entire display and return `true` if it was successful or an `error` if it wasn't.

 
![image](https://user-images.githubusercontent.com/15827604/236480684-2b660788-ff80-4905-b8b9-30301c64c9a8.png)


