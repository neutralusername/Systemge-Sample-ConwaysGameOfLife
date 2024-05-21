export class Cell extends React.Component {
    constructor(props) {
        super(props);
        this.state = {}
    }

    render() {
        return React.createElement(
            "div", {
                id: "cell",
                key: this.props.indexRow * this.props.cols + this.props.indexCol,
                style: {
                    width: this.props.SQUARESIZE + "px",
                    height: this.props.SQUARESIZE + "px",
                    backgroundColor: this.props.cellState ? "black" : "white",
                    border: "1px solid black",
                    boxSizing: "border-box",
                },
                onClick: () =>
                    this.props.WS_CONNECTION.send(
                        this.props.constructMessage(
                            "gridChange",
                            JSON.stringify({
                                row: this.props.indexRow,
                                column: this.props.indexCol,
                                state: (this.props.cellState + 1) % 2,
                            })
                        )
                    ),
                onMouseOver: (e) => {
                    if (e.buttons === 1) {
                        this.props.WS_CONNECTION.send(
                            this.props.constructMessage(
                                "gridChange",
                                JSON.stringify({
                                    row: this.props.indexRow,
                                    column: this.props.indexCol,
                                    state: (this.props.cellState + 1) % 2,
                                })
                            )
                        );
                    }
                },
            },
        )
    }
}