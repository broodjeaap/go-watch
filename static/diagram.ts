class DiagramNode {
    x: number;
    y: number;
    label: string;
    width: number;
    height: number;

    hover: boolean = false;
    inputHover: boolean = false;
    outputHover: boolean = false;

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

    pointInNode(x: number, y: number){
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
            return false;
        }
        return true;
    }

    pointNearNode(x: number, y: number){
        // including the input/output circles
        if (x < this.x - this.height / 3){
            return false;
        }
        if (y < this.y){
            return false;
        }
        if (x > this.x + this.width + this.height / 3){
            return false;
        }
        if (y > this.y + this.height) {
            return false;
        }
        return true;
    }

    pointInInputCircle(x: number, y: number){
        let circleX = this.x;
        let circleY = this.y + this.height / 3;
        let radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
    }

    pointInOutputCircle(x: number, y: number){
        let circleX = this.x + this.width;
        let circleY = this.y + this.height / 2;
        let radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
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

    nodes: Array<DiagramNode> = new Array();

    connections: Array<[DiagramNode, DiagramNode]> = new Array();

    cameraX: number = 0;
    cameraY: number = 0;

    panning: boolean = false;

    nodeDragging: DiagramNode | null = null;
    nodeHover: DiagramNode | null = null;

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
        this.ctx.font = "30px Helvetica";
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        window.onresize = diargramOnResize;
    }

    onmousemove(ev: MouseEvent){
        let canvasRect = this.canvas.getBoundingClientRect();
        let mouseX = ev.x - canvasRect.left;
        let mouseY = ev.y - canvasRect.top;
        let worldX = mouseX - this.cameraX;
        let worldY = mouseY - this.cameraY;

        if (this.nodeHover != null){
            this.nodeHover.hover = false;
            this.nodeHover.inputHover = false
            this.nodeHover.outputHover = false
            this.nodeHover = null;
        }

        if (this.panning){
            this.cameraX += ev.movementX;
            this.cameraY += ev.movementY;
        }
        else if (this.nodeDragging != null){
            // this.nodeDragging.x = worldX - this.nodeDragging.width / 2;
            // this.nodeDragging.y = worldY - this.nodeDragging.height / 2;
            this.nodeDragging.x += ev.movementX;
            this.nodeDragging.y += ev.movementY;
        } else {
            for (let node of this.nodes){
                if (node.pointNearNode(worldX, worldY)){
                    if (node.pointInInputCircle(worldX, worldY)) {
                        node.inputHover = true;
                        this.nodeHover = node;
                        break;
                    } else if (node.pointInOutputCircle(worldX, worldY)){
                        node.outputHover = true;
                        this.nodeHover = node;
                        break;
                    }else if (node.pointInNode(worldX, worldY)){
                        node.hover = true;
                        this.nodeHover = node;
                        break;
                    }
                }
            }
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
        let worldX = mouseX - this.cameraX;
        let worldY = mouseY - this.cameraY;
        for (let node of this.nodes){
            if (node.pointNearNode(worldX, worldY)){
                if (node.pointInInputCircle(worldX, worldY)) {
                    
                } else if (node.pointInOutputCircle(worldX, worldY)){

                }
            }
            if (node.pointInNode(worldX, worldY)) {
                this.nodeDragging = node;
                return;
            }
        }

        this.panning = true;
    }

    onmouseup(ev: MouseEvent){
        this.panning = false;
        this.nodeDragging = null;
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
        let fullCircleRadians = Math.PI + (Math.PI * 3);
        for (let node of this.nodes){
            this.ctx.fillStyle = node.hover ? "#303030" : "#161616";
            this.ctx.fillRect(node.x + this.cameraX, node.y + this.cameraY, node.width, node.height);
            this.ctx.fillStyle = "#D3D3D3";
            this.ctx.font = "30px Helvetica";
            this.ctx.fillText(
                node.label, 
                node.x + this.cameraX + node.height / 2, 
                node.y + this.cameraY + node.height / 1.5
            );
            this.ctx.fillStyle = node.inputHover ? "red" : "green";
            this.ctx.beginPath()
            this.ctx.arc(node.x + this.cameraX, node.y + node.height / 2 + this.cameraY, node.height / 3, 0, fullCircleRadians);
            this.ctx.fill();
            this.ctx.closePath();
            this.ctx.beginPath()
            this.ctx.fillStyle = node.outputHover ? "red" : "green";
            this.ctx.moveTo(node.x + node.width + this.cameraX, node.y + node.height / 2 + this.cameraY);
            this.ctx.arc(node.x + node.width + this.cameraX, node.y + node.height / 2 + this.cameraY, node.height / 3, 0, fullCircleRadians);
            this.ctx.fill();
            this.ctx.closePath();
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