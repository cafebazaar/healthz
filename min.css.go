package healthz

import "net/http"

var minCSS = []byte(`/* Copyright 2014 Owen Versteeg; MIT licensed */
body,textarea,input,select{background:0;border-radius:0;font:16px sans-serif;margin:0}.addon,.btn-sm,.nav,textarea,input,select{outline:0;font-size:14px}.smooth{transition:all .2s}.btn,.nav a{text-decoration:none}.container{margin:0 20px;width:auto}@media(min-width:1310px){.container{margin:auto;width:1270px}}.btn,h2{font-size:2em}h1{font-size:3em}.row{margin:1% 0;overflow:auto}.col{float:left}.table,.c12{width:100%}.c11{width:91.66%}.c10{width:83.33%}.c9{width:75%}.c8{width:66.66%}.c7{width:58.33%}.c6{width:50%}.c5{width:41.66%}.c4{width:33.33%}.c3{width:25%}.c2{width:16.66%}.c1{width:8.33%}@media(max-width:870px){.row .col{width:100%}}.msg{background:#def;border-left:5px solid #59d;padding:1.5em}

.overall {padding: 15px 20px 5px;border-radius: 15px;margin: 1em 0;}
.row {padding: 5px 20px;border-radius: 5px;margin: 0.2em 0;}
.service,.hostname{font-family:monospace;background: #ecf0eb; color: #cc0000; padding: 2px 5px}
.hl2{background: #73d216}
.hl1{background: #729fcf}
.hl0{background: #d3d7cf}
.hl-1{background: #fcaf3e}
.hl-2{background: #ef2929}
`)

func (h *handler) reportMinCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "max-age:2592000, public")
	w.Write(minCSS)
}
