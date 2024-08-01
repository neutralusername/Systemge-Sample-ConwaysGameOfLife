# what is this 
- a simple implementation of conways game of life as a "multiplayer" variant on a 90x140 grid
- dimensions can be changed by changing the corresponding values in appGameOfLife/app.go/New()
- a new generation is generated on button press
- there is a client-side implementation for a new generation loop with a variable interval (time in ms)  
- additionally there is a field which automatically updates with the current grid state
- the content of this field can be copied and changed to save and load a grid state by pressing the set button  
- reset clears the grid

changes to the grid are broadcasted automatically to all other websocket clients.

![chrome_8N6fII47kP](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/ea29055f-c9f9-404b-94a5-56fb0a07e051)


## requirements
- golang (version 1.22.3)
- systemge library on your device

## how to use:  
- make sure to import the Systemge library into the project using the correct path (go.mod)
- locate /main and enter "go run ." into the terminal to launch the command line interface
- type "start" and press enter 
- to interact with the grid open "localhost:8080" in your browser after starting the system
- click cells to change their states
- changes will be propagated to everyone else currently on the website and persist reloads
- enter "randomize" into the command line interface to randomize every cells state
- enter "invert" to invert all cells

## notes
if you intend to use this for production change the tls key and certificate


![60QV4Ftw9N](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/2f5b2d0c-65b4-4045-99da-b73d5727f160)


![chrome_i6yvUFMgJH](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/e220437f-a2c5-483f-a086-fb810827f419)

