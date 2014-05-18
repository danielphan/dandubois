/** @jsx React.DOM */
'use strict'

var
	ContactPage = require('./contact/ContactPage'),
	NotFoundPage = require('./notfound/NotFoundPage'),
	PaintingPage = require('./painting/PaintingPage'),
	PaintingsPage = require('./painting/PaintingsPage'),
	React = require('react'),
	Router = require('react-router-component')
;

var
	Locations = Router.Locations,
	Location = Router.Location,
	NotFound = Router.NotFound
;

var App = React.createClass({
	render: function() {
		return <Locations>
			<Location path='/' handler={PaintingsPage} />
			<Location path='/paintings/?' handler={PaintingsPage} />
			<Location path='/paintings/:id' handler={PaintingPage} />
			<Location path='/contact' handler={ContactPage} />
			<NotFound handler={NotFoundPage} />
		</Locations>;
	}
});

module.exports = App;
