import { 
    Cell 
} from "./cell.js";

export class Grid extends React.Component {
    constructor(props) {
        super(props);
        this.state = {}
    }


    render() {
        let cells = [];
        this.props.grid.grid.forEach((row, indexRow) => {
            row.forEach((cell, indexCol) => {
                cells.push(
                   React.createElement(Cell, {
                        cellState : cell,
                        indexRow : indexRow,
                        indexCol : indexCol,
                        cols : this.props.grid.cols,
                        WS_CONNECTION : this.props.WS_CONNECTION,
                        constructMessage : this.props.constructMessage,
                    })
                );
            })
        });
        return React.createElement(
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
            cells
        );
    }
}