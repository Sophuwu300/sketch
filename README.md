## Sketch
This program allows you to view images in your terminal.

### Features
- View images in your terminal
- Show as many images as you want at once
- Images always scale to fit the terminal size
- Supports png and jpg
- Adjustable scaling with flags

### Usage
```bash
sketch </path/to/images>
# or for 60% size
sketch -6 </path/to/images>
# you can even change size in the middle of the images
sketch cat.png -4 dog.jpg
# will output cat.png with 100% scale and dog.jpg at 40% size
```

### Installation Instructions
1. clone it
```bash
git clone <url ^ that way>
cd sketch
```
2. go build it
```bash
go build -o sketch -ldflags "-s -w" -trimpath .
```
3. install it
```bash
chmod +x sketch
sudo install sketch /usr/local/bin/sketch
```