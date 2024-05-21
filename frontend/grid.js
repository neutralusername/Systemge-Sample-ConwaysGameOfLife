
export class Grid extends React.Component {
    constructor(props) {
        super(props);
        this.state = {}
    }


    render() {
        let gridElements = [];
        if (this.props.grid) {
            this.props.grid.grid.forEach((row, indexRow) => {
                row.forEach((cell, indexCol) => {
                    gridElements.push(
                        React.createElement("div", {
                            key: indexRow*this.props.grid.cols+indexCol,
                            style: {
                                width: this.props.SQUARESIZE + "px",
                                height: this.props.SQUARESIZE + "px",
                                backgroundColor: cell ? "black" : "white",
                                border: "1px solid black",
                                boxSizing: "border-box",
                            },
                            onClick: () =>
                                this.props.WS_CONNECTION.send(
                                    this.props.constructMessage(
                                        "gridChange",
                                        JSON.stringify({
                                            row :indexRow,
                                            column: indexCol,
                                            state: (cell+1)%2,
                                        })
                                    )
                                ),
                            onMouseOver: (e) => {
                                if (e.buttons === 1) {
                                    this.props.WS_CONNECTION.send(
                                        this.props.constructMessage(
                                            "gridChange",
                                            JSON.stringify({
                                                row: indexRow,
                                                column: indexCol,
                                                state: (cell+1)%2,
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
        return  this.props.grid ? React.createElement(
            "div", {
                id: "grid",
                style: {
                    display: "grid",
                    gridTemplateColumns: "repeat(" + this.props.grid.cols + ", " + this.props.SQUARESIZE + "px)",
                    gridTemplateRows: "repeat(" + this.props.grid.rows + ", " + this.props.SQUARESIZE + "px)",
                    width: this.props.grid.cols * this.props.SQUARESIZE + "px",
                    height: this.props.grid.rows * this.props.SQUARESIZE + "px",
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
        ) : null;
    }
}