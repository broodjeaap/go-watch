abstract class CanvasObject {
    x: number;
    y: number;
    width: number;
    height: number;

    hover: boolean

    constructor(x: number, y: number, width: number, height: number){
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.hover = false;
    }

    abstract update(ms: MouseState): void;
    abstract draw(ctx: CanvasRenderingContext2D, ms: MouseState): void;

    pointInObject(p: Point): boolean{
        if (p.x < this.x){
            return false;
        }
        if (p.y < this.y){
            return false;
        }
        if (p.x > this.x + this.width) {
            return false;
        }
        if (p.y > this.y + this.height) {
            return false;
        }
        return true;
    }
}

class Button extends CanvasObject {
    label: string;
    labelWidth: number;
    labelHeight: number;

    constructor(x: number, y: number, label: string, ctx: CanvasRenderingContext2D){
        super(x, y, 0, 0);
        this.label = label;
        this.resize(ctx);
    }

    update(ms: MouseState){
        this.hover = this.pointInObject(new Point(ms.world.x + ms.offset.x, ms.world.y + ms.offset.y));
    }

    draw(ctx: CanvasRenderingContext2D, ms: MouseState){
        ctx.fillStyle = this.hover ? "black" : "#6B6B6B";
        ctx.font = "15px Helvetica";
        ctx.fillText(this.label, this.x + 3, this.y + this.labelHeight + 3);
    }

    resize(ctx: CanvasRenderingContext2D){
        ctx.font = "15px Helvetica";
        let labelSize = ctx.measureText(this.label);
        this.labelWidth = labelSize.width;
        this.width = this.labelWidth + 6;
        this.labelHeight = labelSize.actualBoundingBoxAscent + labelSize.actualBoundingBoxDescent;
        this.height = this.labelHeight + 6;
    }
}

const circleTopRadians = Math.PI / 2;
const circleRightRadians = (Math.PI * 3) / 2;
const circleBottomRadians = Math.PI + (Math.PI * 3);
const circleLeftRadians = Math.PI;
class NodeIO extends CanvasObject {
    node: DiagramNode;
    input: boolean = false;
    radius: number = 15;

    constructor(node: DiagramNode, input: boolean){
        super(0,0,0,0);
        this.input = input
        this.node = node;
        this.reposition();
    }
    update(ms: MouseState): void {
        
    }

    draw(ctx: CanvasRenderingContext2D, ms: MouseState): void {
        ctx.fillStyle = this.input ? "red" : "blue";
        ctx.beginPath();
        ctx.arc(ms.offset.x + this.x, ms.offset.y + this.y, this.radius, circleRightRadians, circleTopRadians, this.input);
        ctx.fill();
    }

    reposition(){
        if (this.input){
            this.x = this.node.x;
            this.y = this.node.y + this.node.height / 2;
        } else {
            this.x = this.node.x + this.node.width;
            this.y = this.node.y + this.node.height / 2;
        }
    }
}

class NodeConnection extends CanvasObject {
    output: DiagramNode;
    input: DiagramNode;

    constructor(output: DiagramNode, input: DiagramNode){
        super(0, 0, 0, 0);
        this.output = output;
        this.input = input;
    }

    update(ms: MouseState): void {
        
    }
    draw(ctx: CanvasRenderingContext2D, ms: MouseState): void {
        let outputX = ms.offset.x + this.output.output.x;
        let outputY = ms.offset.y + this.output.output.y;
        let inputX = ms.offset.x + this.input.input.x;
        let inputY = ms.offset.y + this.input.input.y;
        let dX = Math.abs(outputX - inputX);
        ctx.beginPath();
        ctx.moveTo(outputX, outputY);
        ctx.strokeStyle = "#757575";
        ctx.lineWidth = 5;
        let cp1x = (outputX + dX);
        let cp1y = outputY;
        let cp2x = (inputX - dX);
        let cp2y = inputY;
        ctx.bezierCurveTo(
            cp1x, 
            cp1y, 
            cp2x, 
            cp2y, 
            inputX, 
            inputY
        );
        ctx.stroke();
        ctx.closePath();
        /*
        let halfway = getBezierXY(0.5, outputX, outputY, cp1x, cp1y, cp2x, cp2y, inputX, inputY)
        let mouseOnHalfway = Math.pow(this.mouseX - halfway.x, 2) + Math.pow(this.mouseY - halfway.y, 2) <= 10*10
        if (mouseOnHalfway){
            this.removeConnection(output, input);
            break;
        }
        */
    }

    reposition(){
        
    }
}

class DiagramNode extends CanvasObject {
    id: number;
    label: string;
    type: string;
    
    labelWidth: number;
    labelHeight: number;

    typeWidth: number;
    typeHeight: number;

    deleteButton: Button;
    editButton: Button;
    logButton: Button;

    dragging: boolean = false;
    dragOrigin: Point = new Point();

    input: NodeIO;
    output: NodeIO;

    parents: Array<DiagramNode>;
    children: Array<DiagramNode>;

    meta: Object = {};

    constructor(
        id: number,
        x: number, 
        y: number, 
        label: string,
        meta: Object = {},
        ctx: CanvasRenderingContext2D
    ){
        super(x, y, 0, 0)
        this.id = id;
        this.label = label;
        this.meta = meta;
        // @ts-ignore
        this.type = this.meta.type
        if (["math", "condition"].indexOf(this.type) >= 0 ){
            // @ts-ignore
            this.type = this.meta.var1
        }
        this.resize(ctx);

        this.deleteButton = new Button(0, 0, "Del", ctx);
        this.editButton = new Button(0, 0, "Edit", ctx);
        this.logButton = new Button(0, 0, "Log", ctx);

        this.input = new NodeIO(this, true);
        this.output = new NodeIO(this, false);

    }

    update(ms: MouseState) {
        this.hover = (!ms.draggingNode || this.dragging) && this.pointInObject(ms.world);
        if (this.hover){
            this.deleteButton.update(ms);
            this.editButton.update(ms);
            this.logButton.update(ms)
            let onButtons = this.deleteButton.hover || this.editButton.hover || this.logButton.hover;
            if (!this.dragging && ms.leftDown && !ms.draggingNode && !ms.draggingConnection && !onButtons){
                this.dragging = true;
                ms.draggingNode = true;
                this.dragOrigin.x = this.x - ms.world.x;
                this.dragOrigin.y = this.y - ms.world.y;
            }
        } else {
            this.deleteButton.hover = false;
            this.editButton.hover = false;
            this.logButton.hover = false;
        }

        if (!ms.leftDown){
            this.dragging = false;
            ms.draggingNode = false;
        }
        if (this.dragging){
            this.x = ms.world.x + this.dragOrigin.x;
            this.y = ms.world.y + this.dragOrigin.y;
            this.input.reposition();
            this.output.reposition();
        }
        this.input.update(ms);
        this.output.update(ms);
    }

    draw(ctx: CanvasRenderingContext2D, ms: MouseState){
        ctx.fillStyle = this.hover ? "#DDDDDD" : "#BFBFBF";
        ctx.fillRect(ms.offset.x + this.x, ms.offset.y + this.y, this.width, this.height);
        
        ctx.font = "20px Helvetica";
        ctx.fillStyle = "black";
        let labelX = ms.offset.x + this.x + this.width / 2 - this.labelWidth / 2;
        let labelY = ms.offset.y +this.y + 3 * 2 + this.labelHeight;
        ctx.fillText(this.label, labelX, labelY);
        
        ctx.font = "15px Helvetica";
        let typeX = ms.offset.x + this.x + this.width / 2 - this.typeWidth / 2;
        let typeY = ms.offset.y + this.y + 3 * 3 + this.typeHeight + this.labelHeight;
        ctx.fillText(this.type, typeX, typeY);
        
        this.deleteButton.x = ms.offset.x + this.x;
        this.deleteButton.y = ms.offset.y + this.y + this.height - this.deleteButton.height;
        this.deleteButton.draw(ctx, ms);

        this.editButton.x = ms.offset.x + this.x + this.width - this.editButton.width;
        this.editButton.y = ms.offset.y + this.y + this.height - this.editButton.height;
        this.editButton.draw(ctx, ms);

        this.logButton.x = ms.offset.x + this.x + this.width / 2 - this.logButton.width / 2;
        this.logButton.y = ms.offset.y + this.y + this.height - this.logButton.height;
        this.logButton.draw(ctx, ms);
        
        this.input.draw(ctx, ms);
        this.output.draw(ctx, ms);
        
        ctx.strokeStyle = "#8E8E8E";
        ctx.lineWidth = 3;
        ctx.strokeRect(ms.offset.x + this.x, ms.offset.y + this.y, this.width, this.height);
    }

    resize(ctx: CanvasRenderingContext2D){
        ctx.font = "20px Helvetica";
        let labelSize = ctx.measureText(this.label);
        this.labelWidth = labelSize.width;
        this.labelHeight = labelSize.actualBoundingBoxAscent + labelSize.actualBoundingBoxDescent;
        this.height = 70;

        ctx.font = "15px Helvetica";
        let typeSize = ctx.measureText(this.type);
        this.typeWidth = typeSize.width;
        this.typeHeight = typeSize.actualBoundingBoxAscent + typeSize.actualBoundingBoxDescent;
        
        this.width = Math.max(150, this.labelWidth, this.typeWidth);
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
        return [this.x, this.y + this.height / 2]
    }

    getOutputCircleXY(){
        return [this.x + this.width, this.y + this.height / 2]
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
function tick(){
    _diagram.tick();
    setTimeout(() => {
        tick(), 1000/60;
    });
}
function diagramOnResize(){
    _diagram.onresize();
}
function diagramOnMouseDown(ev: MouseEvent){
    _diagram.onmousedown(ev)
}
function diagramOnMouseUp(ev: MouseEvent){
    _diagram.onmouseup(ev);
}
function diagramOnMouseMove(ev: MouseEvent){
    _diagram.onmousemove(ev)
}
function diagramOnContext(ev: MouseEvent){
    ev.preventDefault();
}

class Point {
    x: number = 0;
    y: number = 0;
    constructor(x: number = 0, y: number = 0){
        this.x = x;
        this.y = y;
    }
}
class MouseState {
    canvas: Point = new Point();
    world: Point = new Point();
    offset: Point = new Point();
    delta: Point = new Point();
    leftDown: boolean = false;
    leftUp: boolean = false;
    panning: boolean = false;
    draggingNode: boolean = false;
    draggingConnection: boolean = false;
}

class Diagrams {
    canvas: HTMLCanvasElement;
    ctx: CanvasRenderingContext2D;
    shouldTick: boolean = true;

    nodes: Map<number, DiagramNode> = new Map();

    connections: Array<NodeConnection> = new Array();

    mouseState: MouseState = new MouseState();


    panning: boolean = false;

    nodeDragging: DiagramNode | null = null;
    nodeHover: DiagramNode | null = null;

    makingConnectionNode: DiagramNode | null = null;

    scale: number = 1.0;

    editNodeCallback: (node: DiagramNode) => void = function (){};
    deleteNodeCallback: (node: DiagramNode) => void = function (){};

    constructor(
            canvasId: string, 
            editNodeCallback: (node: DiagramNode) => void = function (){},
            deleteNodeCallback: (node: DiagramNode) => void = function (){}
        ){
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
        this.editNodeCallback = editNodeCallback;
        this.deleteNodeCallback = deleteNodeCallback;

        this.canvas.onmousemove = diagramOnMouseMove;
        this.canvas.onmousedown = diagramOnMouseDown;
        this.canvas.onmouseup = diagramOnMouseUp;
        window.onresize = diagramOnResize;

        tick();
    }

    tick(){
        this.drawBackground();
        for (let node of this.nodes.values()){
            node.update(this.mouseState);
        }
        for (let connection of this.connections){
            connection.update(this.mouseState);
        }
        for (let connection of this.connections){
            connection.draw(this.ctx, this.mouseState);
        }
        for (let node of this.nodes.values()){
            node.draw(this.ctx, this.mouseState);
        }
        if (this.mouseState.leftUp){
            console.log("Click");
        }
        this.mouseState.leftUp = false;
    }

    onmousemove(ev: MouseEvent){
        let canvasRect = this.canvas.getBoundingClientRect();
        this.mouseState.canvas.x = ev.x - canvasRect.left;
        this.mouseState.canvas.y = ev.y - canvasRect.top;
        this.mouseState.delta.x = ev.movementX;
        this.mouseState.delta.y = ev.movementY;

        if (this.mouseState.panning){
            this.mouseState.offset.x += this.mouseState.delta.x;
            this.mouseState.offset.y += this.mouseState.delta.y;
        }

        this.mouseState.world.x = this.mouseState.canvas.x - this.mouseState.offset.x;
        this.mouseState.world.y = this.mouseState.canvas.y - this.mouseState.offset.y;
    }

    onmousedown(ev: MouseEvent){
        if (ev.button != 0){
            return;
        }
        this.mouseState.leftDown = true;
        for (let object of this.nodes.values()){
            if (object.pointInObject(this.mouseState.world)) {
                return;
            }
        }
        this.mouseState.panning = true;
    }

    onmouseup(ev: MouseEvent){
        this.mouseState.leftDown = false;
        this.mouseState.panning = false;
        this.mouseState.leftUp = true;
    }

    drawBackground(){
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0,0,this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    }

    draw(){

    }

    addNode(id: number, x: number, y: number, label: string, meta: Object = {}){
        let node = new DiagramNode(id, x, y, label, meta, this.ctx);
        this.nodes.set(id, node);
    }

    addConnection(A: DiagramNode, B: DiagramNode){
        this.connections.push(new NodeConnection(A, B));
    }
    addConnectionById(a: number, b: number){
        let A = this.nodes.get(a);
        if (A === undefined){
            console.error(`No node with ID: ${a}`);
            return;
        }
        let B = this.nodes.get(b);
        if (B === undefined){
            console.error(`No node with ID: ${b}`);
            return;
        }
        this.connections.push(new NodeConnection(A, B))
    }
    removeConnection(A: DiagramNode, B: DiagramNode){
        let index = 0;
        for (let connection of this.connections){
            let output = connection.output;
            let input = connection.input;
            if (output.id == A.id && input.id == B.id) {
                this.connections.splice(index, 1);
            }
            index++;
        }
    }

    onresize(){
        this.fillParent();
    }

    fillParent(){
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        //this.draw();
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