const WS_CONNECTION = new WebSocket("ws://localhost:8443/ws")
const DELIMITER1 = "\x02"
const DELIMITER2 = "\x03"

const GRIDSIZE = 75
const SQUARESIZE = 12.5

function constructMessage(type, ...args) {
    let msg = type + DELIMITER1
    for (let arg of args) {
        msg += (arg + DELIMITER2 + DELIMITER1)
    }
    return msg
}
function parseMessage(message) {
    let data = message.split(DELIMITER1);
    let type = data[0]
    let payload = data.slice(1, -1)
    payload.forEach((pl, i) => {
        payload[i] = pl.split(DELIMITER2).slice(0, -1)
    })
    return {type: type, payload: payload}
}

WS_CONNECTION.onclose = function() {
    setTimeout(function() {
        if (WS_CONNECTION.readyState === WebSocket.CLOSED) {}
            window.location.reload()
    }, 2000)
}
WS_CONNECTION.onopen = function() {
    let myLoop = function() {
        WS_CONNECTION.send(constructMessage("heartbeat"))
        setTimeout(myLoop, 15*1000)
    }
    setTimeout(myLoop, 15*1000)
}

WS_CONNECTION.onmessage = function(event) {
    let message = parseMessage(event.data)
    switch (message.type) {
        case "getGrid":
            addOrReplace(getGridElement(message.payload[0][0]))
            break
        case "getGridChange":
            let grid = document.getElementById("grid")
            if (grid) {
                let cell = grid.children[Number(message.payload[0][0])*GRIDSIZE + Number(message.payload[0][1])]
                cell.style.backgroundColor = message.payload[0][2] === "true" ? "black" : "white"
            }
            break
        default:
            console.log("Unknown message type: " + message.type)
            console.log(message)    
            break
    }
}

function addOrReplace(element) {
    let existing = document.getElementById(element.id)
    if (existing) {
        existing.replaceWith(element)
    } else {
        document.body.appendChild(element)
    }
}

function getGridElement(grid) {
    let gridElement = document.createElement("div")
    gridElement.id = "grid"
    gridElement.style.display = "grid"
    gridElement.style.gridTemplateColumns = "repeat("+GRIDSIZE+", 1fr)"
    gridElement.style.gridTemplateRows = "repeat("+GRIDSIZE+", 1fr)"
    gridElement.style.width = SQUARESIZE*GRIDSIZE+"px"
    gridElement.style.height = SQUARESIZE*GRIDSIZE+"px"
    gridElement.style.border = "1px solid black"
    gridElement.style.margin = "auto"
    gridElement.style.marginTop = "10px"
    gridElement.style.marginBottom = "10px"
    gridElement.style.backgroundColor = "white"
    gridElement.style.padding = "0px"
    gridElement.style.boxSizing = "border-box"
    gridElement.style.position = "relative"
    gridElement.style.overflow = "hidden"
    gridElement.style.borderRadius = "5px"
    for (let i = 0; i < GRIDSIZE*GRIDSIZE; i++) {
        let cell = document.createElement("div")
        cell.style.width = SQUARESIZE+"px"
        cell.style.height = SQUARESIZE+"px"
        cell.style.backgroundColor = grid[i] === "1" ? "black" : "white"
        cell.style.border = "1px solid black"
        cell.style.boxSizing = "border-box"
        gridElement.appendChild(cell)
        cell.onclick = function() {
            WS_CONNECTION.send(constructMessage("requestGridChange", Math.floor(i/GRIDSIZE), i%GRIDSIZE))
        }
    }
    return gridElement
}