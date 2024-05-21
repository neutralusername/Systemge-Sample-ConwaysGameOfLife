import {
    Grid
} from "./grid.js";
import { 
    Menu 
} from "./menu.js";

export class root extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            WS_CONNECTION: new WebSocket("ws://localhost:8443/ws"),

            SQUARESIZE: 10,
            autoNextGenDelay_ms: 100,

            grid: null,
            nextGenerationLoop: null,
            stateInput : "",
            constructMessage : this.constructMessage,
            setStateRoot : (state) => {
                this.setState(state)
            }
        },
        (this.state.WS_CONNECTION.onmessage = (event) => {
            let message = JSON.parse(event.data);
            switch (message.topic) {
                case "getGrid":
                    console.log(event.data)
                    let grid = JSON.parse(message.body);
                    let newStateInput = ""
                    grid.grid.forEach((row) => {
                        row.forEach((cell) => {
                            newStateInput += cell
                        })
                    })
                    this.setState({
                        grid: grid,
                        stateInput: newStateInput,
                    });
                    break;
                case "getGridChange": {
                    let gridChange = JSON.parse(message.body);
                    this.state.grid.grid[gridChange.row][gridChange.column] = gridChange.state;
                    let newStateInput = this.state.stateInput
                    newStateInput = newStateInput.substring(0, gridChange.row*this.state.grid.cols+gridChange.column) + gridChange.state + newStateInput.substring(gridChange.row*this.state.grid.cols+gridChange.column+1)
                    this.setState({
                        grid: this.state.grid,
                        stateInput: newStateInput,
                    });
                    break;
                }
                default:
                    console.log("Unknown message topic: " + event.data);
                    break;
            }
        });
        this.state.WS_CONNECTION.onclose = () => {
            setTimeout(() => {
                if (this.state.WS_CONNECTION.readyState === WebSocket.CLOSED) {}
                window.location.reload();
            }, 2000);
        };
        this.state.WS_CONNECTION.onopen = () => {
            let myLoop = () => {
                this.state.WS_CONNECTION.send(this.constructMessage("heartbeat", ""));
                setTimeout(myLoop, 15 * 1000);
            };
            setTimeout(myLoop, 15 * 1000);
        };
    }

    constructMessage = (topic, body) => {
        return JSON.stringify({
            topic: topic,
            body: body,
        });
    };

    render() {
        return React.createElement(
            "div", {
                id: "root",
                onContextMenu: (e) => {
                    e.preventDefault();
                },
                style: {
                    fontFamily: "sans-serif",
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                    alignItems: "center",
                    touchAction : "none",
					userSelect : "none",
                },
            },
            this.state.grid ? React.createElement(Menu, this.state) : null,
            this.state.grid ? React.createElement(Grid, this.state) : null
        );
    }
}