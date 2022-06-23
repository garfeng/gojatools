goja_imports -p js -i os -go js/os.go -ja js/os.ts
goja_imports -p js -i io/ioutil -go js/io.ioutil.go -ja js/io.ioutil.ts

goja_imports -p js -i image -go js/image.go -ja js/image.ts
goja_imports -p js -i image/png -go js/image.png.go -ja js/image.png.ts
goja_imports -p js -i image/color -go js/image.color.go -ja js/image.color.ts
goja_imports -p js -i bytes -go js/bytes.go -ja js/bytes.ts

