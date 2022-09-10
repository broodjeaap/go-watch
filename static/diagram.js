var DiagramNode = /** @class */ (function () {
    function DiagramNode(x, y, width, height, label) {
        this.x = x;
        this.y = y;
        this.width = width;
        this.height = height;
        this.label = label;
    }
    DiagramNode.prototype.pointInDiagram = function (x, y) {
        if (x < this.x) {
            return false;
        }
        if (y < this.y) {
            return false;
        }
        if (x > this.x + this.width) {
            return false;
        }
        if (y > this.y + this.height) {
            return false;
        }
        return true;
    };
    return DiagramNode;
}());
var _diagram;
function diargramOnResize() {
    _diagram.fillParent();
}
function diagramOnMouseDown(ev) {
    _diagram.onmousedown(ev);
}
function diagramOnMouseUp(ev) {
    _diagram.onmouseup(ev);
}
function diagramOnMouseMove(ev) {
    _diagram.onmousemove(ev);
}
var Diagrams = /** @class */ (function () {
    function Diagrams(canvasId) {
        this.canvas = document.getElementById(canvasId);
        if (this.canvas === null) {
            throw "Could not getElementById " + canvasId;
        }
        var ctx = this.canvas.getContext("2d");
        if (ctx === null) {
            throw "Could not get 2d rendering context";
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
    Diagrams.prototype.onmousemove = function (ev) {
        if (this.panning) {
            this.cameraX += ev.movementX;
            this.cameraY += ev.movementY;
        }
        if (this.nodeDrag) {
            if (this.nodeDragged === null) {
                console.error("nodeDrag==true but nodeDragged==null");
                return;
            }
            this.nodeDragged.x += ev.movementX;
            this.nodeDragged.y += ev.movementY;
        }
        this.draw();
    };
    Diagrams.prototype.onmousedown = function (ev) {
        if (ev.button != 0) {
            return;
        }
        var canvasRect = this.canvas.getBoundingClientRect();
        var mouseX = ev.x - canvasRect.left;
        var mouseY = ev.y - canvasRect.top;
        for (var _i = 0, _a = this.nodes; _i < _a.length; _i++) {
            var node = _a[_i];
            if (node.pointInDiagram(mouseX, mouseY)) {
                this.nodeDrag = true;
                this.nodeDragged = node;
                return;
            }
        }
        this.panning = true;
    };
    Diagrams.prototype.onmouseup = function (ev) {
        this.panning = false;
        this.nodeDrag = false;
        this.nodeDragged = null;
    };
    Diagrams.prototype.drawBackground = function () {
        this.ctx.fillStyle = "#D8D8D8";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.strokeStyle = "#888";
        this.ctx.lineWidth = 5;
        this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
    };
    Diagrams.prototype.draw = function () {
        this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
        this.drawBackground();
        for (var _i = 0, _a = this.nodes; _i < _a.length; _i++) {
            var node = _a[_i];
            this.ctx.fillStyle = "gray";
            this.ctx.fillRect(node.x + this.cameraX, node.y + this.cameraY, node.width, node.height);
            this.ctx.fillStyle = "black";
            this.ctx.font = "30px Arial";
            this.ctx.fillText(node.label, node.x + this.cameraX + node.height / 2, node.y + this.cameraY + node.height / 1.5);
        }
    };
    Diagrams.prototype.addNode = function (x, y, label) {
        var textSize = this.ctx.measureText(label);
        var textHeight = 2 * (textSize.actualBoundingBoxAscent + textSize.actualBoundingBoxDescent);
        this.nodes.push(new DiagramNode(x, y, textSize.width + textHeight, textHeight, label));
    };
    Diagrams.prototype.addConnection = function (A, B) {
        this.connections.push([A, B]);
    };
    Diagrams.prototype.fillParent = function () {
        this.canvas.width = this.canvas.clientWidth;
        this.canvas.height = this.canvas.clientHeight;
        this.draw();
    };
    return Diagrams;
}());
