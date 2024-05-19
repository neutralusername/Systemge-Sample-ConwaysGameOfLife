# what is this
a (very) simple implementation of conways game of life as a "multiplayer" variant on a 75x75 grid.
every 5 seconds a new generation is calculated and broadcasted to every client.
during these 5 seconds everyone on the website can make changes, which will affect the next generation.

![60QV4Ftw9N](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/2f5b2d0c-65b4-4045-99da-b73d5727f160)


## how to use:  
- make sure to properly import the Systemge library to the SampleApp.  
- locate /main and enter "go run ." which will start the system.  
- to see the grid open "localhost:8080" in your browser while the system is running.  
- you can now click grids to change the square color. changes will be propagated to everyone else currently on this website and persist reloads.  

![image](https://github.com/NeutralUsername/Systemge-Sample-ConwaysGameOfLife/assets/39095721/8ee7d8b1-de7c-4d9a-91ed-5b7c478a92f2)
