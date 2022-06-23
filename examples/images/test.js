import {image} from "./js/image"
import {png} from "./js/image/png"
import {bytes} from "./js/bytes";
import {ioutil} from "./js/io/ioutil";
import {os} from "./js/os";

const size = image.Rect(0, 0, 100, 100);

let src = image.NewRGBA(size);

for (let j = size.Min.Y; j < size.Max.Y; j++) {
    for (let i = size.Min.X; i < size.Max.X; i++) {
        if (i > 40 && i < 60 && j > 40 && j < 60) {
            src.SetRGBA(i, j, {
                R: 255,
                A: 255,
            });
        } else {
            src.SetRGBA(i,j, {
                A: 255
            });
        }
    }
}

let w = bytes.NewBuffer(null);
png.Encode(w, src)

try {
    os.Mkdir("./tmp", 0775);
} catch (error) {
    
}

let err = ioutil.WriteFile("./tmp/test.png", w.Bytes(), 0664);

if (err != null) {
    console.log(err);
}

