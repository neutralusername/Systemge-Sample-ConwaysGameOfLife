import { Cell } from "./cell.js";
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

                cells : [],

                grid: null,
                nextGenerationLoop: null,
                stateInput: "",
                constructMessage: (topic, payload) => {
                    return JSON.stringify({
                        topic: topic,
                        payload: payload,
                    });
                },
                setStateRoot: (state) => {
                    this.setState(state)
                }
            },
            (this.state.WS_CONNECTION.onmessage = (event) => {
                let message = JSON.parse(event.data);
                switch (message.topic) {
                    case "getGridSync":
                    case "getGrid":

                        let grid = JSON.parse(message.payload);
                        let newStateInput = ""
                        let cells = []
                        grid.grid.forEach((row, indexRow) => {
                            row.forEach((cell, indexCol) => {
                                newStateInput += cell
                                cells.push(
                                    React.createElement(Cell, {
                                        cellState: cell,
                                        indexRow: indexRow,
                                        indexCol: indexCol,
                                        cols: grid.cols,
                                        WS_CONNECTION: this.state.WS_CONNECTION,
                                        constructMessage: this.state.constructMessage,
                                    })
                                );
                            })
                        })
                        this.setState({
                            grid: grid,
                            stateInput: newStateInput,
                            cells: cells,
                        });
                        break;
                    case "getGridChange": {
                        let gridChange = JSON.parse(message.payload);
                        this.state.grid.grid[gridChange.row][gridChange.column] = gridChange.state;
                        let newStateInput = this.state.stateInput
                        this.state.cells[gridChange.row * this.state.grid.cols + gridChange.column] = React.createElement(Cell, {
                            cellState: gridChange.state,
                            indexRow: gridChange.row,
                            indexCol: gridChange.column,
                            cols: this.state.grid.cols,
                            WS_CONNECTION: this.state.WS_CONNECTION,
                            constructMessage: this.state.constructMessage,
                        });
                        newStateInput = newStateInput.substring(0, gridChange.row * this.state.grid.cols + gridChange.column) + gridChange.state + newStateInput.substring(gridChange.row * this.state.grid.cols + gridChange.column + 1)
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
                this.state.WS_CONNECTION.send(this.state.constructMessage("heartbeat", ""));
                setTimeout(myLoop, 15 * 1000);
            };
            setTimeout(myLoop, 15 * 1000);
        };
    }

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
                    touchAction: "none",
                    userSelect: "none",
                },
            },
            this.state.grid ? React.createElement(Menu, this.state) : null,
            this.state.grid ? React.createElement(
                "div", {
                    id: "grid",
                    style: {
                        display: "grid",
                        gridTemplateColumns: "repeat(" + this.state.grid.cols + ", " + this.state.SQUARESIZE + "px)",
                        gridTemplateRows: "repeat(" + this.state.grid.rows + ", " + this.state.SQUARESIZE + "px)",
                        width: this.state.grid.cols * this.state.SQUARESIZE + "px",
                        height: this.state.grid.rows * this.state.SQUARESIZE + "px",
                        border: "1px solid black",
                        margin: "auto",
                        marginTop: "10px",
                        marginBottom: "10px",
                        backgroundColor: "white",
                        padding: "0px",
                        boxSizing: "border-box",
                        position: "relative",
                        overflow: "hidden",
                        borderRadius: "5px",
                    },
                },
                this.state.cells
            ) : null
        );
    }
}