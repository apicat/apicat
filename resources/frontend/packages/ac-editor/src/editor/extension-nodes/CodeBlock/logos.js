import "codemirror/mode/javascript/javascript"
import "codemirror/mode/php/php"
import "codemirror/mode/xml/xml"
import "codemirror/mode/htmlmixed/htmlmixed"
import "codemirror/mode/css/css"
import "codemirror/mode/sass/sass"

export const logos = {
  none: null,
  json: {name:"javascript",json:true},
  html: "htmlmixed",
  xml: "xml",
  javascript: "javascript",
  css: "css",
  scss: "sass",
  sass: "sass",
  php: "php"
};
export default Object.keys(logos);
