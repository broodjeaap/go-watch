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

    callback: (node: DiagramNode) => void = function (){};
    node: DiagramNode;

    constructor(
            x: number, 
            y: number, 
            label: string, 
            ctx: CanvasRenderingContext2D, 
            callback: (node: DiagramNode) => void,
            node: DiagramNode,
        ){
        super(x, y, 0, 0);
        this.label = label;
        this.callback = callback;
        this.node = node;
        this.resize(ctx);
    }

    update(ms: MouseState){
        this.hover = this.pointInObject(new Point(ms.world.x + ms.offset.x, ms.world.y + ms.offset.y));
        if (ms.click && this.hover){
            this.callback(this.node);
            ms.click = false;
        }
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
        if (!ms.draggingConnection && !this.input && this.pointInObject(ms.world) && ms.leftDown){
            ms.draggingConnection = true;
            _diagram.newConnection = new NewConnection(this.node);
        }
    }

    draw(ctx: CanvasRenderingContext2D, ms: MouseState): void {
        ctx.fillStyle = this.input ? "#ED575A" : "#66A7C5";
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

    pointInObject(p: Point): boolean {
        let inCircle = Math.pow(p.x - this.x, 2) + Math.pow(p.y - this.y, 2) <= this.radius * this.radius;
        if (!inCircle){
            this.hover = false;
        } else {
            this.hover = this.input ? p.x < this.x : p.x > this.x;
        }
        return this.hover;
    }
}

class NodeConnection extends CanvasObject {
    output: DiagramNode;
    input: DiagramNode;
    
    controlPoints = {
        dX: 0,
        outputX: 0,
        outputY: 0,
        inputX: 0,
        inputY: 0,
        cp1x: 0,
        cp1y: 0,
        cp2x: 0,
        cp2y: 0,
    }

    halfWayPoint: Point = new Point();

    constructor(output: DiagramNode, input: DiagramNode){
        super(0, 0, 0, 0);
        this.output = output;
        this.input = input;
    }

    update(ms: MouseState): void {
        this.controlPoints.outputX = ms.offset.x + this.output.output.x;
        this.controlPoints.outputY = ms.offset.y + this.output.output.y;
        this.controlPoints.inputX = ms.offset.x + this.input.input.x;
        this.controlPoints.inputY = ms.offset.y + this.input.input.y;
        this.controlPoints.dX = Math.abs(this.controlPoints.outputX - this.controlPoints.inputX);
        
        this.controlPoints.cp1x = (this.controlPoints.outputX + this.controlPoints.dX);
        this.controlPoints.cp1y = this.controlPoints.outputY;
        this.controlPoints.cp2x = (this.controlPoints.inputX - this.controlPoints.dX);
        this.controlPoints.cp2y = this.controlPoints.inputY;
        
        this.halfWayPoint = getBezierXY(
            0.5, 
            this.controlPoints.outputX, 
            this.controlPoints.outputY, 
            this.controlPoints.cp1x, 
            this.controlPoints.cp1y, 
            this.controlPoints.cp2x, 
            this.controlPoints.cp2y, 
            this.controlPoints.inputX, 
            this.controlPoints.inputY
        );
        this.hover = Math.pow(this.halfWayPoint.x - ms.canvas.x, 2) + Math.pow(this.halfWayPoint.y - ms.canvas.y, 2) <= 15*15;
        if (this.hover && ms.click){
            _diagram.removeConnection(this.output, this.input);
            ms.click = false;
        }
    }
    draw(ctx: CanvasRenderingContext2D, ms: MouseState): void {
        ctx.beginPath();
        ctx.moveTo(this.controlPoints.outputX, this.controlPoints.outputY);
        ctx.strokeStyle = "#757575";
        ctx.lineWidth = 5;
        ctx.bezierCurveTo(
            this.controlPoints.cp1x, 
            this.controlPoints.cp1y, 
            this.controlPoints.cp2x, 
            this.controlPoints.cp2y, 
            this.controlPoints.inputX, 
            this.controlPoints.inputY
        );
        ctx.stroke();
        ctx.closePath();
        
        ctx.beginPath();
        ctx.strokeStyle = this.hover ? "red" : "rgba(200, 200, 200, 0.8)";
        ctx.moveTo(this.halfWayPoint.x - 10, this.halfWayPoint.y - 10);
        ctx.lineTo(this.halfWayPoint.x + 10, this.halfWayPoint.y + 10);
        ctx.moveTo(this.halfWayPoint.x + 10, this.halfWayPoint.y - 10);
        ctx.lineTo(this.halfWayPoint.x - 10, this.halfWayPoint.y + 10);
        ctx.stroke();
        ctx.closePath();
    }

    reposition(){
        
    }
}

class NewConnection extends CanvasObject {
    output: DiagramNode;
    input: DiagramNode | null;
    
    controlPoints = {
        dX: 0,
        outputX: 0,
        outputY: 0,
        inputX: 0,
        inputY: 0,
        cp1x: 0,
        cp1y: 0,
        cp2x: 0,
        cp2y: 0,
    }
    constructor(output: DiagramNode){
        super(0, 0, 0, 0);
        this.output = output;
    }

    update(ms: MouseState): void {
        this.input = null;
        for (let node of _diagram.nodes.values()){
            if (this.output.id != node.id && node.pointNearNode(ms.world)){
                this.input = node;
            }
        }
        
        if (this.input == null){
            this.controlPoints.outputX = ms.offset.x + this.output.output.x;
            this.controlPoints.outputY = ms.offset.y + this.output.output.y;
            this.controlPoints.inputX = ms.offset.x + ms.world.x;
            this.controlPoints.inputY = ms.offset.y + ms.world.y;
            this.controlPoints.dX = Math.abs(this.controlPoints.outputX - this.controlPoints.inputX);
        } else {
            this.controlPoints.outputX = ms.offset.x + this.output.output.x;
            this.controlPoints.outputY = ms.offset.y + this.output.output.y;
            this.controlPoints.inputX = ms.offset.x + this.input.input.x;
            this.controlPoints.inputY = ms.offset.y + this.input.input.y;
            this.controlPoints.dX = Math.abs(this.controlPoints.outputX - this.controlPoints.inputX);
        }
        
        this.controlPoints.cp1x = (this.controlPoints.outputX + this.controlPoints.dX);
        this.controlPoints.cp1y = this.controlPoints.outputY;
        this.controlPoints.cp2x = (this.controlPoints.inputX - this.controlPoints.dX);
        this.controlPoints.cp2y = this.controlPoints.inputY;
    }
    draw(ctx: CanvasRenderingContext2D, ms: MouseState): void {
        ctx.beginPath();
        ctx.moveTo(this.controlPoints.outputX, this.controlPoints.outputY);
        ctx.strokeStyle = "#7575A5";
        ctx.lineWidth = 5;
        ctx.bezierCurveTo(
            this.controlPoints.cp1x, 
            this.controlPoints.cp1y, 
            this.controlPoints.cp2x, 
            this.controlPoints.cp2y, 
            this.controlPoints.inputX, 
            this.controlPoints.inputY
        );
        ctx.stroke();
        ctx.closePath();
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

    dragging: boolean = false;
    dragOrigin: Point = new Point();

    input: NodeIO;
    output: NodeIO;

    parents: Array<DiagramNode>;
    children: Array<DiagramNode>;

    meta: Object = {};
    results: Array<string>;
    logs: Array<string>;

    constructor(
        id: number,
        x: number, 
        y: number, 
        label: string,
        meta: Object = {},
        ctx: CanvasRenderingContext2D,
        results: Array<string> = new Array(),
        logs: Array<string> = new Array(),
    ){
        super(x, y, 0, 0)
        this.id = id;
        this.label = label;
        this.meta = meta;
        this.fixType();
        this.resize(ctx);

        this.deleteButton = new Button(0, 0, "Del", ctx, _diagram.deleteNodeCallback, this);
        this.editButton = new Button(0, 0, "Edit", ctx, _diagram.editNodeCallback, this);

        this.input = new NodeIO(this, true);
        this.output = new NodeIO(this, false);
        this.results = results;
        this.logs = logs;
    }

    update(ms: MouseState) {
        if (this.pointNearNode(ms.world)){
            this.input.update(ms);
            this.output.update(ms);
        }
        this.hover = (!ms.draggingNode || this.dragging) && super.pointInObject(ms.world);
        if (this.hover){
            this.deleteButton.update(ms);
            this.editButton.update(ms);
            let onButtons = this.deleteButton.hover || this.editButton.hover;
            if (!this.dragging && ms.leftDown && !ms.draggingNode && !ms.draggingConnection && !onButtons){
                this.dragging = true;
                ms.draggingNode = true;
                this.dragOrigin.x = this.x - ms.world.x;
                this.dragOrigin.y = this.y - ms.world.y;
            }
        } else {
            this.deleteButton.hover = false;
            this.editButton.hover = false;
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
        ctx.fillStyle = "#898989";
        let typeX = ms.offset.x + this.x + this.width / 2 - this.typeWidth / 2;
        let typeY = ms.offset.y + this.y + this.height - 3;
        ctx.fillText(this.type, typeX, typeY);

        let resultCount = `${this.results.length}`
        let resultCountSize = ctx.measureText(resultCount);
        let resultCountWidth = resultCountSize.width;
        let resultCountHeight = resultCountSize.actualBoundingBoxAscent + resultCountSize.actualBoundingBoxDescent;
        let resultCountX = ms.offset.x + this.x + this.width - resultCountWidth - 3 * 3;
        let resultCountY = ms.offset.y + this.y + resultCountHeight + 3 * 3;
        ctx.fillText(resultCount, resultCountX, resultCountY)

        this.deleteButton.x = ms.offset.x + this.x;
        this.deleteButton.y = ms.offset.y + this.y + this.height - this.deleteButton.height;
        this.deleteButton.draw(ctx, ms);

        this.editButton.x = ms.offset.x + this.x + this.width - this.editButton.width;
        this.editButton.y = ms.offset.y + this.y + this.height - this.editButton.height;
        this.editButton.draw(ctx, ms);
        
        this.input.draw(ctx, ms);
        this.output.draw(ctx, ms);

        if(this.logs.length > 0){
            ctx.moveTo(ms.offset.x + this.x + 21, ms.offset.y + this.y + 6);
            ctx.fillStyle = "orange";
            ctx.beginPath();
            ctx.lineTo(ms.offset.x + this.x + 23, ms.offset.y + this.y + 21);
            ctx.lineTo(ms.offset.x + this.x + 6, ms.offset.y + this.y + 21);
            ctx.lineTo(ms.offset.x + this.x + 14, ms.offset.y + this.y + 6);
            ctx.fill();            
        }
        
        ctx.strokeStyle = "#8E8E8E";
        ctx.lineWidth = 3;
        ctx.strokeRect(ms.offset.x + this.x, ms.offset.y + this.y, this.width, this.height);
    }

    fixType() {
        // @ts-ignore
        this.type = this.meta.type
        if (["math", "condition"].indexOf(this.type) >= 0 ){
            // @ts-ignore
            this.type = this.meta.var1
        }
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
        
        this.width = Math.max(130, this.labelWidth * 1.5, this.typeWidth * 1.2);
    }

    pointInObject(p: Point): boolean {
        return this.pointNearNode(p) && (super.pointInObject(p) || this.input.pointInObject(p) || this.output.pointInObject(p));
    }


    pointNearNode(p: Point){
        // including the input/output circles
        if (p.x < this.x - this.input.radius){
            return false;
        }
        if (p.y < this.y){
            return false;
        }
        if (p.x > this.x + this.width + this.output.radius){
            return false;
        }
        if (p.y > this.y + this.height) {
            return false;
        }
        return true;
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
function diagramOnWheel(ev: WheelEvent){
    _diagram.onwheel(ev);   
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
    click: boolean = true;
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

    newConnection: NewConnection | null = null;

    scaleLevel: number = 0;
    scaleMax: number = 3;
    scaleMin: number = -1;
    get scale(): number {return 1 - (1 / this.scaleMax) * this.scaleLevel;}

    editNodeCallback: (node: DiagramNode) => void = function (){};
    deleteNodeCallback: (node: DiagramNode) => void = function (){};

    constructor(
            canvasId: string, 
            editNodeCallback: (node: DiagramNode) => void = function (){},
            deleteNodeCallback: (node: DiagramNode) => void = function (){},
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
        this.canvas.onwheel = diagramOnWheel;
        window.onresize = diagramOnResize;
        tick();
    }
    
    tick(){
        this.drawBackground();
        if (this.mouseState.leftUp && !this.mouseState.panning && !this.mouseState.draggingNode && !this.mouseState.draggingConnection){
            this.mouseState.click = true;
        }
        for (let node of this.nodes.values()){
            node.update(this.mouseState);
        }
        for (let connection of this.connections){
            connection.update(this.mouseState);
        }
        if (this.newConnection != null){
            this.newConnection.update(this.mouseState);
        }
        for (let connection of this.connections){
            connection.draw(this.ctx, this.mouseState);
        }
        if (this.newConnection != null){
            this.newConnection.draw(this.ctx, this.mouseState);
        }
        for (let node of this.nodes.values()){
            node.draw(this.ctx, this.mouseState);
        }
        this.drawWarning();
        this.mouseState.leftUp = false;
        this.mouseState.click = false;
    }

    onmousemove(ev: MouseEvent){
        let canvasRect = this.canvas.getBoundingClientRect();
        let scale = 1 / this.scale;
        this.mouseState.canvas.x = (ev.x - canvasRect.left) * scale;
        this.mouseState.canvas.y = (ev.y - canvasRect.top) * scale;
        this.mouseState.delta.x = ev.movementX * scale;
        this.mouseState.delta.y = ev.movementY * scale;

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
        if (this.newConnection != null){
            if (this.newConnection.input != null){
                this.addConnection(this.newConnection.output, this.newConnection.input);
            }
            this.mouseState.draggingConnection = false;
        }
        this.newConnection = null;
    }

    onwheel(ev: WheelEvent) {
        ev.preventDefault();
        let sign = Math.sign(ev.deltaY);
        let zoomOut = sign > 0;
        if (zoomOut && this.scaleLevel >= this.scaleMax-1) {
            return;
        } 
        let zoomIn = !zoomOut
        if (zoomIn && this.scaleLevel <= this.scaleMin) {
            return;
        }
        
        // undo previous scaling
        let currentScale = this.scale;
        let unscale = 1 / currentScale;
        this.ctx.scale(unscale, unscale);
        
        this.scaleLevel += sign;
        
        // scale with new value
        let scale = this.scale;
        this.ctx.scale(scale, scale);
    }

    drawBackground(){
        this.ctx.fillStyle = "#D8D8D8";
        let scale = 1 / this.scale;
        this.ctx.fillRect(0,0,this.canvas.width * scale, this.canvas.height * scale);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5 * scale;
        this.ctx.strokeRect(0, 0, this.canvas.width * scale, this.canvas.height * scale);
    }

    drawWarning(){
        let nodeWithLogs: DiagramNode | null = null;
        for (let node of this.nodes.values()){
            if (node.logs.length > 0) {
                nodeWithLogs = node;
                break;
            }
        }
        if (nodeWithLogs == null){
            return
        }
        let warningString = `Check log of '${nodeWithLogs.label}' Filter!`
        this.ctx.font = "30px Helvetica";
        let warningSize = this.ctx.measureText(warningString);

        this.ctx.fillStyle = "orange";
        this.ctx.fillRect(this.canvas.width - warningSize.width - 30, 0, warningSize.width + 30, 50);
        this.ctx.fillStyle = "#000";
        this.ctx.fillText(warningString, this.canvas.width - warningSize.width - 15, 35)

    }

    addNode(
            id: number, 
            x: number, 
            y: number, 
            label: string, 
            meta: Object = {}, 
            results: Array<string> = new Array(),
            logs: Array<string> = new Array()
        ){
        let node = new DiagramNode(id, x, y, label, meta, this.ctx, results, logs);
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
    return new Point(
      Math.pow(1-t,3) * sx + 3 * t * Math.pow(1 - t, 2) * cp1x 
        + 3 * t * t * (1 - t) * cp2x + t * t * t * ex,
      Math.pow(1-t,3) * sy + 3 * t * Math.pow(1 - t, 2) * cp1y 
        + 3 * t * t * (1 - t) * cp2y + t * t * t * ey
    );
  }