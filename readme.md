# what is this 
- a simple implementation of conways game of life as a "multiplayer" variant on a 90x140 grid
- dimensions can be changed by changing the corresponding constants in appGameOfLife/app.go
- a new generation is generated on button press
- there is a client-side implementation for a new generation loop with a variable interval (time in ms)  
- additionally there is a field which automatically updates with the current grid state
- the content of this field can be copied and changed to save and load a grid state by pressing the set button  
- reset clears the grid

changes to the grid are broadcasted automatically to all other websocket clients.

![Screenshot from 2024-06-08 22-13-33](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/304513a9-7659-47b7-a83b-1174476d41cf)


![systemge-game-of-life(10)](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/b6f9c94c-f8e6-4d5b-9c43-b8b044626413)

## requirements
- golang (version 1.22.3)

## how to use:  
- make sure to import the Systemge library into the project using the correct path (go.mod)
- locate /main and enter "go run ." into the terminal to launch the command line interface
- type "start" and press enter 
- to interact with the grid open "localhost:8080" in your browser after starting the system
- click cells to change their states
- changes will be propagated to everyone else currently on the website and persist reloads
- enter "randomize" into the command line interface to randomize every cells state
- enter "invert" to invert all cells


![60QV4Ftw9N](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/2f5b2d0c-65b4-4045-99da-b73d5727f160)


![chrome_i6yvUFMgJH](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/e220437f-a2c5-483f-a086-fb810827f419)

