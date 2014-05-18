/** @jsx React.DOM */
'use strict'

var
	Page = require('../page/Page'),
	React = require('react')
;

var PaintingPage = React.createClass({
	render: function() {
		return <Page>
			<div className='PaintingPage'>
				{this.props.id}
			</div>
		</Page>;
	}
});

module.exports = PaintingPage;
