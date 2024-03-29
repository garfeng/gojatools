import gocv from "./js/gocv.io.x.gocv"
import image from "./js/image";
import color from "./js/image.color";

let mat = gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3);

let rect = image.Rect(30, 30, 70, 70);
// or:
/*
let rect = image._newRectangle({
    Min: {
        X: 30, Y: 30
    }, 
    Max : {
        X: 70, Y: 70
    }
})
*/

// Color Red
let c = color._newRGBA({
    R: 255,
    A: 255,
});

// draw rectangle
gocv.Rectangle(mat, rect, c, 10);

// save image
gocv.IMWrite("tmp.png", mat);

// show image
let w = gocv.NewWindow("gocv test")
w.IMShow(mat)
w.WaitKey(-1);

let matAddr = addr(mat)
matAddr.Close();
