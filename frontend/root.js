import { Cell } from "./cell.js";
import { Menu } from "./menu.js";

export class root extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            SQUARESIZE: 10,
            cells: [],
            grid: null,
            stateInput: "",
        };
        
        this.WS_CONNECTION = new WebSocket("ws://localhost:8443/ws");
        this.WS_CONNECTION.onmessage = this.handleMessage.bind(this);
        this.WS_CONNECTION.onclose = this.handleClose.bind(this);
        this.WS_CONNECTION.onopen = this.handleOpen.bind(this);
    }

    handleMessage(event) {
        const message = JSON.parse(event.data);
        switch (message.topic) {
            case "getGrid":
            case "propagateGrid":
                this.updateGrid(JSON.parse(message.payload));
                break;
            case "propagateGridChange":
                this.updateGridChange(JSON.parse(message.payload));
                break;
            default:
                console.log("Unknown message topic: " + event.data);
                break;
        }
    }

    handleClose() {
        setTimeout(() => {
            if (this.WS_CONNECTION.readyState === WebSocket.CLOSED) {
                window.location.reload();
            }
        }, 2000);
    }

    handleOpen() {
        const sendHeartbeat = () => {
            this.WS_CONNECTION.send(this.constructMessage("heartbeat", ""));
            setTimeout(sendHeartbeat, 15 * 1000);
        };
        setTimeout(sendHeartbeat, 15 * 1000);
    }

    constructMessage(topic, payload) {
        return JSON.stringify({ topic, payload });
    }

    updateGrid(grid) {
        const newStateInput = grid.grid.flat().join('');
        const cells = grid.grid.flatMap((row, indexRow) =>
            row.map((cell, indexCol) =>
                React.createElement(Cell, {
                    key: `${indexRow}-${indexCol}`,
                    cellState: cell,
                    indexRow,
                    indexCol,
                    cols: grid.cols,
                    WS_CONNECTION: this.WS_CONNECTION,
                    constructMessage: this.constructMessage,
                })
            )
        );

        this.setState({ grid, stateInput: newStateInput, cells });
    }

    updateGridChange(gridChange) {
        const { row, column, state } = gridChange;
        this.setState((prevState) => {
            const newGrid = { ...prevState.grid };
            newGrid.grid[row][column] = state;
            const newCells = [...prevState.cells];
            const cellIndex = row * newGrid.cols + column;
            newCells[cellIndex] = React.cloneElement(newCells[cellIndex], { cellState: state });
            const newStateInput = prevState.stateInput.substring(0, cellIndex) + state + prevState.stateInput.substring(cellIndex + 1);
            return { grid: newGrid, stateInput: newStateInput, cells: newCells };
        });
    }
    
    render() {
        const { grid, cells, SQUARESIZE } = this.state;
        return React.createElement(
            "div",
            {
                id: "root",
                onContextMenu: (e) => e.preventDefault(),
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
            grid ? React.createElement(Menu, { 
                ...this.state, 
                setStateRoot: this.setState.bind(this),
                WS_CONNECTION: this.WS_CONNECTION,
                constructMessage: this.constructMessage,
            }) : null,
            grid
                ? React.createElement(
                      "div",
                      {
                          id: "grid",
                          style: {
                              display: "grid",
                              gridTemplateColumns: `repeat(${grid.cols}, ${SQUARESIZE}px)`,
                              gridTemplateRows: `repeat(${grid.rows}, ${SQUARESIZE}px)`,
                              width: grid.cols * SQUARESIZE,
                              height: grid.rows * SQUARESIZE,
                              border: "1px solid black",
                              margin: "10px auto",
                              backgroundColor: "white",
                              padding: "0",
                              boxSizing: "border-box",
                              borderRadius: "5px",
                          },
                      },
                      cells
                  )
                : null
        );
    }
}