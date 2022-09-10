class DiagramNode {
    x: number;
    y: number;
    label: string;
    width: number;
    height: number;

    parents: Array<DiagramNode>;
    children: Array<DiagramNode>;

    constructor(
        x: number, 
        y: number, 
        width: number,
        height: number,
        label: string,
    ){
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.label = label;
    }

    pointInDiagram(x: number, y: number){
        if (x < this.x){
            return false;
        }
        if (y < this.y){
            return false;
        }
        if (x > this.x + this.width) {
            return false;
        }
        if (y > this.y + this.height) {
            return false
        }
        return true;
    }
}

let _diagram: Diagrams;
function diargramOnResize(){
    _diagram.fillParent();
}
function diagramOnMouseDown(ev: MouseEvent){
    _diagram.onmousedown(ev);
}
function diagramOnMouseUp(ev: MouseEvent){
    _diagram.onmouseup(ev);
}
function diagramOnMouseMove(ev: MouseEvent){
    _diagram.onmousemove(ev);
}

class Diagrams {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;

    nodes: Array<DiagramNode>;

    connections: Array<[DiagramNode, DiagramNode]>;

    cameraX: number;
    cameraY: number;

    panning: boolean;

    nodeDrag: boolean;
    nodeDragged: DiagramNode | null;

    constructor(canvasId: string){
        this.canvas = document.getElementById(canvasId) as HTMLCanvasElement;
        if (this.canvas === null){
            throw `Could not getElementById ${canvasId}`;
        }
        let ctx = this.canvas.getContext("2d");
        if (ctx === null){
            throw `Could not get 2d rendering context`
        }
        _diagram = this;
        this.ctx = ctx;
        this.ctx.font = "30px Arial";
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        window.onresize = diargramOnResize;

        this.nodes = new Array();
        this.connections = new Array();

        this.cameraX = 0;
        this.cameraY = 0;
    }

    onmousemove(ev: MouseEvent){
        if (this.panning){
            this.cameraX += ev.movementX;
            this.cameraY += ev.movementY;
        }
        if (this.nodeDrag){
            if (this.nodeDragged === null){
                console.error("nodeDrag==true but nodeDragged==null");
                return
            }
            this.nodeDragged.x += ev.movementX;
            this.nodeDragged.y += ev.movementY;
        }
        this.draw();
    }

    onmousedown(ev: MouseEvent){
        if (ev.button != 0){
            return;
        }
        let canvasRect = this.canvas.getBoundingClientRect();
        let mouseX = ev.x - canvasRect.left;
        let mouseY = ev.y - canvasRect.top;
        for (let node of this.nodes){
            if (node.pointInDiagram(mouseX, mouseY)) {
                this.nodeDrag = true;
                this.nodeDragged = node;
                return;
            }
        }

        this.panning = true;
    }

    onmouseup(ev: MouseEvent){
        this.panning = false;
        this.nodeDrag = false;
        this.nodeDragged = null;
    }

    drawBackground(){
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0,0,this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    }

    draw(){
        this.ctx.clearRect(0,0, this.canvas.width, this.canvas.height);
        this.drawBackground();
        for (let node of this.nodes){
            this.ctx.fillStyle = "gray";
            this.ctx.fillRect(node.x + this.cameraX, node.y + this.cameraY, node.width, node.height);
            this.ctx.fillStyle = "black";
            this.ctx.font = "30px Arial";
            this.ctx.fillText(
                node.label, 
                node.x + this.cameraX + node.height / 2, 
                node.y + this.cameraY + node.height / 1.5
            );
        }
    }

    addNode(x: number, y: number, label: string){
        let textSize = this.ctx.measureText(label);
        let textHeight = 2 * (textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent);
        this.nodes.push(new DiagramNode(x, y, textSize.width + textHeight, textHeight, label));
    }

    addConnection(A: DiagramNode, B: DiagramNode){
        this.connections.push([A, B]);
    }

    fillParent(){
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        this.draw();
    }
}