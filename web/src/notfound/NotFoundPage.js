/** @jsx React.DOM */
'use strict'

var
	React = require('react')
;

var NotFoundPage = React.createClass({
	render: function() {
		return <div className='NotFoundPage'>404 Not Found</div>;
	}
});

module.exports = NotFoundPage;
