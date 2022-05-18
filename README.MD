#### Description
A tool used to convert standard text into 'art',
using only characters available in the ASCII table.

#### Usage
In order to use this tool on your own machine, please clone the repository and run main.go.
After running main.go and allowing firewall privileges, please head to http://localhost:8080/
where you will be able to generate your text using the different banners available.

#### Implementation Details
The frontend is written in HTML/CSS, backend is written in Go.
The user enters a text/banner in the HTML form, and that data is sent via POST to /ascii-art, 
where we pass it to a modified version of GenerateArt() from my original ASCII Art project. 
The modification I have made is so that it generates a text file instead of a slice of strings,
with the output (the actual ASCII art end result) being stored inside that file.

The HTML file makes use of the \<pre> tag, which allows us to then display the file contents to the user
after they click the 'Generate Art' button on the web page. 
