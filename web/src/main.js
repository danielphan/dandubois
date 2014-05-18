/** @jsx React.DOM */
'use strict'

var
	App = require('./app'),
	jQuery = require('jquery'),
	React = require('react'),
	underscore = require('underscore')
;

underscore.once(function() {
	window.React = React;
})();

jQuery(document).ready(function() {
	React.renderComponent(<App />, document.getElementsByClassName('App')[0]);
});
