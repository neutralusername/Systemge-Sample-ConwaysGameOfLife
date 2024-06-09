# what is this 
- a simple implementation of conways game of life as a "multiplayer" variant on a 90x140 grid
- dimensions can be changed by changing the corresponding constants in appGameOfLife/app.go
- a new generation is generated on button press
- there is a client-side implementation for a new generation loop with a variable interval (time in ms)  
- additionally there is a field which automatically updates with the current grid state
- the content of this field can be copied and changed to save and load a state by pressing the set button  
- reset clears the grid

changes to the grid are broadcasted automatically to all other websocket clients.

![Screenshot from 2024-06-08 22-13-33](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/304513a9-7659-47b7-a83b-1174476d41cf)


![systemge-game-of-life(10)](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/b6f9c94c-f8e6-4d5b-9c43-b8b044626413)



## how to use:  
- make sure to import the Systemge library to the SampleApp
- locate /main and enter "go run ." which will start the system
- type "start" and press enter 
- to see the grid open "localhost:8080" in your browser while the system is running
- you can now click grids to change the square color. changes will be propagated to everyone else currently on this website and persist reloads
- you can type "randomize" in the console to randomize the entire grid
- you can type "invert" to invert all grid cells


![60QV4Ftw9N](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/2f5b2d0c-65b4-4045-99da-b73d5727f160)


![chrome_i6yvUFMgJH](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/e220437f-a2c5-483f-a086-fb810827f419)

