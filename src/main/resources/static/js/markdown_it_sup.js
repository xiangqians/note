/*! markdown-it-sup 2.0.0 https://github.com/markdown-it/markdown-it-sup @license MIT */
!function(e,s){"object"==typeof exports&&"undefined"!=typeof module?module.exports=s():"function"==typeof define&&define.amd?define(s):(e="undefined"!=typeof globalThis?globalThis:e||self).markdownitSup=s()}(this,(function(){"use strict";const e=/\\([ \\!"#$%&'()*+,./:;<=>?@[\]^_`{|}~-])/g;function s(s,o){const n=s.posMax,p=s.pos;if(94!==s.src.charCodeAt(p))return!1;if(o)return!1;if(p+2>=n)return!1;s.pos=p+1;let t=!1;for(;s.pos<n;){if(94===s.src.charCodeAt(s.pos)){t=!0;break}s.md.inline.skipToken(s)}if(!t||p+1===s.pos)return s.pos=p,!1;const r=s.src.slice(p+1,s.pos);if(r.match(/(^|[^\\])(\\\\)*\s/))return s.pos=p,!1;s.posMax=s.pos,s.pos=p+1;s.push("sup_open","sup",1).markup="^";s.push("text","",0).content=r.replace(e,"$1");return s.push("sup_close","sup",-1).markup="^",s.pos=s.posMax+1,s.posMax=n,!0}return function(e){e.inline.ruler.after("emphasis","sup",s)}}));
