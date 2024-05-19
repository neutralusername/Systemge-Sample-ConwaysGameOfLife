
export class root extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			WS_CONNECTION : new WebSocket("ws://localhost:8443/ws"),

			GRIDSIZE : 75,
            SQUARESIZE : 12.5,
            AUTONEXTGENERATIONTEMPO : 100,

            grid : [],
            nextGenerationLoop : null,
		},
        this.state.WS_CONNECTION.onmessage = (event) => {
            let message = JSON.parse(event.data)
            console.log(message)
            switch (message.type) {
                case "getGrid":
                    let newGrid = []
                    message.body.split('').forEach((digit) => {
                        newGrid.push(Number(digit))
                    })
                    this.setState({
                        grid : newGrid
                    })
                    break
                case "getGridChange":
                    let gridChange = JSON.parse(message.body)
                    this.state.grid[gridChange.row*this.state.GRIDSIZE + gridChange.column] = gridChange.state ? 1 : 0
                    this.setState({
                        grid : this.state.grid
                    })
                    break
                default:
                    console.log("Unknown message type: " + event.data)
                    break
            }
        }
        this.state.WS_CONNECTION.onclose = function() {
            setTimeout(function() {
                if (WS_CONNECTION.readyState === WebSocket.CLOSED) {}
                    window.location.reload()
            }, 2000)
        }
        this.state.WS_CONNECTION.onopen = function() {
            let myLoop = function() {
                WS_CONNECTION.send(this.constructMessage("heartbeat", ""))
                setTimeout(myLoop, 15*1000)
            }
            setTimeout(myLoop, 15*1000)
        }
    }

    constructMessage = (type, body) => {
        return JSON.stringify({
             type: type,
             body: body  
        })
    }

	render() {
        let gridElements = []
        this.state.grid.forEach((cell, index) => {
            gridElements.push(React.createElement('div', {
                key : index,
                style : {
                    width : this.state.SQUARESIZE+"px",
                    height : this.state.SQUARESIZE+"px",
                    backgroundColor : cell ? "black" : "white",
                    border : "1px solid black",
                    boxSizing : "border-box",
                },
                onClick : () => this.state.WS_CONNECTION.send(this.constructMessage("gridChange", JSON.stringify({row:Math.floor(index/this.state.GRIDSIZE), column:index%this.state.GRIDSIZE, state:cell ? false : true}))),
                onMouseOver : (e) => {
                    if (e.buttons === 1) {
                        this.state.WS_CONNECTION.send(this.constructMessage("gridChange", JSON.stringify({row:Math.floor(index/this.state.GRIDSIZE), column:index%this.state.GRIDSIZE, state:cell ? false : true})))
                    }
                }
            }))
        })
		return React.createElement('div', {
				id : "root",
				onContextMenu: e => {
					e.preventDefault()
				},
				style : {
					fontFamily : "sans-serif",
					display : "flex",
					flexDirection : "column",
					justifyContent : "center",
					alignItems : "center",
				}
			},
            React.createElement('button', {
                    id : "nextGeneration",
                    style : {
                        position : "absolute",
                        top : "10px",
                        left : "10px",
                        padding : "5px",
                        border : "1px solid black",
                        borderRadius : "5px",
                        backgroundColor : "white",
                        color : "black",
                        fontFamily : "Arial",
                        fontSize : "16px",
                        cursor : "pointer",
                    },
                    onClick : () => this.state.WS_CONNECTION.send(this.constructMessage("nextGeneration", "")),
                    innerHTML : "Next Generation"
                },
                "Next Generation"
            ),
            React.createElement('button', {
                    id : "nextGenerationLoop",
                    style : {
                        position : "absolute",
                        top : "50px",
                        left : "10px",
                        padding : "5px",
                        border : "1px solid black",
                        borderRadius : "5px",
                        backgroundColor : "white",
                        color : "black",
                        fontFamily : "Arial",
                        fontSize : "16px",
                        cursor : "pointer",
                    },
                    onClick : () => {
                        if (this.state.nextGenerationLoop === null) {
                            this.setState({
                                nextGenerationLoop : setInterval(() => {
                                    this.state.WS_CONNECTION.send(this.constructMessage("nextGeneration", ""))
                                }, this.state.AUTONEXTGENERATIONTEMPO)
                            })
                        } else {
                            clearInterval(this.state.nextGenerationLoop)
                            this.setState({
                                nextGenerationLoop : null
                            })
                        }
                    },
                },
                this.state.nextGenerationLoop === null ? "Start Loop" : "Stop Loop"
            ),
			React.createElement('div', {
                    id : "grid",
                    style : {
                        display : "grid",
                        gridTemplateColumns : "repeat("+this.state.GRIDSIZE+", 1fr)",
                        gridTemplateRows : "repeat("+this.state.GRIDSIZE+", 1fr)",
                        width : this.state.SQUARESIZE*this.state.GRIDSIZE+"px",
                        height : this.state.SQUARESIZE*this.state.GRIDSIZE+"px",
                        border : "1px solid black",
                        margin : "auto",
                        marginTop : "10px",
                        marginBottom : "10px",
                        backgroundColor : "white",
                        padding : "0px",
                        boxSizing : "border-box",
                        position : "relative",
                        overflow : "hidden",
                        borderRadius : "5px",
                    }
                },
                gridElements
            )      
		) 
	}
}