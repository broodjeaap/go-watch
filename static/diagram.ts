class DiagramNode {
    id: number;
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
        id: number,
        x: number, 
        y: number, 
        width: number,
        height: number,
        label: string,
    ){
        this.id = id;
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

    getInputCircleXY(){
        return [this.x, this.y + this.height / 3]
    }

    getOutputCircleXY(){
        return [this.x + this.width, this.y + this.height / 3]
    }

    pointInInputCircle(x: number, y: number){
        let [circleX, circleY] = this.getInputCircleXY()
        let radiusSqrd = Math.pow(this.height / 3, 2);
        return Math.pow(x - circleX, 2) + Math.pow(y - circleY, 2) <= radiusSqrd;
    }

    pointInOutputCircle(x: number, y: number){
        let [circleX, circleY] = this.getOutputCircleXY()
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
function diagramOnWheel(ev: WheelEvent){
    //_diagram.onwheel(ev);
}

class Diagrams {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;

    nodes: Array<DiagramNode> = new Array();

    connections: Array<[DiagramNode, DiagramNode]> = new Array();

    cameraX: number = 0; // camera position
    cameraY: number = 0;
    mouseX: number = 0; // mouse position on the canvas
    mouseY: number = 0;
    worldX: number = 0; // relative mouse position
    worldY: number = 0;

    panning: boolean = false;

    nodeDragging: DiagramNode | null = null;
    nodeHover: DiagramNode | null = null;

    makingConnectionNode: DiagramNode | null = null;

    scale: number = 1.0;

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
        this.ctx.font = "20px Helvetica";
        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        this.canvas.onwheel = diagramOnWheel;
        window.onresize = diargramOnResize;

    }

    onmousemove(ev: MouseEvent){
        let canvasRect = this.canvas.getBoundingClientRect();
        this.mouseX = ev.x - canvasRect.left;
        this.mouseY = ev.y - canvasRect.top;
        this.worldX = this.mouseX - this.cameraX;
        this.worldY = this.mouseY - this.cameraY;

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
                if (node.pointNearNode(this.worldX, this.worldY)){
                    if (node.pointInInputCircle(this.worldX, this.worldY)) {
                        node.inputHover = true;
                        this.nodeHover = node;
                        break;
                    } else if (this.makingConnectionNode == null && node.pointInOutputCircle(this.worldX, this.worldY)){
                        node.outputHover = true;
                        this.nodeHover = node;
                        break;
                    }else if (node.pointInNode(this.worldX, this.worldY)){
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
        this.mouseX = ev.x - canvasRect.left;
        this.mouseY = ev.y - canvasRect.top;
        this.worldX = this.mouseX - this.cameraX;
        this.worldY = this.mouseY - this.cameraY;
        for (let node of this.nodes){
            if (node.pointNearNode(this.worldX, this.worldY)){
                if (node.pointInInputCircle(this.worldX, this.worldY)) {
                    // no dragging from inputs ?
                } else if (node.pointInOutputCircle(this.worldX, this.worldY)){
                    this.makingConnectionNode = node;
                    return;
                }
            }
            if (node.pointInNode(this.worldX, this.worldY)) {
                this.nodeDragging = node;
                return;
            }
        }

        this.panning = true;
    }

    onmouseup(ev: MouseEvent){
        this.panning = false;
        this.nodeDragging = null;
        if (this.makingConnectionNode !== null){
            for (let node of this.nodes){
                if (node == this.makingConnectionNode){
                    continue;
                }
                if(node.pointInInputCircle(this.worldX, this.worldY)){
                    this.addConnection(this.makingConnectionNode, node);
                }
            }
            this.makingConnectionNode = null;
        }

        for (let [output, input] of this.connections){
            let [outputX, outputY] = output.getOutputCircleXY();
            outputX += this.cameraX;
            outputY += this.cameraY;
            let [inputX, inputY] = input.getInputCircleXY();
            inputX += this.cameraX;
            inputY += this.cameraY;
            let dX = Math.abs(outputX - inputX);
            this.ctx.beginPath();
            this.ctx.moveTo(outputX, outputY);
            this.ctx.strokeStyle = "black";
            let cp1x = (outputX + dX / 2);
            let cp1y = outputY;
            let cp2x = (inputX - dX / 2);
            let cp2y = inputY;
            this.ctx.bezierCurveTo(
                cp1x, 
                cp1y, 
                cp2x, 
                cp2y, 
                inputX, 
                inputY
            );
            this.ctx.stroke();
            this.ctx.closePath();
            let halfway = getBezierXY(0.5, outputX, outputY, cp1x, cp1y, cp2x, cp2y, inputX, inputY)
            let mouseOnHalfway = Math.pow(this.mouseX - halfway.x, 2) + Math.pow(this.mouseY - halfway.y, 2) <= 10*10
            if (mouseOnHalfway){
                this.connections.splice(this.connections.indexOf([output, input]), 1)
            }
        }
        
        this.draw();
    }

    onwheel(ev: WheelEvent){
        if(ev.deltaY > 0){
            return;
        }
        this.scale = Math.min(Math.max(this.scale - 0.1, 0.1), 1.0);
        this.ctx.scale(this.scale, this.scale);
    }

    drawBackground(){
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0,0,this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    }

    draw(){
        let scale = 1 / this.scale;
        this.ctx.clearRect(0,0, this.canvas.width * scale, this.canvas.height * scale);
        this.drawBackground();
        let fullCircleRadians = Math.PI + (Math.PI * 3);
        if (this.makingConnectionNode != null){
            let [circleX, circleY] = this.makingConnectionNode.getOutputCircleXY();
            let dX = Math.abs((circleX + this.cameraX) - this.mouseX);
            this.ctx.beginPath();
            this.ctx.moveTo(circleX + this.cameraX, circleY + this.cameraY);
            this.ctx.strokeStyle = "black";
            let cp1x = (circleX + dX / 2) + this.cameraX;
            let cp1y = circleY + this.cameraY;
            let cp2x = (this.mouseX - dX / 2);
            let cp2y = this.mouseY;
            this.ctx.bezierCurveTo(
                cp1x, 
                cp1y, 
                cp2x, 
                cp2y, 
                this.mouseX, 
                this.mouseY
            );
            this.ctx.stroke();
            this.ctx.closePath();
        }
        for (let [output, input] of this.connections){
            let [outputX, outputY] = output.getOutputCircleXY();
            outputX += this.cameraX;
            outputY += this.cameraY;
            let [inputX, inputY] = input.getInputCircleXY();
            inputX += this.cameraX;
            inputY += this.cameraY;
            let dX = Math.abs(outputX - inputX);
            this.ctx.beginPath();
            this.ctx.moveTo(outputX, outputY);
            this.ctx.strokeStyle = "black";
            let cp1x = (outputX + dX / 2);
            let cp1y = outputY;
            let cp2x = (inputX - dX / 2);
            let cp2y = inputY;
            this.ctx.bezierCurveTo(
                cp1x, 
                cp1y, 
                cp2x, 
                cp2y, 
                inputX, 
                inputY
            );
            this.ctx.stroke();
            this.ctx.closePath();
            let halfway = getBezierXY(0.5, outputX, outputY, cp1x, cp1y, cp2x, cp2y, inputX, inputY)
            let mouseOnHalfway = Math.pow(this.mouseX - halfway.x, 2) + Math.pow(this.mouseY - halfway.y, 2) <= 10*10
            this.ctx.beginPath();
            this.ctx.strokeStyle = mouseOnHalfway ? "red" : "rgba(200, 200, 200, 0.8)";
            this.ctx.moveTo(halfway.x - 10, halfway.y - 10);
            this.ctx.lineTo(halfway.x + 10, halfway.y + 10);
            this.ctx.moveTo(halfway.x + 10, halfway.y - 10);
            this.ctx.lineTo(halfway.x - 10, halfway.y + 10);
            this.ctx.stroke();
            this.ctx.closePath();
        }
        for (let node of this.nodes){
            this.ctx.fillStyle = node.hover ? "#303030" : "#161616";
            this.ctx.fillRect(node.x + this.cameraX, node.y + this.cameraY, node.width, node.height);
            this.ctx.fillStyle = "#D3D3D3";
            this.ctx.font = "20px Helvetica";
            this.ctx.fillText(
                node.label, 
                node.x + this.cameraX + node.height / 2, 
                node.y + this.cameraY + node.height / 1.5
            );
            this.ctx.fillStyle = node.inputHover ? "red" : "green";
            this.ctx.beginPath()
            this.ctx.arc(node.x + this.cameraX, node.y + node.height / 2 + this.cameraY, node.height / 3, 0, fullCircleRadians);
            this.ctx.fill();
            this.ctx.beginPath()
            this.ctx.fillStyle = node.outputHover ? "red" : "green";
            this.ctx.moveTo(node.x + node.width + this.cameraX, node.y + node.height / 2 + this.cameraY);
            this.ctx.arc(node.x + node.width + this.cameraX, node.y + node.height / 2 + this.cameraY, node.height / 3, 0, fullCircleRadians);
            this.ctx.fill();
        }
    }

    addNode(id: number, x: number, y: number, label: string){
        let textSize = this.ctx.measureText(label);
        let textHeight = 2 * (textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent);
        this.nodes.push(new DiagramNode(id, x, y, textSize.width + textHeight, textHeight, label));
    }

    addConnection(A: DiagramNode, B: DiagramNode){
        this.connections.push([A, B]);
    }

    drawDiagramNode(node: DiagramNode){
        
    }

    fillParent(){
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        this.draw();
    }
}

// http://www.independent-software.com/determining-coordinates-on-a-html-canvas-bezier-curve.html
function getBezierXY(t, sx, sy, cp1x, cp1y, cp2x, cp2y, ex, ey) {
    return {
      x: Math.pow(1-t,3) * sx + 3 * t * Math.pow(1 - t, 2) * cp1x 
        + 3 * t * t * (1 - t) * cp2x + t * t * t * ex,
      y: Math.pow(1-t,3) * sy + 3 * t * Math.pow(1 - t, 2) * cp1y 
        + 3 * t * t * (1 - t) * cp2y + t * t * t * ey
    };
  }