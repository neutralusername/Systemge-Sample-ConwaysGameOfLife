export class root extends React.Component {
    constructor(props) {
        super(props);
        (this.state = {
            WS_CONNECTION: new WebSocket("ws://localhost:8443/ws"),

            SQUARESIZE: 12.5,
            autoNextGenDelay_ms: 100,

            grid: null,
            nextGenerationLoop: null,
        }),
        (this.state.WS_CONNECTION.onmessage = (event) => {
            let message = JSON.parse(event.data);
            console.log(message);
            switch (message.type) {
                case "getGrid":
                    let grid = JSON.parse(message.body);
                    this.setState({
                        grid: grid,
                    });
                    break;
                case "getGridChange":
                    let gridChange = JSON.parse(message.body);
                    this.state.grid.grid[gridChange.row][gridChange.column] = gridChange.state;
                    this.setState({
                        grid: this.state.grid,
                    });
                    break;
                default:
                    console.log("Unknown message type: " + event.data);
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

    constructMessage = (type, body) => {
        return JSON.stringify({
            type: type,
            body: body,
        });
    };

    startNextGenerationLoop = () => {
        this.state.WS_CONNECTION.send(
            this.constructMessage("nextGeneration", "")
        );
        this.setState({
            nextGenerationLoop: setTimeout(this.startNextGenerationLoop, this.state.autoNextGenDelay_ms),
        });
    }

    render() {
        let gridElements = [];
        if (this.state.grid) {
            this.state.grid.grid.forEach((row, indexRow) => {
                row.forEach((cell, indexCol) => {
                    gridElements.push(
                        React.createElement("div", {
                            key: indexRow*this.state.grid.cols+indexCol,
                            style: {
                                width: this.state.SQUARESIZE + "px",
                                height: this.state.SQUARESIZE + "px",
                                backgroundColor: cell ? "black" : "white",
                                border: "1px solid black",
                                boxSizing: "border-box",
                            },
                            onClick: () =>
                                this.state.WS_CONNECTION.send(
                                    this.constructMessage(
                                        "gridChange",
                                        JSON.stringify({
                                            row :indexRow,
                                            column: indexCol,
                                            state: cell ? 0 : 1,
                                        })
                                    )
                                ),
                            onMouseOver: (e) => {
                                if (e.buttons === 1) {
                                    this.state.WS_CONNECTION.send(
                                        this.constructMessage(
                                            "gridChange",
                                            JSON.stringify({
                                                row: indexRow,
                                                column: indexCol,
                                                state: cell ? 0 : 1,
                                            })
                                        )
                                    );
                                }
                            },
                        })
                    );
                })
            });
        }
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
            React.createElement(
                "button", {
                    id: "nextGeneration",
                    style: {
                        position: "absolute",
                        top: "10px",
                        left: "10px",
                        padding: "5px",
                        backgroundColor: "white",
                        color: "black",
                        fontFamily: "Arial",
                        fontSize: "16px",
                        cursor: "pointer",
                    },
                    onClick: () =>
                        this.state.WS_CONNECTION.send(
                            this.constructMessage("nextGeneration", "")
                        ),
                    innerHTML: "Next Generation",
                },
                "Next Generation"
            ),
            React.createElement("div", {
                    style: {
                        position: "absolute",
                        display: "flex",
                        flexDirection: "row",
                        top: "50px",
                        left: "10px",
                        padding: "5px",
                        border: "1px solid black",
                        borderRadius: "5px",
                        backgroundColor: "white",
                        color: "black",
                        fontFamily: "Arial",
                        fontSize: "16px",
                        gap: "10px",    
                    },
                },
                React.createElement(
                    "button", {
                        id: "nextGenerationLoop",
                        style: {
                            backgroundColor: "white",
                            color: "black",
                            fontFamily: "Arial",
                            fontSize: "16px",
                            cursor: "pointer",
                        },
                        onClick: () => {
                            if (this.state.nextGenerationLoop === null) {
                                this.startNextGenerationLoop ()
                            } else {
                                clearTimeout(this.state.nextGenerationLoop);
                                this.setState({
                                    nextGenerationLoop: null,
                                });
                            }
                        },
                    },
                    this.state.nextGenerationLoop === null ? "Start Loop" : "Stop Loop"
                ),
                React.createElement("input", {
                    id: "autoNextGenDelay",
                    type: "number",
                    style: {
                        width: "65px",
                        height: "20px",
                        padding: "5px",
                        border: "1px solid black",
                        borderRadius: "5px",
                        backgroundColor: "white",
                        color: "black",
                        fontFamily: "Arial",
                        fontSize: "16px",
                    },
                    value: this.state.autoNextGenDelay_ms,
                    onChange: (e) => {
                        this.setState({
                            autoNextGenDelay_ms: e.target.value,
                        });
                    },
                }),
            ),
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
                gridElements
            ) : null
        );
    }
}