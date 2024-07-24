export class Menu extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            nextGenerationLoop: null,
            autoNextGenDelay_ms: 100,
        }
    }

    nextGeneratioLoop = () => {
        this.props.WS_CONNECTION.send(
            this.props.constructMessage("nextGeneration", "")
        );
        this.setState({
            nextGenerationLoop: setTimeout(this.nextGeneratioLoop, this.state.autoNextGenDelay_ms),
        });
    }

    render() {
        return React.createElement(
            "div", {
                id: "menu",
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
                        this.props.WS_CONNECTION.send(
                            this.props.constructMessage("nextGeneration", "")
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
                                this.nextGeneratioLoop()
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
            React.createElement("div", {
                    style: {
                        position: "absolute",
                        display: "flex",
                        flexDirection: "row",
                        top: "101px",
                        left: "10px",
                        padding: "5px",
                        border: "1px solid black",
                        borderRadius: "5px",
                        backgroundColor: "white",
                        color: "black",
                        fontFamily: "Arial",
                        fontSize: "16px",
                        gap: "10px",
                    }
                },
                React.createElement("input", {
                    style: {
                        width: "74px",
                        height: "20px",
                        border: "1px solid black",
                        borderRadius: "5px",
                        backgroundColor: "white",
                        color: "black",
                        fontFamily: "Arial",
                        fontSize: "16px",
                    },
                    onChange: (e) => {
                        this.props.setStateRoot({
                            stateInput: e.target.value,
                        })
                    },
                    value: this.props.stateInput,
                }),
                React.createElement("button", {
                        onClick: () => {
                            this.props.WS_CONNECTION.send(
                                this.props.constructMessage("setGrid", this.props.stateInput)
                            )
                        }
                    },
                    "set"
                ),
                React.createElement("button", {
                        onClick: () => {
                            this.props.WS_CONNECTION.send(
                                this.props.constructMessage("setGrid", (() => {
                                    let str = ""
                                    for (let i = 0; i < this.props.grid.rows; i++) {
                                        for (let j = 0; j < this.props.grid.cols; j++) {
                                            str += "0"
                                        }
                                    }
                                    return str
                                })())
                            )
                        }
                    },
                    "reset"
                )
            ),
        )
    }
}