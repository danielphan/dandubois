"use strict";

var _ = require("underscore");
var jQuery = require("jquery");
var React = require("react");

require("./main.less");

jQuery(document).ready(function() {
  jQuery("head").append(
    jQuery("<link></link>")
      .attr("rel", "shortcut icon")
      .attr("href", require("./favicon.ico"))
  );

  jQuery.get(process.env.API_URL + "/paintings", function(data) {
    console.log(data);
  });
});
