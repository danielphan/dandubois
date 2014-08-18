/** @jsx React.DOM */
"use strict";

var _ = require("underscore");
var $ = require("jquery");
var App = require("./app/App.jsx");
var React = require("react");

_.once(function() {
  require("es5-shim");
  require("es5-shim/es5-sham.js");

  if ($("html").is(".ie6, .ie7, .ie8")) {
    require("html5shiv", function(shiv) {});
  }

  window.React = React;
})();

$(document).ready(function() {
  $("head").append(
    $("<link></link>")
      .attr("rel", "shortcut icon")
      .attr("href", require("./favicon.ico"))
  );
  React.renderComponent(<App />, $(".App").get(0));
});
