### Introduction

This project is to generate a pdf file from images which is similar to Avery label generator.
All you have to do it define a label page settings and image information in JSON.


# Page Settings

Define your label paper format in JSON, say, page.json,

```

[
{
	"model" : "EU30147WJ",
	"description" : "onlinelabels.com EU30147WJ 70mm x 50mm",
	"paper" : { "name" : "A4", "unit" : "mm", "width": 210, "height": 297, "orientation": "P", 
  		    "top": 14, "bottom": 13, "left": 33, "right": 32, "columns" : 2, "rows" : 5,
		    "label_width": 70, "label_height": 50, "hspace": 5, "vspace": 5}
},
{
        "model" : "TAOBAO70x50",
        "description" : "TAOBAO 70mm x 50mm",
        "paper" : { "name" : "A4", "unit" : "mm", "width": 210, "height": 297, "orientation": "P",
        	    "columns" : 3, "rows" : 5, "top": 24, "bottom": 20, "left": 0, "right": 0,
		    "label_width": 70, "label_height": 50, "hspace": 1, "vspace": 1}
}
]

```

# Image Settings in JSON

Define the images you want to print in JSON format (e.g. label.json).  Image path with empty string will be skipped and leave a blank space on the pdf.

```
{
	"model": "EU30147WJ",
	"image_type": "PNG",
	"images" : [["image1", "image2"], 
			["image3", "image4"],
			["image5", "image6"],
			["image7", "image8"],
			["image9", "image19"]
		]
}

```

# Create the pdf file by running the command line ```label2pdf```

```
% label2pdf page.json label.json out.pdf
```
