class DragPosition {
    x: number;
    y: number;
    oldX: number;
    oldY: number;
    elem: HTMLElement;

    constructor(e: HTMLElement){
        this.x = 0;
        this.y = 0;
        this.oldX = 0;
        this.oldY = 0;
        this.elem = e;
    }
}
var position: DragPosition;

function startDrag(e: MouseEvent){
    e = e || window.Event
    e.preventDefault()
    var elem = e.target as HTMLElement;
    while (!elem.classList.contains("node-card")){
        if (elem.parentElement === null){
            return
        }
        elem = elem.parentElement;
    }
    position = new DragPosition(elem);
    position.oldX = e.clientX;
    position.oldY = e.clientY;
    document.onmouseup = stopDrag;
    document.onmousemove = dragging;
}

function dragging(e: MouseEvent){
    e = e || window.Event
    e.preventDefault();

    position.x = position.oldX - e.clientX;
    position.y = position.oldY - e.clientY;
    position.oldX = e.clientX;
    position.oldY = e.clientY;
    position.elem.style.top = `${(position.elem.offsetTop - position.y)}px`;
    position.elem.style.left = `${(position.elem.offsetLeft - position.x)}px`;
}

function stopDrag(e: MouseEvent){
    document.onmouseup = null;
    document.onmousemove = null;
}

function log(e: MouseEvent){
    e.stopPropagation();
    console.log(e.target);
}

document.addEventListener("DOMContentLoaded", function(event) { 
    let node_cards = document.getElementsByClassName("node-card");
    [].forEach.call(node_cards, function(node_card: HTMLElement) {
        node_card.onmousedown = startDrag;
    });

    let inputs = document.getElementsByClassName("node-input");
    [].forEach.call(inputs, function(input: HTMLElement){
        input.onmousedown = log;
    });

    let c = document.getElementById("node-canvas") as HTMLCanvasElement;
    let ctx = c.getContext("2d");
    if (ctx === null){
        return
    }
    ctx.beginPath();
    ctx.moveTo(50, 50);
    ctx.lineTo(1000, 1000);
    ctx.lineWidth = 10;

    // set line color
    ctx.strokeStyle = '#ff0000';
    ctx.stroke();
});